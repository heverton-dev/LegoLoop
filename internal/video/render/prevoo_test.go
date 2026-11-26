package render

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	jProjeto = `{"schema_version":1,"slug":"demo","formatos":["9:16"],"estilo_visual":"codigo"}`
	jMarca   = `{"schema_version":1,"nome":"M","formatos":["16:9"],"voz":{"motor":"edge-tts","voz":"pt-BR-ThalitaMultilingualNeural","velocidade":"+0%"}}`
	jScript  = `{"schema_version":1,"title":"t","lines":[{"text":"a"},{"text":"b"}]}`
	jCenas   = `{"schema_version":1,"projeto_slug":"demo","total_linhas":2,"cenas":[
{"indice_linha":0,"texto":"a","cenario_id":"cen_01","metafora_visual":"m","elementos_moveis":[{"id":"e0","tipo":"personagem","sprite_id":"narrador_pose_1","acao":"gesto_enfatiza"}]},
{"indice_linha":1,"texto":"b","cenario_id":"cen_01","metafora_visual":"m","elementos_moveis":[{"id":"e1","tipo":"camera","acao":"travelling_suave"}]}]}`
	jSprite = `{"id":"narrador_pose_1","personagem_id":"n","tipo":"pose","arquivo_png":"sprites/narrador_pose_1.png","largura":100,"altura":100,"alfa_valido":true,"bounding_box":{"x":0,"y":0,"w":100,"h":100}}`
)

func projetoFixture(t *testing.T, sobrescreve map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	arquivos := map[string]string{
		"projeto.json": jProjeto, "marca.json": jMarca, "script.json": jScript,
		"cenas.json": jCenas, "sprites/narrador_pose_1.json": jSprite,
	}
	for k, v := range sobrescreve {
		arquivos[k] = v
	}
	for nome, conteudo := range arquivos {
		if conteudo == "" {
			continue
		}
		caminho := filepath.Join(dir, filepath.FromSlash(nome))
		if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func depsOK() DepsPreVoo {
	return DepsPreVoo{
		ProvedorRegistrado: func(m string) bool { return m == "edge-tts" },
		Ferramenta:         func(string) bool { return true },
	}
}

func TestPreVooAprovaProjetoValido(t *testing.T) {
	dir := projetoFixture(t, nil)
	rel, c, err := ExecutarPreVoo(context.Background(), dir, depsOK())
	if err != nil || rel.Bloqueado {
		t.Fatalf("esperava aprovacao, obteve err=%v rel=%+v", err, rel)
	}
	if len(c.Script.Lines) != 2 || len(c.Storyboard.Cenas) != 2 || c.Sprites["narrador_pose_1"].ID == "" {
		t.Fatalf("contratos nao carregados: %+v", c)
	}
	if c.Config.Resolucoes[0] != "1080p" || c.Config.ReprovasWERMax != 3 {
		t.Fatalf("render.json ausente deve usar padrao: %+v", c.Config)
	}
}

func TestPreVooUsaRenderJSON(t *testing.T) {
	dir := projetoFixture(t, map[string]string{"render.json": `{"schema_version":1,"resolucoes":["720p","4k"],"reprovas_wer_max":2}`})
	_, c, err := ExecutarPreVoo(context.Background(), dir, depsOK())
	if err != nil || len(c.Config.Resolucoes) != 2 || c.Config.ReprovasWERMax != 2 {
		t.Fatalf("render.json nao aplicado: err=%v cfg=%+v", err, c.Config)
	}
}

func TestPreVooBloqueios(t *testing.T) {
	casos := []struct {
		nome   string
		arqs   map[string]string
		deps   func(*DepsPreVoo)
		alvo   error
		trecho string
	}{
		{"script ausente", map[string]string{"script.json": "{"}, nil, ErrPreVooBloqueado, "script.json"},
		{"cenas cobrem menos linhas", map[string]string{"script.json": `{"schema_version":1,"lines":[{"text":"a"},{"text":"b"},{"text":"c"}]}`}, nil, ErrPreVooBloqueado, "cobertura"},
		{"indice duplicado", map[string]string{"cenas.json": strings.Replace(jCenas, `"indice_linha":1`, `"indice_linha":0`, 1)}, nil, ErrPreVooBloqueado, "indice_linha"},
		{"campo desconhecido em cenas", map[string]string{"cenas.json": strings.Replace(jCenas, `"total_linhas":2`, `"total_linhas":2,"extra":1`, 1)}, nil, ErrPreVooBloqueado, "cenas.json"},
		{"sprite inexistente", map[string]string{"sprites/narrador_pose_1.json": strings.Replace(jSprite, `narrador_pose_1`, `outro`, 1)}, nil, ErrPreVooBloqueado, "sprite_id"},
		{"motor nao registrado", map[string]string{"marca.json": strings.Replace(jMarca, `edge-tts`, `desconhecido`, 1)}, nil, ErrProvedorDesconhecido, "desconhecido"},
		{"ffmpeg ausente", nil, func(d *DepsPreVoo) { d.Ferramenta = func(n string) bool { return n != "ffmpeg" } }, ErrPreVooBloqueado, "ffmpeg"},
		{"whisper ausente", nil, func(d *DepsPreVoo) { d.Ferramenta = func(n string) bool { return n != "faster-whisper" } }, ErrPreVooBloqueado, "faster-whisper"},
		{"resolucao invalida", map[string]string{"render.json": `{"schema_version":1,"resolucoes":["480p"]}`}, nil, ErrPreVooBloqueado, "render.json"},
		{"render.json com campo desconhecido", map[string]string{"render.json": `{"schema_version":1,"x":1}`}, nil, ErrPreVooBloqueado, "render.json"},
		{"pasta de saida sem escrita", map[string]string{"audio": "arquivo no lugar da pasta"}, nil, ErrPreVooBloqueado, "audio"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			dir := projetoFixture(t, c.arqs)
			d := depsOK()
			if c.deps != nil {
				c.deps(&d)
			}
			rel, _, err := ExecutarPreVoo(context.Background(), dir, d)
			if !errors.Is(err, c.alvo) {
				t.Fatalf("esperava %v, obteve %v", c.alvo, err)
			}
			if !strings.Contains(err.Error(), c.trecho) {
				t.Fatalf("erro deveria citar %q: %v", c.trecho, err)
			}
			if !rel.Bloqueado {
				t.Fatalf("relatorio deveria estar bloqueado: %+v", rel)
			}
		})
	}
}

func TestPreVooNaoAcionaAdaptadoresComContratoInvalido(t *testing.T) {
	dir := projetoFixture(t, map[string]string{"script.json": "{"})
	chamou := false
	d := DepsPreVoo{
		ProvedorRegistrado: func(string) bool { chamou = true; return true },
		Ferramenta:         func(string) bool { chamou = true; return true },
	}
	if _, _, err := ExecutarPreVoo(context.Background(), dir, d); err == nil || chamou {
		t.Fatalf("adaptador acionado antes do contrato valido (err=%v chamou=%v)", err, chamou)
	}
}

func TestPreVooDepsPadrao(t *testing.T) {
	d := normalizarDepsPreVoo(DepsPreVoo{})
	if !d.ProvedorRegistrado("edge-tts") || d.ProvedorRegistrado("inexistente") {
		t.Fatal("provedor padrao deve conhecer apenas edge-tts")
	}
	if d.Ferramenta("ferramenta-que-nao-existe-xyz") {
		t.Fatal("ferramenta inexistente nao pode estar disponivel")
	}
}

func TestPreVooPathJail(t *testing.T) {
	if _, err := resolverSeguro(t.TempDir(), "../fora.json"); err == nil {
		t.Fatal("caminho fora da raiz deve ser recusado")
	}
	if _, err := CarregarContratos(t.TempDir()); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("pasta vazia deve bloquear: %v", err)
	}
}

