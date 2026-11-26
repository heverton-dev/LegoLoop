package roteiro

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ItemPreVoo struct {
	Nome    string `json:"nome"`
	Nivel   string `json:"nivel"` // "OK" | "FALTA" | "AVISO"
	Detalhe string `json:"detalhe,omitempty"`
}

type RelatorioPreVoo struct {
	Bloqueado bool         `json:"bloqueado"`
	Itens     []ItemPreVoo `json:"itens"`
}

type DepsPreVoo struct {
	ChecarLegoloop func(context.Context) error
	ChecarCofreRea func() bool
	ChecarOpcional func(nome string) bool
	HTTPClient     *http.Client
}

// TotalAvisos retorna a quantidade de itens apenas informativos (AVISO).
func (r RelatorioPreVoo) TotalAvisos() int {
	n := 0
	for _, it := range r.Itens {
		if it.Nivel == "AVISO" {
			n++
		}
	}
	return n
}

// ExecutarPreVoo verifica os 5 bloqueantes do Processo 1 e registra itens de P2/P3 como AVISO.
func ExecutarPreVoo(ctx context.Context, p Projeto, cat CatalogoFontes, projDir, acervoDir string, deps DepsPreVoo) RelatorioPreVoo {
	deps = preencherDepsPadrao(deps)
	var rel RelatorioPreVoo
	adicionarItem(&rel, checarBinarioLegoloop(ctx, deps.ChecarLegoloop))
	adicionarItem(&rel, checarCatalogoPreVoo(cat))
	adicionarItem(&rel, checarPermissoesPastas(projDir, acervoDir))
	adicionarItem(&rel, checarConectividadeHTTP(ctx, p, cat, deps.HTTPClient))
	if strings.TrimSpace(p.VideoReferencia) != "" {
		adicionarItem(&rel, checarCofreVideoRea(deps.ChecarCofreRea))
	}
	for _, opt := range []string{"chrome_flow_9333", "ffmpeg", "chromium", "edge-tts", "gitleaks"} {
		adicionarItem(&rel, checarOpcionalAviso(opt, deps.ChecarOpcional))
	}
	return rel
}

func preencherDepsPadrao(d DepsPreVoo) DepsPreVoo {
	if d.ChecarLegoloop == nil {
		d.ChecarLegoloop = func(ctx context.Context) error {
			return exec.CommandContext(ctx, "legoloop", "version").Run()
		}
	}
	if d.ChecarCofreRea == nil {
		d.ChecarCofreRea = func() bool {
			_, err := os.Stat(filepath.Join("processos", "legoloop-video-rea"))
			return err == nil
		}
	}
	if d.ChecarOpcional == nil {
		d.ChecarOpcional = func(nome string) bool {
			_, err := exec.LookPath(nome)
			return err == nil
		}
	}
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{Timeout: 8 * time.Second}
	}
	return d
}

func adicionarItem(rel *RelatorioPreVoo, it ItemPreVoo) {
	rel.Itens = append(rel.Itens, it)
	if it.Nivel == "FALTA" {
		rel.Bloqueado = true
	}
}

func checarBinarioLegoloop(ctx context.Context, fn func(context.Context) error) ItemPreVoo {
	if err := fn(ctx); err != nil {
		return ItemPreVoo{Nome: "legoloop_version", Nivel: "FALTA", Detalhe: err.Error()}
	}
	return ItemPreVoo{Nome: "legoloop_version", Nivel: "OK"}
}

func checarCatalogoPreVoo(cat CatalogoFontes) ItemPreVoo {
	if err := cat.Validar(); err != nil {
		return ItemPreVoo{Nome: "catalogo_fontes", Nivel: "FALTA", Detalhe: err.Error()}
	}
	return ItemPreVoo{Nome: "catalogo_fontes", Nivel: "OK"}
}

func checarPermissoesPastas(projDir, acervoDir string) ItemPreVoo {
	for _, dir := range []string{projDir, acervoDir} {
		if err := testarEscritaLeitura(dir); err != nil {
			return ItemPreVoo{Nome: "permissoes_pastas", Nivel: "FALTA", Detalhe: fmt.Sprintf("%s: %v", dir, err)}
		}
	}
	return ItemPreVoo{Nome: "permissoes_pastas", Nivel: "OK"}
}

func testarEscritaLeitura(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("diretorio vazio")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".prevoo_*")
	if err != nil {
		return err
	}
	nome := tmp.Name()
	if err := tmp.Close(); err != nil {
		return err
	}
	_, readErr := os.ReadFile(nome)
	remErr := os.Remove(nome)
	if readErr != nil {
		return readErr
	}
	return remErr
}

func checarConectividadeHTTP(ctx context.Context, p Projeto, cat CatalogoFontes, client *http.Client) ItemPreVoo {
	ativos := cat.ConectoresParaProjeto(p)
	if len(ativos) == 0 {
		return ItemPreVoo{Nome: "conectividade_http", Nivel: "FALTA", Detalhe: "nenhum conector compativel com perfis_fonte"}
	}
	for _, c := range ativos {
		if sondaHTTP200(ctx, client, montarURLConector(c.URL, p.Nicho)) {
			return ItemPreVoo{Nome: "conectividade_http", Nivel: "OK", Detalhe: c.ID}
		}
	}
	return ItemPreVoo{Nome: "conectividade_http", Nivel: "FALTA", Detalhe: "nenhum conector respondeu status 200 em <= 8s"}
}

func sondaHTTP200(ctx context.Context, client *http.Client, alvo string) bool {
	cctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(cctx, http.MethodGet, alvo, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "LegoLoop-Diorama-Roteiro/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func montarURLConector(rawURL, nicho string) string {
	return strings.ReplaceAll(rawURL, "{{QUERY}}", url.QueryEscape(strings.TrimSpace(nicho)))
}

func checarCofreVideoRea(fn func() bool) ItemPreVoo {
	if !fn() {
		return ItemPreVoo{Nome: "cofre_video_rea", Nivel: "FALTA", Detalhe: "legoloop-video-rea ausente no Cofre"}
	}
	return ItemPreVoo{Nome: "cofre_video_rea", Nivel: "OK"}
}

func checarOpcionalAviso(nome string, fn func(string) bool) ItemPreVoo {
	if fn(nome) {
		return ItemPreVoo{Nome: nome, Nivel: "OK"}
	}
	return ItemPreVoo{Nome: nome, Nivel: "AVISO", Detalhe: "opcional no Processo 1"}
}
