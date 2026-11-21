package processo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// harnessesConhecidos são os harnesses com pasta de skills própria. O padrão é ~/.<nome>/skills;
// os que fogem dele estão em layoutsEspeciais.
var harnessesConhecidos = []string{
	"claude", "codex", "gemini", "antigravity", "opencode", "cursor", "kiro",
	"windsurf", "copilot", "trae", "roo", "continue", "qwen", "cline", "hermes",
	"mimocode", "augment", "junie", "grok", "crush", "devin", "goose", "pi",
}

// layoutEspecial guarda a raiz do harness e a pasta de skills, relativas ao HOME.
type layoutEspecial struct{ raiz, skills []string }

var layoutsEspeciais = map[string]layoutEspecial{
	"antigravity": {[]string{".gemini", "config"}, []string{".gemini", "config", "skills"}},
	"opencode":    {[]string{".config", "opencode"}, []string{".config", "opencode", "skills"}},
	"windsurf":    {[]string{".codeium", "windsurf"}, []string{".codeium", "windsurf", "skills"}},
	"crush":       {[]string{".config", "crush"}, []string{".config", "crush", "skills"}},
	"devin":       {[]string{".config", "devin"}, []string{".config", "devin", "skills"}},
	"goose":       {[]string{".config", "goose"}, []string{".config", "goose", "skills"}},
	"pi":          {[]string{".pi", "agent"}, []string{".pi", "agent", "skills"}},
}

// NomesHarnesses devolve os harnesses conhecidos, na ordem de detecção.
func NomesHarnesses() []string { return append([]string(nil), harnessesConhecidos...) }

// PastaSkills devolve a pasta onde o harness lê skills, sob o HOME dado.
func PastaSkills(home, nome string) string {
	_, skills := caminhosHarness(home, nome)
	return skills
}

// Harness é um harness detectado no HOME: o nome e a pasta onde ele lê skills.
type Harness struct {
	Nome   string // ex.: "claude"
	Skills string // ex.: ~/.claude/skills
}

// DetectarHarnesses devolve os harnesses cuja pasta existe no HOME, na ordem de harnessesConhecidos.
func DetectarHarnesses(fsys FS, home string) []Harness {
	if home == "" {
		return nil
	}
	var out []Harness
	for _, n := range harnessesConhecidos {
		if skillsDir, ok := pastaDetectada(fsys, home, n); ok {
			out = append(out, Harness{Nome: n, Skills: skillsDir})
		}
	}
	return out
}

// layoutsLegados guardam o layout antigo de harnesses que mudaram de pasta: quem só tem a pasta
// antiga continua detectado, e a pasta nova vence quando existe.
var layoutsLegados = map[string]layoutEspecial{
	"opencode": {[]string{".opencode"}, []string{".opencode", "skills"}},
}

// pastaDetectada devolve a pasta de skills do harness se a raiz dele existe no HOME.
func pastaDetectada(fsys FS, home, nome string) (string, bool) {
	existe := func(dir string) bool {
		info, err := fsys.Stat(dir)
		return err == nil && info.IsDir()
	}
	raiz, skills := caminhosHarness(home, nome)
	if existe(raiz) {
		return skills, true
	}
	if l, ok := layoutsLegados[nome]; ok && existe(filepath.Join(append([]string{home}, l.raiz...)...)) {
		return filepath.Join(append([]string{home}, l.skills...)...), true
	}
	return "", false
}

func caminhosHarness(home, nome string) (string, string) {
	if l, ok := layoutsEspeciais[nome]; ok {
		return filepath.Join(append([]string{home}, l.raiz...)...), filepath.Join(append([]string{home}, l.skills...)...)
	}
	base := filepath.Join(home, "."+nome)
	return base, filepath.Join(base, "skills")
}

// Vinculo liga a skill do cofre (Origem) à pasta de skills de um harness (Destino).
type Vinculo struct {
	Harness string
	Origem  string
	Destino string
}

