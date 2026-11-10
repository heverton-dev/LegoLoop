package dag

import (
	"strings"
	"testing"
)

// v2 monta um DAG v2 com o trigger dado (JSON literal) e os nós informados.
func v2(trigger string, nodes ...string) string {
	t := ""
	if trigger != "" {
		t = `"trigger":` + trigger + `,`
	}
	return `{"name":"t","version":"2.0.0",` + t + `"nodes":[` + strings.Join(nodes, ",") + `]}`
}

// n2 é um nó script com artefato e campos extras (execution/slots) injetados.
func n2(id, deps, artifact, extra string) string {
	if extra != "" {
		extra = "," + extra
	}
	return `{"id":"` + id + `","type":"script","depends_on":[` + deps + `],
	"command":["go","version"],"routing":{"mode":"manual","model":"none"},
	"output":{"artifact":"` + artifact + `","gate":"exit_code == 0"}` + extra + `}`
}

func parseErr(t *testing.T, name, src, want string) {
	t.Helper()
	_, err := Parse([]byte(src))
	if err == nil {
		t.Errorf("%s: esperado erro", name)
		return
	}
	if want != "" && !strings.Contains(err.Error(), want) {
		t.Errorf("%s: erro %q não contém %q", name, err, want)
	}
}

// R-TRG-1: DAG v2 exige exatamente 1 trigger.
func TestTriggerRequiredInV2(t *testing.T) {
	parseErr(t, "v2 sem trigger", v2("", n2("a", "", "A.md", "")), "trigger")
	parseErr(t, "trigger array", v2(`[{"type":"manual"},{"type":"manual"}]`, n2("a", "", "A.md", "")), "")
}

// §3: v1 sem trigger continua válido, assumido manual, com aviso.
func TestV1WithoutTriggerIsManualWithWarning(t *testing.T) {
	d, err := Parse([]byte(specExample))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.EffectiveTrigger().Type; got != TriggerManual {
		t.Fatalf("trigger efetivo = %q", got)
	}
	w := d.Warnings()
	if len(w) != 1 || !strings.Contains(w[0], "manual") {
		t.Fatalf("avisos = %v", w)
	}
	out, err := d.Marshal()
	if err != nil || strings.Contains(string(out), `"trigger"`) {
		t.Fatalf("v1 não pode ganhar trigger ao serializar: %v\n%s", err, out)
	}
	if v := strings.Replace(specExample, `"1.0.0"`, `"3.0.0"`, 1); true {
		parseErr(t, "versão desconhecida", v, "versão")
	}
}

func TestV2ValidTriggers(t *testing.T) {
	cases := map[string]string{
		"manual":             `{"type":"manual"}`,
		"cron dias úteis":    `{"type":"cron","cron":"0 9 * * 1-5"}`,
		"cron passo":         `{"type":"cron","cron":"*/15 * * * *"}`,
		"cron lista":         `{"type":"cron","cron":"0 0 1,15 * *"}`,
		"cron faixa+passo":   `{"type":"cron","cron":"0-30/10 8-18 * * 0-7"}`,
		"file_change":        `{"type":"file_change","path":"docs"}`,
		"file_change nested": `{"type":"file_change","path":"src/entrada.md"}`,
		"processo_fim":       `{"type":"processo_fim","processo":"clinica-vet"}`,
	}
	for name, trig := range cases {
		d, err := Parse([]byte(v2(trig, n2("a", "", "OUT/a.md", ""))))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if len(d.Warnings()) != 0 {
			t.Errorf("%s: v2 válido não deve ter avisos: %v", name, d.Warnings())
		}
	}
}

func TestInvalidCron(t *testing.T) {
	bad := []string{"", "* * * *", "* * * * * *", "60 * * * *", "* 24 * * *", "* * 0 * *", "* * 32 * *",
		"* * * 13 *", "* * * 0 *", "* * * * 8", "*/0 * * * *", "5-1 * * * *", "a * * * *", "+1 * * * *",
		"1,,2 * * * *", "-1 * * * *", "1-2-3 * * * *", "*/ * * * *", strings.Repeat("1,", 200) + "1 * * * *"}
	for _, c := range bad {
		src := v2(`{"type":"cron","cron":"`+c+`"}`, n2("a", "", "A.md", ""))
		parseErr(t, "cron "+c, src, "cron")
	}
}

