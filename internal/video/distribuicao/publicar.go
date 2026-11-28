package distribuicao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Cliente fala com a API publica do Postiz. A chave nunca aparece em erro, log ou arquivo.
type Cliente struct {
	BaseURL string
	Chave   string
	HTTP    *http.Client
}

// SegredoDoAmbiente le o valor da variavel de ambiente ou de um arquivo .env* da raiz.
func SegredoDoAmbiente(nome, raiz string, getenv func(string) string) (string, bool) {
	if v := strings.TrimSpace(getenv(nome)); v != "" {
		return limparSegredo(v), true
	}
	arquivos, _ := filepath.Glob(filepath.Join(raiz, ".env*"))
	slices.Sort(arquivos)
	for _, a := range arquivos {
		dados, err := os.ReadFile(a)
		if err != nil {
			continue
		}
		for _, linha := range strings.Split(string(dados), "\n") {
			linha = strings.TrimSpace(strings.TrimPrefix(linha, "\xef\xbb\xbf"))
			if v, ok := strings.CutPrefix(linha, nome+"="); ok && !strings.HasPrefix(linha, "#") {
				return limparSegredo(v), true
			}
		}
	}
	return "", false
}

func limparSegredo(v string) string {
	return strings.Trim(strings.TrimSpace(v), "\"'")
}

// NovoCliente monta o cliente com POSTIZ_API_KEY e POSTIZ_URL; falta de qualquer um e erro tipado.
func NovoCliente(raiz string, getenv func(string) string) (Cliente, error) {
	chave, ok := SegredoDoAmbiente("POSTIZ_API_KEY", raiz, getenv)
	if !ok || chave == "" {
		return Cliente{}, fmt.Errorf("%w: POSTIZ_API_KEY", ErrSegredoAusente)
	}
	url, ok := SegredoDoAmbiente("POSTIZ_URL", raiz, getenv)
	if !ok || url == "" {
		return Cliente{}, fmt.Errorf("%w: POSTIZ_URL", ErrSegredoAusente)
	}
	return Cliente{BaseURL: strings.TrimRight(url, "/") + "/api/public/v1", Chave: chave, HTTP: &http.Client{Timeout: 10 * time.Minute}}, nil
}

// Canal e uma rede ligada ao Postiz.
type Canal struct {
	ID         string `json:"id"`
	Rede       string `json:"rede"`
	Nome       string `json:"nome"`
	Desativado bool   `json:"desativado,omitempty"`
}

func (c Cliente) chamar(ctx context.Context, metodo, caminho string, corpo []byte, tipo string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, metodo, c.BaseURL+caminho, bytes.NewReader(corpo))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Chave)
	req.Header.Set("Content-Type", tipo)
	req.Header.Set("User-Agent", "legoloop-video-distribuicao/1.0")
	cli := c.HTTP
	if cli == nil {
		cli = http.DefaultClient
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPostiz, c.ocultar(err.Error()))
	}
	defer resp.Body.Close()
	dados, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: HTTP %d: %s", ErrPostiz, resp.StatusCode, c.ocultar(trecho(string(dados), 300)))
	}
	return dados, nil
}

func (c Cliente) ocultar(s string) string {
	if c.Chave == "" {
		return s
	}
	return strings.ReplaceAll(s, c.Chave, "***")
}

func trecho(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n])
	}
	return s
}

// Canais lista os canais ligados ao Postiz.
func (c Cliente) Canais(ctx context.Context) ([]Canal, error) {
	dados, err := c.chamar(ctx, http.MethodGet, "/integrations", nil, "application/json")
	if err != nil {
		return nil, err
	}
	var brutos []struct {
		ID         string `json:"id"`
		Identifier string `json:"identifier"`
		Provider   string `json:"providerIdentifier"`
		Name       string `json:"name"`
		Disabled   bool   `json:"disabled"`
	}
	if err := json.Unmarshal(dados, &brutos); err != nil {
		var envelope struct {
			Integrations json.RawMessage `json:"integrations"`
		}
		if err2 := json.Unmarshal(dados, &envelope); err2 != nil || json.Unmarshal(envelope.Integrations, &brutos) != nil {
			return nil, fmt.Errorf("%w: resposta de canais ilegivel", ErrPostiz)
		}
	}
	out := make([]Canal, 0, len(brutos))
	for _, b := range brutos {
		rede := b.Identifier
		if rede == "" {
			rede = b.Provider
		}
		out = append(out, Canal{ID: b.ID, Rede: rede, Nome: b.Name, Desativado: b.Disabled})
	}
	return out, nil
}

