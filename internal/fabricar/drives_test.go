package fabricar_test

import (
	"context"
	"errors"
	"io/fs"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
)

// probeFor simula unidades: ok responde nil, "hang" bloqueia até o prazo, "deny" nega acesso;
// qualquer outra letra não existe.
func probeFor(ok []string, hang, deny string) fabricar.DriveProbe {
	return func(ctx context.Context, root string) error {
		switch {
		case slices.Contains(ok, root):
			return nil
		case root == hang:
			<-ctx.Done()
			return ctx.Err()
		case root == deny:
			return errors.New("acesso negado")
		}
		return fs.ErrNotExist
	}
}

func TestDiscoverDrives_ReturnsAccessibleLetters(t *testing.T) {
	probe := probeFor([]string{`C:\`, `D:\`, `E:\`}, "", "")
	got := fabricar.DiscoverDrives(context.Background(), probe)
	if !slices.Equal(got, []string{`C:\`, `D:\`, `E:\`}) {
		t.Fatalf("unidades = %v", got)
	}
}

func TestDiscoverDrives_HangingDriveIsSkippedWithinTheTimeout(t *testing.T) {
	probe := probeFor([]string{`C:\`, `D:\`}, `F:\`, `G:\`)
	start := time.Now()
	ok, skipped := fabricar.DiscoverDrivesDetailed(context.Background(), probe, 200*time.Millisecond)
	if elapsed := time.Since(start); elapsed > 1500*time.Millisecond {
		t.Fatalf("a unidade que trava atrasou tudo: %v", elapsed)
	}
	if !slices.Equal(ok, []string{`C:\`, `D:\`}) {
		t.Errorf("as demais unidades devem voltar: %v", ok)
	}
	reasons := map[string]string{}
	for _, s := range skipped {
		reasons[s.Root] = s.Reason
	}
	if reasons[`F:\`] != "timeout" || reasons[`G:\`] != "inacessivel" || len(reasons) != 2 {
		t.Errorf("puladas incorretas (letras inexistentes não contam): %v", reasons)
	}
}

const procMounts = `proc /proc proc rw 0 0
sysfs /sys sysfs rw 0 0
devtmpfs /dev devtmpfs rw 0 0
tmpfs /run tmpfs rw 0 0
cgroup2 /sys/fs/cgroup cgroup2 rw 0 0
overlay /var/lib/docker/overlay2/x overlay rw 0 0
/dev/loop1 /snap/core squashfs ro 0 0
/dev/sda1 / ext4 rw 0 0
/dev/sda2 /home ext4 rw 0 0
/dev/sdb1 /mnt/dados ext4 rw 0 0
tmpfs /mnt/ram tmpfs rw 0 0
/dev/sdc1 /media/user/USB\040STICK vfat rw 0 0
`

func TestParseMounts_KeepsRootAndExternalMountsOnly(t *testing.T) {
	got := fabricar.ParseMounts(strings.NewReader(procMounts))
	want := []string{"/", "/media/user/USB STICK", "/mnt/dados"}
	if !slices.Equal(got, want) {
		t.Fatalf("ParseMounts = %v; quer %v", got, want)
	}
}

func TestParseMounts_EmptyInputStillHasRoot(t *testing.T) {
	if got := fabricar.ParseMounts(strings.NewReader("")); !slices.Equal(got, []string{"/"}) {
		t.Fatalf("ParseMounts(vazio) = %v", got)
	}
}

func TestMacVolumes_RootPlusVolumes(t *testing.T) {
	list := func(dir string) ([]string, error) { return []string{"Pendrive", "Backup"}, nil }
	if got := fabricar.MacVolumes(list); !slices.Equal(got, []string{"/", "/Volumes/Backup", "/Volumes/Pendrive"}) {
		t.Fatalf("MacVolumes = %v", got)
	}
	failing := func(string) ([]string, error) { return nil, errors.New("sem /Volumes") }
	if got := fabricar.MacVolumes(failing); !slices.Equal(got, []string{"/"}) {
		t.Fatalf("falha na listagem deve manter só a raiz: %v", got)
	}
}

func TestSystemRoots_RealHostHasAtLeastOneRoot(t *testing.T) {
	roots, _ := fabricar.SystemRoots(context.Background())
	if len(roots) == 0 {
		t.Fatal("o host real deve ter ao menos uma raiz")
	}
}