func TestPreVooContextoCancelado(t *testing.T) {
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	rel, _, err := ExecutarPreVoo(ctx, projetoFixture(t, nil), depsOK())
	if !errors.Is(err, ErrPreVooBloqueado) || rel.Bloqueado {
		t.Fatalf("contexto cancelado deve bloquear sem marcar item: err=%v rel=%+v", err, rel)
	}
}

func TestPreVooSpriteInvalido(t *testing.T) {
	for nome, conteudo := range map[string]string{"json quebrado": "{", "sem id": `{"personagem_id":"n"}`} {
		t.Run(nome, func(t *testing.T) {
			dir := projetoFixture(t, map[string]string{"sprites/ruim.json": conteudo})
			if _, err := CarregarContratos(dir); !errors.Is(err, ErrPreVooBloqueado) {
				t.Fatalf("esperava bloqueio, obteve %v", err)
			}
		})
	}
}

func TestPreVooLeituraForaDaRaiz(t *testing.T) {
	if _, err := lerBytes(t.TempDir(), "../x.json"); err == nil {
		t.Fatal("leitura fora da raiz deve falhar")
	}
}

func TestPreVooFerramentaLocal(t *testing.T) {
	if !ferramentaLocal("go") {
		t.Fatal("go deve estar no PATH durante os testes")
	}
	_ = ferramentaLocal("faster-whisper") // resultado depende da maquina; exercita a sonda Python
}

func TestPreVooSondarEscritaForaDaRaiz(t *testing.T) {
	if err := sondarEscrita(t.TempDir(), "../fora"); err == nil {
		t.Fatal("sonda fora da raiz deve falhar")
	}
}
