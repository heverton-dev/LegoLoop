// Package migrar converte um Processo v1 para os nomes canônicos (TICK-01.4):
// <nome>.loop.json vira <nome>.studio.json e pipeline.dag.json vira
// processo.dag.json. Os arquivos .dag.json e .processo.json não mudam. Antes de
// mover, a migração copia cada original para a pasta de backup. A operação é
// idempotente: sem origem pendente, nada muda no disco.
package migrar

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/caminhos"
)

// SufixoV1 é o sufixo do documento do Studio no dialeto antigo.
const SufixoV1 = ".loop.json"

// SufixoStudio é o sufixo do documento do Studio no dialeto novo.
const SufixoStudio = ".studio.json"

// DirBackup é a pasta, sob a raiz, que guarda a cópia dos originais.
const DirBackup = ".migrar-backup"

var (
	// ErrDestinoExiste indica que o destino de um movimento já existe; nada é alterado.
	ErrDestinoExiste = errors.New("migrar: destino já existe")

	// ErrSymlink indica um symlink no lugar de um arquivo a migrar (Path Jail).
	ErrSymlink = errors.New("migrar: symlink recusado")

	// ErrRaizInvalida indica raiz vazia, com "..", inexistente ou que não é pasta.
	ErrRaizInvalida = errors.New("migrar: raiz inválida")
)

// Movimento descreve uma renomeação, com caminhos relativos à raiz e barra "/".
type Movimento struct {
	Origem string

	Destino string
}

// Relatorio resume a execução: os movimentos aplicados e a pasta de backup.
type Relatorio struct {
	Movimentos []Movimento

	Backup string
}

// Migrar planeja, confere, copia o backup e só então renomeia. Qualquer conflito
// detectado na conferência retorna erro antes de tocar o disco.
func Migrar(raiz string) (Relatorio, error) {
	raiz, err := validarRaiz(raiz)
	if err != nil {
		return Relatorio{}, err
	}
	plano, err := Planejar(raiz)
	if err != nil || len(plano) == 0 {
		return Relatorio{}, err
	}
	backup := filepath.Join(raiz, DirBackup)
	if err := conferir(raiz, backup, plano); err != nil {
		return Relatorio{}, err
	}
	if err := copiarBackup(raiz, backup, plano); err != nil {
		return Relatorio{}, err
	}
	if err := aplicar(raiz, plano); err != nil {
		return Relatorio{}, err
	}
	return Relatorio{Movimentos: plano, Backup: backup}, nil
}

// validarRaiz recusa raiz vazia, com "..", symlink ou que não seja pasta.
func validarRaiz(raiz string) (string, error) {
	if strings.TrimSpace(raiz) == "" {
		return "", fmt.Errorf("%w: vazia", ErrRaizInvalida)
	}
	for _, parte := range strings.Split(filepath.ToSlash(raiz), "/") {
		if parte == ".." {
			return "", fmt.Errorf("%w: %q contém ..", ErrRaizInvalida, raiz)
		}
	}
	abs, err := filepath.Abs(raiz)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRaizInvalida, err)
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRaizInvalida, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("%w: %q", ErrSymlink, raiz)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: %q não é pasta", ErrRaizInvalida, raiz)
	}
	return abs, nil
}

// Planejar lista os movimentos pendentes sob a raiz, sem alterar o disco.
// Pastas symlink não são seguidas; a pasta de backup é ignorada.
func Planejar(raiz string) ([]Movimento, error) {
	entradas, err := os.ReadDir(raiz)
	if err != nil {
		return nil, fmt.Errorf("migrar: ler %q: %w", raiz, err)
	}
	var plano []Movimento
	for _, e := range entradas {
		m, ok, err := planejarEntrada(raiz, e)
		if err != nil {
			return nil, err
		}
		if ok {
			plano = append(plano, m)
		}
	}
	return plano, nil
}

func planejarEntrada(raiz string, e os.DirEntry) (Movimento, bool, error) {
	nome := e.Name()
	symlink := e.Type()&os.ModeSymlink != 0
	switch {
	case nome == DirBackup:
		return Movimento{}, false, nil
	case symlink && (strings.HasSuffix(nome, SufixoV1) || nome == caminhos.ArquivoDAG):
		return Movimento{}, false, fmt.Errorf("%w: %q", ErrSymlink, nome)
	case symlink:
		return Movimento{}, false, nil
	case e.Type().IsRegular() && strings.HasSuffix(nome, SufixoV1) && len(nome) > len(SufixoV1):
		return Movimento{Origem: nome, Destino: strings.TrimSuffix(nome, SufixoV1) + SufixoStudio}, true, nil
	case e.Type().IsRegular() && nome == caminhos.ArquivoDAG:
		return Movimento{Origem: nome, Destino: caminhos.ArquivoProcessoDAG}, true, nil
	case e.IsDir():
		return planejarPasta(raiz, nome)
	}
	return Movimento{}, false, nil
}

// planejarPasta procura pipeline.dag.json dentro de uma pasta de Processo.
func planejarPasta(raiz, pasta string) (Movimento, bool, error) {
	origem := pasta + "/" + caminhos.ArquivoDAG
	info, err := os.Lstat(filepath.Join(raiz, filepath.FromSlash(origem)))
	if errors.Is(err, os.ErrNotExist) {
		return Movimento{}, false, nil
	}
	if err != nil {
		return Movimento{}, false, fmt.Errorf("migrar: %q: %w", origem, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return Movimento{}, false, fmt.Errorf("%w: %q", ErrSymlink, origem)
	}
	if !info.Mode().IsRegular() {
		return Movimento{}, false, nil
	}
	return Movimento{Origem: origem, Destino: pasta + "/" + caminhos.ArquivoProcessoDAG}, true, nil
}

// conferir garante que nenhum destino nem cópia de backup já exista.
func conferir(raiz, backup string, plano []Movimento) error {
	for _, m := range plano {
		alvos := []string{
			filepath.Join(raiz, filepath.FromSlash(m.Destino)),
			filepath.Join(backup, filepath.FromSlash(m.Origem)),
		}
		for _, alvo := range alvos {
			_, err := os.Lstat(alvo)
			if err == nil {
				return fmt.Errorf("%w: %q", ErrDestinoExiste, alvo)
			}
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("migrar: %q: %w", alvo, err)
			}
		}
	}
	return nil
}

func copiarBackup(raiz, backup string, plano []Movimento) error {
	for _, m := range plano {
		origem := filepath.Join(raiz, filepath.FromSlash(m.Origem))
		destino := filepath.Join(backup, filepath.FromSlash(m.Origem))
		if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
			return fmt.Errorf("migrar: backup: %w", err)
		}
		if err := copiarArquivo(origem, destino); err != nil {
			return fmt.Errorf("migrar: backup de %q: %w", m.Origem, err)
		}
	}
	return nil
}

// copiarArquivo copia sem sobrescrever (O_EXCL), preservando as permissões.
func copiarArquivo(origem, destino string) error {
	in, err := os.Open(origem)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(destino, os.O_WRONLY|os.O_CREATE|os.O_EXCL, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func aplicar(raiz string, plano []Movimento) error {
	for _, m := range plano {
		origem := filepath.Join(raiz, filepath.FromSlash(m.Origem))
		destino := filepath.Join(raiz, filepath.FromSlash(m.Destino))
		if err := os.Rename(origem, destino); err != nil {
			return fmt.Errorf("migrar: renomear %q: %w", m.Origem, err)
		}
	}
	return nil
}
