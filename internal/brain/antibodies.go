package brain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"legoloop/internal/escrita"
)

// AntibodiesDBFile é o nome do arquivo persistido por processo (schema anticorpos.v1).
const AntibodiesDBFile = "anticorpos.json"

// SchemaAnticorpos é a versão de schema do banco persistido (fonte: schemas/anticorpos.v1.json).
const SchemaAnticorpos = "anticorpos.v1"

// sufixoBackup marca a cópia de segurança criada ao detectar banco corrompido.
const sufixoBackup = ".legoloop-bak"

// ErrAnticorposCorrompidos sinaliza anticorpos.json inválido no carregamento (AC-E02.2).
var ErrAnticorposCorrompidos = errors.New("brain: anticorpos.json corrompido")

// ErrAnticorpoInvalido sinaliza errSig ou pacote vazio no registro (AC-E02.4).
var ErrAnticorpoInvalido = errors.New("brain: anticorpo inválido")

// ErrGateAntiReincidencia bloqueia erro conhecido sem antídoto aplicável (AC-E02.7).
var ErrGateAntiReincidencia = errors.New("brain: erro conhecido sem antídoto aplicável")

// Escopo define onde a vacina é persistida (AC-E04.3).
type Escopo int

const (
	// EscopoDominio guarda regras de domínio locais do processo.
	EscopoDominio Escopo = iota
	// EscopoGlobal guarda falhas de infraestrutura/LLM compartilhadas entre processos.
	EscopoGlobal
)

// AcaoGate descreve a decisão do gate anti-reincidência (AC-E02.7/AC-E02.8).
type AcaoGate int

const (
	// AcaoNormal: erro desconhecido, execução segue sem intervenção.
	AcaoNormal AcaoGate = iota
	// AcaoAntidoto: erro conhecido com antídoto aplicável.
	AcaoAntidoto
	// AcaoBloqueado: erro conhecido sem antídoto aplicável, execução bloqueia.
	AcaoBloqueado
)

// Antibody captura a tripla imunológica de um erro resolvido (Motor III).
type Antibody struct {
	Package          string `json:"package"`
	ErrorSignature   string `json:"error_signature"`
	RootCause        string `json:"root_cause"`
	DiffSolution     string `json:"diff_solution"`
	Timestamp        string `json:"timestamp"`
	Recorrencias     int    `json:"recorrencias"`      // 1 no registro; +1 a cada reincidência
	AntidotoAplicado bool   `json:"antidoto_aplicado"` // true após vacinação de reincidência
}

// vacinaJSON é a forma persistida da vacina no schema anticorpos.v1.
type vacinaJSON struct {
	Assinatura       string `json:"assinatura"`
	Erro             string `json:"erro"`
	Antidoto         string `json:"antidoto"`
	Recorrencias     int    `json:"recorrencias"`
	Pacote           string `json:"pacote,omitempty"`
	CausaRaiz        string `json:"causa_raiz,omitempty"`
	Timestamp        string `json:"timestamp,omitempty"`
	AntidotoAplicado bool   `json:"antidoto_aplicado,omitempty"`
}

// bancoJSON é o envelope persistido em anticorpos.json (schema anticorpos.v1).
type bancoJSON struct {
	SchemaVersion string       `json:"schema_version"`
	Processo      string       `json:"processo"`
	Globais       []vacinaJSON `json:"globais"`
	Dominio       []vacinaJSON `json:"dominio"`
}

// AntibodiesDB armazena e consulta anticorpos imunológicos para injeção preventiva.
type AntibodiesDB struct {
	mu            sync.RWMutex
	globais       []Antibody
	dominio       []Antibody
	caminho       string // anticorpos.json do processo; vazio desliga a persistência
	caminhoGlobal string // anticorpos.json compartilhado (globais); vazio = só local
	processo      string
}

// NewAntibodiesDB instancia um banco de anticorpos em memória.
func NewAntibodiesDB() *AntibodiesDB {
	return &AntibodiesDB{}
}

