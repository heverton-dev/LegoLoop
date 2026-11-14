package triagem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type resolvedorFalso struct {
	chamado  bool
	resposta Decisao
	err      error
}

func (r *resolvedorFalso) Resolver(_ context.Context, _ string, _ []TriagemSpec) (Decisao, error) {
	r.chamado = true
	return r.resposta, r.err
}

func TestValidarTriagemSpec(t *testing.T) {
	valida := GerarTriagemPadrao("video-rea", []string{"ingest", "forensics"})
	if err := valida.Validar(); err != nil {
		t.Fatalf("esperava spec válida, obteve erro: %v", err)
	}

	casosInvalidos := []struct {
		nome string
		mut  func(*TriagemSpec)
	}{
		{"schema_version inválido", func(s *TriagemSpec) { s.SchemaVersion = "triagem.v0" }},
		{"processo vazio", func(s *TriagemSpec) { s.Processo = "" }},
		{"pontos_entrada vazio", func(s *TriagemSpec) { s.PontosEntrada = nil }},
		{"regras vazias", func(s *TriagemSpec) { s.RegrasDeterministicas = nil }},
	}

	for _, tc := range casosInvalidos {
		t.Run(tc.nome, func(t *testing.T) {
			copia := GerarTriagemPadrao("video-rea", []string{"ingest"})
			tc.mut(&copia)
			if err := copia.Validar(); err == nil {
				t.Fatalf("esperava erro para %s", tc.nome)
			}
		})
	}
}

func TestDecidirDeterministicoPorComandoESlug(t *testing.T) {
	catalogo := []TriagemSpec{
		GerarTriagemPadrao("video-rea", []string{"ingest", "forensics"}),
		GerarTriagemPadrao("cartoon-shorts", []string{"roteiro", "render"}),
	}
	llm := &resolvedorFalso{}

	casos := []struct {
		entrada  string
		esperado string
	}{
		{"/legoloop-video-rea analisar canal", "video-rea"},
		{"legoloop cartoon-shorts --theme gato", "cartoon-shorts"},
		{"video-rea", "video-rea"},
	}

	for _, tc := range casos {
		dec, err := Decidir(context.Background(), Entrada{Texto: tc.entrada}, catalogo, llm)
		if err != nil {
			t.Fatalf("erro inesperado para %q: %v", tc.entrada, err)
		}
		if dec.Processo != tc.esperado {
			t.Errorf("entrada %q: esperado %s, obtido %s", tc.entrada, tc.esperado, dec.Processo)
		}
		if dec.Camada != "deterministica" {
			t.Errorf("esperava camada deterministica, obteve %s", dec.Camada)
		}
		if llm.chamado {
			t.Errorf("LLM não deveria ser chamado para %q", tc.entrada)
		}
	}
}

