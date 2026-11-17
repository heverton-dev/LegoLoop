package fechar

import (
	"encoding/json"
	"os"
	"path/filepath"

	"legoloop/internal/triagem"
)

func ehDiretorio(caminho string) bool {
	info, err := os.Stat(caminho)
	return err == nil && info.IsDir()
}

// resolverDirLivro prefere docs/books/<canon> quando a pasta já existe;
// caso contrário mantém docs/books/<slug>.
func resolverDirLivro(repoDir, slug, canon string) string {
	base := filepath.Join(repoDir, "docs", "books")
	if canon != "" && ehDiretorio(filepath.Join(base, canon)) {
		return filepath.Join(base, canon)
	}
	return filepath.Join(base, slug)
}

type dagNos struct {
	Nodes []struct {
		ID string `json:"id"`
	} `json:"nodes"`
}

// garantirQuarteto cria os arquivos padrão ausentes do Processo (DAG, triagem,
// memória) e aplica o zeramento seletivo da memória: zera só a sessão dona e
// limpa os órfãos.
func garantirQuarteto(dir, dagPath, sessaoID string) error {
	dados, err := os.ReadFile(dagPath)
	if err != nil {
		return err
	}
	nome := filepath.Base(dir)
	if err := garantirTriagem(dir, nome, dados); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(dir, triagem.ArquivoMemoria)); os.IsNotExist(err) {
		if err := triagem.SalvarMemoria(dir, triagem.MemoriaExecucao{Processo: nome}); err != nil {
			return err
		}
	}
	return zerarMemoriaSessao(dir, sessaoID)
}

// zerarMemoriaSessao zera a sessão dona do fechar e limpa sessão e lock órfãos.
func zerarMemoriaSessao(dir, sessaoID string) error {
	if sessaoID != "" {
		if err := triagem.ZerarSessao(dir, sessaoID); err != nil {
			return err
		}
	}
	_, err := triagem.LimparOrfaos(dir)
	return err
}

// garantirTriagem aplica a regra: um nó só não tem triagem; dois ou mais exigem uma triagem que cubra todos.
// A ausente nasce, a existente ganha só as regras que faltam, e no fim o Conferir prova a cobertura.
func garantirTriagem(dir, nome string, dagJSON []byte) error {
	var dag dagNos
	_ = json.Unmarshal(dagJSON, &dag)
	var etapas []string
	for _, n := range dag.Nodes {
		if n.ID != "" {
			etapas = append(etapas, n.ID)
		}
	}
	if !triagem.Exigida(len(etapas)) {
		return nil
	}
	path := filepath.Join(dir, triagem.ArquivoTriagem)
	var atual triagem.TriagemSpec
	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &atual); err != nil {
			return err
		}
	}
	saida, err := json.MarshalIndent(triagem.Completar(atual, nome, etapas), "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append(saida, '\n'), 0o644); err != nil {
		return err
	}
	return triagem.Conferir(dir, etapas)
}
