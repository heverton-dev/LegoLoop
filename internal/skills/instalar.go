// Package skills distribui as skills universais e as skills de Processo para
// o Cofre canônico e para todos os harnesses detectados, sem passo manual.
package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"legoloop/internal/caminhos"
	"legoloop/internal/processo"
	nativas "legoloop/skills"
)

// Opcao descreve a instalação de skills.
type Opcao struct {
	Vault        string
	Home         string
	DirProcessos string
	DirSkills    string
	DryRun       bool
	Eager        bool // true = vincula toda skill em toda pasta de harness (comportamento anterior a 20/11)
}

// politica devolve a regra de vínculo: lazy por padrão, com as exceções de ~/.legoloop/skills-eager.txt.
func (o Opcao) politica() processo.Politica {
	return processo.Politica{Eager: o.Eager, Excecoes: processo.LerExcecoes(o.Home)}
}

// Relatorio resume a instalação executada.
type Relatorio struct {
	Vault     string   `json:"vault"`
	Skills    []string `json:"skills"`
	Harnesses []string `json:"harnesses"`
	Vinculos  []string `json:"vinculos"`
	Comandos  []string `json:"comandos"`
	DryRun    bool     `json:"dry_run"`
	Status    string   `json:"status"`
}

// Instalar valida todas as skills e as distribui para o Cofre e os harnesses.
// A validação roda antes de qualquer gravação: skill inválida aborta tudo.
func Instalar(op Opcao) (*Relatorio, error) {
	op = completarOpcao(op)
	skProc, cmds, err := carregarSkillsProcesso(op)
	if err != nil {
		return nil, err
	}
	nat, err := nativas.TodasAsSkills()
	if err != nil {
		return nil, fmt.Errorf("carregar skills nativas: %w", err)
	}
	if err := validarNativas(nat); err != nil {
		return nil, err
	}
	hs := processo.DetectarHarnesses(processo.OSFS(), op.Home)
	return distribuir(op, hs, nat, skProc, cmds)
}

func completarOpcao(op Opcao) Opcao {
	if op.DirProcessos == "" {
		op.DirProcessos = caminhos.DirProcessos
	}
	if op.DirSkills == "" {
		op.DirSkills = filepath.Join(filepath.Dir(op.DirProcessos), "skills")
	}
	return op
}

// processosComDAG lista os Processos com contrato DAG no disco.
func processosComDAG(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("ler %s: %w", dir, err)
	}
	var slugs []string
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		procDir := filepath.Join(dir, e.Name())
		if _, err := os.Stat(caminhos.CaminhoDAGLeitura(procDir)); err == nil {
			slugs = append(slugs, e.Name())
		}
	}
	sort.Strings(slugs)
	return slugs, nil
}

// carregarSkillsProcesso reúne a skill e o comando de cada Processo. A skill
// vem do disco quando o usuário já a tem; senão é gerada. Tudo é validado
// antes de qualquer gravação.
func carregarSkillsProcesso(op Opcao) (map[string][]byte, map[string][]byte, error) {
	slugs, err := processosComDAG(op.DirProcessos)
	if err != nil {
		return nil, nil, err
	}
	sk := make(map[string][]byte, len(slugs))
	cmds := make(map[string][]byte, len(slugs))
	for _, slug := range slugs {
		camSk := filepath.Join(op.DirSkills, slug, "SKILL.md")
		dados, err := lerOuGerar(camSk, slug, gerarSkillProcesso)
		if err != nil {
			return nil, nil, err
		}
		if err := validarFrontmatter(camSk, dados); err != nil {
			return nil, nil, err
		}
		sk[slug] = dados
		cmds["legoloop-"+slug] = gerarComandoProcesso(slug)
	}
	return sk, cmds, nil
}

func lerOuGerar(caminho, slug string, gerar func(string) []byte) ([]byte, error) {
	if dados, err := os.ReadFile(caminho); err == nil {
		return dados, nil
	}
	return gerar(slug), nil
}

// validarNativas confere o front matter de todas as skills universais.
func validarNativas(nat map[string][]byte) error {
	for _, fn := range nativas.FuncoesNativas {
		if err := validarFrontmatter(fn, nat[fn]); err != nil {
			return err
		}
	}
	return nil
}

// validarFrontmatter exige front matter fechado com name e description.
func validarFrontmatter(caminho string, dados []byte) error {
	linhas := strings.Split(string(dados), "\n")
	if len(linhas) < 2 || strings.TrimSpace(linhas[0]) != "---" {
		return fmt.Errorf("skill inválida em %s: front matter ausente", caminho)
	}
	fim := -1
	for i := 1; i < len(linhas); i++ {
		if strings.TrimSpace(linhas[i]) == "---" {
			fim = i
			break
		}
	}
	if fim < 0 {
		return fmt.Errorf("skill inválida em %s: front matter não fechado", caminho)
	}
	cab := strings.Join(linhas[1:fim], "\n")
	if !chavePreenchida(cab, "name") {
		return fmt.Errorf("skill inválida em %s: name ausente", caminho)
	}
	if !chavePreenchida(cab, "description") {
		return fmt.Errorf("skill inválida em %s: description ausente", caminho)
	}
	return nil
}