func TestDecidirDeterministicoPorArtefatoWorkspace(t *testing.T) {
	ws := t.TempDir()
	artefato := filepath.Join(ws, "forensics.json")
	if err := os.WriteFile(artefato, []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	specVideo := GerarTriagemPadrao("video-rea", []string{"ingest", "forensics"})
	specVideo.RegrasDeterministicas[0].ArtefatosAtivos = []string{"forensics.json"}

	specCartoon := GerarTriagemPadrao("cartoon-shorts", []string{"roteiro"})
	specCartoon.RegrasDeterministicas[0].ArtefatosAtivos = []string{"storyboard.json"}

	llm := &resolvedorFalso{}
	dec, err := Decidir(context.Background(), Entrada{
		Texto:     "continuar trabalho",
		Workspace: ws,
	}, []TriagemSpec{specVideo, specCartoon}, llm)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if dec.Processo != "video-rea" {
		t.Errorf("esperava video-rea, obteve %s", dec.Processo)
	}
	if llm.chamado {
		t.Errorf("LLM não deveria ser chamado quando há artefato único ativo")
	}
}

func TestDecidirProbabilisticoEmLinguagemNaturalOuEmpate(t *testing.T) {
	specA := GerarTriagemPadrao("video-rea", []string{"ingest"})
	specB := GerarTriagemPadrao("cartoon-shorts", []string{"roteiro"})
	catalogo := []TriagemSpec{specA, specB}

	llm := &resolvedorFalso{
		resposta: Decisao{
			Processo:           "cartoon-shorts",
			EtapaInicial:       "roteiro",
			Confianca:          0.88,
			Justificativa:      "pedido aberto de animação",
			PrecisaConfirmacao: false,
		},
	}

	dec, err := Decidir(context.Background(), Entrada{
		Texto: "quero criar um vídeo animado curto sobre gatos",
	}, catalogo, llm)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !llm.chamado {
		t.Fatalf("esperava chamada ao ResolvedorLLM")
	}
	if dec.Processo != "cartoon-shorts" || dec.Camada != "probabilistica" {
		t.Errorf("decisão inesperada: %+v", dec)
	}

	// Empate entre dois processos com o mesmo comando/artefato deve acionar LLM
	ws := t.TempDir()
	_ = os.WriteFile(filepath.Join(ws, "comum.json"), []byte(`{}`), 0o644)
	specA.RegrasDeterministicas[0].ArtefatosAtivos = []string{"comum.json"}
	specB.RegrasDeterministicas[0].ArtefatosAtivos = []string{"comum.json"}

	llmEmpate := &resolvedorFalso{
		resposta: Decisao{Processo: "video-rea", EtapaInicial: "ingest", Confianca: 0.75},
	}
	decEmp, err := Decidir(context.Background(), Entrada{
		Texto:     "prosseguir",
		Workspace: ws,
	}, []TriagemSpec{specA, specB}, llmEmpate)
	if err != nil {
		t.Fatalf("erro inesperado no empate: %v", err)
	}
	if !llmEmpate.chamado || decEmp.Camada != "probabilistica" {
		t.Errorf("empate deveria acionar camada probabilística, obteve %+v", decEmp)
	}
}

func TestDecidirSemResolvedorEmCasoAmbiguoRetornaErro(t *testing.T) {
	catalogo := []TriagemSpec{GerarTriagemPadrao("video-rea", []string{"ingest"})}
	_, err := Decidir(context.Background(), Entrada{Texto: "algo totalmente vago"}, catalogo, nil)
	if !errors.Is(err, ErrResolvedorAusente) {
		t.Errorf("esperava ErrResolvedorAusente, obteve %v", err)
	}
}

func TestDecidirDeterministicoPorEtapa(t *testing.T) {
	spec := GerarTriagemPadrao("video-rea", []string{"ingest", "forensics"})
	catalogo := []TriagemSpec{spec}

	casos := []struct {
		nome  string
		texto string
		etapa string
	}{
		{"slash", "/video-rea-forensics", "forensics"},
		{"flag", "legoloop video-rea --etapa forensics", "forensics"},
		{"palavra_chave", "rodar forensics agora", "forensics"},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			dec, err := Decidir(context.Background(), Entrada{Texto: c.texto}, catalogo, nil)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if dec.Processo != "video-rea" || dec.EtapaInicial != c.etapa || dec.Camada != "deterministica" {
				t.Errorf("decisão inesperada: %+v", dec)
			}
		})
	}
}

func TestMemoriaExecucaoCicloCompleto(t *testing.T) {
	dir := t.TempDir()
	mem := MemoriaExecucao{
		Processo:     "video-rea",
		EtapaAtual:   "ingest",
		ProximoPasso: "forensics",
	}

	if err := SalvarMemoria(dir, mem); err != nil {
		t.Fatalf("erro ao salvar memória: %v", err)
	}

	lida, err := LerMemoria(dir)
	if err != nil {
		t.Fatalf("erro ao ler memória: %v", err)
	}
	if lida.Processo != mem.Processo || lida.EtapaAtual != mem.EtapaAtual || lida.ProximoPasso != mem.ProximoPasso {
		t.Errorf("memória lida diverge: %+v", lida)
	}
	if lida.AtualizadoEm.IsZero() {
		t.Errorf("AtualizadoEm deveria ser preenchido automaticamente")
	}

	if err := LimparMemoria(dir); err != nil {
		t.Fatalf("erro ao limpar memória: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ArquivoMemoria)); !os.IsNotExist(err) {
		t.Errorf("arquivo de memória deveria ter sido removido, err=%v", err)
	}
	if err := LimparMemoria(dir); err != nil {
		t.Errorf("limpar memória inexistente não deveria falhar: %v", err)
	}
}
