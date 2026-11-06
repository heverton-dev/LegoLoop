package ste_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/ste"
)

func TestSyncHarnesses(t *testing.T) {
	t.Parallel()

	tmpHome := t.TempDir()

	// Cria estrutura prévia de teste para Gemini e Claude
	geminiRulesDir := filepath.Join(tmpHome, ".gemini", "config", "rules")
	claudeRulesDir := filepath.Join(tmpHome, ".claude", "rules")
	if err := os.MkdirAll(geminiRulesDir, 0755); err != nil {
		t.Fatalf("falha ao criar pasta gemini rules: %v", err)
	}
	if err := os.MkdirAll(claudeRulesDir, 0755); err != nil {
		t.Fatalf("falha ao criar pasta claude rules: %v", err)
	}

	// Cria arquivo do usuário pré-existente para provar não-invasividade
	userOtherRule := filepath.Join(geminiRulesDir, "other-user-rule.md")
	if err := os.WriteFile(userOtherRule, []byte("# Minha Regra"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo alheio do usuario: %v", err)
	}

	cfg := ste.NovoConfig()
	cfg.Padrao = ste.ConfigPerfil{
		Modo:     ste.ModoHibrido,
		Pausa:    ste.PausaDestrutivo,
		Thinking: "caveman",
	}

	// Sincroniza (Ativação)
	statusList, err := ste.Sync(cfg, tmpHome, false)
	if err != nil {
		t.Fatalf("sync falhou: %v", err)
	}

	geminiSteptRule := filepath.Join(geminiRulesDir, "stept-rules.md")
	if _, err := os.Stat(geminiSteptRule); os.IsNotExist(err) {
		t.Errorf("arquivo de regras stept nao foi gerado para gemini: %s", geminiSteptRule)
	}

	// Prova não-invasividade: o outro arquivo do usuário continua intacto
	if _, err := os.Stat(userOtherRule); os.IsNotExist(err) {
		t.Errorf("arquivo alheio do usuario foi indevidamente deletado!")
	}

	// Verifica se OMP não instalado não teve diretórios criados (sem createDirs)
	ompDir := filepath.Join(tmpHome, ".omp")
	if _, err := os.Stat(ompDir); !os.IsNotExist(err) {
		t.Errorf("diretorio fantasma foi criado para omp: %s", ompDir)
	}

	// Verifica status reportados
	var geminiStatus ste.HarnessStatus
	for _, s := range statusList {
		if s.ID == ste.HarnessGemini {
			geminiStatus = s
			break
		}
	}
	if !geminiStatus.Ativo || geminiStatus.Perfil.Modo != ste.ModoHibrido {
		t.Errorf("status inesperado para gemini: %+v", geminiStatus)
	}

	// Agora desativação (ModoOff)
	cfg.Padrao.Modo = ste.ModoOff
	statusListOff, err := ste.Sync(cfg, tmpHome, false)
	if err != nil {
		t.Fatalf("sync off falhou: %v", err)
	}

	// Prova remoção limpa: stept-rules.md foi deletado
	if _, err := os.Stat(geminiSteptRule); !os.IsNotExist(err) {
		t.Errorf("stept-rules.md deveria ter sido removido no modo off")
	}

	// Prova que pasta pai e arquivo alheio continuam intactos
	if _, err := os.Stat(userOtherRule); os.IsNotExist(err) {
		t.Errorf("arquivo alheio do usuario sumiu na desativacao!")
	}

	for _, s := range statusListOff {
		if s.ID == ste.HarnessGemini && s.Ativo {
			t.Errorf("status de gemini deveria ser Ativo=false apos modo off")
		}
	}
}
