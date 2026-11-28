package distribuicao

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const chaveTeste = "chave-super-secreta-123"

type postizFalso struct {
	posts      []map[string]any
	auth       []string
	canais     string
	falhaPosts bool
	uploadRuim bool
}

func (p *postizFalso) servir(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/public/v1/integrations", func(w http.ResponseWriter, r *http.Request) {
		p.auth = append(p.auth, r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, p.canais)
	})
	mux.HandleFunc("/api/public/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		if p.uploadRuim {
			_, _ = io.WriteString(w, `{}`)
			return
		}
		_, _ = io.WriteString(w, `{"id":"up1","path":"/uploads/v.mp4"}`)
	})
	mux.HandleFunc("/api/public/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		if p.falhaPosts {
			http.Error(w, "erro com "+chaveTeste+" no corpo", http.StatusBadRequest)
			return
		}
		var corpo map[string]any
		_ = json.NewDecoder(r.Body).Decode(&corpo)
		p.posts = append(p.posts, corpo)
		_, _ = io.WriteString(w, `[{"postId":"p1"}]`)
	})
	return httptest.NewServer(mux)
}

const canaisJSON = `[{"id":"c1","identifier":"youtube","name":"Bilhoes em Jogo"},{"id":"c2","providerIdentifier":"instagram","name":"Insta","disabled":true}]`

func clienteDe(srv *httptest.Server) Cliente {
	return Cliente{BaseURL: srv.URL + "/api/public/v1", Chave: chaveTeste, HTTP: srv.Client()}
}

func pedidoOK(t *testing.T) PedidoEnvio {
	t.Helper()
	dir := t.TempDir()
	tocar(t, dir, "v.mp4")
	return PedidoEnvio{Canal: "youtube", Video: filepath.Join(dir, "v.mp4"), Titulo: "Titulo", Descricao: "desc", Visibilidade: VisPrivado, RotuloIA: true, Tags: []string{"a", " ", "b"}}
}

func TestPedidoEnvioValidar(t *testing.T) {
	casos := map[string]struct {
		mudar func(*PedidoEnvio)
		quer  error
	}{
		"sem canal":         {func(p *PedidoEnvio) { p.Canal = "" }, ErrPedidoInvalido},
		"titulo vazio":      {func(p *PedidoEnvio) { p.Titulo = " " }, ErrPedidoInvalido},
		"titulo longo":      {func(p *PedidoEnvio) { p.Titulo = strings.Repeat("a", 101) }, ErrPedidoInvalido},
		"visibilidade":      {func(p *PedidoEnvio) { p.Visibilidade = "secreto" }, ErrPedidoInvalido},
		"publico sem trava": {func(p *PedidoEnvio) { p.Visibilidade = VisPublico }, ErrPublicoSemConfirmacao},
		"sem rotulo de IA":  {func(p *PedidoEnvio) { p.RotuloIA = false }, ErrSemRotuloIA},
		"video ausente":     {func(p *PedidoEnvio) { p.Video = filepath.Join(t.TempDir(), "nao.mp4") }, ErrPedidoInvalido},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			p := pedidoOK(t)
			c.mudar(&p)
			if err := p.Validar(); !errors.Is(err, c.quer) {
				t.Fatalf("quer %v, obteve %v", c.quer, err)
			}
		})
	}
	p := pedidoOK(t)
	p.Visibilidade, p.ConfirmoPublico = VisPublico, true
	if err := p.Validar(); err != nil {
		t.Fatalf("publico com confirmacao: %v", err)
	}
}

func TestEnviarDryRunNaoChamaARede(t *testing.T) {
	c := Cliente{BaseURL: "http://127.0.0.1:1/nunca", Chave: chaveTeste}
	p := pedidoOK(t)
	p.DryRun = true
	res, err := c.Enviar(context.Background(), p)
	if err != nil || !res.DryRun || res.Resumo["visibilidade"] != VisPrivado || res.Resumo["modo"] != "agora" {
		t.Fatalf("dry-run: %+v %v", res, err)
	}
	quando := time.Date(2026, 12, 1, 14, 0, 0, 0, time.UTC)
	p.Quando = &quando
	res, _ = c.Enviar(context.Background(), p)
	if !strings.Contains(res.Resumo["modo"].(string), "2026-12-01T14:00:00Z") {
		t.Fatalf("modo agendado: %+v", res.Resumo)
	}
	p.RotuloIA = false
	if _, err := c.Enviar(context.Background(), p); !errors.Is(err, ErrSemRotuloIA) {
		t.Fatalf("pedido invalido nao passa nem em dry-run: %v", err)
	}
}