func chavePreenchida(cab, chave string) bool {
	for _, ln := range strings.Split(cab, "\n") {
		ln = strings.TrimSpace(ln)
		if !strings.HasPrefix(ln, chave+":") {
			continue
		}
		return strings.TrimSpace(strings.TrimPrefix(ln, chave+":")) != ""
	}
	return false
}

// distribuir grava projeto, Cofre e harnesses depois da validação total.
func distribuir(op Opcao, hs []processo.Harness, nat, skProc, cmds map[string][]byte) (*Relatorio, error) {
	res := &Relatorio{Vault: op.Vault, DryRun: op.DryRun, Status: "ok"}
	for _, h := range hs {
		res.Harnesses = append(res.Harnesses, h.Nome)
	}
	for _, fn := range nativas.FuncoesNativas {
		res.Skills = append(res.Skills, fn)
		vs, err := publicar(op, hs, fn, nat[fn])
		if err != nil {
			return nil, err
		}
		res.Vinculos = append(res.Vinculos, vs...)
	}
	for _, slug := range ordenarChaves(skProc) {
		res.Skills = append(res.Skills, slug)
		camSk := filepath.Join(op.DirSkills, slug, "SKILL.md")
		if err := gravarSeNovo(camSk, skProc[slug], op.DryRun); err != nil {
			return nil, err
		}
		vs, err := publicar(op, hs, slug, skProc[slug])
		if err != nil {
			return nil, err
		}
		res.Vinculos = append(res.Vinculos, vs...)
	}
	if err := distribuirComandos(op, hs, cmds, res); err != nil {
		return nil, err
	}
	idx, err := garantirIndice(op, hs)
	if err != nil {
		return nil, err
	}
	res.Vinculos = append(res.Vinculos, idx...)
	return res, nil
}

// publicar grava a skill no Cofre e a vincula em cada harness.
func publicar(op Opcao, hs []processo.Harness, nome string, dados []byte) ([]string, error) {
	dirCofre := filepath.Join(op.Vault, "skills", nome)
	if !op.DryRun {
		if err := gravarArquivo(filepath.Join(dirCofre, "SKILL.md"), dados); err != nil {
			return nil, err
		}
	}
	var vs []processo.Vinculo
	if op.politica().Vincula(nome) {
		vs = processo.PlanejarVinculos(hs, nome, dirCofre)
	}
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.Destino)
	}
	if !op.DryRun && len(vs) > 0 {
		if _, err := processo.Vincular(vs, nil); err != nil {
			return nil, fmt.Errorf("vincular %s: %w", nome, err)
		}
	}
	return out, nil
}

// distribuirComandos grava cada comando no projeto e em todos os harnesses.
func distribuirComandos(op Opcao, hs []processo.Harness, cmds map[string][]byte, res *Relatorio) error {
	for _, nome := range ordenarChaves(cmds) {
		res.Comandos = append(res.Comandos, nome)
		slug := strings.TrimPrefix(nome, "legoloop-")
		camProj := filepath.Join(op.DirSkills, slug, "comandos", nome+".md")
		if err := gravarSeNovo(camProj, cmds[nome], op.DryRun); err != nil {
			return err
		}
		if op.DryRun {
			continue
		}
		for _, h := range hs {
			dest := filepath.Join(filepath.Dir(h.Skills), "commands", nome+".md")
			if err := gravarArquivo(dest, cmds[nome]); err != nil {
				return err
			}
		}
	}
	return nil
}

func gravarSeNovo(dest string, dados []byte, dry bool) error {
	if dry {
		return nil
	}
	if _, err := os.Stat(dest); err == nil {
		return nil
	}
	return gravarArquivo(dest, dados)
}

func gravarArquivo(dest string, dados []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("criar pasta %s: %w", filepath.Dir(dest), err)
	}
	if err := os.WriteFile(dest, dados, 0o644); err != nil {
		return fmt.Errorf("gravar %s: %w", dest, err)
	}
	return nil
}

func ordenarChaves(m map[string][]byte) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// gerarSkillProcesso gera a skill universal de um Processo. O corpo delega a
// execução ao comando canônico legoloop exec.
func gerarSkillProcesso(slug string) []byte {
	corpo := "---\nname: " + slug + "\ndescription: Executa o Processo " + slug +
		" pelo comando legoloop exec.\n---\n\n# " + slug +
		"\n\nExecute o Processo " + slug + " com o comando:\n\n    legoloop exec " + slug + "\n"
	return []byte(corpo)
}

// gerarComandoProcesso gera o comando /legoloop-<slug> de um Processo.
func gerarComandoProcesso(slug string) []byte {
	corpo := "---\nname: legoloop-" + slug + "\ndescription: Executa o Processo " + slug +
		" pelo comando legoloop exec.\n---\n\nExecute o Processo " + slug +
		" com o comando:\n\n    legoloop exec " + slug + "\n"
	return []byte(corpo)
}
