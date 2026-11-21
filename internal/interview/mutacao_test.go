package interview

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestClientDecodeStatus300EURLInvalida(t *testing.T) {
	resp := &http.Response{
		StatusCode: 300,
		Body:       io.NopCloser(strings.NewReader(`{"choices":[{"message":{"content":"ok"}}]}`)),
	}
	if _, err := decodeCompletion(resp); err == nil {
		t.Fatal("esperava erro para status 300")
	}
	c := NewRouterClient("://url-invalida", "", "m")
	if _, err := c.newChatRequest(context.Background(), nil); err == nil {
		t.Fatal("esperava erro em newChatRequest com URL inválida")
	}
	if _, err := c.Chat(context.Background(), nil); err == nil {
		t.Fatal("esperava erro em Chat com URL inválida")
	}
}

func TestContextEmptyCutUTF8ERender(t *testing.T) {
	if (Context{Pieces: []string{"p1"}}).Empty() {
		t.Fatal("Context com Pieces não pode ser Empty")
	}
	if (Context{Docs: []Snippet{{Name: "a"}}}).Empty() {
		t.Fatal("Context com Docs não pode ser Empty")
	}
	if got, cut := cutUTF8("abcd", 4); got != "abcd" || cut {
		t.Fatalf("cutUTF8 limite exato falhou: got=%q cut=%v", got, cut)
	}
	if got, cut := cutUTF8("abcde", 3); got != "abc" || !cut {
		t.Fatalf("cutUTF8 ASCII falhou: got=%q cut=%v", got, cut)
	}
	if got, cut := cutUTF8("\x80\x80", 1); got != "" || !cut {
		t.Fatalf("cutUTF8 continuation bytes falhou: got=%q cut=%v", got, cut)
	}
	if strings.Contains((Context{Goal: "g"}).Render(), "Peças aptas já catalogadas:") {
		t.Fatal("Render não deve incluir seção de peças quando Pieces está vazio")
	}
}

type mockLLMErroSeq struct {
	resps []string
	errs  []error
	idx   int
}

func (m *mockLLMErroSeq) Chat(_ context.Context, _ []Message) (string, error) {
	i := m.idx
	m.idx++
	if i < len(m.errs) && m.errs[i] != nil {
		return "", m.errs[i]
	}
	if i < len(m.resps) {
		return m.resps[i], nil
	}
	return "", nil
}

func TestEngineErrosStartExchangeMaterialize(t *testing.T) {
	engErr := NewEngine(&mockLLMErroSeq{errs: []error{errors.New("falha llm")}})
	if _, _, err := engErr.Start(context.Background(), "ideia"); err == nil {
		t.Fatal("Start deveria propagar erro da LLM")
	}
	if _, err := engErr.MaterializeIntent("intv_999", t.TempDir(), ""); err == nil {
		t.Fatal("MaterializeIntent deveria falhar para sessão inexistente")
	}
	engImpasse := NewEngine(&mockLLMErroSeq{
		resps: []string{"CONFIANÇA: 96%\nHIPÓTESE: ok"},
		errs:  []error{nil, errors.New("falha premonicao")},
	})
	if _, _, err := engImpasse.Start(context.Background(), "ideia"); err == nil {
		t.Fatal("exchange deveria propagar erro de pedirPremonicao")
	}
}

func TestEngineIntentTargetCaptureParseESanitize(t *testing.T) {
	dir := t.TempDir()
	p, _ := intentTarget(dir, "x.md", "tema", "2026-01-01")
	if filepath.Dir(p) != filepath.Join(dir, "docs", "intents") {
		t.Fatalf("esperava docs/intents, obteve %q", p)
	}
	dirArq := t.TempDir()
	_ = os.WriteFile(filepath.Join(dirArq, "docs"), []byte("x"), 0o644)
	p2, _ := intentTarget(dirArq, "x.md", "tema", "2026-01-01")
	if filepath.Dir(p2) != dirArq {
		t.Fatalf("esperava fallback para baseDir, obteve %q", p2)
	}
	reSemGrupo := regexp.MustCompile("abc")
	if _, ok := capture(reSemGrupo, "abc"); ok {
		t.Fatal("capture sem grupo deveria retornar false")
	}
	if v, ok := capture(reQuestion, "PERGUNTA: qual?"); !ok || v != "qual?" {
		t.Fatalf("capture falhou: %q %v", v, ok)
	}
	if _, ok := capture(reQuestion, "nada"); ok {
		t.Fatal("capture sem match deveria retornar false")
	}
	if _, ok := captureLast(reConfidence, "nada"); ok {
		t.Fatal("captureLast sem match deveria retornar false")
	}
	if _, ok := captureLast(reSemGrupo, "abc"); ok {
		t.Fatal("captureLast sem grupo deveria retornar false")
	}
}

func TestEngineParse95SemTituloESanitize(t *testing.T) {
	raw := "CONFIANÇA: 95%\n- Resultado Esperado: Entregar valor\n"
	turn := ParseTurnResult(raw)
	if turn.Status != StatusCompleted || turn.Mirror == nil {
		t.Fatalf("esperava StatusCompleted com 95%% e espelho, obteve status=%s mirror=%v", turn.Status, turn.Mirror)
	}
	if got := sanitizeFileName("!!!"); got != "NOVA_INTENCAO" {
		t.Fatalf("esperava NOVA_INTENCAO, obteve %q", got)
	}
	if got := sanitizeFileName("abc"); got != "ABC" {
		t.Fatalf("esperava ABC, obteve %q", got)
	}
}