// R-TRG-2: caminho do file_change relativo, sem links/8.3/dirs sensíveis e sem tocar nos artefatos.
func TestFileChangeTriggerValidation(t *testing.T) {
	node := n2("a", "", "out/a.md", "")
	cases := map[string]string{
		"vazio":              "",
		"traversal":          "../x",
		"absoluto":           "/etc",
		"drive":              "C:x",
		"git":                ".git/HEAD",
		"claude":             ".claude/settings.json",
		"8.3":                "GIT~1/x",
		"raiz":               ".",
		"o próprio artefato": "out/a.md",
		"pasta do artefato":  "out",
	}
	for name, p := range cases {
		src := v2(`{"type":"file_change","path":"`+p+`"}`, node)
		parseErr(t, "file_change "+name, src, "")
	}
	parseErr(t, "file_change sem path", v2(`{"type":"file_change"}`, node), "path")
}

func TestProcessoFimTriggerValidation(t *testing.T) {
	node := n2("a", "", "A.md", "")
	parseErr(t, "sem processo", v2(`{"type":"processo_fim"}`, node), "processo")
	parseErr(t, "processo inválido", v2(`{"type":"processo_fim","processo":"a b"}`, node), "processo")
}

// R-TRG-3: tipos da etapa 2 e desconhecidos são recusados com mensagem clara.
func TestTriggerPhase2AndUnknown(t *testing.T) {
	node := n2("a", "", "A.md", "")
	parseErr(t, "webhook", v2(`{"type":"webhook"}`, node), "etapa 2")
	parseErr(t, "git_event", v2(`{"type":"git_event"}`, node), "etapa 2")
	parseErr(t, "desconhecido", v2(`{"type":"telepatia"}`, node), "telepatia")
	parseErr(t, "manual com campo de cron", v2(`{"type":"manual","cron":"* * * * *"}`, node), "manual")
	parseErr(t, "cron com path", v2(`{"type":"cron","cron":"* * * * *","path":"x"}`, node), "cron")
	parseErr(t, "file_change com cron", v2(`{"type":"file_change","path":"x","cron":"* * * * *"}`, node), "file_change")
	parseErr(t, "processo_fim com path", v2(`{"type":"processo_fim","processo":"x","path":"y"}`, node), "path")
}

// R-NODE-7: execution.mode válido; inline proibido em camada paralela com artefatos sobrepostos.
func TestExecutionMode(t *testing.T) {
	trig := `{"type":"manual"}`
	for _, m := range []string{"", "inline", "isolated"} {
		extra := ""
		if m != "" {
			extra = `"execution":{"mode":"` + m + `"}`
		}
		if _, err := Parse([]byte(v2(trig, n2("a", "", "A.md", extra)))); err != nil {
			t.Errorf("modo %q: %v", m, err)
		}
	}
	parseErr(t, "modo inválido", v2(trig, n2("a", "", "A.md", `"execution":{"mode":"yolo"}`)), "execution")
	parseErr(t, "execution sem mode", v2(trig, n2("a", "", "A.md", `"execution":{}`)), "execution")
}

func TestExecutionEffectiveModeDefaultsToInline(t *testing.T) {
	d, err := Parse([]byte(v2(`{"type":"manual"}`,
		n2("a", "", "A.md", ""), n2("b", "", "B.md", `"execution":{"mode":"isolated"}`))))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Node("a").EffectiveMode(); got != ExecInline {
		t.Fatalf("a = %q", got)
	}
	if got := d.Node("b").EffectiveMode(); got != ExecIsolated {
		t.Fatalf("b = %q", got)
	}
}

func TestInlineParallelOverlapRejected(t *testing.T) {
	trig := `{"type":"manual"}`
	iso := `"execution":{"mode":"isolated"}`
	parseErr(t, "paralelo inline sobreposto",
		v2(trig, n2("a", "", "out/a.md", ""), n2("b", "", "out", "")), "inline")
	// um deles isolado resolve
	if _, err := Parse([]byte(v2(trig, n2("a", "", "out/a.md", ""), n2("b", "", "out", iso)))); err != nil {
		t.Errorf("com b isolado deveria passar: %v", err)
	}
	// em sequência (não paralelo) a sobreposição é permitida
	if _, err := Parse([]byte(v2(trig, n2("a", "", "out/a.md", ""), n2("b", `"a"`, "out", "")))); err != nil {
		t.Errorf("sequencial deveria passar: %v", err)
	}
	// artefatos disjuntos em paralelo, ambos inline: permitido
	if _, err := Parse([]byte(v2(trig, n2("a", "", "out/a.md", ""), n2("b", "", "outro/b.md", "")))); err != nil {
		t.Errorf("disjuntos deveriam passar: %v", err)
	}
}

