package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// passo é a saída JSON de cada passo da entrevista conduzida pelo Chat.
type passo struct {
	Sessao    string   `json:"sessao"`
	Status    string   `json:"status"`
	Pergunta  string   `json:"pergunta"`
	Palpite   string   `json:"palpite"`
	Opcoes    []string `json:"opcoes"`
	Confianca int      `json:"confianca"`
	Meta      int      `json:"meta"`
	Turnos    int      `json:"turnos"`
	Espelho   *struct {
		Resultado    string   `json:"resultado_esperado"`
		Quem         string   `json:"quem_se_beneficia"`
		PorQueAgora  string   `json:"por_que_agora"`
		Criterio     string   `json:"criterio_de_sucesso"`
		Restricoes   []string `json:"restricoes_inegociaveis"`
		ForaDeEscopo []string `json:"fora_de_escopo"`
	} `json:"espelho"`
	Lacunas []string `json:"lacunas"`
}

func (r *interRig) entrevista(t *testing.T, args ...string) (passo, int) {
	t.Helper()
	r.out.Reset()
	r.errb.Reset()
	full := append([]string{"fabricar", "entrevista"}, args...)
	full = append(full, "--vault", r.vault, "--docs", r.doc, "--model", "sonnet")
	code := run(t, r.d, full...)
	var p passo
	if code == 0 {
		if err := json.Unmarshal(r.out.Bytes(), &p); err != nil {
			t.Fatalf("a saída não é JSON puro: %v\n%s", err, r.out)
		}
	}
	return p, code
}

func TestEntrevistaPassoAPassoIniciarEResponder(t *testing.T) {
	r := newInterRig(t, "", lowReply, doneReply)
	p, code := r.entrevista(t, "iniciar", "extrair áudio")
	if code != 0 {
		t.Fatalf("iniciar: exit %d\n%s", code, r.errb)
	}
	if p.Sessao == "" || p.Pergunta != "Qual o formato?" || p.Palpite != "mp3" || p.Confianca != 50 || p.Meta != 95 || p.Espelho != nil {
		t.Errorf("primeiro passo inesperado: %+v", p)
	}
	if len(p.Opcoes) < 2 || p.Opcoes[0] != "Aceitar o palpite: mp3" {
		t.Errorf("opções inesperadas: %v", p.Opcoes)
	}
	if _, err := os.Stat(filepath.Join(r.vault, "sessions", p.Sessao+".json")); err != nil {
		t.Errorf("a sessão deve persistir em sessions/: %v", err)
	}
	// segunda chamada = outro "processo": o motor só conhece o arquivo.
	p2, code := r.entrevista(t, "responder", p.Sessao, "mp3")
	if code != 0 {
		t.Fatalf("responder: exit %d\n%s", code, r.errb)
	}
	if p2.Sessao != p.Sessao || p2.Confianca != 98 || p2.Status != "COMPLETED" || p2.Espelho == nil {
		t.Fatalf("segundo passo inesperado: %+v", p2)
	}
	if p2.Espelho.Resultado != "extrair áudio em mp3" || p2.Espelho.Quem != "editores" || len(p2.Espelho.ForaDeEscopo) != 1 {
		t.Errorf("espelho incompleto: %+v", p2.Espelho)
	}
}

