package canon

import (
	"os"
	"path/filepath"
	"sort"
)

// Group é um conjunto de duplicatas exatas (mesmo SHA-256).
type Group struct {
	SHA256 string     `json:"sha256"`
	Files  []FileInfo `json:"files"`
}

// GroupDuplicates agrupa por hash e devolve apenas grupos com 2+ arquivos, ordenados por hash.
func GroupDuplicates(files []FileInfo) []Group {
	byHash := map[string][]FileInfo{}
	for _, f := range files {
		byHash[f.SHA256] = append(byHash[f.SHA256], f)
	}
	var out []Group
	for h, fs := range byHash {
		if len(fs) > 1 {
			out = append(out, Group{SHA256: h, Files: fs})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SHA256 < out[j].SHA256 })
	return out
}

// Elect escolhe a versão canônica: mais recente; empate ⇒ maior cobertura de
// frontmatter; empate ⇒ caminho mais curto; empate ⇒ ordem lexical.
func Elect(files []FileInfo) FileInfo {
	best := files[0]
	for _, f := range files[1:] {
		if better(f, best) {
			best = f
		}
	}
	return best
}

func better(a, b FileInfo) bool {
	if !a.ModTime.Equal(b.ModTime) {
		return a.ModTime.After(b.ModTime)
	}
	if a.FrontmatterKeys != b.FrontmatterKeys {
		return a.FrontmatterKeys > b.FrontmatterKeys
	}
	if len(a.Path) != len(b.Path) {
		return len(a.Path) < len(b.Path)
	}
	return a.Path < b.Path
}

// GroupPlan descreve o que será feito com um grupo de duplicatas.
type GroupPlan struct {
	SHA256     string   `json:"sha256"`
	Canonical  FileInfo `json:"canonical"`
	VaultPath  string   `json:"vault_path"`
	Links      []string `json:"links"`
	SavedBytes int64    `json:"saved_bytes"`
}

// Plan é o resultado do --dry-run: nenhuma escrita em disco.
type Plan struct {
	Vault      string      `json:"vault"`
	Groups     []GroupPlan `json:"groups"`
	SavedBytes int64       `json:"saved_bytes"`
}

// BuildPlan elege canônicos e calcula a economia. Arquivos que já são o mesmo
// arquivo físico do cofre (execução anterior) são ignorados ⇒ idempotente.
func BuildPlan(files []FileInfo, vault string) Plan {
	p := Plan{Vault: vault}
	for _, g := range GroupDuplicates(files) {
		c := Elect(g.Files)
		gp := GroupPlan{SHA256: g.SHA256, Canonical: c,
			VaultPath: filepath.Join(vault, g.SHA256[:2], g.SHA256, filepath.Base(c.Path))}
		vaultInfo, vaultErr := os.Stat(gp.VaultPath)
		for _, f := range g.Files {
			if fi, err := os.Stat(f.Path); vaultErr == nil && err == nil && os.SameFile(fi, vaultInfo) {
				continue
			}
			gp.Links = append(gp.Links, f.Path)
		}
		gp.SavedBytes = c.Size * int64(len(gp.Links))
		if vaultErr != nil {
			gp.SavedBytes -= c.Size
		}
		if len(gp.Links) > 0 && gp.SavedBytes > 0 {
			p.Groups = append(p.Groups, gp)
			p.SavedBytes += gp.SavedBytes
		}
	}
	return p
}