func TestEnviarCriaPostagemPrivadaComRotuloDeIA(t *testing.T) {
	falso := &postizFalso{canais: canaisJSON}
	srv := falso.servir(t)
	defer srv.Close()
	p := pedidoOK(t)
	quando := time.Date(2026, 12, 1, 14, 0, 0, 0, time.UTC)
	p.Quando = &quando
	res, err := clienteDe(srv).Enviar(context.Background(), p)
	if err != nil || !strings.Contains(string(res.Resposta), "p1") {
		t.Fatalf("Enviar: %+v %v", res, err)
	}
	if len(falso.posts) != 1 || falso.posts[0]["type"] != "schedule" {
		t.Fatalf("postagem: %+v", falso.posts)
	}
	corpo, _ := json.Marshal(falso.posts[0])
	for _, quer := range []string{`"containsSyntheticMedia":true`, `"type":"private"`, `"__type":"youtube"`, `"id":"c1"`, `"path":"/uploads/v.mp4"`} {
		if !strings.Contains(string(corpo), quer) {
			t.Fatalf("corpo sem %s: %s", quer, corpo)
		}
	}
	if strings.Contains(string(corpo), `"label":" "`) {
		t.Fatalf("tag vazia nao deve ir: %s", corpo)
	}
	if falso.auth[0] != chaveTeste {
		t.Fatalf("cabecalho Authorization ausente")
	}
	p.Quando = nil
	if _, err := clienteDe(srv).Enviar(context.Background(), p); err != nil || falso.posts[1]["type"] != "now" {
		t.Fatalf("envio imediato: %v %+v", err, falso.posts)
	}
}

func TestEnviarOcultaAChaveNosErros(t *testing.T) {
	falso := &postizFalso{canais: canaisJSON, falhaPosts: true}
	srv := falso.servir(t)
	defer srv.Close()
	_, err := clienteDe(srv).Enviar(context.Background(), pedidoOK(t))
	if !errors.Is(err, ErrPostiz) || strings.Contains(err.Error(), chaveTeste) || !strings.Contains(err.Error(), "***") {
		t.Fatalf("a chave nao pode vazar: %v", err)
	}
	morto := httptest.NewServer(http.NotFoundHandler())
	cli := clienteDe(morto)
	morto.Close()
	if _, err := cli.Canais(context.Background()); !errors.Is(err, ErrPostiz) || strings.Contains(err.Error(), chaveTeste) {
		t.Fatalf("servidor fora do ar: %v", err)
	}
	semChave := Cliente{BaseURL: "http://127.0.0.1:1"}
	if semChave.ocultar("abc") != "abc" {
		t.Fatal("sem chave nada muda")
	}
}

func TestCanaisEAchar(t *testing.T) {
	falso := &postizFalso{canais: canaisJSON}
	srv := falso.servir(t)
	defer srv.Close()
	c := clienteDe(srv)
	canais, err := c.Canais(context.Background())
	if err != nil || len(canais) != 2 || canais[1].Rede != "instagram" || !canais[1].Desativado {
		t.Fatalf("canais: %+v %v", canais, err)
	}
	falso.canais = `{"integrations":[{"id":"x","identifier":"tiktok","name":"T"}]}`
	canais, err = c.Canais(context.Background())
	if err != nil || len(canais) != 1 || canais[0].Rede != "tiktok" {
		t.Fatalf("envelope: %+v %v", canais, err)
	}
	falso.canais = `lixo`
	if _, err := c.Canais(context.Background()); !errors.Is(err, ErrPostiz) {
		t.Fatalf("resposta ilegivel: %v", err)
	}
	falso.canais = `{"integrations":"x"}`
	if _, err := c.Canais(context.Background()); !errors.Is(err, ErrPostiz) {
		t.Fatalf("envelope ruim: %v", err)
	}
	falso.canais = canaisJSON
	for _, q := range []string{"zzz", "Insta"} {
		if _, err := c.achar(context.Background(), q); !errors.Is(err, ErrPedidoInvalido) {
			t.Fatalf("canal %q inexistente ou desativado: %v", q, err)
		}
	}
	if k, err := c.achar(context.Background(), "bilhoes em jogo"); err != nil || k.ID != "c1" {
		t.Fatalf("achar por nome: %+v %v", k, err)
	}
}

func TestSubirErros(t *testing.T) {
	falso := &postizFalso{canais: canaisJSON, uploadRuim: true}
	srv := falso.servir(t)
	defer srv.Close()
	p := pedidoOK(t)
	if _, err := clienteDe(srv).Subir(context.Background(), p.Video); !errors.Is(err, ErrPostiz) {
		t.Fatalf("upload sem id: %v", err)
	}
	if _, err := clienteDe(srv).Subir(context.Background(), filepath.Join(t.TempDir(), "nao.mp4")); err == nil {
		t.Fatal("arquivo ausente deve falhar")
	}
	if _, err := clienteDe(srv).Enviar(context.Background(), p); !errors.Is(err, ErrPostiz) {
		t.Fatalf("Enviar com upload ruim: %v", err)
	}
	p.Canal = "nao-existe"
	if _, err := clienteDe(srv).Enviar(context.Background(), p); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("canal desconhecido: %v", err)
	}
}

