package assets

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// root é um diretório de ativos com o escopo e a origem que ele confere.
type root struct{ dir, scope, origin string }

type cacheEntry struct {
	mtime, size int64
	items       []Item
	gen         uint64
}

// Scanner varre a máquina. A re-varredura só relê arquivos cujo (caminho, mtime, tamanho) mudou.
type Scanner struct {
	cfg Config

	scanMu sync.Mutex // serializa Scan: a geração do cache precisa ser exclusiva
	mu     sync.Mutex
	cache  map[string]cacheEntry
	gen    uint64

	onParse func(path string) // gancho de teste: conta parses reais (cache miss)
}

// NewScanner cria um scanner com cache vazio.
func NewScanner(cfg Config) *Scanner {
	return &Scanner{cfg: cfg.withDefaults(), cache: map[string]cacheEntry{}}
}

// Scan inventaria todas as categorias em paralelo. Nunca falha: o que não puder ser lido é ignorado.
func (s *Scanner) Scan(ctx context.Context) Snapshot {
	s.scanMu.Lock()
	defer s.scanMu.Unlock()
	s.mu.Lock()
	s.gen++
	s.mu.Unlock()

	scans := []func(context.Context) []Item{
		s.scanSkills, s.scanMCPs, s.scanScripts, s.scanHooks, s.scanGates, s.scanHarnesses, s.scanBinaries, s.scanProcessos,
	}
	parts := make([][]Item, len(scans))
	var wg sync.WaitGroup
	for i, fn := range scans {
		wg.Go(func() { parts[i] = fn(ctx) })
	}
	wg.Wait()
	s.prune()

	var items []Item
	for _, p := range parts {
		items = append(items, p...)
	}
	flagDuplicateSkills(items)
	for i := range items {
		items[i].ID = itemID(items[i])
	}
	sortItems(items)
	counts := make(map[Kind]int, len(Kinds()))
	for _, k := range Kinds() {
		counts[k] = 0
	}
	for _, it := range items {
		counts[it.Kind]++
	}
	return Snapshot{Items: items, Counts: counts, ScannedAt: time.Now()}
}

// Roots devolve os diretórios de ativos que existem (base da canonização).
func (s *Scanner) Roots() []string {
	var out []string
	seen := map[string]bool{}
	for _, group := range [][]root{s.skillRoots(), s.scriptRoots(context.Background())} {
		for _, r := range group {
			if st, err := os.Stat(r.dir); err == nil && st.IsDir() && !seen[r.dir] {
				seen[r.dir] = true
				out = append(out, r.dir)
			}
		}
	}
	return out
}

// load devolve os itens do arquivo, usando o cache quando mtime e tamanho não mudaram. A chave
// inclui o tipo do parser: o mesmo arquivo pode alimentar vários (settings.json tem hooks e MCPs).
func (s *Scanner) load(kind Kind, path string, parse func(string) []Item) []Item {
	key := string(kind) + "\x00" + path
	st, err := os.Stat(path)
	if err != nil || !st.Mode().IsRegular() {
		return nil
	}
	mtime, size := st.ModTime().UnixNano(), st.Size()
	s.mu.Lock()
	if e, ok := s.cache[key]; ok && e.mtime == mtime && e.size == size {
		e.gen = s.gen
		s.cache[key] = e
		s.mu.Unlock()
		return cloneItems(e.items)
	}
	gen := s.gen
	s.mu.Unlock()
	if s.onParse != nil {
		s.onParse(path)
	}
	items := parse(path)
	s.mu.Lock()
	s.cache[key] = cacheEntry{mtime: mtime, size: size, items: cloneItems(items), gen: gen}
	s.mu.Unlock()
	return items
}

// prune descarta do cache o que não foi visto na varredura atual (arquivo removido).
func (s *Scanner) prune() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for p, e := range s.cache {
		if e.gen != s.gen {
			delete(s.cache, p)
		}
	}
}

func sortItems(items []Item) {
	order := map[Kind]int{}
	for i, k := range Kinds() {
		order[k] = i
	}
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.Kind != b.Kind {
			return order[a.Kind] < order[b.Kind]
		}
		if an, bn := strings.ToLower(a.Name), strings.ToLower(b.Name); an != bn {
			return an < bn
		}
		return a.Path < b.Path
	})
}

// flagDuplicateSkills marca skills de conteúdo idêntico (mesmo SHA-256) com a contagem de cópias.
func flagDuplicateSkills(items []Item) {
	count := map[string]int{}
	for _, it := range items {
		if it.Kind == KindSkill {
			count[it.Meta["sha256"]]++
		}
	}
	for i, it := range items {
		if n := count[it.Meta["sha256"]]; it.Kind == KindSkill && n > 1 {
			items[i].Meta["duplicates"] = strconv.Itoa(n)
		}
	}
}

// joinRoots monta raízes relativas ao HOME (global) ou a cada projeto (project).
func (s *Scanner) joinRoots(home, project []root) []root {
	var out []root
	for _, r := range home {
		out = append(out, root{filepath.Join(s.cfg.Home, filepath.FromSlash(r.dir)), "global", r.origin})
	}
	for _, p := range s.cfg.Projects {
		for _, r := range project {
			out = append(out, root{filepath.Join(p, filepath.FromSlash(r.dir)), "project", r.origin})
		}
	}
	return out
}

// skipDirName: diretórios que nunca contêm ativos do usuário.
func skipDirName(name string) bool { return name == "node_modules" || name == ".git" }
