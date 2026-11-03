package canon

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// DefaultWorkers é o tamanho padrão do worker pool (scanner de 32 threads do Symbiont).
const DefaultWorkers = 32

// maxMetaBytes limita a leitura usada para extrair frontmatter/chaves JSON.
const maxMetaBytes = 1 << 20

var frontmatterKeyRe = regexp.MustCompile(`^[A-Za-z0-9_-]+\s*:`)

// FileInfo são os metadados de um ativo.
type FileInfo struct {
	Path            string    `json:"path"`
	Size            int64     `json:"size"`
	ModTime         time.Time `json:"mod_time"`
	SHA256          string    `json:"sha256"`
	FrontmatterKeys int       `json:"frontmatter_keys"`
}

// ScanOptions configura a varredura.
type ScanOptions struct {
	Roots   []string
	Exclude []string
	Workers int
	// IncludeGit inclui working trees git (padrão: ignoradas, pois um link
	// commitado quebra o repositório para colaboradores).
	IncludeGit bool
}

// HashTree varre as raízes (sem seguir symlinks, ignorando .git e node_modules) e
// calcula SHA-256 + metadados de cada .md/.json num worker pool.
func HashTree(ctx context.Context, opts ScanOptions) ([]FileInfo, error) {
	paths := make(chan string, 256)
	results := make(chan FileInfo, 256)
	var wg sync.WaitGroup
	for range max(opts.Workers, 1) {
		wg.Go(func() { hashWorker(ctx, paths, results) })
	}
	walkErr := make(chan error, 1)
	go func() {
		walkErr <- walkRoots(ctx, opts, paths)
		close(paths)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	var out []FileInfo
	for fi := range results {
		out = append(out, fi)
	}
	if err := <-walkErr; err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, ctx.Err()
}

func walkRoots(ctx context.Context, opts ScanOptions, paths chan<- string) error {
	for _, root := range opts.Roots {
		abs, err := filepath.Abs(root)
		if err != nil {
			return err
		}
		if _, err := os.Stat(abs); err != nil {
			return err
		}
		if _, inRepo := GitRoot(abs, homeDir()); inRepo && !opts.IncludeGit {
			continue
		}
		if err := filepath.WalkDir(abs, visitor(ctx, opts, paths)); err != nil {
			return err
		}
	}
	return nil
}

func visitor(ctx context.Context, opts ScanOptions, paths chan<- string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if err != nil || isExcluded(path, opts.Exclude) {
			return skipEntry(d)
		}
		if d.IsDir() {
			if skipDir(path, d.Name(), opts.IncludeGit) {
				return fs.SkipDir
			}
			return nil
		}
		if d.Type().IsRegular() && isAsset(path) {
			paths <- path
		}
		return nil
	}
}

// GitRoot devolve o working tree git que contém dir (dir ou um ancestral com .git).
// A subida para em stop (exclusivo), ex.: o HOME — um ~/.git de dotfiles não conta.
func GitRoot(dir, stop string) (string, bool) {
	for p := dir; !samePath(p, stop); p = filepath.Dir(p) {
		if _, err := os.Lstat(filepath.Join(p, ".git")); err == nil {
			return p, true
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return "", false
}

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

func samePath(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// skipDir ignora .git, node_modules e (sem IncludeGit) qualquer working tree git.
func skipDir(path, name string, includeGit bool) bool {
	if name == ".git" || name == "node_modules" {
		return true
	}
	if includeGit {
		return false
	}
	_, err := os.Lstat(filepath.Join(path, ".git"))
	return err == nil
}

func skipEntry(d fs.DirEntry) error {
	if d != nil && d.IsDir() {
		return fs.SkipDir
	}
	return nil
}

func isAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".json"
}

func isExcluded(path string, exclude []string) bool {
	for _, ex := range exclude {
		if abs, err := filepath.Abs(ex); err == nil && isUnder(path, abs) {
			return true
		}
	}
	return false
}

// isUnder compara caminhos respeitando a insensibilidade a caixa do Windows.
func isUnder(path, root string) bool {
	if runtime.GOOS == "windows" {
		path, root = strings.ToLower(path), strings.ToLower(root)
	}
	rel, err := filepath.Rel(root, path)
	return err == nil && filepath.IsLocal(rel)
}

func hashWorker(ctx context.Context, paths <-chan string, results chan<- FileInfo) {
	for p := range paths {
		if ctx.Err() != nil {
			continue // drena o canal para não bloquear o walker
		}
		if fi, err := hashFile(p); err == nil {
			results <- fi
		}
	}
}

func hashFile(path string) (FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileInfo{}, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return FileInfo{}, err
	}
	h := sha256.New()
	var head bytes.Buffer
	if _, err := io.Copy(io.MultiWriter(h, &limitWriter{buf: &head, n: maxMetaBytes}), f); err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Path: path, Size: st.Size(), ModTime: st.ModTime(),
		SHA256:          hex.EncodeToString(h.Sum(nil)),
		FrontmatterKeys: metaKeys(path, head.Bytes()),
	}, nil
}

// metaKeys conta chaves do frontmatter YAML (.md) ou chaves de topo (.json).
func metaKeys(path string, data []byte) int {
	if strings.EqualFold(filepath.Ext(path), ".json") {
		var m map[string]json.RawMessage
		if json.Unmarshal(data, &m) != nil {
			return 0
		}
		return len(m)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	if !sc.Scan() || strings.TrimSpace(sc.Text()) != "---" {
		return 0
	}
	n := 0
	for sc.Scan() && strings.TrimSpace(sc.Text()) != "---" {
		if frontmatterKeyRe.MatchString(sc.Text()) {
			n++
		}
	}
	return n
}

// limitWriter guarda apenas os primeiros n bytes e descarta o resto sem erro.
type limitWriter struct {
	buf *bytes.Buffer
	n   int
}

func (l *limitWriter) Write(p []byte) (int, error) {
	if l.n > 0 {
		k := min(len(p), l.n)
		l.buf.Write(p[:k])
		l.n -= k
	}
	return len(p), nil
}