// NewAntibodiesDBPersistente instancia o banco em dir com persistência em anticorpos.json.
// Carrega o arquivo ao iniciar; JSON inválido devolve ErrAnticorposCorrompidos com backup
// .legoloop-bak e banco vazio funcional (o processo segue sem vacina, AC-E04.4).
func NewAntibodiesDBPersistente(dir string) (*AntibodiesDB, error) {
	return NewAntibodiesDBPersistenteComGlobal(dir, "")
}

// NewAntibodiesDBPersistenteComGlobal instancia o banco do processo com globais
// compartilhadas via dirGlobal (AC-E04.3).
func NewAntibodiesDBPersistenteComGlobal(dir, dirGlobal string) (*AntibodiesDB, error) {
	db := &AntibodiesDB{
		caminho:  filepath.Join(dir, AntibodiesDBFile),
		processo: filepath.Base(dir),
	}
	if dirGlobal != "" {
		db.caminhoGlobal = filepath.Join(dirGlobal, AntibodiesDBFile)
	}
	globais, dominio, err := carregarBanco(db.caminho)
	if err != nil {
		return db, err
	}
	db.globais, db.dominio = globais, dominio
	db.fundirGlobais()
	return db, nil
}

// fundirGlobais absorve as vacinas globais compartilhadas entre processos.
func (db *AntibodiesDB) fundirGlobais() {
	if db.caminhoGlobal == "" {
		return
	}
	globais, _, err := carregarBanco(db.caminhoGlobal)
	if err != nil {
		return
	}
	for _, g := range globais {
		if _, ok := db.procurar(g.Package, g.ErrorSignature); !ok {
			db.globais = append(db.globais, g)
		}
	}
}

// carregarBanco interpreta o formato anticorpos.v1 ou o legado (mapa pacote→lista).
func carregarBanco(caminho string) ([]Antibody, []Antibody, error) {
	raw, err := os.ReadFile(caminho)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("brain: lendo anticorpos: %w", err)
	}
	return parseBanco(caminho, raw)
}

// parseBanco decodifica o banco no schema v1 ou migra o formato legado sem perda.
func parseBanco(caminho string, raw []byte) ([]Antibody, []Antibody, error) {
	var campos map[string]json.RawMessage
	if err := json.Unmarshal(raw, &campos); err != nil {
		return nil, nil, backupAnticorposCorrompidos(caminho, raw)
	}
	if _, ok := campos["schema_version"]; ok {
		return parseBancoV1(caminho, raw)
	}
	var legado map[string][]Antibody
	if err := json.Unmarshal(raw, &legado); err != nil {
		return nil, nil, backupAnticorposCorrompidos(caminho, raw)
	}
	var dominio []Antibody
	for _, lista := range legado {
		dominio = append(dominio, lista...)
	}
	return nil, dominio, nil
}

// parseBancoV1 decodifica o envelope anticorpos.v1 em vacinas de domínio e globais.
func parseBancoV1(caminho string, raw []byte) ([]Antibody, []Antibody, error) {
	var banco bancoJSON
	if err := json.Unmarshal(raw, &banco); err != nil {
		return nil, nil, backupAnticorposCorrompidos(caminho, raw)
	}
	return paraAntibodies(banco.Globais), paraAntibodies(banco.Dominio), nil
}

// backupAnticorposCorrompidos preserva o conteúdo inválido e devolve a sentinela.
func backupAnticorposCorrompidos(caminho string, raw []byte) error {
	_ = escrita.GravarAtomico(caminho+sufixoBackup, raw)
	return ErrAnticorposCorrompidos
}

// paraAntibodies converte as vacinas persistidas para o modelo interno.
func paraAntibodies(lista []vacinaJSON) []Antibody {
	saida := make([]Antibody, 0, len(lista))
	for _, v := range lista {
		saida = append(saida, Antibody{
			Package:          v.Pacote,
			ErrorSignature:   v.Erro,
			RootCause:        v.CausaRaiz,
			DiffSolution:     v.Antidoto,
			Timestamp:        v.Timestamp,
			Recorrencias:     v.Recorrencias,
			AntidotoAplicado: v.AntidotoAplicado,
		})
	}
	return saida
}

