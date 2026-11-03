package studio

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func TestFabricarModalStatic(t *testing.T) {
	html := readStatic(t, "index.html")

	t.Run("o HTML tem os ids do modal", func(t *testing.T) {
		for _, id := range []string{"btn-novo-processo", "modal-fabricar", "fabricar-name", "fabricar-model", "fabricar-start", "fabricar-error"} {
			if !strings.Contains(html, `id="`+id+`"`) {
				t.Errorf("index.html sem id=%q", id)
			}
		}
		if !regexp.MustCompile(`<dialog id="modal-fabricar"`).MatchString(html) {
			t.Error("o modal deve ser um <dialog> da própria página, como o do cofre")
		}
		if !regexp.MustCompile(`(?s)<dialog id="modal-fabricar".*?id="fabricar-error".*?</dialog>`).MatchString(html) {
			t.Error("os erros devem aparecer no próprio modal (fabricar-error dentro do dialog), nunca em diálogo nativo")
		}
	})

	t.Run("sem script inline, sem diálogos nativos, sem innerHTML, sem emoji", func(t *testing.T) {
		for _, m := range regexp.MustCompile(`(?s)<script\b([^>]*)>`).FindAllStringSubmatch(html, -1) {
			if !strings.Contains(m[1], "src=") {
				t.Errorf("script inline proibido pela CSP: <script%s>", m[1])
			}
		}
		native := regexp.MustCompile(`(^|[^.\w])(alert|confirm|prompt)\s*\(`)
		emoji := regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{1F000}-\x{1F2FF}\x{FE00}-\x{FE0F}]`)
		for _, f := range []string{"fabricar.js", "app.js", "index.html", "style.css"} {
			for i, line := range strings.Split(readStatic(t, f), "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "//") {
					continue
				}
				if f != "style.css" && native.MatchString(line) {
					t.Errorf("%s:%d diálogo nativo: %s", f, i+1, strings.TrimSpace(line))
				}
				if f == "fabricar.js" && (strings.Contains(line, "innerHTML") || strings.Contains(line, "insertAdjacentHTML")) {
					t.Errorf("%s:%d HTML dinâmico: %s", f, i+1, strings.TrimSpace(line))
				}
				if emoji.MatchString(line) {
					t.Errorf("%s:%d emoji: %s", f, i+1, strings.TrimSpace(line))
				}
			}
		}
	})

	t.Run("o modal nao usa estilo inline", func(t *testing.T) {
		block := regexp.MustCompile(`(?s)<dialog id="modal-fabricar".*?</dialog>`).FindString(html)
		if strings.Contains(block, "style=") {
			t.Error("o modal tem style=: a CSP style-src 'self' bloqueia, use classe no style.css")
		}
		if !strings.Contains(readStatic(t, "style.css"), "modal-fabricar") && !strings.Contains(readStatic(t, "style.css"), ".fabricar-") {
			t.Error("style.css sem regras do modal de fabricação")
		}
	})

	t.Run("o modal tem o botao de sobrescrever e o app.js pergunta antes de agir", func(t *testing.T) {
		if !regexp.MustCompile(`(?s)<dialog id="modal-fabricar".*?id="fabricar-overwrite".*?</dialog>`).MatchString(html) {
			t.Error("o botão fabricar-overwrite deve estar dentro do modal")
		}
		app := readStatic(t, "app.js")
		if strings.Count(app, "confirmDialog(") < 2 || !strings.Contains(app, "fabricarNeedsReplaceConfirm") || !strings.Contains(app, "fabricarCanOverwrite") {
			t.Error("app.js deve confirmar (confirmDialog) antes de substituir o Canvas e antes de sobrescrever")
		}
		if !strings.Contains(readStatic(t, "style.css"), ".fabricar-overwrite[hidden]") {
			t.Error("o botão oculto precisa de regra [hidden], pois a classe define display")
		}
	})

	t.Run("elementos com hidden continuam ocultos apesar do display do CSS", func(t *testing.T) {
		// display:grid/flex numa classe vence o estilo hidden do navegador: o painel apareceria antes da hora.
		css := readStatic(t, "style.css")
		if !strings.Contains(css, ".fabricar-interview[hidden]") {
			t.Error("style.css precisa de .fabricar-interview[hidden] { display: none }")
		}
	})

	t.Run("fabricar.js é servido, carregado antes do app.js, e a CSP segue intacta", func(t *testing.T) {
		s, _ := newServer(t)
		rec := do(t, s.Handler(), "GET", "/fabricar.js", "")
		if rec.Code != 200 {
			t.Fatalf("/fabricar.js = %d", rec.Code)
		}
		if got := rec.Header().Get("Content-Security-Policy"); !strings.HasPrefix(got, "default-src 'self'; script-src 'self'") {
			t.Errorf("CSP alterada: %q", got)
		}
		i, j := strings.Index(html, `src="fabricar.js"`), strings.Index(html, `src="app.js"`)
		if i < 0 || j < 0 || i > j {
			t.Errorf("fabricar.js deve ser carregado antes do app.js (fabricar@%d, app@%d)", i, j)
		}
	})

	t.Run("o wiring usa o inventário, a entrevista e a API de fabricação", func(t *testing.T) {
		app := readStatic(t, "app.js")
		for _, want := range []string{"btn-novo-processo", "modal-fabricar", "/api/interview/start", "/api/interview/message", "/api/fabricar/run"} {
			if !strings.Contains(app, want) {
				t.Errorf("app.js não referencia %q", want)
			}
		}
	})
}

type fabModelOpt struct {
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
	Label    string `json:"label"`
}

func TestFabricarModalLogic(t *testing.T) {
	out := runNode(t, `(() => {
	  const models = [{id:"sonnet",source:"claude-cli"},{id:"code-fast",source:"9router"},{id:"qwen",source:"ollama"},{id:"code-pro",source:"9router"}];
	  const opts = fabricarModelOptions(models);
	  return {
	    opts,
	    picked: fabricarDefaultModel(opts),
	    none: fabricarDefaultModel(fabricarModelOptions([{id:"sonnet",source:"claude-cli"}])),
	    payloadFromMirror: fabricarPayload({name:"  Extrator  ", category:"build", model:"code-fast"}, {mirror:{outcome:"extrair áudio"}, hypothesis:"h"}),
	    payloadFromHyp: fabricarPayload({name:"x", category:"web", model:"m"}, {hypothesis:"hipótese"}),
	    payloadNoTurn: fabricarPayload({name:"x", category:"build", model:"m"}, null),
	    e400: fabricarErrorMessage(400, {error:"nome inválido: use letras e números"}),
	    e409: fabricarErrorMessage(409, {error:"já existe um processo com este nome"}),
	    e422: fabricarErrorMessage(422, {error:"x"}),
	    e502: fabricarErrorMessage(502, {error:"x"}),
	    e500: fabricarErrorMessage(500, {}),
	    eNet: fabricarErrorMessage(0, null),
	    forced: fabricarPayload({name:"x", category:"build", model:"m"}, null, true),
	    unforced: fabricarPayload({name:"x", category:"build", model:"m"}, null, false),
	    ow: [fabricarCanOverwrite(409, {code:"exists"}), fabricarCanOverwrite(409, {code:"busy"}), fabricarCanOverwrite(400, {code:"exists"}), fabricarCanOverwrite(409, null)],
	    ask: [fabricarNeedsReplaceConfirm(0), fabricarNeedsReplaceConfirm(1), fabricarNeedsReplaceConfirm(7)],
	    valid: fabricarNameError(""), validOk: fabricarNameError("Meu fluxo"), validBad: fabricarNameError("!!!"),
	  };
	})()`, "fabricar.js")
	var r struct {
		Opts                               []fabModelOpt
		Picked, None                       string
		PayloadFromMirror, PayloadFromHyp  map[string]any
		PayloadNoTurn                      map[string]any
		E400, E409, E422, E502, E500, ENet string
		Valid, ValidOk, ValidBad           string
		Forced, Unforced                   map[string]any
		Ow, Ask                            []bool
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	byID := map[string]fabModelOpt{}
	for _, o := range r.Opts {
		byID[o.ID] = o
	}
	if len(r.Opts) != 4 || byID["code-fast"].Disabled || byID["code-pro"].Disabled || !byID["sonnet"].Disabled || !byID["qwen"].Disabled {
		t.Errorf("só modelos do 9Router podem conduzir a fabricação; os demais ficam desabilitados: %+v", r.Opts)
	}
	if !strings.Contains(byID["sonnet"].Label, "9Router") {
		t.Errorf("o modelo desabilitado deve explicar o motivo: %q", byID["sonnet"].Label)
	}
	if r.Picked != "code-fast" || r.None != "" {
		t.Errorf("padrão = %q (esperado code-fast, o primeiro do 9Router); sem 9Router = %q (esperado vazio)", r.Picked, r.None)
	}
	if r.PayloadFromMirror["name"] != "Extrator" || r.PayloadFromMirror["goal"] != "extrair áudio" || r.PayloadFromMirror["model"] != "code-fast" {
		t.Errorf("payload do espelho: %v", r.PayloadFromMirror)
	}
	if r.PayloadFromHyp["goal"] != "hipótese" || r.PayloadFromHyp["category"] != "web" {
		t.Errorf("payload da hipótese: %v", r.PayloadFromHyp)
	}
	if g, ok := r.PayloadNoTurn["goal"]; ok && g != "" {
		t.Errorf("sem entrevista o goal fica vazio (o servidor usa o nome): %v", r.PayloadNoTurn)
	}
	for _, k := range []string{"force", "vault", "docs"} {
		if _, ok := r.PayloadFromMirror[k]; ok {
			t.Errorf("o cliente não envia %q", k)
		}
	}
	for name, msg := range map[string]string{"400": r.E400, "409": r.E409, "422": r.E422, "502": r.E502, "500": r.E500, "rede": r.ENet} {
		if strings.TrimSpace(msg) == "" {
			t.Errorf("mensagem de erro %s vazia", name)
		}
	}
	if !strings.Contains(r.E400, "nome inválido") || !strings.Contains(strings.ToLower(r.E409), "já existe") {
		t.Errorf("400/409 devem repetir o motivo do servidor: %q | %q", r.E400, r.E409)
	}
	if r.Forced["force"] != true {
		t.Errorf("sobrescrever envia force=true: %v", r.Forced)
	}
	if _, ok := r.Unforced["force"]; ok {
		t.Errorf("sem sobrescrita o payload não leva force: %v", r.Unforced)
	}
	if len(r.Ow) != 4 || !r.Ow[0] || r.Ow[1] || r.Ow[2] || r.Ow[3] {
		t.Errorf("só 409 com code=exists oferece a sobrescrita: %v", r.Ow)
	}
	if len(r.Ask) != 3 || r.Ask[0] || !r.Ask[1] || !r.Ask[2] {
		t.Errorf("canvas vazio planta direto; com nós, pergunta antes de substituir: %v", r.Ask)
	}
	if r.Valid == "" || r.ValidOk != "" || r.ValidBad == "" {
		t.Errorf("validação do nome: vazio=%q ok=%q inválido=%q", r.Valid, r.ValidOk, r.ValidBad)
	}
}