// Upload e o arquivo ja guardado no Postiz.
type Upload struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// Subir envia o arquivo ao armazenamento do Postiz; nao publica nada.
func (c Cliente) Subir(ctx context.Context, arquivo string) (Upload, error) {
	conteudo, err := os.ReadFile(arquivo)
	if err != nil {
		return Upload{}, err
	}
	var corpo bytes.Buffer
	mw := multipart.NewWriter(&corpo)
	parte, err := mw.CreateFormFile("file", filepath.Base(arquivo))
	if err != nil {
		return Upload{}, err
	}
	if _, err := parte.Write(conteudo); err != nil {
		return Upload{}, err
	}
	if err := mw.Close(); err != nil {
		return Upload{}, err
	}
	dados, err := c.chamar(ctx, http.MethodPost, "/upload", corpo.Bytes(), mw.FormDataContentType())
	if err != nil {
		return Upload{}, err
	}
	var up Upload
	if err := json.Unmarshal(dados, &up); err != nil || up.ID == "" || up.Path == "" {
		return Upload{}, fmt.Errorf("%w: resposta do upload sem id e path", ErrPostiz)
	}
	return up, nil
}

// PedidoEnvio descreve uma postagem. Padrao: privada. Publica exige confirmacao explicita.
type PedidoEnvio struct {
	Canal           string     `json:"canal"`
	Video           string     `json:"video"`
	Titulo          string     `json:"titulo"`
	Descricao       string     `json:"descricao"`
	Visibilidade    string     `json:"visibilidade"`
	ConfirmoPublico bool       `json:"confirmo_publico"`
	Tags            []string   `json:"tags,omitempty"`
	Quando          *time.Time `json:"quando,omitempty"`
	RotuloIA        bool       `json:"rotulo_ia"`
	DryRun          bool       `json:"dry_run"`
}

// Validar aplica as travas: titulo, visibilidade, confirmacao de publico, rotulo de IA e video no disco.
func (p PedidoEnvio) Validar() error {
	n := len([]rune(strings.TrimSpace(p.Titulo)))
	switch {
	case strings.TrimSpace(p.Canal) == "":
		return fmt.Errorf("%w: canal obrigatorio", ErrPedidoInvalido)
	case n == 0 || n > MaxTitulo:
		return fmt.Errorf("%w: titulo de 1 a %d caracteres", ErrPedidoInvalido, MaxTitulo)
	case !slices.Contains([]string{VisPrivado, "unlisted", VisPublico}, p.Visibilidade):
		return fmt.Errorf("%w: visibilidade %q", ErrPedidoInvalido, p.Visibilidade)
	case p.Visibilidade == VisPublico && !p.ConfirmoPublico:
		return ErrPublicoSemConfirmacao
	case !p.RotuloIA:
		return ErrSemRotuloIA
	}
	if _, err := os.Stat(p.Video); err != nil {
		return fmt.Errorf("%w: video %q ausente", ErrPedidoInvalido, p.Video)
	}
	return nil
}

// ResultadoEnvio resume o que foi (ou seria) enviado.
type ResultadoEnvio struct {
	DryRun   bool            `json:"dry_run"`
	Resumo   map[string]any  `json:"resumo"`
	Resposta json.RawMessage `json:"resposta,omitempty"`
}

