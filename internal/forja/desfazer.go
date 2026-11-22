package forja

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/jail"
)

func executarDesfazer(ctx context.Context, op Opcao, rel Relatorio) (Relatorio, error) {
	if err := jail.ValidarNomeProcesso(op.Desfazer); err != nil {
		rel.Veredito = VereditoFalha
		return rel, err
	}
	dirData := filepath.Join(op.Home, ".aidd", "backups", "forja_"+op.Desfazer)
	entradasDir, err := os.ReadDir(dirData)
	if err != nil {
		rel.Veredito = VereditoFalha
		if os.IsNotExist(err) {
			return rel, fmt.Errorf("%w: %s", ErrBackupInexistente, op.Desfazer)
		}
		return rel, err
	}
	harnesses := listarSubpastasOrdenadas(entradasDir)
	for _, h := range harnesses {
		if err := ctx.Err(); err != nil {
			rel.Veredito = VereditoFalha
			return rel, err
		}
		etapa, err := restaurarHarness(op.Home, filepath.Join(dirData, h), h)
		rel.Etapas = append(rel.Etapas, etapa)
		if err != nil {
			rel.Veredito = VereditoFalha
			return rel, err
		}
	}
	limparPastaBackup(op.Home, dirData)
	return rel, nil
}

func listarSubpastasOrdenadas(entradas []os.DirEntry) []string {
	var out []string
	for _, e := range entradas {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func restaurarHarness(home, dirBackup, harness string) (EtapaResultado, error) {
	etapa := EtapaResultado{Nome: "desfazer:" + harness, Harness: harness, Backup: dirBackup}
	entradas, err := lerManifesto(filepath.Join(dirBackup, "MANIFESTO.tsv"))
	if err != nil {
		etapa.Estado = EstadoFalha
		etapa.Mensagem = err.Error()
		return etapa, err
	}
	var restaurados []string
	for i := len(entradas) - 1; i >= 0; i-- {
		if err := aplicarRestauracao(home, dirBackup, entradas[i]); err != nil {
			etapa.Estado = EstadoFalha
			etapa.Mensagem = err.Error()
			return etapa, err
		}
		if entradas[i].acao != acaoDiretorioCriado {
			restaurados = append(restaurados, entradas[i].origem)
		}
	}
	etapa.Estado = EstadoOK
	etapa.Alterados = restaurados
	etapa.Mensagem = "restaurado"
	return etapa, nil
}

func lerManifesto(caminho string) ([]entradaManifesto, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	linhas := strings.Split(strings.TrimSpace(string(dados)), "\n")
	if len(linhas) == 0 || strings.TrimSpace(linhas[0]) != strings.TrimSpace(cabecalhoManifesto) {
		return nil, fmt.Errorf("%w: cabeçalho ausente em %s", ErrManifestoInvalido, caminho)
	}
	var out []entradaManifesto
	for _, linha := range linhas[1:] {
		linha = strings.TrimRight(linha, "\r")
		if strings.TrimSpace(linha) == "" {
			continue
		}
		cols := strings.Split(linha, "\t")
		if len(cols) != 4 {
			return nil, fmt.Errorf("%w: linha com %d colunas", ErrManifestoInvalido, len(cols))
		}
		out = append(out, entradaManifesto{acao: cols[0], origem: cols[1], backup: cols[2], modo: cols[3]})
	}
	return out, nil
}

func aplicarRestauracao(home, dirBackup string, ent entradaManifesto) error {
	destino, err := jail.DentroDaRaiz(home, ent.origem)
	if err != nil {
		return err
	}
	switch ent.acao {
	case acaoModificado:
		return restaurarArquivoModificado(dirBackup, destino, ent)
	case acaoCriado:
		if err := os.Remove(destino); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	case acaoDiretorioCriado:
		if err := os.Remove(destino); err != nil && !os.IsNotExist(err) {
			return nil
		}
		return nil
	default:
		return fmt.Errorf("%w: ação desconhecida %q", ErrManifestoInvalido, ent.acao)
	}
}

func restaurarArquivoModificado(dirBackup, destino string, ent entradaManifesto) error {
	origemBk, err := jail.DentroDaRaiz(dirBackup, ent.backup)
	if err != nil {
		return err
	}
	dados, err := os.ReadFile(origemBk)
	if err != nil {
		return err
	}
	modo := fs.FileMode(0o644)
	if parsed, err := strconv.ParseUint(ent.modo, 8, 32); err == nil && parsed != 0 {
		modo = fs.FileMode(parsed)
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destino, dados, modo)
}

func limparPastaBackup(home, dirData string) {
	if err := os.RemoveAll(dirData); err != nil {
		return
	}
	if err := os.Remove(filepath.Join(home, ".aidd", "backups")); err != nil {
		return
	}
	if err := os.Remove(filepath.Join(home, ".aidd")); err != nil {
		return
	}
}
