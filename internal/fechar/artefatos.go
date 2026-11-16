package fechar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// minLinhasSkill separa a Skill real do stub de três linhas que o fechar gravava.
	minLinhasSkill = 10
	// minChavesLivro separa o livro.json real do stub de duas chaves (nome_base e titulo).
	minChavesLivro = 3
)

func relativo(repoDir, caminho string) string {
	if rel, err := filepath.Rel(repoDir, caminho); err == nil {
		return filepath.ToSlash(rel)
	}
	return caminho
}

func erroAusente(repoDir, rotulo, caminho string) error {
	return fmt.Errorf("%w: %s em %s", ErrArtefatoAusente, rotulo, relativo(repoDir, caminho))
}

func erroStub(repoDir, rotulo, caminho, motivo string) error {
	return fmt.Errorf("%w: %s em %s (%s)", ErrArtefatoStub, rotulo, relativo(repoDir, caminho), motivo)
}

// exigirManual aceita só o manual deste processo: o nome termina em _MANUAL[_FUNCAO]_<SLUG>[_LEGOLOOP].md.
func exigirManual(repoDir, slug string) (string, error) {
	dir := filepath.Join(repoDir, "docs", "manuals")
	padrao := `_MANUAL(_FUNCAO)?_` + regexp.QuoteMeta(upperSlug(slug)) + `(_LEGOLOOP)?\.md$`
	re := regexp.MustCompile(padrao)
	entradas, _ := os.ReadDir(dir)
	for _, e := range entradas {
		if !e.IsDir() && re.MatchString(e.Name()) {
			return filepath.Join(dir, e.Name()), nil
		}
	}
	alvo := filepath.Join(dir, "<data>_MANUAL_"+upperSlug(slug)+".md")
	return alvo, erroAusente(repoDir, "manual operacional", alvo)
}

// exigirLivro recusa o livro ausente e o livro.json que só tem o par nome_base e titulo.
func exigirLivro(repoDir, dir string) (string, error) {
	manifesto := filepath.Join(dir, "livro.json")
	dados, err := os.ReadFile(manifesto)
	if err != nil {
		return dir, erroAusente(repoDir, "livro didático (livro.json)", manifesto)
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(dados, &obj); err != nil {
		return dir, erroStub(repoDir, "livro didático", manifesto, "JSON inválido")
	}
	if len(obj) < minChavesLivro {
		return dir, erroStub(repoDir, "livro didático", manifesto, fmt.Sprintf("%d chaves, mínimo %d", len(obj), minChavesLivro))
	}
	return dir, nil
}

// exigirSkill recusa a Skill ausente e a Skill com menos de minLinhasSkill linhas.
func exigirSkill(repoDir, skillPath string) (string, error) {
	dados, err := os.ReadFile(skillPath)
	if err != nil {
		return skillPath, erroAusente(repoDir, "skill canônica (SKILL.md)", skillPath)
	}
	linhas := len(strings.Split(strings.TrimSpace(string(dados)), "\n"))
	if linhas < minLinhasSkill {
		return skillPath, erroStub(repoDir, "skill canônica", skillPath, fmt.Sprintf("%d linhas, mínimo %d", linhas, minLinhasSkill))
	}
	return skillPath, nil
}

// exigirNosDoDAG recusa o DAG ausente e o DAG sem nenhum nó.
func exigirNosDoDAG(repoDir, dagPath string) error {
	dados, err := os.ReadFile(dagPath)
	if err != nil {
		return erroAusente(repoDir, "DAG do processo", dagPath)
	}
	var dag dagNos
	if err := json.Unmarshal(dados, &dag); err != nil {
		return erroStub(repoDir, "DAG do processo", dagPath, "JSON inválido")
	}
	if len(dag.Nodes) == 0 {
		return erroStub(repoDir, "DAG do processo", dagPath, "nodes vazio")
	}
	return nil
}
