package provider

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Source identifica o executor capaz de rodar um modelo.
type Source string

const (
	SourceRouter    Source = "9router"
	SourceOllama    Source = "ollama"
	SourceClaudeCLI Source = "claude-cli"
	SourceAgy       Source = "agy"
	SourceOmp       Source = "omp"
	SourceMimo      Source = "mimo"
	SourceOpencode  Source = "opencode"
)

// harnessSources são executores cujos modelos têm prefixo "<harness>:<modelo>"
// ("<harness>:default" = modelo padrão configurado no próprio harness). Listar os
// modelos deles leva de 3 a 40 s, então o inventário passivo só exige o binário no PATH;
// o próprio harness valida o modelo na execução.
func harnessSources() []Source {
	return []Source{SourceAgy, SourceOmp, SourceMimo, SourceOpencode}
}

// Model é um modelo "verde" no inventário: existe E tem executor real.
type Model struct {
	ID     string `json:"id"`
	Source Source `json:"source"`
}

// Config injeta dependências do scanner (testável, sem estado global).
type Config struct {
	HTTP        *http.Client
	OllamaURL   string
	RouterURLs  []string
	RouterKey   string
	LookPath    func(string) (string, error)
	Getenv      func(string) string
	HarnessCLIs []string
}

// DefaultConfig usa as portas oficiais: Ollama :11434, 9Router :20128 e :9000.
// NINEROUTER_URL, se definida, passa a ser a única base sondada.
func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	routers := []string{"http://127.0.0.1:20128", "http://127.0.0.1:9000"}
	if explicit := os.Getenv("NINEROUTER_URL"); explicit != "" {
		if u := normalizeRouterURL(explicit); safeRouterURL(u) {
			routers = []string{u}
		}
		return Config{
			HTTP:        &http.Client{Timeout: RemoteScanTimeout},
			OllamaURL:   "http://127.0.0.1:11434",
			RouterURLs:  routers,
			RouterKey:   routerKeyFromEnv(os.Getenv),
			LookPath:    exec.LookPath,
			Getenv:      os.Getenv,
			HarnessCLIs: defaultHarnessCLIs(),
		}
	}
	envURL, envKey := findRouterCreds(home, os.Getenv)
	if u := normalizeRouterURL(envURL); safeRouterURL(u) {
		routers = append([]string{u}, routers...)
	}
	return Config{
		HTTP:        &http.Client{Timeout: RemoteScanTimeout},
		OllamaURL:   "http://127.0.0.1:11434",
		RouterURLs:  routers,
		RouterKey:   envKey,
		LookPath:    exec.LookPath,
		Getenv:      os.Getenv,
		HarnessCLIs: defaultHarnessCLIs(),
	}
}

// normalizeRouterURL tolera a forma usada pelos harnesses (".../v1").
func normalizeRouterURL(raw string) string {
	u := strings.TrimRight(strings.TrimSpace(raw), "/")
	return strings.TrimSuffix(u, "/v1")
}

// routerKeyFromEnv aceita NINEROUTER_API_KEY ou a convenção NINEROUTER_KEY (omp/aidd-9router).
func routerKeyFromEnv(getenv func(string) string) string {
	if k := getenv("NINEROUTER_API_KEY"); k != "" {
		return k
	}
	return getenv("NINEROUTER_KEY")
}

// safeRouterURL aceita https ou http somente em loopback (a credencial não trafega em claro).
func safeRouterURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return false
	}
	if u.Scheme == "https" {
		return true
	}
	host := u.Hostname()
	ip := net.ParseIP(host)
	return u.Scheme == "http" && (host == "localhost" || (ip != nil && ip.IsLoopback()))
}

// TrustedRouterURL é a única base que recebe RouterKey (a primeira configurada).
func (c Config) TrustedRouterURL() string {
	if len(c.RouterURLs) == 0 {
		return ""
	}
	return c.RouterURLs[0]
}

func defaultHarnessCLIs() []string {
	return []string{"claude", "agy", "omp", "mimo", "opencode", "codex", "gemini", "cursor-agent"}
}

// Inventory é a Matriz Viva de Modelos da máquina.
type Inventory struct {
	RouterURL    string   `json:"router_url,omitempty"`
	RouterModels []string `json:"router_models"`
	Ollama       []string `json:"ollama"`
	CLIs         []CLI    `json:"clis"`
	APIKeys      []string `json:"api_keys"`
	Errors       []string `json:"errors,omitempty"`
}

// Discover roda as sondas em paralelo; falhas viram Errors, nunca pânico.
func Discover(ctx context.Context, c Config) Inventory {
	var inv Inventory
	var mu sync.Mutex
	var wg sync.WaitGroup
	fail := func(err error) {
		mu.Lock()
		inv.Errors = append(inv.Errors, err.Error())
		mu.Unlock()
	}
	wg.Go(func() {
		models, err := ScanOllama(ctx, c.HTTP, c.OllamaURL)
		if err != nil {
			fail(err)
		}
		inv.Ollama = models
	})
	wg.Go(func() {
		base, models, err := ScanNineRouter(ctx, c.HTTP, c.RouterURLs, c.RouterKey)
		if err != nil {
			fail(err)
		}
		inv.RouterURL, inv.RouterModels = base, models
	})
	names := c.HarnessCLIs
	if len(names) == 0 {
		names = defaultHarnessCLIs()
	}
	inv.CLIs = DetectCLIs(c.LookPath, names)
	inv.APIKeys = DetectAPIKeys(c.Getenv)
	wg.Wait()
	return inv
}

// claudeCLIModels são os aliases aceitos por `claude --model`.
func claudeCLIModels() []string { return []string{"haiku", "sonnet", "opus"} }

// Models lista os modelos executáveis; combos do 9Router primeiro (recomendados).
func (inv Inventory) Models() []Model {
	var out []Model
	for _, m := range inv.RouterModels {
		out = append(out, Model{ID: m, Source: SourceRouter})
	}
	for _, m := range inv.Ollama {
		out = append(out, Model{ID: m, Source: SourceOllama})
	}
	if inv.HasCLI("claude") {
		for _, m := range claudeCLIModels() {
			out = append(out, Model{ID: m, Source: SourceClaudeCLI})
		}
	}
	for _, h := range harnessSources() {
		if inv.HasCLI(string(h)) {
			out = append(out, Model{ID: string(h) + ":default", Source: h})
		}
	}
	return out
}

// Source devolve o primeiro executor (por prioridade) que roda o modelo; IDs com
// prefixo de harness ("agy:<modelo>") valem quando o harness está no PATH.
func (inv Inventory) Source(model string) (Source, bool) {
	for _, m := range inv.Models() {
		if m.ID == model {
			return m.Source, true
		}
	}
	prefix, name, ok := strings.Cut(model, ":")
	for _, h := range harnessSources() {
		if ok && name != "" && prefix == string(h) && inv.HasCLI(prefix) {
			return h, true
		}
	}
	return "", false
}

// Has aplica a Regra de Ouro: só modelos do inventário são permitidos.
func (inv Inventory) Has(model string) bool {
	_, ok := inv.Source(model)
	return ok
}

// HasCLI informa se o harness está no PATH.
func (inv Inventory) HasCLI(name string) bool {
	for _, c := range inv.CLIs {
		if c.Name == name {
			return true
		}
	}
	return false
}

// CLIPath devolve o caminho absoluto resolvido do harness.
func (inv Inventory) CLIPath(name string) string {
	for _, c := range inv.CLIs {
		if c.Name == name {
			return c.Path
		}
	}
	return ""
}
