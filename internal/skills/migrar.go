package skills

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// limiteCopiaSkill é o tamanho acima do qual só o SKILL.md vai ao Cofre (as skills enormes
// ficam completas no backup).
const limiteCopiaSkill = 3 * 1024 * 1024

var frontmatterValido = regexp.MustCompile(`(?s)^---\s*\n.*?\n---`)

// MigrarOpcao configura a migração das skills eager de cada harness para o Cofre.
type MigrarOpcao struct {
	Home    string
	Vault   string
	Data    string // sufixo do diretório de backup, ex.: 20261121
	Aplicar bool   // false = só simula
}

// MigrarHarness resume o que mudou (ou mudaria) em um harness.
type MigrarHarness struct {
	Nome      string `json:"nome"`
	Pasta     string `json:"pasta"`
	Mover     int    `json:"mover"`
	Copiar    int    `json:"copiar_ao_cofre"`
	Invalidas int    `json:"invalidas"`
}

// MigrarRelatorio é o resultado da migração.
type MigrarRelatorio struct {
	Backup    string          `json:"backup"`
	Harnesses []MigrarHarness `json:"harnesses"`
}

type planoMigracao struct {
	mover     []string          // entradas da pasta do harness que saem dela
	copiar    map[string]string // nome da skill -> pasta de origem, para quem ainda não está no Cofre
	invalidas int
}

// Migrar tira as skills eager de cada harness detectado: copia ao Cofre as que faltam, move o resto
// para um backup com manifesto e deixa só a skill-índice. Com Aplicar=false, apenas conta.
func Migrar(op MigrarOpcao) (MigrarRelatorio, error) {
	cofre := filepath.Join(op.Vault, "skills")
	rel := MigrarRelatorio{Backup: filepath.Join(op.Home, ".aidd", "backups", "skills_"+op.Data)}
	existentes := nomesNoCofre(cofre)
	for _, h := range processo.DetectarHarnesses(processo.OSFS(), op.Home) {
		if info, err := os.Stat(h.Skills); err != nil || !info.IsDir() {
			continue
		}
		plano := planejarMigracao(h.Skills, existentes)
		rel.Harnesses = append(rel.Harnesses, MigrarHarness{Nome: h.Nome, Pasta: h.Skills, Mover: len(plano.mover), Copiar: len(plano.copiar), Invalidas: plano.invalidas})
		if !op.Aplicar {
			continue
		}
		if err := aplicarMigracao(op, h, cofre, rel.Backup, plano, existentes); err != nil {
			return rel, fmt.Errorf("migrar %s: %w", h.Nome, err)
		}
	}
	return rel, nil
}

func nomesNoCofre(cofre string) map[string]bool {
	out := map[string]bool{}
	entradas, _ := os.ReadDir(cofre)
	for _, e := range entradas {
		out[e.Name()] = true
	}
	return out
}

func planejarMigracao(pasta string, existentes map[string]bool) planoMigracao {
	plano := planoMigracao{copiar: map[string]string{}}
	entradas, _ := os.ReadDir(pasta)
	for _, e := range entradas {
		nome := e.Name()
		dir := filepath.Join(pasta, nome)
		if nome == processo.IndiceSkills || strings.HasPrefix(nome, ".") || !dirOuLinkParaDir(dir) {
			continue
		}
		fontes := skillsDe(dir)
		if len(fontes) == 0 {
			continue
		}
		plano.mover = append(plano.mover, nome)
		if strings.HasSuffix(nome, ".disabled") || strings.Contains(nome, ".backup") {
			continue
		}
		classificarFontes(&plano, fontes, existentes)
	}
	return plano
}

func classificarFontes(plano *planoMigracao, fontes []string, existentes map[string]bool) {
	for _, src := range fontes {
		nome := filepath.Base(src)
		if !skillValida(filepath.Join(src, "SKILL.md")) {
			plano.invalidas++
			continue
		}
		if _, ja := plano.copiar[nome]; !existentes[nome] && !ja {
			plano.copiar[nome] = src
		}
	}
}

func dirOuLinkParaDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

// marcadoresDeFerramenta indicam que uma pasta é repositório de ferramenta (código), não pacote de skills.
var marcadoresDeFerramenta = []string{".git", "package.json", "Cargo.toml", "go.mod", "pyproject.toml", "setup.py", "Makefile"}

func pareceFerramenta(dir string) bool {
	for _, m := range marcadoresDeFerramenta {
		if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
			return true
		}
	}
	return false
}

// skillsDe devolve a própria pasta se tiver SKILL.md; senão, as skills aninhadas até 3 níveis (pacote).
// Repositórios de ferramenta sem SKILL.md no topo não são pacotes: ficam onde estão.
func skillsDe(dir string) []string {
	if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err == nil {
		return []string{dir}
	}
	if pareceFerramenta(dir) {
		return nil
	}
	var out []string
	base := strings.Count(dir, string(os.PathSeparator))
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if strings.Count(p, string(os.PathSeparator))-base > 3 {
			return filepath.SkipDir
		}
		if _, e := os.Stat(filepath.Join(p, "SKILL.md")); e == nil {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

func skillValida(arquivo string) bool {
	b, err := os.ReadFile(arquivo)
	if err != nil {
		return false
	}
	return frontmatterValido.Match([]byte(strings.TrimPrefix(string(b), string(rune(0xFEFF)))))
}

func aplicarMigracao(op MigrarOpcao, h processo.Harness, cofre, backup string, plano planoMigracao, existentes map[string]bool) error {
	for nome, src := range plano.copiar {
		if err := copiarSkillAoCofre(src, filepath.Join(cofre, nome)); err != nil {
			return err
		}
		existentes[nome] = true
	}
	bk := filepath.Join(backup, h.Nome)
	if err := os.MkdirAll(bk, 0o755); err != nil {
		return err
	}
	var manifesto strings.Builder
	for _, nome := range plano.mover {
		de, para := filepath.Join(h.Skills, nome), filepath.Join(bk, nome)
		if err := os.Rename(de, para); err != nil {
			return err
		}
		manifesto.WriteString(para + "\t" + de + "\n")
	}
	if err := os.WriteFile(filepath.Join(bk, "MANIFESTO.tsv"), []byte(manifesto.String()), 0o644); err != nil {
		return err
	}
	_, err := garantirIndice(Opcao{Home: op.Home}, []processo.Harness{h})
	return err
}

func copiarSkillAoCofre(src, dest string) error {
	if tamanhoPasta(src) > limiteCopiaSkill {
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return err
		}
		return copiarArquivo(filepath.Join(src, "SKILL.md"), filepath.Join(dest, "SKILL.md"))
	}
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(src, p)
		alvo := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(alvo, 0o755)
		}
		return copiarArquivo(p, alvo)
	})
}

func tamanhoPasta(dir string) int64 {
	var total int64
	_ = filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if info, e := d.Info(); e == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total
}

func copiarArquivo(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return nil // link quebrado ou arquivo sumido: não derruba a migração
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
