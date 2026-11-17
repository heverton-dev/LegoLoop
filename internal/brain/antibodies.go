package brain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ErrAnticorposCorrompidos sinaliza anticorpos.json inválido no carregamento (AC-E02.2).
var ErrAnticorposCorrompidos = errors.New("brain: anticorpos.json corrompido")

// ErrAnticorpoInvalido sinaliza errSig ou pacote vazio no registro (AC-E02.4).
var ErrAnticorpoInvalido = errors.New("brain: anticorpo inválido")

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

// AntibodiesDB armazena e consulta anticorpos imunológicos para injeção preventiva.
type AntibodiesDB struct {
	mu         sync.RWMutex
	antibodies map[string][]Antibody
	caminho    string // anticorpos.json; vazio desliga a persistência
}

// NewAntibodiesDB instancia um banco de anticorpos em memória.
func NewAntibodiesDB() *AntibodiesDB {
	return &AntibodiesDB{
		antibodies: make(map[string][]Antibody),
	}
}

// NewAntibodiesDBPersistente instancia o banco em dir com persistência em anticorpos.json.
// Carrega o arquivo ao iniciar; JSON inválido devolve ErrAnticorposCorrompidos sem sobrescrever.
func NewAntibodiesDBPersistente(dir string) (*AntibodiesDB, error) {
	db := NewAntibodiesDB()
	db.caminho = filepath.Join(dir, "anticorpos.json")

	raw, err := os.ReadFile(db.caminho)
	if errors.Is(err, os.ErrNotExist) {
		return db, nil
	}
	if err != nil {
		return nil, fmt.Errorf("brain: lendo anticorpos: %w", err)
	}
	if err := json.Unmarshal(raw, &db.antibodies); err != nil {
		return nil, ErrAnticorposCorrompidos
	}
	return db, nil
}

// Record grava uma tripla imunológica associada a um pacote.
// Reincidência do errSig deduplica e vacina com antídoto automático (AC-E02.3).
func (db *AntibodiesDB) Record(pkg, errSig, rootCause, diffSol string) error {
	if strings.TrimSpace(pkg) == "" || strings.TrimSpace(errSig) == "" {
		return ErrAnticorpoInvalido
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	lista := db.antibodies[pkg]
	for i := range lista {
		if strings.EqualFold(strings.TrimSpace(lista[i].ErrorSignature), strings.TrimSpace(errSig)) {
			lista[i].Recorrencias++
			lista[i].AntidotoAplicado = true
			db.antibodies[pkg] = lista
			return db.persistir()
		}
	}

	ab := Antibody{
		Package:        pkg,
		ErrorSignature: errSig,
		RootCause:      rootCause,
		DiffSolution:   diffSol,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Recorrencias:   1,
	}
	db.antibodies[pkg] = append(lista, ab)
	return db.persistir()
}

// persistir grava anticorpos.json; exige lock de escrita já retido.
func (db *AntibodiesDB) persistir() error {
	if db.caminho == "" {
		return nil
	}
	raw, err := json.MarshalIndent(db.antibodies, "", "  ")
	if err != nil {
		return fmt.Errorf("brain: serializando anticorpos: %w", err)
	}
	if err := os.WriteFile(db.caminho, raw, 0o644); err != nil {
		return fmt.Errorf("brain: gravando anticorpos: %w", err)
	}
	return nil
}

// Find localiza um anticorpo pela assinatura do erro dentro de um pacote.
func (db *AntibodiesDB) Find(pkg, errSig string) (*Antibody, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, ab := range db.antibodies[pkg] {
		if strings.EqualFold(strings.TrimSpace(ab.ErrorSignature), strings.TrimSpace(errSig)) {
			return &ab, true
		}
	}
	return nil, false
}

// AntidotoAplicado devolve o antídoto do errSig quando a vacinação de reincidência ocorreu (AC-E02.3).
func (db *AntibodiesDB) AntidotoAplicado(pkg, errSig string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, ab := range db.antibodies[pkg] {
		if strings.EqualFold(strings.TrimSpace(ab.ErrorSignature), strings.TrimSpace(errSig)) && ab.AntidotoAplicado {
			return ab.DiffSolution, true
		}
	}
	return "", false
}

// InjectPreventiveContext gera um texto de advertência preventiva com os anticorpos do pacote.
func (db *AntibodiesDB) InjectPreventiveContext(pkg string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := db.antibodies[pkg]
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
