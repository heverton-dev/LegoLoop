package renomear

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type linhaManifesto struct {
	origem string
	backup string
	acao   string
}

type backupCtx struct {
	dir    string
	linhas []linhaManifesto
	seq    int
}

func resolverBackupDir(op Opcoes) string {
	if strings.TrimSpace(op.BackupDir) != "" {
		return op.BackupDir
	}
	ts := op.Agora().Format("20060102-150405")
	nome := fmt.Sprintf("renomear-%s-%s-%s", op.Antigo, op.Novo, ts)
	return filepath.Join(op.Raiz, ".backups", nome)
}

func novoBackup(dir string) (*backupCtx, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &backupCtx{dir: dir}, nil
}

func (b *backupCtx) salvar(origem, acao string) error {
	info, err := os.Lstat(origem)
	if err != nil {
		return err
	}
	b.seq++
	destino := filepath.Join(b.dir, fmt.Sprintf("%03d_%s", b.seq, filepath.Base(origem)))
	if info.Mode()&os.ModeSymlink != 0 {
		if err := copiarSymlink(origem, destino); err != nil {
			return err
		}
	} else if info.IsDir() {
		if err := copiarPasta(origem, destino); err != nil {
			return err
		}
	} else {
		if err := copiarArquivo(origem, destino, info.Mode()); err != nil {
			return err
		}
	}
	b.linhas = append(b.linhas, linhaManifesto{origem: origem, backup: destino, acao: acao})
	return nil
}

func (b *backupCtx) gravarManifesto() (string, error) {
	var sb strings.Builder
	sb.WriteString("origem\tbackup\tacao\n")
	for _, l := range b.linhas {
		sb.WriteString(l.origem)
		sb.WriteByte('\t')
		sb.WriteString(l.backup)
		sb.WriteByte('\t')
		sb.WriteString(l.acao)
		sb.WriteByte('\n')
	}
	caminho := filepath.Join(b.dir, "MANIFESTO.tsv")
	if err := os.WriteFile(caminho, []byte(sb.String()), 0o644); err != nil {
		return "", err
	}
	return caminho, nil
}

func copiarSymlink(origem, destino string) error {
	alvo, err := os.Readlink(origem)
	if err != nil {
		return err
	}
	return os.WriteFile(destino+".symlink.txt", []byte(alvo+"\n"), 0o644)
}

func copiarPasta(origem, destino string) error {
	return filepath.WalkDir(origem, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(origem, path)
		if err != nil {
			return err
		}
		alvo := filepath.Join(destino, rel)
		if d.IsDir() {
			return os.MkdirAll(alvo, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return copiarArquivo(path, alvo, info.Mode())
	})
}

func copiarArquivo(origem, destino string, modo fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	in, err := os.Open(origem)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destino, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, modo.Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