// paraVacinas converte o modelo interno para a forma persistida com assinatura SHA-256.
func paraVacinas(lista []Antibody) []vacinaJSON {
	saida := make([]vacinaJSON, 0, len(lista))
	for _, ab := range lista {
		saida = append(saida, vacinaJSON{
			Assinatura:       assinaturaErro(ab.ErrorSignature),
			Erro:             ab.ErrorSignature,
			Antidoto:         ab.DiffSolution,
			Recorrencias:     ab.Recorrencias,
			Pacote:           ab.Package,
			CausaRaiz:        ab.RootCause,
			Timestamp:        ab.Timestamp,
			AntidotoAplicado: ab.AntidotoAplicado,
		})
	}
	return saida
}

// assinaturaErro deriva a assinatura SHA-256 do erro para detecção de reincidência.
func assinaturaErro(errSig string) string {
	soma := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(errSig))))
	return "sha256:" + hex.EncodeToString(soma[:])
}

// mesmaAssinatura compara erros pela assinatura normalizada (AC-E02.3).
func mesmaAssinatura(a, b string) bool {
	return assinaturaErro(a) == assinaturaErro(b)
}

// Record grava uma tripla imunológica no escopo de domínio do processo.
// Reincidência do errSig deduplica e vacina com antídoto automático (AC-E02.3).
func (db *AntibodiesDB) Record(pkg, errSig, rootCause, diffSol string) error {
	return db.RecordComEscopo(EscopoDominio, pkg, errSig, rootCause, diffSol)
}