// PlanejarVinculos monta a prévia: um vínculo por harness, sem tocar o disco.
func PlanejarVinculos(hs []Harness, nome, origem string) []Vinculo {
	out := make([]Vinculo, 0, len(hs))
	for _, h := range hs {
		out = append(out, Vinculo{Harness: h.Nome, Origem: origem, Destino: filepath.Join(h.Skills, nome)})
	}
	return out
}

// LinkFunc cria um vínculo; é a borda que os testes trocam para provar o desfazer.
type LinkFunc func(origem, destino string) error

// SymlinkOuCopia cria destino como symlink para a pasta da skill; se o sistema não permitir, copia o SKILL.md.
func SymlinkOuCopia(origem, destino string) error {
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	if err := os.Symlink(origem, destino); err == nil {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(origem, "SKILL.md"))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destino, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(destino, "SKILL.md"), data, 0o644)
}

// Vincular cria todos os vínculos ou nenhum: na primeira falha desfaz os que ele mesmo criou.
// Um symlink que já aponta à Origem é reaproveitado e nunca é desfeito (não foi criado aqui).
// Devolve os destinos que ficaram ligados. link nil usa SymlinkOuCopia.
func Vincular(vs []Vinculo, link LinkFunc) ([]string, error) {
	if link == nil {
		link = SymlinkOuCopia
	}
	var ligados, criados []string
	for _, v := range vs {
		ja, err := jaLigado(v)
		if err != nil {
			desfazer(criados)
			return nil, err
		}
		if !ja {
			if err := link(v.Origem, v.Destino); err != nil {
				desfazer(append(criados, v.Destino))
				return nil, fmt.Errorf("vínculo com %s: %w", v.Harness, err)
			}
			criados = append(criados, v.Destino)
		}
		ligados = append(ligados, v.Destino)
	}
	return ligados, nil
}

// jaLigado diz se o destino já é o symlink certo; qualquer outra coisa ocupando o lugar é erro.
func jaLigado(v Vinculo) (bool, error) {
	info, err := os.Lstat(v.Destino)
	if err != nil {
		return false, nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		if alvo, err := os.Readlink(v.Destino); err == nil && filepath.Clean(alvo) == filepath.Clean(v.Origem) {
			return true, nil
		}
	}
	if copiaIdentica(v) {
		return true, nil
	}
	if copiaDesatualizada(v) {
		return false, nil
	}
	return false, fmt.Errorf("vínculo com %s: já existe %s", v.Harness, v.Destino)
}

func copiaDesatualizada(v Vinculo) bool {
	a, errA := os.ReadFile(filepath.Join(v.Destino, "SKILL.md"))
	b, errB := os.ReadFile(filepath.Join(v.Origem, "SKILL.md"))
	if errA != nil || errB != nil {
		return false
	}
	nomeSkill := filepath.Base(v.Origem)
	tag := "name: " + nomeSkill
	return strings.Contains(string(a), tag) && strings.Contains(string(b), tag)
}

// copiaIdentica reconhece o vínculo feito por cópia (sistema sem symlink): o SKILL.md é o mesmo do cofre.
func copiaIdentica(v Vinculo) bool {
	a, errA := os.ReadFile(filepath.Join(v.Destino, "SKILL.md"))
	b, errB := os.ReadFile(filepath.Join(v.Origem, "SKILL.md"))
	return errA == nil && errB == nil && string(a) == string(b)
}

// desfazer remove os vínculos criados. RemoveAll em symlink remove só o link, nunca o alvo.
func desfazer(destinos []string) {
	for i := len(destinos) - 1; i >= 0; i-- {
		_ = os.RemoveAll(destinos[i])
	}
}

// HarnessesVinculados devolve os nomes dos harnesses que têm a pasta skills/<nome>.
func HarnessesVinculados(fsys FS, hs []Harness, nome string) []string {
	if fsys == nil {
		fsys = OSFS()
	}
	out := []string{}
	for _, h := range hs {
		if _, err := fsys.Stat(filepath.Join(h.Skills, nome)); err == nil {
			out = append(out, h.Nome)
		}
	}
	return out
}