func TestSegredoDoAmbienteENovoCliente(t *testing.T) {
	raiz := t.TempDir()
	vazio := func(string) string { return "" }
	if _, err := NovoCliente(raiz, vazio); !errors.Is(err, ErrSegredoAusente) {
		t.Fatalf("sem chave: %v", err)
	}
	so := func(k string) string {
		if k == "POSTIZ_API_KEY" {
			return chaveTeste
		}
		return ""
	}
	if _, err := NovoCliente(raiz, so); !errors.Is(err, ErrSegredoAusente) {
		t.Fatalf("sem URL: %v", err)
	}
	_ = os.WriteFile(filepath.Join(raiz, ".env-teste"), []byte("\xef\xbb\xbf# comentario\nOUTRA=1\nPOSTIZ_URL=\"https://postiz.exemplo.org/\"\n"), 0o644)
	c, err := NovoCliente(raiz, so)
	if err != nil || c.BaseURL != "https://postiz.exemplo.org/api/public/v1" || c.Chave != chaveTeste {
		t.Fatalf("cliente: %+v %v", c, err)
	}
	if v, ok := SegredoDoAmbiente("OUTRA", raiz, vazio); !ok || v != "1" {
		t.Fatalf("arquivo .env: %q %v", v, ok)
	}
	if _, ok := SegredoDoAmbiente("AUSENTE", raiz, vazio); ok {
		t.Fatal("variavel ausente")
	}
	comBOM := t.TempDir()
	_ = os.WriteFile(filepath.Join(comBOM, ".env"), []byte("\xef\xbb\xbfPOSTIZ_API_KEY=zzz\n"), 0o644)
	if v, ok := SegredoDoAmbiente("POSTIZ_API_KEY", comBOM, vazio); !ok || v != "zzz" {
		t.Fatalf("BOM na primeira linha: %q %v", v, ok)
	}
	_ = os.Mkdir(filepath.Join(raiz, ".env-dir"), 0o755)
	if _, ok := SegredoDoAmbiente("AUSENTE", raiz, vazio); ok {
		t.Fatal("diretorio .env-dir deve ser ignorado")
	}
}

func TestPublicarDoAcervo(t *testing.T) {
	falso := &postizFalso{canais: canaisJSON}
	srv := falso.servir(t)
	defer srv.Close()
	c := clienteDe(srv)
	p := pedidoOK(t)
	it := itemBase("a", "tema", TipoPerene)
	a := Acervo{Itens: []Item{it}}

	if _, err := c.PublicarDoAcervo(context.Background(), &a, "zz", p); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("inexistente: %v", err)
	}
	a.Itens[0].Visibilidade = VisPublico
	if _, err := c.PublicarDoAcervo(context.Background(), &a, "a", p); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("visibilidade diferente: %v", err)
	}
	a.Itens[0].Visibilidade = VisPrivado
	a.Itens[0].Autorizado = false
	if _, err := c.PublicarDoAcervo(context.Background(), &a, "a", p); !errors.Is(err, ErrSemAutorizacao) {
		t.Fatalf("sem autorizacao: %v", err)
	}
	a.Itens[0].Autorizado, a.Itens[0].Status = true, StatusExpirado
	if _, err := c.PublicarDoAcervo(context.Background(), &a, "a", p); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("status expirado: %v", err)
	}
	a.Itens[0].Status = StatusPronto
	p.DryRun = true
	if _, err := c.PublicarDoAcervo(context.Background(), &a, "a", p); err != nil || a.Itens[0].Status != StatusPronto || len(falso.posts) != 0 {
		t.Fatalf("dry-run nao muda nada: %v %+v", err, a.Itens[0])
	}
	p.DryRun = false
	if _, err := c.PublicarDoAcervo(context.Background(), &a, "a", p); err != nil || a.Itens[0].Status != StatusPublic || len(a.Itens[0].Links) != 1 {
		t.Fatalf("publicar agora: %v %+v", err, a.Itens[0])
	}
	b := Acervo{Itens: []Item{itemBase("b", "outro", TipoPerene)}}
	quando := time.Now().Add(time.Hour)
	p.Quando = &quando
	if _, err := c.PublicarDoAcervo(context.Background(), &b, "b", p); err != nil || b.Itens[0].Status != StatusPronto || len(b.Itens[0].Links) != 1 {
		t.Fatalf("agendado no Postiz mantem status: %v %+v", err, b.Itens[0])
	}
	falso.falhaPosts = true
	cc := Acervo{Itens: []Item{itemBase("c", "terceiro", TipoPerene)}}
	if _, err := c.PublicarDoAcervo(context.Background(), &cc, "c", p); err == nil || cc.Itens[0].Status != StatusPronto {
		t.Fatalf("falha nao muda o acervo: %v", err)
	}
}
