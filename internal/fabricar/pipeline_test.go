package fabricar_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/reputation"
)

func validRequest() fabricar.Request {
	return fabricar.Request{RawName: "Extrator de Áudio", Goal: "extrair áudio", Category: fabricar.CategoryBuild, ModelChoice: "sonnet"}
}

func TestRequestValidate(t *testing.T) {
	if err := validRequest().Validate(); err != nil {
		t.Fatalf("pedido válido rejeitado: %v", err)
	}
	cases := []struct {
		name   string
		mutate func(*fabricar.Request)
		want   error
	}{
		{"nome vazio", func(r *fabricar.Request) { r.RawName = "" }, fabricar.ErrInvalidName},
		{"nome ideograma", func(r *fabricar.Request) { r.RawName = "日本語" }, fabricar.ErrInvalidName},
		{"categoria desconhecida", func(r *fabricar.Request) { r.Category = "fisica" }, fabricar.ErrInvalidRequest},
		{"categoria vazia", func(r *fabricar.Request) { r.Category = "" }, fabricar.ErrInvalidRequest},
		{"modelo vazio", func(r *fabricar.Request) { r.ModelChoice = "  " }, fabricar.ErrInvalidRequest},
	}
	for _, c := range cases {
		r := validRequest()
		c.mutate(&r)
		if err := r.Validate(); !errors.Is(err, c.want) {
			t.Errorf("%s: erro = %v; quer %v", c.name, err, c.want)
		}
	}
	for _, c := range []fabricar.Category{fabricar.CategoryBuild, fabricar.CategoryReverse, fabricar.CategoryWeb, fabricar.CategoryMobile} {
		r := validRequest()
		r.Category = c
		if err := r.Validate(); err != nil {
			t.Errorf("categoria %q deve ser válida: %v", c, err)
		}
	}
}

// scriptedLLM responde em sequência e conta as chamadas.
type scriptedLLM struct {
	replies []string
	calls   int
}

func (s *scriptedLLM) Chat(context.Context, []interview.Message) (string, error) {
	i := min(s.calls, len(s.replies)-1)
	s.calls++
	return s.replies[i], nil
}

const (
	lowConfidence = "HIPÓTESE: algo\nCONFIANÇA: 50%\nPERGUNTA: Qual o formato?\nPALPITE: mp3"
	doneReply     = "HIPÓTESE: extrator\nCONFIANÇA: 98%\nESPELHO CANÔNICO:\n- Resultado Esperado: extrair áudio em mp3\n- Quem se Beneficia: editores\n- Por que Agora: retrabalho\n- Critério de Sucesso: arquivo gerado\n- Restrições Inegociáveis: sem rede\n- Fora de Escopo: edição de vídeo"
)

func TestClarifier_AutonomousDoesNotCallTheLLMAndUsesGoalOrName(t *testing.T) {
	llm := &scriptedLLM{replies: []string{lowConfidence}}
	c := fabricar.NewClarifier(llm)
	got, err := c.Clarify(context.Background(), validRequest(), nil)
	if err != nil || llm.calls != 0 {
		t.Fatalf("autônomo não pergunta: %v (chamadas: %d)", err, llm.calls)
	}
	if got.Mirror.Outcome != "extrair áudio" || got.Confidence != 0 || got.Reached {
		t.Errorf("deve usar o Goal sem inventar confiança: %+v", got)
	}
	req := validRequest()
	req.Goal = ""
	got, _ = c.Clarify(context.Background(), req, nil)
	if got.Mirror.Outcome != "Extrator de Áudio" {
		t.Errorf("sem Goal usa o próprio nome: %+v", got)
	}
}

func TestClarifier_InteractiveFunnelsUntilConfident(t *testing.T) {
	llm := &scriptedLLM{replies: []string{lowConfidence, lowConfidence, doneReply}}
	var questions []string
	ask := func(q, guess string) (string, error) {
		questions = append(questions, q+"|"+guess)
		return "resposta", nil
	}
	req := validRequest()
	req.Interactive = true
	got, err := fabricar.NewClarifier(llm).Clarify(context.Background(), req, ask)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Reached || got.Confidence < 95 || got.Mirror.Outcome != "extrair áudio em mp3" || got.Turns != 2 {
		t.Errorf("clarificação incorreta: %+v", got)
	}
	if len(questions) != 2 || questions[0] != "Qual o formato?|mp3" {
		t.Errorf("perguntas entregues ao usuário: %v", questions)
	}
}

func TestClarifier_StopsAtTheTurnLimitWithoutLooping(t *testing.T) {
	llm := &scriptedLLM{replies: []string{lowConfidence}}
	asked := 0
	ask := func(string, string) (string, error) { asked++; return "x", nil }
	c := fabricar.NewClarifier(llm)
	c.MaxTurns = 3
	req := validRequest()
	req.Interactive = true
	got, err := c.Clarify(context.Background(), req, ask)
	if err != nil || got.Reached || asked != 3 {
		t.Fatalf("deve parar no limite sem ≥95%%: %v reached=%v perguntas=%d", err, got.Reached, asked)
	}
	if got.Confidence != 50 {
		t.Errorf("a confiança real deve ser devolvida: %d", got.Confidence)
	}
}