func TestEspelhoCamposFronteiras(t *testing.T) {
	raw := "intro ignorada\n- Resultado Esperado: Valor A\nCONFIANÇA: 90%\nresto ignorado\n- Critério de Sucesso:\n\n  - item 1\n"
	c := camposEspelho(raw)
	if len(c["resultado esperado"]) != 1 || c["resultado esperado"][0] != "Valor A" {
		t.Fatalf("resultado esperado incorreto: %v", c["resultado esperado"])
	}
	if len(c["critério de sucesso"]) != 1 || c["critério de sucesso"][0] != "item 1" {
		t.Fatalf("critério de sucesso após linha em branco incorreto: %v", c["critério de sucesso"])
	}
}

func TestImpasseSemEspelhoEPedirPremonicao(t *testing.T) {
	if !semEspelho(&TurnResult{Status: StatusInProgress, Mirror: nil, Confidence: 95}) {
		t.Fatal("semEspelho deve ser true em 95%%")
	}
	if semEspelho(&TurnResult{Status: StatusInProgress, Mirror: &Mirror{Outcome: "ok"}, Confidence: 95}) {
		t.Fatal("semEspelho deve ser false com Mirror")
	}
	if semEspelho(&TurnResult{Status: StatusInProgress, Mirror: nil, Confidence: 50}) {
		t.Fatal("semEspelho deve ser false com confiança 50")
	}
	engFim := NewEngine(&mockLLMErroSeq{resps: []string{"ESPELHO CANÔNICO\n- Resultado Esperado: X\nCONFIANÇA: 96%"}})
	turnFim, err := engFim.pedirPremonicao(context.Background(), &Session{}, &TurnResult{Confidence: 96})
	if err != nil || turnFim.Status != StatusCompleted || turnFim.Question != "" {
		t.Fatalf("conclusão não deve preencher Question: err=%v status=%s q=%q", err, turnFim.Status, turnFim.Question)
	}
	engSemQ := NewEngine(&mockLLMErroSeq{resps: []string{"CONFIANÇA: 50%\nHIPÓTESE: h"}})
	turnSemQ, err := engSemQ.pedirPremonicao(context.Background(), &Session{}, &TurnResult{Confidence: 80})
	if err != nil || turnSemQ.Question != perguntaDeConfirmacao {
		t.Fatalf("esperava fallback sem pergunta: err=%v q=%q", err, turnSemQ.Question)
	}
	engSemEspQ := NewEngine(&mockLLMErroSeq{resps: []string{"CONFIANÇA: 97%\nPERGUNTA: extra?"}})
	turnSemEspQ, err := engSemEspQ.pedirPremonicao(context.Background(), &Session{}, &TurnResult{Confidence: 96})
	if err != nil || turnSemEspQ.Question != perguntaSemEspelho {
		t.Fatalf("esperava fallback perguntaSemEspelho: err=%v q=%q", err, turnSemEspQ.Question)
	}
}

func TestIntentRenderETextoLimites(t *testing.T) {
	if intentStatus(95) != "UNDER_REVIEW" {
		t.Fatalf("esperava UNDER_REVIEW em 95, obteve %q", intentStatus(95))
	}
	if !ehSoftware("") || !ehSoftware("software") || ehSoftware("juridico") {
		t.Fatal("ehSoftware falhou nos casos vazio, software e juridico")
	}
	if openRisks(100) != "Nenhum registrado." {
		t.Fatalf("openRisks(100) inesperado: %q", openRisks(100))
	}
	if specBase("/a.md") != "a.md" || specBase("docs/specs/a.md") != "a.md" {
		t.Fatal("specBase falhou")
	}
	if got := tituloDoTema(". segunda frase"); got != "." {
		t.Fatalf("tituloDoTema('. segunda frase') = %q", got)
	}
	if got := tituloDoTema("Primeira. Segunda"); got != "Primeira." {
		t.Fatalf("tituloDoTema('Primeira. Segunda') = %q", got)
	}
	exato := strings.Repeat("a", tituloMax)
	if got := tituloDoTema(exato); got != exato {
		t.Fatalf("tituloDoTema com %d chars não deve truncar: %q", tituloMax, got)
	}
	if got := criteriosNumerados("unico"); got != "unico" {
		t.Fatalf("criteriosNumerados('unico') = %q", got)
	}
	if got := criteriosNumerados("a; b"); got != "(1) a (2) b" {
		t.Fatalf("criteriosNumerados('a; b') = %q", got)
	}
}

func TestStoreSaveELoadErros(t *testing.T) {
	arqPai := filepath.Join(t.TempDir(), "arquivo")
	_ = os.WriteFile(arqPai, []byte("x"), 0o644)
	err := (Store{Dir: filepath.Join(arqPai, "sub")}).Save(&Session{ID: "intv_1"})
	var pe *os.PathError
	if err == nil || !errors.As(err, &pe) || pe.Op != "mkdir" {
		t.Fatalf("Save deveria falhar com PathError mkdir, obteve %v", err)
	}
	s := Store{Dir: t.TempDir()}
	sessInvalida := &Session{ID: "intv_1", CreatedAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := s.Save(sessInvalida); err == nil {
		t.Fatal("Save deveria falhar ao serializar ano 10000")
	}
	if err := os.Mkdir(filepath.Join(s.Dir, "intv_2.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(&Session{ID: "intv_2"}); err == nil {
		t.Fatal("Save deveria falhar quando destino é diretório")
	}
	if _, err := s.Load("intv_2"); err == nil || strings.Contains(err.Error(), "corrompida") {
		t.Fatalf("Load deveria falhar com erro de leitura, obteve %v", err)
	}
}
