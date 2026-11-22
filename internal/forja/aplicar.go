package forja

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	acaoModificado      = "modificado"
	acaoCriado          = "criado"
	acaoDiretorioCriado = "diretorio_criado"
	cabecalhoManifesto  = "acao\torigem\tbackup\tmodo\n"
)

type entradaManifesto struct {
	acao   string
	origem string
	backup string
	modo   string
}

func pastaBackupHarness(home, data, harness string) string {
	return filepath.Join(home, ".aidd", "backups", "forja_"+data, harness)
}

func aplicarEtapa(op Opcao, h string, pendentes []AlteracaoPlanejada) (EtapaResultado, error) {
	etapa := EtapaResultado{Nome: "forjar:" + h, Harness: h, Estado: EstadoOK}
	if len(pendentes) == 0 {
		etapa.Mensagem = "sem alterações"
		return etapa, nil
	}
	dirBackup := pastaBackupHarness(op.Home, op.Data, h)
	if err := os.MkdirAll(dirBackup, 0o755); err != nil {
		etapa.Estado = EstadoFalha
		etapa.Mensagem = err.Error()
		return etapa, err
	}
	entradas, err := gravarComBackup(op.Home, dirBackup, pendentes)
	if err != nil {
		etapa.Estado = EstadoFalha
		etapa.Mensagem = err.Error()
		return etapa, err
	}
	if err := escreverManifesto(dirBackup, entradas); err != nil {
		etapa.Estado = EstadoFalha
		etapa.Mensagem = err.Error()
		return etapa, err
	}
	etapa.Alterados = caminhosRelativos(op.Home, pendentes)
	etapa.Backup = dirBackup
	etapa.Mensagem = "aplicado"
	return etapa, nil
}

func gravarComBackup(home, dirBackup string, pendentes []AlteracaoPlanejada) ([]entradaManifesto, error) {
	var entradas []entradaManifesto
	seq := 0
	for _, alt := range pendentes {
		dirsCriados, err := garantirPaisRegistrando(home, filepath.Dir(alt.Caminho))
		if err != nil {
			return nil, err
		}
		for _, d := range dirsCriados {
			entradas = append(entradas, entradaManifesto{
				acao: acaoDiretorioCriado, origem: d, backup: "-", modo: "0",
			})
		}
		seq++
		ent, err := salvarEGravarArquivo(home, dirBackup, seq, alt)
		if err != nil {
			return nil, err
		}
		entradas = append(entradas, ent)
	}
	return entradas, nil
}

func garantirPaisRegistrando(home, dir string) ([]string, error) {
	absHome, err := filepath.Abs(home)
	if err != nil {
		return nil, err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	var faltantes []string
	for atual := absDir; atual != absHome && len(atual) > len(absHome); atual = filepath.Dir(atual) {
		if _, err := os.Stat(atual); err == nil {
			break
		}
		rel, err := filepath.Rel(absHome, atual)
		if err != nil {
			return nil, err
		}
		faltantes = append([]string{filepath.ToSlash(rel)}, faltantes...)
	}
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return nil, err
	}
	return faltantes, nil
}

func salvarEGravarArquivo(home, dirBackup string, seq int, alt AlteracaoPlanejada) (entradaManifesto, error) {
	rel, err := filepath.Rel(home, alt.Caminho)
	if err != nil {
		return entradaManifesto{}, err
	}
	relSlash := filepath.ToSlash(rel)
	info, statErr := os.Stat(alt.Caminho)
	ent := entradaManifesto{acao: acaoCriado, origem: relSlash, backup: "-", modo: "0"}
	if statErr == nil {
		nomeBk := fmt.Sprintf("%03d_%s", seq, filepath.Base(alt.Caminho))
		dadosOrig, err := os.ReadFile(alt.Caminho)
		if err != nil {
			return entradaManifesto{}, err
		}
		if err := os.WriteFile(filepath.Join(dirBackup, nomeBk), dadosOrig, 0o644); err != nil {
			return entradaManifesto{}, err
		}
		ent = entradaManifesto{
			acao:   acaoModificado,
			origem: relSlash,
			backup: nomeBk,
			modo:   fmt.Sprintf("%04o", info.Mode().Perm()),
		}
	}
	if err := os.WriteFile(alt.Caminho, alt.Conteudo, alt.Modo.Perm()); err != nil {
		return entradaManifesto{}, err
	}
	return ent, nil
}

func escreverManifesto(dirBackup string, entradas []entradaManifesto) error {
	var sb strings.Builder
	sb.WriteString(cabecalhoManifesto)
	for _, e := range entradas {
		sb.WriteString(e.acao)
		sb.WriteByte('\t')
		sb.WriteString(e.origem)
		sb.WriteByte('\t')
		sb.WriteString(e.backup)
		sb.WriteByte('\t')
		sb.WriteString(e.modo)
		sb.WriteByte('\n')
	}
	return os.WriteFile(filepath.Join(dirBackup, "MANIFESTO.tsv"), []byte(sb.String()), 0o644)
}