func TestClarifier_PropagatesUserAbortAndLLMFailure(t *testing.T) {
	req := validRequest()
	req.Interactive = true
	boom := errors.New("usuário cancelou")
	_, err := fabricar.NewClarifier(&scriptedLLM{replies: []string{lowConfidence}}).Clarify(context.Background(), req, func(string, string) (string, error) { return "", boom })
	if !errors.Is(err, boom) {
		t.Errorf("o cancelamento do usuário deve abortar: %v", err)
	}
	if _, err := fabricar.NewClarifier(nil).Clarify(context.Background(), req, func(string, string) (string, error) { return "x", nil }); err == nil {
		t.Error("modo interativo sem LLM deve falhar")
	}
}

func piece(name, path string, score float64, untrusted bool) fabricar.DiscoveredPiece {
	return fabricar.DiscoveredPiece{Name: name, Kind: fabricar.KindSkill, Drive: "D:\\", Path: path, Score: score, Tier: reputation.TierOuro, Untrusted: untrusted}
}

func stepNames(m fabricar.Merged) []string {
	var out []string
	for _, s := range m.Steps {
		out = append(out, s.Piece)
	}
	return out
}

func TestArchitectMerger_SelectsByScoreThenPathAndExcludesUntrustedByDefault(t *testing.T) {
	pieces := []fabricar.DiscoveredPiece{
		piece("b", "D:\\b", 0.9, false), piece("a", "D:\\a", 0.9, false),
		piece("c", "D:\\c", 0.8, false), piece("pendrive", "E:\\p", 0.99, true),
	}
	m := fabricar.ArchitectMerge(pieces, fabricar.MergeOptions{})
	if !slices.Equal(stepNames(m), []string{"a", "b", "c"}) {
		t.Errorf("ordem = %v; quer Score desc, Path asc, sem Untrusted", stepNames(m))
	}
	if len(m.Excluded) != 1 || m.Excluded[0].Name != "pendrive" || len(m.Warnings) == 0 {
		t.Errorf("a peça de origem não confiável deve ficar de fora, com aviso: %+v %v", m.Excluded, m.Warnings)
	}
	trusted := fabricar.ArchitectMerge(pieces, fabricar.MergeOptions{TrustRemovable: true})
	if !slices.Contains(stepNames(trusted), "pendrive") {
		t.Errorf("com TrustRemovable a peça entra: %v", stepNames(trusted))
	}
}

func TestArchitectMerger_NeverInheritsMCPsHooksOrScripts(t *testing.T) {
	var pieces []fabricar.DiscoveredPiece
	for _, k := range []fabricar.PieceKind{fabricar.KindMCP, fabricar.KindHook, fabricar.KindScript, fabricar.KindLoop} {
		pieces = append(pieces, fabricar.DiscoveredPiece{Name: "x-" + string(k), Kind: k, Score: 1, Tier: reputation.TierOuro})
	}
	m := fabricar.ArchitectMerge(pieces, fabricar.MergeOptions{})
	for _, name := range stepNames(m) {
		if strings.HasPrefix(name, "x-") {
			t.Errorf("só skills viram passos (default-deny para MCP, hook, script e loop): %v", stepNames(m))
		}
	}
	if len(m.Used) != 0 {
		t.Errorf("nenhuma dessas peças pode ser usada: %+v", m.Used)
	}
}

func TestArchitectMerger_DropsToxicDuplicatesAndCapsSteps(t *testing.T) {
	toxic := piece("t", "D:\\t", 0.1, false)
	toxic.Tier = reputation.TierToxica
	pieces := []fabricar.DiscoveredPiece{toxic, piece("dup", "D:\\1", 0.9, false), piece("dup", "D:\\2", 0.7, false)}
	for _, n := range []string{"e", "f", "g", "h", "i", "j"} {
		pieces = append(pieces, piece(n, "D:\\"+n, 0.5, false))
	}
	m := fabricar.ArchitectMerge(pieces, fabricar.MergeOptions{})
	names := stepNames(m)
	if slices.Contains(names, "t") || len(names) != 5 {
		t.Errorf("tóxica fora e teto de 5 passos: %v", names)
	}
	if names[0] != "dup" || countName(names, "dup") != 1 {
		t.Errorf("duplicada: fica só a melhor: %v", names)
	}
}

func TestArchitectMerger_FallsBackToDefaultStepsWithAWarning(t *testing.T) {
	m := fabricar.ArchitectMerge(nil, fabricar.MergeOptions{})
	if !slices.Equal(stepNames(m), []string{"aidd-spec", "aidd-tickets"}) || len(m.Warnings) == 0 {
		t.Errorf("sem peças curadas usa os passos padrão e avisa: %v %v", stepNames(m), m.Warnings)
	}
}

func countName(names []string, want string) int {
	n := 0
	for _, s := range names {
		if s == want {
			n++
		}
	}
	return n
}