// Enviar valida o pedido; em dry-run so devolve o resumo. Sem dry-run sobe o video e cria a postagem.
func (c Cliente) Enviar(ctx context.Context, p PedidoEnvio) (ResultadoEnvio, error) {
	if err := p.Validar(); err != nil {
		return ResultadoEnvio{}, err
	}
	info, _ := os.Stat(p.Video)
	modo := "agora"
	if p.Quando != nil {
		modo = "agendado " + p.Quando.UTC().Format(time.RFC3339)
	}
	res := ResultadoEnvio{DryRun: p.DryRun, Resumo: map[string]any{"canal": p.Canal, "video": filepath.Base(p.Video), "bytes": info.Size(),
		"titulo": p.Titulo, "visibilidade": p.Visibilidade, "tags": len(p.Tags), "modo": modo, "rotulo_ia": true}}
	if p.DryRun {
		return res, nil
	}
	canal, err := c.achar(ctx, p.Canal)
	if err != nil {
		return res, err
	}
	up, err := c.Subir(ctx, p.Video)
	if err != nil {
		return res, err
	}
	corpo, err := json.Marshal(montarPostagem(canal, up, p))
	if err != nil {
		return res, err
	}
	res.Resposta, err = c.chamar(ctx, http.MethodPost, "/posts", corpo, "application/json")
	return res, err
}

func (c Cliente) achar(ctx context.Context, q string) (Canal, error) {
	canais, err := c.Canais(ctx)
	if err != nil {
		return Canal{}, err
	}
	var achados []Canal
	for _, k := range canais {
		if q == k.ID || q == k.Rede || strings.EqualFold(q, k.Nome) {
			achados = append(achados, k)
		}
	}
	if len(achados) != 1 || achados[0].Desativado {
		return Canal{}, fmt.Errorf("%w: --canal %q casou com %d canais ativos; use o id exato", ErrPedidoInvalido, q, len(achados))
	}
	return achados[0], nil
}

func montarPostagem(canal Canal, up Upload, p PedidoEnvio) map[string]any {
	tags := []map[string]string{}
	for _, t := range p.Tags {
		if t = strings.TrimSpace(t); t != "" {
			tags = append(tags, map[string]string{"value": t, "label": t})
		}
	}
	tipo, data := "now", time.Now().UTC().Format(time.RFC3339)
	if p.Quando != nil {
		tipo, data = "schedule", p.Quando.UTC().Format(time.RFC3339)
	}
	return map[string]any{"type": tipo, "shortLink": false, "date": data, "tags": []string{},
		"posts": []map[string]any{{
			"integration": map[string]string{"id": canal.ID},
			"value":       []map[string]any{{"content": p.Descricao, "image": []map[string]string{{"id": up.ID, "path": up.Path}}}},
			"settings":    map[string]any{"__type": canal.Rede, "title": p.Titulo, "type": p.Visibilidade, "tags": tags, "containsSyntheticMedia": true},
		}}}
}

// PublicarDoAcervo exige autorizacao gravada, visibilidade igual a do item e atualiza o acervo.
func (c Cliente) PublicarDoAcervo(ctx context.Context, a *Acervo, id string, p PedidoEnvio) (ResultadoEnvio, error) {
	it := a.Buscar(id)
	switch {
	case it == nil:
		return ResultadoEnvio{}, fmt.Errorf("%w: item %q inexistente", ErrAcervoInvalido, id)
	case !it.Autorizado:
		return ResultadoEnvio{}, ErrSemAutorizacao
	case it.Status != StatusPronto && it.Status != StatusAgendado:
		return ResultadoEnvio{}, fmt.Errorf("%w: item %q com status %q", ErrAcervoInvalido, id, it.Status)
	case p.Visibilidade != it.Visibilidade:
		return ResultadoEnvio{}, fmt.Errorf("%w: visibilidade do pedido difere da do item", ErrPedidoInvalido)
	}
	res, err := c.Enviar(ctx, p)
	if err != nil || p.DryRun {
		return res, err
	}
	it.Links = append(it.Links, "postiz:"+trecho(strings.TrimSpace(string(res.Resposta)), 120))
	if p.Quando == nil {
		it.Status = StatusPublic
	}
	return res, nil
}