func TestEntrevistaPassoAPassoIniciarLeAgentsDoProjeto(t *testing.T) {
	r := newInterRig(t, "", lowReply)
	spy := &sysLLM{inner: r.llm}
	r.d.Fabricar.LLM = spy
	_ = os.MkdirAll(r.doc, 0o755)
	if err := os.WriteFile(filepath.Join(r.doc, "AGENTS.md"), []byte("Go, stdlib only"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := r.entrevista(t, "iniciar", "extrair áudio"); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(spy.system, "Go, stdlib only") {
		t.Errorf("a entrevista do Chat deve ser contextual:\n%s", spy.system)
	}
}

func TestEntrevistaPassoAPassoErros(t *testing.T) {
	r := newInterRig(t, "", lowReply, doneReply)
	for _, id := range []string{"../fora", "intv_1/../x", "qualquer"} {
		if _, code := r.entrevista(t, "responder", id, "oi"); code != 2 {
			t.Errorf("id %q: exit %d, queria 2", id, code)
		}
	}
	if _, code := r.entrevista(t, "responder", "intv_999", "oi"); code != 1 || !strings.Contains(r.errb.String(), "não encontrada") {
		t.Errorf("sessão ausente: exit %d\n%s", code, r.errb)
	}
	if _, code := r.entrevista(t, "iniciar"); code != 2 {
		t.Errorf("iniciar sem ideia deve sair com 2, saiu %d", code)
	}
	if _, code := r.entrevista(t, "voar"); code != 2 {
		t.Errorf("subcomando desconhecido deve sair com 2, saiu %d", code)
	}
	p, _ := r.entrevista(t, "iniciar", "x")
	if _, code := r.entrevista(t, "responder", p.Sessao, "ok"); code != 0 { // conclui
		t.Fatal("responder deveria concluir")
	}
	if _, code := r.entrevista(t, "responder", p.Sessao, "mais"); code != 1 || !strings.Contains(r.errb.String(), "finalizada") {
		t.Errorf("responder numa sessão concluída: exit %d\n%s", code, r.errb)
	}
}

// fabricarComEntrevista roda a fabricação a partir de uma sessão do Chat.
func (r *interRig) fabricarComEntrevista(t *testing.T, sessao string, extra ...string) int {
	t.Helper()
	r.out.Reset()
	r.errb.Reset()
	args := append([]string{"--yes", "--entrevista", sessao}, extra...)
	return run(t, r.d, r.args("novo processo", args...)...)
}

func TestEntrevistaPassoAPassoFabricarCarregaOEspelho(t *testing.T) {
	r := newInterRig(t, "", lowReply, doneReply)
	p, _ := r.entrevista(t, "iniciar", "extrair áudio")
	r.entrevista(t, "responder", p.Sessao, "mp3")
	before := r.llm.interview
	if code := r.fabricarComEntrevista(t, p.Sessao); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, r.out, r.errb)
	}
	if r.llm.interview != before {
		t.Error("a fabricação não pode entrevistar de novo")
	}
	intent := firstGlob(t, filepath.Join(r.doc, "docs", "intents", "*.md"))
	for _, want := range []string{"editores", "retrabalho", "sem rede", "edição de vídeo"} {
		if !strings.Contains(intent, want) {
			t.Errorf("INTENT sem %q:\n%s", want, intent)
		}
	}
	if strings.Contains(intent, "Não informado") {
		t.Errorf("espelho completo não deixa 'Não informado':\n%s", intent)
	}
}

func TestEntrevistaPassoAPassoAbaixoDaMetaExigeAceite(t *testing.T) {
	r := newInterRig(t, "", lowReply)
	p, _ := r.entrevista(t, "iniciar", "extrair áudio")
	if code := r.fabricarComEntrevista(t, p.Sessao); code != 1 {
		t.Fatalf("exit %d, queria 1\n%s", code, r.errb)
	}
	if !strings.Contains(r.errb.String(), "--aceitar-rascunho") || wroteDocs(t, r) != 0 {
		t.Errorf("deve explicar o aceite e não gravar (gravou %d)\n%s", r.wrote(t), r.errb)
	}
	if code := r.fabricarComEntrevista(t, p.Sessao, "--aceitar-rascunho"); code != 0 {
		t.Fatalf("com o aceite: exit %d\n%s", code, r.errb)
	}
	intent := firstGlob(t, filepath.Join(r.doc, "docs", "intents", "*.md"))
	if !strings.Contains(intent, "**Status:** DRAFT") || !strings.Contains(intent, "50%") {
		t.Errorf("o INTENT deve sair DRAFT com a confiança medida:\n%s", intent)
	}
}

