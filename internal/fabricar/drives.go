package fabricar

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/fs"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

// probeTimeout é o prazo de cada sondagem de unidade (unidades de rede desconectadas e leitores
// vazios bloqueiam as chamadas de arquivo).
const probeTimeout = 2 * time.Second

// SkippedRoot registra uma raiz que não foi varrida e por quê.
type SkippedRoot struct {
	Root   string `json:"root"`
	Reason string `json:"reason"` // "timeout", "inacessivel"
}

// DriveProbe sonda uma raiz. nil significa acessível; fs.ErrNotExist significa que a unidade não existe.
type DriveProbe func(ctx context.Context, root string) error

func osProbe(_ context.Context, root string) error {
	_, err := os.Stat(root)
	return err
}

// DiscoverDrives devolve as letras de unidade (A:\ a Z:\) que respondem à sondagem, em ordem alfabética.
// Uma unidade que trava é pulada no prazo e não atrasa as demais.
func DiscoverDrives(ctx context.Context, probe DriveProbe) []string {
	roots, _ := DiscoverDrivesDetailed(ctx, probe, probeTimeout)
	return roots
}

type probeResult struct {
	root string
	err  error
}

// DiscoverDrivesDetailed é DiscoverDrives com o prazo por unidade e a lista de unidades puladas.
// Letras inexistentes não são unidades e não aparecem como puladas.
func DiscoverDrivesDetailed(ctx context.Context, probe DriveProbe, timeout time.Duration) ([]string, []SkippedRoot) {
	results := make([]probeResult, 26)
	var wg sync.WaitGroup
	for i := range results {
		root := string(rune('A'+i)) + `:\`
		wg.Go(func() { results[i] = probeOne(ctx, probe, root, timeout) })
	}
	wg.Wait()
	return classify(results)
}

// probeOne roda a sondagem em goroutine própria. Se ela travar para sempre, a goroutine vaza, mas o
// chamador segue no prazo (o canal tem buffer, então ela nunca bloqueia ao terminar).
func probeOne(ctx context.Context, probe DriveProbe, root string, timeout time.Duration) probeResult {
	pctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- probe(pctx, root) }()
	select {
	case err := <-done:
		return probeResult{root, err}
	case <-pctx.Done():
		return probeResult{root, pctx.Err()}
	}
}

func classify(results []probeResult) ([]string, []SkippedRoot) {
	var ok []string
	var skipped []SkippedRoot
	for _, r := range results {
		switch {
		case r.err == nil:
			ok = append(ok, r.root)
		case errors.Is(r.err, fs.ErrNotExist):
		case errors.Is(r.err, context.DeadlineExceeded):
			skipped = append(skipped, SkippedRoot{r.root, "timeout"})
		default:
			skipped = append(skipped, SkippedRoot{r.root, "inacessivel"})
		}
	}
	return ok, skipped
}

// virtualFS são os sistemas de arquivos que não guardam peças do usuário.
var virtualFS = map[string]bool{
	"proc": true, "sysfs": true, "devtmpfs": true, "devpts": true, "tmpfs": true, "cgroup": true,
	"cgroup2": true, "overlay": true, "squashfs": true, "securityfs": true, "debugfs": true,
	"tracefs": true, "configfs": true, "fusectl": true, "mqueue": true, "hugetlbfs": true,
	"pstore": true, "bpf": true, "autofs": true, "binfmt_misc": true,
}

var mountUnescape = strings.NewReplacer(`\040`, " ", `\011`, "\t", `\012`, "\n", `\134`, `\`)

// ParseMounts lê o formato de /proc/mounts e devolve "/" mais os pontos de montagem externos
// (/media, /mnt e /run/media), sem sistemas de arquivos virtuais.
func ParseMounts(r io.Reader) []string {
	set := map[string]bool{"/": true}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if mp, ok := externalMount(sc.Text()); ok {
			set[mp] = true
		}
	}
	return slices.Sorted(maps.Keys(set))
}

func externalMount(line string) (string, bool) {
	f := strings.Fields(line)
	if len(f) < 3 || virtualFS[f[2]] {
		return "", false
	}
	mp := mountUnescape.Replace(f[1])
	for _, prefix := range []string{"/media/", "/mnt/", "/run/media/"} {
		if strings.HasPrefix(mp, prefix) {
			return mp, true
		}
	}
	return "", false
}

// MacVolumes devolve "/" mais cada volume em /Volumes. list é a borda de listagem de diretório.
func MacVolumes(list func(dir string) ([]string, error)) []string {
	roots := []string{"/"}
	names, err := list("/Volumes")
	if err != nil {
		return roots
	}
	for _, n := range slices.Sorted(slices.Values(names)) {
		roots = append(roots, "/Volumes/"+n)
	}
	return roots
}

// SystemRoots devolve as raízes de varredura do sistema atual e as unidades puladas.
func SystemRoots(ctx context.Context) ([]string, []SkippedRoot) {
	switch runtime.GOOS {
	case "windows":
		return DiscoverDrivesDetailed(ctx, osProbe, probeTimeout)
	case "linux":
		return linuxRoots(), nil
	case "darwin":
		return MacVolumes(listVolumes), nil
	}
	return []string{"/"}, nil
}

func linuxRoots() []string {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return []string{"/"}
	}
	defer f.Close()
	return ParseMounts(f)
}

// listVolumes lista /Volumes sem os symlinks (o volume de boot aparece como link para "/").
func listVolumes(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.Type()&fs.ModeSymlink == 0 {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
