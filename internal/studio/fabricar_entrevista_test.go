package studio

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/provider"
)

type mockStudioEntrevistaLLM struct {
	mu        sync.Mutex
	callCount int
}

func (m *mockStudioEntrevistaLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	for _, msg := range msgs {
		if strings.Contains(msg.Content, "Arquiteto de Especificações") {
			return studioSpec, nil
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.callCount == 1 {
		return "HIPÓTESE: extrator de áudio\nCONFIANÇA: 60%\nPERGUNTA: Qual formato?\nPALPITE: mp3", nil
	}
	return "HIPÓTESE: extrator de áudio\nCONFIANÇA: 98%\nESPELHO CANÔNICO:\n- Resultado Esperado: extrair áudio em mp3\n- Quem se Beneficia: editores sobrecarregados\n- Por que Agora: retrabalho manual diário\n- Critério de Sucesso: arquivo mp3 gerado com sucesso\n- Restrições Inegociáveis: sem conexão com a rede\n- Fora de Escopo: edição e corte de vídeo\n", nil
}

func TestStudioEntrevista(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	llm := &mockStudioEntrevistaLLM{}
	s.Interview = interview.NewEngine(llm)
	s.Discover = func(context.Context) provider.Inventory {
		return provider.Inventory{
			RouterURL:    "http://127.0.0.1:8000",
			RouterModels: []string{"code-pro"},
			CLIs:         []provider.CLI{{Name: "claude"}},
		}
	}
	s.FabricarMatriz = studioMatriz()

	vault, err := s.vault()
	if err != nil {
		t.Fatal(err)
	}

	h := s.Handler()

	// 1. Inicia entrevista
	startRec := do(t, h, http.MethodPost, "/api/interview/start", `{"idea":"extrator de audio"}`)
	if startRec.Code != http.StatusOK {
		t.Fatalf("esperava 200 no start, obteve %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		Session struct {
			ID string `json:"id"`
		} `json:"session"`
	}
	decode(t, startRec, &startResp)
	sessID := startResp.Session.ID
	if sessID == "" {
		t.Fatal("sessID não gerado")
	}

	// 2. Envia mensagem para completar entrevista
	msgRec := do(t, h, http.MethodPost, "/api/interview/message", `{"session_id":"`+sessID+`","answer":"tudo certo"}`)
	if msgRec.Code != http.StatusOK {
		t.Fatalf("esperava 200 no message, obteve %d: %s", msgRec.Code, msgRec.Body.String())
	}

	// 3. Configura s.Fabricar para associar a sessão carregada
	s.Fabricar = func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		deps := fabricar.Deps{
			LLM:    llm,
			Matriz: studioMatriz(),
			Clock:  func() time.Time { return studioNow },
			Roots:  fabricar.Roots{Vault: vault, Docs: s.ProjectDir},
		}
		if sess, err := s.Interview.GetSession(sessID); err == nil && sess != nil {
			clar := fabricar.ClarificationFromSession(sess)
			deps.Interview = &clar
		}
		return fabricar.Fabricar(ctx, req, deps)
	}

	// 4. Envia fabricação com o identificador da entrevista
	fabRec := do(t, h, http.MethodPost, "/api/fabricar/run", `{
		"name": "extrator de audio",
		"model": "code-pro",
		"category": "build",
		"entrevista": "`+sessID+`"
	}`)
	if fabRec.Code != http.StatusOK {
		t.Fatalf("esperava 200 no fabricar/run, obteve %d: %s", fabRec.Code, fabRec.Body.String())
	}

	// 5. Verifica o INTENT gerado em docs/intents/
	intentsDir := filepath.Join(s.ProjectDir, "docs", "intents")
	entries, err := os.ReadDir(intentsDir)
	if err != nil {
		t.Fatalf("falha ao ler docs/intents: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("nenhum arquivo INTENT gerado")
	}

	intentContent, err := os.ReadFile(filepath.Join(intentsDir, entries[0].Name()))
	if err != nil {
		t.Fatalf("falha ao ler INTENT: %v", err)
	}
	intentStr := string(intentContent)

	for _, want := range []string{
		"editores sobrecarregados",
		"retrabalho manual diário",
		"sem conexão com a rede",
		"edição e corte de vídeo",
	} {
		if !strings.Contains(intentStr, want) {
			t.Errorf("INTENT não contém %q do espelho:\n%s", want, intentStr)
		}
	}
	if strings.Contains(intentStr, "Não informado") {
		t.Errorf("INTENT não deve conter 'Não informado' quando vem do espelho:\n%s", intentStr)
	}
}
