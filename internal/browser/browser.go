package browser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const PortaPadrao = 9222

var (
	ErrNavegadorFechado       = errors.New("browser: nenhuma porta CDP ativa (9222/9333)")
	ErrNavegadorNaoEncontrado = errors.New("browser: Chrome ou Edge nao encontrado no sistema")
)

type OpcoesAbertura struct {
	Home            string
	Porta           int
	URLInicial      string
	DetectarExe     func() (string, error)
	IniciarProcesso func(ctx context.Context, exe string, args []string) error
}

// PortasPadrao retorna as URLs CDP sondadas por ordem de prioridade.
func PortasPadrao() []string {
	return []string{"http://127.0.0.1:9222", "http://127.0.0.1:9333"}
}

// ResolverCDPURL sonda as URLs CDP informadas (ou 9222 e 9333) e retorna a primeira ativa.
func ResolverCDPURL(ctx context.Context, cli *http.Client, candidatas []string) (string, error) {
	if cli == nil {
		cli = &http.Client{Timeout: 2 * time.Second}
	}
	if len(candidatas) == 0 {
		candidatas = PortasPadrao()
	}
	for _, base := range candidatas {
		if sondarCDP(ctx, cli, base) == nil {
			return strings.TrimRight(base, "/"), nil
		}
	}
	return "", ErrNavegadorFechado
}

func sondarCDP(ctx context.Context, cli *http.Client, baseURL string) error {
	alvo := strings.TrimRight(baseURL, "/") + "/json/version"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, alvo, nil)
	if err != nil {
		return err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	return nil
}

// PerfilPadrao reutiliza ~/.chrome-automation-profile se existir ou cria ~/.legoloop/chrome-profile.
func PerfilPadrao(home string) string {
	if strings.TrimSpace(home) == "" {
		h, err := os.UserHomeDir()
		if err == nil {
			home = h
		}
	}
	legado := filepath.Join(home, ".chrome-automation-profile")
	if info, err := os.Stat(legado); err == nil && info.IsDir() {
		return legado
	}
	return filepath.Join(home, ".legoloop", "chrome-profile")
}

// ConstruirArgs gera os argumentos de linha de comando para sessao persistente via CDP.
func ConstruirArgs(perfilDir string, porta int, urlInicial string) []string {
	if porta <= 0 {
		porta = PortaPadrao
	}
	args := []string{
		"--user-data-dir=" + perfilDir,
		fmt.Sprintf("--remote-debugging-port=%d", porta),
		"--remote-allow-origins=*",
	}
	if strings.TrimSpace(urlInicial) != "" {
		args = append(args, urlInicial)
	}
	return args
}

// DetectarNavegador localiza Chrome ou Edge no PATH ou nos caminhos padrao do sistema.
func DetectarNavegador(lookPath func(string) (string, error), existe func(string) bool) (string, error) {
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	if existe == nil {
		existe = func(p string) bool {
			info, err := os.Stat(p)
			return err == nil && !info.IsDir()
		}
	}
	for _, bin := range []string{"chrome", "google-chrome", "chromium", "msedge"} {
		if p, err := lookPath(bin); err == nil && p != "" {
			return p, nil
		}
	}
	for _, fixo := range caminhosFixosNavegador() {
		if existe(fixo) {
			return fixo, nil
		}
	}
	return "", ErrNavegadorNaoEncontrado
}

func caminhosFixosNavegador() []string {
	return []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`,
	}
}

// AbrirPersistente abre o navegador com perfil isolado persistente em modo desacoplado.
func AbrirPersistente(ctx context.Context, opts OpcoesAbertura) error {
	detectar := opts.DetectarExe
	if detectar == nil {
		detectar = func() (string, error) {
			return DetectarNavegador(nil, nil)
		}
	}
	exe, err := detectar()
	if err != nil {
		return err
	}
	perfil := PerfilPadrao(opts.Home)
	if err := os.MkdirAll(perfil, 0o755); err != nil {
		return err
	}
	iniciar := opts.IniciarProcesso
	if iniciar == nil {
		iniciar = func(_ context.Context, bin string, a []string) error {
			return exec.Command(bin, a...).Start()
		}
	}
	return iniciar(ctx, exe, ConstruirArgs(perfil, opts.Porta, opts.URLInicial))
}