// RecordComEscopo grava a vacina no escopo pedido: domínio local ou global compartilhado.
func (db *AntibodiesDB) RecordComEscopo(escopo Escopo, pkg, errSig, rootCause, diffSol string) error {
	if strings.TrimSpace(pkg) == "" || strings.TrimSpace(errSig) == "" {
		return ErrAnticorpoInvalido
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	lista, nome := &db.dominio, "dominio"
	if escopo == EscopoGlobal {
		lista, nome = &db.globais, "globais"
	}
	for i := range *lista {
		if mesmaAssinatura((*lista)[i].ErrorSignature, errSig) {
			(*lista)[i].Recorrencias++
			(*lista)[i].AntidotoAplicado = true
			return db.persistir(nome)
		}
	}
	*lista = append(*lista, Antibody{
		Package:        pkg,
		ErrorSignature: errSig,
		RootCause:      rootCause,
		DiffSolution:   diffSol,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Recorrencias:   1,
	})
	return db.persistir(nome)
}

// persistir grava anticorpos.json no schema v1; exige lock de escrita já retido.
// O escopo alterado grava também no arquivo global compartilhado quando configurado.
func (db *AntibodiesDB) persistir(escopoAlterado string) error {
	if db.caminho == "" {
		return nil
	}
	if err := gravarBanco(db.caminho, db.processo, paraVacinas(db.globais), paraVacinas(db.dominio)); err != nil {
		return err
	}
	if db.caminhoGlobal != "" && escopoAlterado == "globais" {
		return gravarBanco(db.caminhoGlobal, "compartilhado", paraVacinas(db.globais), nil)
	}
	return nil
}

// gravarBanco serializa o envelope anticorpos.v1 e grava com escrita atômica.
func gravarBanco(caminho, processo string, globais, dominio []vacinaJSON) error {
	banco := bancoJSON{SchemaVersion: SchemaAnticorpos, Processo: processo, Globais: globais, Dominio: dominio}
	raw, err := json.MarshalIndent(banco, "", "  ")
	if err != nil {
		return fmt.Errorf("brain: serializando anticorpos: %w", err)
	}
	if err := escrita.GravarAtomico(caminho, raw); err != nil {
		return fmt.Errorf("brain: gravando anticorpos: %w", err)
	}
	caminhoMD := filepath.Join(filepath.Dir(caminho), AntibodiesMarkdownFile)
	md := RenderizarMarkdown(processo, globais, dominio)
	if err := escrita.GravarAtomico(caminhoMD, md); err != nil {
		return fmt.Errorf("brain: gravando %s: %w", AntibodiesMarkdownFile, err)
	}
	return nil
}

// procurar localiza a vacina: primeiro domínio do pacote, depois globais.
func (db *AntibodiesDB) procurar(pkg, errSig string) (*Antibody, bool) {
	if ab := procurarLista(db.dominio, pkg, errSig, true); ab != nil {
		return ab, true
	}
	if ab := procurarLista(db.globais, pkg, errSig, false); ab != nil {
		return ab, true
	}
	return nil, false
}

// procurarLista varre uma lista de vacinas pela assinatura do erro.
func procurarLista(lista []Antibody, pkg, errSig string, exigirPacote bool) *Antibody {
	for i := range lista {
		if exigirPacote && lista[i].Package != pkg {
			continue
		}
		if mesmaAssinatura(lista[i].ErrorSignature, errSig) {
			copia := lista[i]
			return &copia
		}
	}
	return nil
}

// Find localiza um anticorpo pela assinatura do erro dentro de um pacote.
func (db *AntibodiesDB) Find(pkg, errSig string) (*Antibody, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.procurar(pkg, errSig)
}

// AntidotoAplicado devolve o antídoto do errSig quando a vacinação de reincidência ocorreu (AC-E02.3).
func (db *AntibodiesDB) AntidotoAplicado(pkg, errSig string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	ab, ok := db.procurar(pkg, errSig)
	if !ok || !ab.AntidotoAplicado {
		return "", false
	}
	return ab.DiffSolution, true
}

// GateAntiReincidencia decide o destino de um erro de runtime (AC-E02.7/AC-E02.8):
// desconhecido segue normal; conhecido com antídoto aplica; sem antídoto bloqueia citando o errSig.
func (db *AntibodiesDB) GateAntiReincidencia(pkg, errSig string) (AcaoGate, string, error) {
	ab, ok := db.Find(pkg, errSig)
	if !ok {
		return AcaoNormal, "", nil
	}
	if ab.DiffSolution != "" {
		return AcaoAntidoto, ab.DiffSolution, nil
	}
	return AcaoBloqueado, "", fmt.Errorf("%w: %s", ErrGateAntiReincidencia, strings.TrimSpace(errSig))
}

// InjectPreventiveContext gera um texto de advertência preventiva com os anticorpos do pacote.
func (db *AntibodiesDB) InjectPreventiveContext(pkg string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := append(append([]Antibody{}, db.dominioFiltrado(pkg)...), db.globais...)
	if len(list) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n🛡️ [ANTICORPO PREVENTIVO // MEMÓRIA IMUNOLÓGICA]:\n")
	for _, ab := range list {
		sb.WriteString(fmt.Sprintf("· Causa Raiz Conhecida: %s\n· Solução Canônica: %s\n", ab.RootCause, ab.DiffSolution))
	}
	return sb.String()
}

// dominioFiltrado devolve as vacinas de domínio do pacote informado.
func (db *AntibodiesDB) dominioFiltrado(pkg string) []Antibody {
	var saida []Antibody
	for _, ab := range db.dominio {
		if ab.Package == pkg {
			saida = append(saida, ab)
		}
	}
	return saida
}

// InjecaoPreventiva vira passo do processo: injeta o contexto preventivo antes do Nó 1 (AC-E02.5).
// Devolve os contextos de nó com a advertência prefixada ao primeiro. Sem anticorpo no pacote,
// injeta vazio e o contexto do Nó 1 segue inalterado (AC-E02.6).
func (db *AntibodiesDB) InjecaoPreventiva(pkg string, contextos []string) []string {
	saida := make([]string, len(contextos))
	copy(saida, contextos)
	if len(saida) == 0 {
		return saida
	}
	saida[0] = db.InjectPreventiveContext(pkg) + saida[0]
	return saida
}
