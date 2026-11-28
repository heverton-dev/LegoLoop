package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func escreverDist(t *testing.T, caminho string, v any) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	dados, ok := v.([]byte)
	if !ok {
		var err error
		if dados, err = json.Marshal(v); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(caminho, dados, 0o644); err != nil {
		t.Fatal(err)
	}
	return caminho
}

func itemDist(id, tema string) map[string]any {
	return map[string]any{"id": id, "tema": tema, "tipo": "perene", "canais": []string{"youtube"}, "formatos": []string{"16:9"},
		"status": "pronto", "visibilidade": "private", "autorizado": true}
}

func TestVideoDistribuicaoAjudaEErros(t *testing.T) {
	if r := rodarCLI(t, "video-distribuicao"); r.code != 0 || !strings.Contains(r.stdout, "revisao folhas") {
		t.Fatalf("ajuda: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "nao-existe"); r.code != 2 {
		t.Fatalf("subcomando desconhecido: %+v", r)
	}
	for _, sub := range []string{"prevoo", "capa", "agenda", "aprendizado"} {
		if r := rodarCLI(t, "video-distribuicao", sub); r.code != 2 || !strings.Contains(r.stderr, "--projeto") {
			t.Errorf("%s sem --projeto: %+v", sub, r)
		}
	}
	for _, grupo := range []string{"revisao", "vitrine", "acervo", "publicar", "comentarios"} {
		if r := rodarCLI(t, "video-distribuicao", grupo); r.code != 2 {
			t.Errorf("%s sem subcomando: %+v", grupo, r)
		}
		if r := rodarCLI(t, "video-distribuicao", grupo, "x"); r.code != 2 {
			t.Errorf("%s x sem --projeto: %+v", grupo, r)
		}
		if r := rodarCLI(t, "video-distribuicao", grupo, "x", "--projeto", t.TempDir()); r.code == 0 {
			t.Errorf("%s x deveria falhar: %+v", grupo, r)
		}
	}
	if r := rodarCLI(t, "video-distribuicao", "prevoo", "--nao-existe"); r.code != 2 {
		t.Fatalf("flag invalida: %+v", r)
	}
}

func TestVideoDistribuicaoPrevoo(t *testing.T) {
	dir := t.TempDir()
	r := rodarCLI(t, "video-distribuicao", "prevoo", "--projeto", dir, "--raiz", t.TempDir())
	if r.code != 0 || !strings.Contains(r.stdout, "OK: pasta acervo/") || !strings.Contains(r.stdout, "pre-voo aprovado") {
		t.Fatalf("prevoo: %+v", r)
	}
	arq := t.TempDir()
	_ = os.WriteFile(filepath.Join(arq, "acervo"), []byte("x"), 0o644)
	if r := rodarCLI(t, "video-distribuicao", "prevoo", "--projeto", arq); r.code != 1 || !strings.Contains(r.stdout, "FALTA") {
		t.Fatalf("acervo que e arquivo: %+v", r)
	}
}

func TestVideoDistribuicaoAcervoEAgenda(t *testing.T) {
	dir := t.TempDir()
	p := escreverDist(t, filepath.Join(t.TempDir(), "i1.json"), itemDist("a", "juros compostos"))
	if r := rodarCLI(t, "video-distribuicao", "acervo", "adicionar", "--projeto", dir, "--arquivo", p); r.code != 0 {
		t.Fatalf("adicionar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "adicionar", "--projeto", dir, "--arquivo", p); r.code != 1 {
		t.Fatalf("id repetido deve falhar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "adicionar", "--projeto", dir, "--arquivo", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("arquivo ausente: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "autorizar", "--projeto", dir, "--id", "a", "--revogar"); r.code != 0 {
		t.Fatalf("revogar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "autorizar", "--projeto", dir, "--id", "zz"); r.code != 1 {
		t.Fatalf("id inexistente: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "listar", "--projeto", dir); r.code != 0 || !strings.Contains(r.stdout, `"autorizado": false`) {
		t.Fatalf("listar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "autorizar", "--projeto", dir, "--id", "a"); r.code != 0 {
		t.Fatalf("autorizar: %+v", r)
	}
	agora := "2026-11-27T10:00:00Z"
	r := rodarCLI(t, "video-distribuicao", "agenda", "--projeto", dir, "--agora", agora)
	if r.code != 0 || !strings.Contains(r.stdout, `"acao": "agendar"`) || !strings.Contains(r.stderr, "aplicado: false") {
		t.Fatalf("agenda sem aplicar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "agenda", "--projeto", dir, "--agora", agora, "--aplicar"); r.code != 0 {
		t.Fatalf("agenda aplicando: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "acervo", "listar", "--projeto", dir); !strings.Contains(r.stdout, `"status": "agendado"`) {
		t.Fatalf("status agendado esperado: %+v", r)
	}
	cfg := escreverDist(t, filepath.Join(t.TempDir(), "cfg.json"), map[string]any{"horarios": []string{"09:30"}})
	if r := rodarCLI(t, "video-distribuicao", "agenda", "--projeto", dir, "--agora", agora, "--config", cfg); r.code != 0 {
		t.Fatalf("agenda com config: %+v", r)
	}
	for _, args := range [][]string{{"--agora", "ontem"}, {"--config", filepath.Join(dir, "nao.json")}} {
		if r := rodarCLI(t, append([]string{"video-distribuicao", "agenda", "--projeto", dir}, args...)...); r.code != 1 {
			t.Fatalf("agenda %v deveria falhar: %+v", args, r)
		}
	}
	ruim := t.TempDir()
	escreverDist(t, distAcervoPath(ruim), []byte("{"))
	for _, sub := range []string{"agenda", "acervo"} {
		args := []string{"video-distribuicao", sub, "--projeto", ruim}
		if sub == "acervo" {
			args = []string{"video-distribuicao", "acervo", "listar", "--projeto", ruim}
		}
		if r := rodarCLI(t, args...); r.code != 1 {
			t.Fatalf("acervo corrompido em %s: %+v", sub, r)
		}
	}
}

func distAcervoPath(dir string) string { return filepath.Join(dir, "acervo", "acervo.json") }

func TestVideoDistribuicaoPublicarEhDryRunPorPadrao(t *testing.T) {
	dir := t.TempDir()
	escreverDist(t, distAcervoPath(dir), map[string]any{"schema_version": 1, "itens": []any{itemDist("a", "tema")}})
	video := filepath.Join(dir, "v.mp4")
	escreverDist(t, video, []byte("x"))
	pedido := escreverDist(t, filepath.Join(dir, "p.json"), map[string]any{"canal": "youtube", "video": video, "titulo": "T", "descricao": "d", "visibilidade": "private", "rotulo_ia": true})
	r := rodarCLI(t, "video-distribuicao", "publicar", "enviar", "--projeto", dir, "--id", "a", "--arquivo", pedido)
	if r.code != 0 || !strings.Contains(r.stdout, `"dry_run": true`) {
		t.Fatalf("dry-run por padrao: %+v", r)
	}
	t.Setenv("POSTIZ_API_KEY", "")
	t.Setenv("POSTIZ_URL", "")
	r = rodarCLI(t, "video-distribuicao", "publicar", "enviar", "--projeto", dir, "--id", "a", "--arquivo", pedido, "--executar", "--raiz", t.TempDir())
	if r.code != 1 || !strings.Contains(r.stderr, "POSTIZ_API_KEY") {
		t.Fatalf("--executar sem segredo: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "publicar", "canais", "--projeto", dir, "--raiz", t.TempDir()); r.code != 1 {
		t.Fatalf("canais sem segredo: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "publicar", "enviar", "--projeto", dir, "--id", "a", "--arquivo", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("pedido ausente: %+v", r)
	}
	escreverDist(t, distAcervoPath(dir), []byte("{"))
	if r := rodarCLI(t, "video-distribuicao", "publicar", "enviar", "--projeto", dir, "--id", "a", "--arquivo", pedido); r.code != 1 {
		t.Fatalf("acervo corrompido: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "publicar", "enviar", "--projeto", t.TempDir(), "--id", "zz", "--arquivo", pedido); r.code != 1 {
		t.Fatalf("item inexistente: %+v", r)
	}
}

func TestVideoDistribuicaoRevisaoEVitrine(t *testing.T) {
	dir := t.TempDir()
	parecer := escreverDist(t, filepath.Join(t.TempDir(), "p.json"), map[string]any{"versao": 1, "veredito": "ENVIAR", "revisor": "agente-isolado"})
	if r := rodarCLI(t, "video-distribuicao", "revisao", "liberar", "--projeto", dir, "--versao", "1"); r.code != 1 {
		t.Fatalf("sem parecer nao libera: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "revisao", "parecer", "--projeto", dir, "--arquivo", parecer); r.code != 0 {
		t.Fatalf("parecer: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "revisao", "parecer", "--projeto", dir, "--arquivo", parecer); r.code != 1 {
		t.Fatalf("sobrescrever deve falhar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "revisao", "parecer", "--projeto", dir, "--arquivo", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("arquivo ausente: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "revisao", "liberar", "--projeto", dir, "--versao", "1"); r.code != 0 {
		t.Fatalf("ENVIAR libera: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "revisao", "folhas", "--projeto", dir, "--video", "renders/v1.mp4", "--versao", "1"); r.code != 1 {
		t.Fatalf("folhas sem video: %+v", r)
	}
	r := rodarCLI(t, "video-distribuicao", "vitrine", "pontuar", "--projeto", dir, "--palavra", "juros", "Juros: 5 erros", "outro")
	if r.code != 0 || !strings.Contains(r.stdout, "Juros: 5 erros") {
		t.Fatalf("pontuar: %+v", r)
	}
	vitrine := escreverDist(t, filepath.Join(t.TempDir(), "v.json"), map[string]any{})
	if r := rodarCLI(t, "video-distribuicao", "vitrine", "gravar", "--projeto", dir, "--arquivo", vitrine, "--titulo", "0"); r.code != 1 {
		t.Fatalf("vitrine vazia: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "vitrine", "gravar", "--projeto", dir, "--arquivo", vitrine, "--titulo", "x"); r.code != 1 {
		t.Fatalf("titulo nao numerico: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "vitrine", "gravar", "--projeto", dir, "--arquivo", filepath.Join(dir, "nao.json"), "--titulo", "0"); r.code != 1 {
		t.Fatalf("arquivo ausente: %+v", r)
	}
	capa := escreverDist(t, filepath.Join(t.TempDir(), "c.json"), map[string]any{"linha1": "A", "fundo": "nao.jpg", "saida": "c.png"})
	if r := rodarCLI(t, "video-distribuicao", "capa", "--projeto", dir, "--arquivo", capa); r.code != 1 {
		t.Fatalf("capa sem fundo: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "capa", "--projeto", dir, "--arquivo", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("pedido de capa ausente: %+v", r)
	}
}

func TestVideoDistribuicaoComentariosEAprendizado(t *testing.T) {
	dir := t.TempDir()
	entrada := escreverDist(t, filepath.Join(t.TempDir(), "c.json"), []map[string]any{{"id": "1", "rede": "youtube", "texto": "oi"}, {"id": "2", "texto": "ja respondido", "respondido": true}})
	r := rodarCLI(t, "video-distribuicao", "comentarios", "listar", "--projeto", dir, "--entrada", entrada)
	if r.code != 0 || !strings.Contains(r.stdout, `"id": "1"`) || strings.Contains(r.stdout, `"id": "2"`) {
		t.Fatalf("listar pendentes: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "rascunhar", "--projeto", dir, "--entrada", entrada, "--id", "1", "--texto", "Obrigado!"); r.code != 0 {
		t.Fatalf("rascunhar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "rascunhar", "--projeto", dir, "--entrada", entrada, "--id", "9", "--texto", "x"); r.code != 1 {
		t.Fatalf("comentario fora da entrada: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "rascunhar", "--projeto", dir, "--entrada", filepath.Join(dir, "nao.json"), "--id", "1", "--texto", "x"); r.code != 1 {
		t.Fatalf("entrada ausente: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "decidir", "--projeto", dir, "--id", "1"); r.code != 1 {
		t.Fatalf("decidir sem escolha: %+v", r)
	}
	r = rodarCLI(t, "video-distribuicao", "comentarios", "decidir", "--projeto", dir, "--id", "1", "--aprovar", "--editado", "Valeu!")
	if r.code != 0 || !strings.Contains(r.stdout, "aprovadas: 1") {
		t.Fatalf("aprovar: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "listar", "--projeto", dir, "--entrada", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("listar sem entrada: %+v", r)
	}
	ruim := t.TempDir()
	escreverDist(t, filepath.Join(ruim, "acervo", "respostas.json"), []byte("{"))
	if r := rodarCLI(t, "video-distribuicao", "comentarios", "listar", "--projeto", ruim, "--entrada", entrada); r.code != 1 {
		t.Fatalf("respostas corrompidas: %+v", r)
	}
	ms := []map[string]any{
		{"item_id": "a", "tema": "juros", "formato": "9:16", "duracao_s": 40, "retencao": 0.7, "alcance": 100, "hora_utc": 14, "novos_seguidores": 1},
		{"item_id": "b", "tema": "cofre", "formato": "16:9", "duracao_s": 300, "retencao": 0.4, "alcance": 50, "hora_utc": 19, "novos_seguidores": 2},
		{"item_id": "c", "tema": "juros", "formato": "9:16", "duracao_s": 20, "retencao": 0.9, "alcance": 300, "hora_utc": 14, "novos_seguidores": 3},
	}
	metricas := escreverDist(t, filepath.Join(t.TempDir(), "m.json"), ms)
	if r := rodarCLI(t, "video-distribuicao", "aprendizado", "--projeto", dir, "--metricas", metricas); r.code != 0 || !strings.Contains(r.stdout, "aprendizado.md") {
		t.Fatalf("aprendizado: %+v", r)
	}
	if r := rodarCLI(t, "video-distribuicao", "aprendizado", "--projeto", dir, "--metricas", filepath.Join(dir, "nao.json")); r.code != 1 {
		t.Fatalf("metricas ausentes: %+v", r)
	}
	poucas := escreverDist(t, filepath.Join(t.TempDir(), "p.json"), ms[:1])
	if r := rodarCLI(t, "video-distribuicao", "aprendizado", "--projeto", dir, "--metricas", poucas); r.code != 1 {
		t.Fatalf("poucas medidas: %+v", r)
	}
}