// R-NODE-1/3 + SEC-01/03: slots passam pela mesma validação de segurança dos nós.
func TestSlotsValid(t *testing.T) {
	slots := `"slots":{"gates":["file_exists('A.md')"],"scripts":[["go","version"]],
	"hooks":{"pre":[["go","version"]],"post":[["go","env"]]}}`
	d, err := Parse([]byte(v2(`{"type":"manual"}`, n2("a", "", "A.md", slots))))
	if err != nil {
		t.Fatal(err)
	}
	s := d.Node("a").Slots
	if s == nil || len(s.Gates) != 1 || len(s.Scripts) != 1 || len(s.Hooks.Pre) != 1 || len(s.Hooks.Post) != 1 {
		t.Fatalf("slots mal decodificados: %+v", s)
	}
	if !d.Node("a").HasSlots() || (&Node{}).HasSlots() {
		t.Fatal("HasSlots incorreto")
	}
}

func TestSlotsInvalid(t *testing.T) {
	trig := `{"type":"manual"}`
	cases := map[string]string{
		"gate inválido":     `"slots":{"gates":["exit_code =="]}`,
		"script relativo":   `"slots":{"scripts":[["tools/x.exe"]]}`,
		"script vazio":      `"slots":{"scripts":[[]]}`,
		"hook pre relativo": `"slots":{"hooks":{"pre":[["bin\\x.exe"]]}}`,
		"hook post vazio":   `"slots":{"hooks":{"post":[[""]]}}`,
		"hook drive-rel":    `"slots":{"hooks":{"pre":[["C:x.exe"]]}}`,
		"gates demais":      `"slots":{"gates":[` + strings.TrimSuffix(strings.Repeat(`"true",`, MaxSlotEntries+1), ",") + `]}`,
	}
	for name, extra := range cases {
		parseErr(t, name, v2(trig, n2("a", "", "A.md", extra)), "slot")
	}
	// campo desconhecido é barrado pelo decoder estrito, sem a palavra "slot" na mensagem
	parseErr(t, "campo desconhecido", v2(trig, n2("a", "", "A.md", `"slots":{"evil":[]}`)), "evil")
}

// R-MCP-3: mcp não executa. O nó subprocesso exige referência válida.
func TestMCPNodeRejectedExplicitly(t *testing.T) {
	src := v2(`{"type":"manual"}`, `{"id":"a","type":"mcp","ref":"rea",
	"routing":{"mode":"manual","model":"m"},"output":{"artifact":"A.md","gate":"true"}}`)
	parseErr(t, "mcp", src, "MCP")
}

func TestSubprocessoNode(t *testing.T) {
	mk := func(ref string) string {
		r := ""
		if ref != "" {
			r = `"ref":"` + ref + `",`
		}
		return v2(`{"type":"manual"}`, `{"id":"a","type":"subprocesso",`+r+
			`"routing":{"mode":"manual","model":"m"},"output":{"artifact":"A.md","gate":"true"}}`)
	}
	d, err := Parse([]byte(mk("clinica-vet")))
	if err != nil || d.Node("a").Type != TypeSubloop {
		t.Fatalf("subprocesso válido: %v", err)
	}
	parseErr(t, "subprocesso sem ref", mk(""), "subprocesso")
	parseErr(t, "subprocesso ref inválida", mk("../fora"), "subprocesso")
}

// Slots são opcionais: ausente, vazio ou listas vazias são escolhas válidas do usuário,
// sem erro e sem aviso, e nunca são preenchidos automaticamente.
func TestSlotsAreOptional(t *testing.T) {
	trig := `{"type":"manual"}`
	cases := map[string]string{
		"sem slots":        "",
		"slots vazio":      `"slots":{}`,
		"listas vazias":    `"slots":{"gates":[],"scripts":[],"hooks":{"pre":[],"post":[]}}`,
		"só hooks vazio":   `"slots":{"hooks":{}}`,
		"só um slot usado": `"slots":{"gates":["file_exists('A.md')"]}`,
	}
	for name, extra := range cases {
		d, err := Parse([]byte(v2(trig, n2("a", "", "A.md", extra))))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if len(d.Warnings()) != 0 {
			t.Errorf("%s: slot vazio não gera aviso: %v", name, d.Warnings())
		}
		wantUsed := name == "só um slot usado"
		if d.Node("a").HasSlots() != wantUsed {
			t.Errorf("%s: HasSlots = %v", name, !wantUsed)
		}
	}
}
