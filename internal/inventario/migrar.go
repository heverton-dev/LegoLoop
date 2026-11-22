package inventario

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
)

// ResultadoMigracao resume a movimentação para o arquivo morto.
type ResultadoMigracao struct {
	Movidos    []string // destinos relativos a docs/
	Protegidos []string // relatórios globais mantidos em docs/reports/
}

// Migrar move a documentação de construção de docs/ para docs/archive/.
// Move com rename. Nunca deleta arquivo (contrato AC-E01.2).
func Migrar(raiz string) (ResultadoMigracao, error) {
	var res ResultadoMigracao
	docs := filepath.Join(raiz, "docs")
	if err := os.MkdirAll(filepath.Join(docs, "archive"), 0o755); err != nil {
		return res, err
	}
	for _, dir := range []string{"videorea", "sessions"} {
		origem := filepath.Join(docs, dir)
		destino := filepath.Join(docs, "archive", dir)
		if err := moverDiretorio(origem, destino, docs, &res); err != nil {
			return res, err
		}
	}
	globais, err := RelatoriosGlobais(raiz)
	if err != nil {
		return res, err
	}
	return res, moverRelatorios(docs, globais, &res)
}

// moverDiretorio move todos os arquivos de origem preservando a estrutura.
func moverDiretorio(origem, destino, raizDocs string, res *ResultadoMigracao) error {
	if _, err := os.Stat(origem); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(origem, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, errRel := filepath.Rel(origem, p)
		if errRel != nil {
			return errRel
		}
		alvo := filepath.Join(destino, rel)
		if err := moverArquivo(p, alvo); err != nil {
			return err
		}
		registrarMovido(res, raizDocs, alvo)
		return nil
	})
}

// moverRelatorios move relatórios históricos e protege os relatórios globais.
func moverRelatorios(docs string, globais map[string]bool, res *ResultadoMigracao) error {
	origem := filepath.Join(docs, "reports")
	if _, err := os.Stat(origem); os.IsNotExist(err) {
		return nil
	}
	destino := filepath.Join(docs, "archive", "reports")
	return filepath.WalkDir(origem, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, errRel := filepath.Rel(origem, p)
		if errRel != nil {
			return errRel
		}
		if globais[filepath.Base(p)] {
			res.Protegidos = append(res.Protegidos, filepath.ToSlash(filepath.Join("reports", rel)))
			return nil
		}
		alvo := filepath.Join(destino, rel)
		if err := moverArquivo(p, alvo); err != nil {
			return err
		}
		registrarMovido(res, docs, alvo)
		return nil
	})
}

// moverArquivo move um arquivo com rename. Nunca deleta.
func moverArquivo(origem, destino string) error {
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	return os.Rename(origem, destino)
}

// registrarMovido registra o destino relativo a docs/.
func registrarMovido(res *ResultadoMigracao, raizDocs, destino string) {
	rel, err := filepath.Rel(raizDocs, destino)
	if err != nil {
		rel = destino
	}
	res.Movidos = append(res.Movidos, filepath.ToSlash(rel))
}

// tiposQuinteto são as pastas de documento do Quinteto saneadas por MigrarQuinteto.
var tiposQuinteto = []string{"intents", "specs", "plans", "routes"}

