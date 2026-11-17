package fechar

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/caminhos"
)

// ErrParidadeSkillDAG sinaliza Nó do DAG não citado na Skill do Processo.
var ErrParidadeSkillDAG = errors.New("fechar: nó do DAG não citado na Skill")

var (
	errParidadeDAGAusente  = errors.New("DAG ausente")
	errParidadeDAGInvalido = errors.New("DAG com JSON inválido")
)

// NoDAG é o Nó do DAG que a Skill deve citar (por id ou por nome).
type NoDAG struct {
	ID   string `json:"id"`
	Nome string `json:"name"`
	Ref  string `json:"ref"`
}

// ParidadeRelato é a linha do modo relatório do Gate de paridade Skill×DAG.
type ParidadeRelato struct {
	Processo string
	Skill    string
	OK       bool
	Faltando []string
	Motivo   string
}

// NosNaoCitados devolve o rótulo de cada Nó do DAG ausente do texto da Skill.
// Um Nó vale como citado quando o id ou o nome aparece no texto, sem diferenciar
// maiúsculas de minúsculas.
func NosNaoCitados(skillTexto string, nos []NoDAG) []string {
	texto := strings.ToLower(skillTexto)
	var faltando []string
	for _, n := range nos {
		if !noCitado(texto, n) {
			faltando = append(faltando, rotuloNo(n))
		}
	}
	return faltando
}

func noCitado(textoBaixo string, n NoDAG) bool {
	if n.ID != "" && strings.Contains(textoBaixo, strings.ToLower(n.ID)) {
		return true
	}
	return n.Nome != "" && strings.Contains(textoBaixo, strings.ToLower(n.Nome))
}

func rotuloNo(n NoDAG) string {
	switch {
	case n.ID != "" && n.Nome != "":
		return fmt.Sprintf("%s (%s)", n.ID, n.Nome)
	case n.ID != "":
		return n.ID
	default:
		return n.Nome
	}
}

func lerNosDoDAG(dagJSON []byte) ([]NoDAG, error) {
	var dag struct {
		Nodes []NoDAG `json:"nodes"`
	}
	if err := json.Unmarshal(dagJSON, &dag); err != nil {
		return nil, err
	}
	return dag.Nodes, nil
}

func carregarNosParidade(dagPath string) ([]NoDAG, error) {
	dagJSON, err := os.ReadFile(dagPath)
	if err != nil {
		return nil, errParidadeDAGAusente
	}
	nos, err := lerNosDoDAG(dagJSON)
	if err != nil {
		return nil, errParidadeDAGInvalido
	}
	return nos, nil
}

// resolverSkillParidade devolve a SKILL.md do Processo: a canônica primeiro,
// a do slug em seguida e as refs dos nós por último.
func resolverSkillParidade(repoDir, slug, canon string, nos []NoDAG) string {
	candidatos := []string{filepath.Join("skills", canon, "SKILL.md")}
	if slug != "" && slug != canon {
		candidatos = append(candidatos, filepath.Join("skills", slug, "SKILL.md"))
	}
	for _, n := range nos {
		if ref := strings.TrimPrefix(n.Ref, "skills/"); ref != n.Ref {
			candidatos = append(candidatos, filepath.Join("skills", ref, "SKILL.md"))
		}
	}
	for _, rel := range candidatos {
		alvo := filepath.Join(repoDir, filepath.FromSlash(rel))
		if info, err := os.Stat(alvo); err == nil && !info.IsDir() {
			return alvo
		}
	}
	return ""
}

// auditParidadeSkillDAG é o Gate de paridade da Lei 11: falha quando a Skill
// do Processo não cita todos os nós do DAG.
func auditParidadeSkillDAG(repoDir, slug, canon, dagPath string) error {
	nos, err := carregarNosParidade(dagPath)
	if err != nil {
		if errors.Is(err, errParidadeDAGAusente) {
			return erroAusente(repoDir, "DAG do processo", dagPath)
		}
		return fmt.Errorf("%w: %s", errParidadeDAGInvalido, relativo(repoDir, dagPath))
	}
	skillPath := resolverSkillParidade(repoDir, slug, canon, nos)
	if skillPath == "" {
		return erroAusente(repoDir, "skill canônica (SKILL.md)", filepath.Join(repoDir, "skills", canon, "SKILL.md"))
	}
	texto, err := os.ReadFile(skillPath)
	if err != nil {
		return erroAusente(repoDir, "skill canônica (SKILL.md)", skillPath)
	}
	if faltando := NosNaoCitados(string(texto), nos); len(faltando) > 0 {
		return fmt.Errorf("%w: %s não cita os nós %s", ErrParidadeSkillDAG, relativo(repoDir, skillPath), strings.Join(faltando, ", "))
	}
	return nil
}

// RelatorioParidadeSkillDAG varre os Processos do repositório e devolve o
// relato de paridade Skill×DAG de cada um, sem alterar o disco.
func RelatorioParidadeSkillDAG(repoDir string) []ParidadeRelato {
	var relatos []ParidadeRelato
	for _, base := range caminhos.DirsLeitura(repoDir) {
		entradas, _ := os.ReadDir(base)
		for _, e := range entradas {
			if !e.IsDir() {
				continue
			}
			relatos = append(relatos, relatarParidadeProcesso(repoDir, filepath.Join(base, e.Name()), e.Name()))
		}
	}
	return relatos
}

func relatarParidadeProcesso(repoDir, procDir, nome string) ParidadeRelato {
	rel := ParidadeRelato{Processo: nome}
	nos, err := carregarNosParidade(caminhos.CaminhoDAGLeitura(procDir))
	if err != nil {
		rel.Motivo = err.Error()
		return rel
	}
	slug, canon := nomesDoProcesso(nome)
	skillPath := resolverSkillParidade(repoDir, slug, canon, nos)
	if skillPath == "" {
		rel.Motivo = "SKILL.md ausente"
		return rel
	}
	rel.Skill = relativo(repoDir, skillPath)
	texto, err := os.ReadFile(skillPath)
	if err != nil {
		rel.Motivo = "SKILL.md ilegível"
		return rel
	}
	rel.Faltando = NosNaoCitados(string(texto), nos)
	rel.OK = len(rel.Faltando) == 0
	return rel
}

func nomesDoProcesso(nome string) (string, string) {
	slug, canon, err := NormalizeSlug(nome)
	if err != nil {
		return nome, nome
	}
	return slug, canon
}