func TestEntrevistaPassoAPassoYesExigeEntrevistaOuOptOut(t *testing.T) {
	r := newInterRig(t, "", lowReply)
	run1 := func(extra ...string) int {
		r.out.Reset()
		r.errb.Reset()
		return run(t, r.d, r.args("novo processo", append([]string{"--yes"}, extra...)...)...)
	}
	if code := run1(); code != 2 {
		t.Fatalf("--yes sem entrevista: exit %d, queria 2", code)
	}
	for _, want := range []string{"INTENT sem entrevista", "--sem-entrevista"} {
		if !strings.Contains(r.errb.String(), want) {
			t.Errorf("o aviso não traz %q:\n%s", want, r.errb)
		}
	}
	if wroteDocs(t, r) != 0 {
		t.Error("nada pode ser gravado")
	}
	if code := run1("--sem-entrevista", "--entrevista", "intv_1"); code != 2 {
		t.Errorf("--sem-entrevista com --entrevista é contraditório: exit %d", code)
	}
	if code := run1("--sem-entrevista"); code != 0 {
		t.Fatalf("com o opt-out: exit %d\n%s", code, r.errb)
	}
	if !strings.Contains(r.errb.String(), "INTENT sem entrevista") {
		t.Errorf("o opt-out ainda deve avisar:\n%s", r.errb)
	}
}

func TestEntrevistaPassoAPassoSkillProibeAtalho(t *testing.T) {
	_, body := readFabricarSkill(t)
	for _, want := range []string{"legoloop fabricar entrevista iniciar", "legoloop fabricar entrevista responder", "--entrevista <sessao>"} {
		if !strings.Contains(body, want) {
			t.Errorf("a skill não ensina %q", want)
		}
	}
	rule := regexp.MustCompile(`(?m)^\d+\. \*\*NUNCA\*\*[^\n]*--yes[^\n]*sem[^\n]*--entrevista[^\n]*$`)
	if !rule.MatchString(body) {
		t.Error("a skill deve ter uma regra NUNCA proibindo --yes sem --entrevista")
	}
	if !strings.Contains(body, "--sem-entrevista") {
		t.Error("a skill deve citar --sem-entrevista para proibi-lo por iniciativa própria")
	}
	proc := body[strings.Index(body, "### Procedimento"):]
	if strings.Index(proc, "entrevista iniciar") > strings.Index(proc, "--dry-run") {
		t.Error("a entrevista pelo motor deve vir antes da prévia")
	}
}

// wroteDocs conta o que a fabricação gravou, sem contar as sessões da entrevista.
func wroteDocs(t *testing.T, r *interRig) int {
	t.Helper()
	n := 0
	for _, f := range r.files(t) {
		if !strings.Contains(filepath.ToSlash(f), "/sessions/") {
			n++
		}
	}
	return n
}

func TestEntrevistaIntentMaterializaINTENT(t *testing.T) {
	r := newInterRig(t, "", lowReply, doneReply)
	p, code := r.entrevista(t, "iniciar", "processo local de vídeo offline")
	if code != 0 {
		t.Fatalf("iniciar: exit %d\n%s", code, r.errb)
	}
	if _, code := r.entrevista(t, "responder", p.Sessao, "FFmpeg e CPU"); code != 0 {
		t.Fatalf("responder: exit %d\n%s", code, r.errb)
	}

	r.out.Reset()
	r.errb.Reset()
	code = run(t, r.d, "fabricar", "entrevista", "intent", p.Sessao, "--vault", r.vault, "--docs", r.doc)
	if code != 0 {
		t.Fatalf("intent: exit %d\n%s", code, r.errb)
	}
	var out struct {
		Arquivo string `json:"arquivo"`
	}
	if err := json.Unmarshal(r.out.Bytes(), &out); err != nil {
		t.Fatalf("saída não é JSON: %v\n%s", err, r.out)
	}
	data, err := os.ReadFile(out.Arquivo)
	if err != nil {
		t.Fatalf("lendo INTENT: %v", err)
	}
	conteudo := string(data)
	if !strings.Contains(conteudo, "UNDER_REVIEW") {
		t.Errorf("esperava status UNDER_REVIEW na INTENT")
	}
	if strings.Contains(conteudo, "{{") {
		t.Errorf("INTENT com placeholder não resolvido")
	}
	if !strings.Contains(conteudo, p.Sessao) {
		t.Errorf("INTENT sem o id da sessão %s", p.Sessao)
	}
	if !strings.Contains(filepath.Base(out.Arquivo), "INTENT_") {
		t.Errorf("nome de arquivo fora do padrão: %s", out.Arquivo)
	}
}