// reDataTipo separa o nome-base do documento do prefixo de data e tipo.
var reDataTipo = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_(INTENT|SPEC|PLAN|ROUTE|DAG|INTERVIEW|REPORT|MANUAL)_`)

// EntradaOrfao registra um documento do Quinteto sem dono (AC-E05.4).
type EntradaOrfao struct {
	Caminho string `json:"caminho"` // origem relativa a docs/
	Destino string `json:"destino"` // destino relativo a docs/
	Tipo    string `json:"tipo"`
	Motivo  string `json:"motivo"`
}

// ResultadoQuinteto resume o saneamento do Quinteto.
type ResultadoQuinteto struct {
	Orfaos     []EntradaOrfao
	Protegidos []string // caminhos relativos a docs/ mantidos na raiz
}

// posse é o mapa de donos: rota exata e token normalizado para a fatia.
type posse struct {
	porRota  map[string]string
	porToken map[string]string
}

// MigrarQuinteto saneia o Quinteto da raiz docs/ (AC-E05.2 revisado).
// Documento com dono ou protegido permanece em docs/<tipo>/ (escrita canônica).
// Sem dono, vai para docs/archive/quinteto/ com registro. Nada vai para
// processos/<slug>/docs/. Nunca deleta.
func MigrarQuinteto(raiz string, protegidos []string) (ResultadoQuinteto, error) {
	var res ResultadoQuinteto
	pos, err := donosDe(raiz)
	if err != nil {
		return res, err
	}
	docs := filepath.Join(raiz, "docs")
	for _, tipo := range tiposQuinteto {
		if err := migrarTipo(docs, tipo, raiz, pos, protegidos, &res); err != nil {
			return res, err
		}
	}
	return res, registrarOrfaos(docs, res.Orfaos)
}

// donosDe lê os manifestos de fatia e devolve a posse por rota e por token.
func donosDe(raiz string) (posse, error) {
	pos := posse{porRota: map[string]string{}, porToken: map[string]string{}}
	base := filepath.Join(raiz, caminhos.DirProcessos)
	entradas, err := os.ReadDir(base)
	if os.IsNotExist(err) {
		return pos, nil
	}
	if err != nil {
		return pos, err
	}
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		slug := e.Name()
		manifesto := filepath.Join(base, slug, slug+caminhos.SufixoStudio)
		rota, errR := rotaDoManifesto(manifesto)
		if errR != nil {
			return pos, errR
		}
		registrarPosse(pos, slug, rota)
	}
	return pos, nil
}

// registrarPosse grava a fatia por rota exata e por token, o primeiro vence.
func registrarPosse(pos posse, slug, rota string) {
	if rota != "" {
		if _, ok := pos.porRota[rota]; !ok {
			pos.porRota[rota] = slug
		}
		if token := normalizarToken(filepath.Base(rota)); token != "" {
			if _, ok := pos.porToken[token]; !ok {
				pos.porToken[token] = slug
			}
		}
	}
	token := normalizarToken(slug)
	if _, ok := pos.porToken[token]; !ok {
		pos.porToken[token] = slug
	}
}

// rotaDoManifesto extrai o caminho da rota upstream (texto ou objeto com path).
func rotaDoManifesto(caminho string) (string, error) {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var m struct {
		Upstream map[string]any `json:"upstream_signatures"`
	}
	if err := json.Unmarshal(dados, &m); err != nil {
		return "", err
	}
	switch v := m.Upstream["route"].(type) {
	case string:
		return filepath.ToSlash(v), nil
	case map[string]any:
		if p, ok := v["path"].(string); ok {
			return filepath.ToSlash(p), nil
		}
	}
	return "", nil
}

// normalizarToken reduz nome de documento ou de fatia a um token comparável.
func normalizarToken(nome string) string {
	base := strings.TrimSuffix(filepath.Base(nome), filepath.Ext(nome))
	base = reDataTipo.ReplaceAllString(base, "")
	base = strings.ToLower(base)
	base = strings.ReplaceAll(base, "_", "-")
	return strings.TrimPrefix(base, "legoloop-")
}

// resolverDono devolve a fatia dona do documento, ou vazio se órfão.
func resolverDono(pos posse, relDoc string) string {
	if slug, ok := pos.porRota[filepath.ToSlash(relDoc)]; ok {
		return slug
	}
	return pos.porToken[normalizarToken(filepath.Base(relDoc))]
}

// migrarTipo move os documentos de um tipo do Quinteto para donos ou arquivo morto.
func migrarTipo(docs, tipo, raiz string, pos posse, protegidos []string, res *ResultadoQuinteto) error {
	origem := filepath.Join(docs, tipo)
	if _, err := os.Stat(origem); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(origem, func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		base := filepath.Base(p)
		if casarAlgum(protegidos, base) {
			rel, _ := filepath.Rel(docs, p)
			res.Protegidos = append(res.Protegidos, filepath.ToSlash(rel))
			return nil
		}
		rel, _ := filepath.Rel(raiz, p)
		if resolverDono(pos, rel) != "" {
			manterNaRaiz(docs, p, res)
			return nil
		}
		return arquivarOrfao(docs, tipo, p, res)
	})
}

// casarAlgum diz se o nome-base casa algum glob de proteção.
func casarAlgum(globs []string, base string) bool {
	for _, g := range globs {
		if ok, err := filepath.Match(g, base); err == nil && ok {
			return true
		}
	}
	return false
}

// manterNaRaiz registra o documento com dono mantido em docs/<tipo>/ (canônico).
func manterNaRaiz(docs, origem string, res *ResultadoQuinteto) {
	rel, _ := filepath.Rel(docs, origem)
	res.Protegidos = append(res.Protegidos, filepath.ToSlash(rel))
}

// arquivarOrfao move o documento sem dono para docs/archive/quinteto/ e registra.
func arquivarOrfao(docs, tipo, origem string, res *ResultadoQuinteto) error {
	relO, _ := filepath.Rel(docs, origem)
	alvo := filepath.Join(docs, "archive", "quinteto", tipo, filepath.Base(origem))
	if err := moverArquivo(origem, alvo); err != nil {
		return err
	}
	relD, _ := filepath.Rel(docs, alvo)
	res.Orfaos = append(res.Orfaos, EntradaOrfao{
		Caminho: filepath.ToSlash(relO), Destino: filepath.ToSlash(relD), Tipo: tipo, Motivo: "sem_dono",
	})
	return nil
}

// registrarOrfaos grava o registro dos órfãos no arquivo morto (AC-E05.4).
func registrarOrfaos(docs string, orfaos []EntradaOrfao) error {
	if len(orfaos) == 0 {
		return nil
	}
	alvo := filepath.Join(docs, "archive", "quinteto", "registro_orfaos.json")
	if err := os.MkdirAll(filepath.Dir(alvo), 0o755); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(orfaos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(alvo, dados, 0o644)
}

// OrfaosDe lê o registro de órfãos do Quinteto no arquivo morto.
func OrfaosDe(raiz string) ([]EntradaOrfao, error) {
	caminho := filepath.Join(raiz, "docs", "archive", "quinteto", "registro_orfaos.json")
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var orfaos []EntradaOrfao
	return orfaos, json.Unmarshal(dados, &orfaos)
}
