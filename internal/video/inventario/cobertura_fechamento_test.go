package inventario

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestVerificarEstiloSemFalhas cobre VerificarEstilo com todos os itens integros (0 falhas).
func TestVerificarEstiloSemFalhas(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	if err := s.NovoEstilo(estiloRascunho("com-itens")); err != nil {
		t.Fatal(err)
	}
	it := itemComArquivo(t, raiz)
	it.Estilo, it.Tipo = "com-itens", "cenario"
	if err := s.GravarItem(it); err != nil {
		t.Fatal(err)
	}
	falhas, err := s.VerificarEstilo("com-itens")
	if err != nil || len(falhas) != 0 {
		t.Fatalf("falhas=%v err=%v, queria 0 falhas", falhas, err)
	}
}

// TestVerificarEstiloComHashDivergente cobre a lista de falhas de VerificarEstilo.
func TestVerificarEstiloComHashDivergente(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	if err := s.NovoEstilo(estiloRascunho("com-falha")); err != nil {
		t.Fatal(err)
	}
	it := itemComArquivo(t, raiz)
	it.Estilo, it.Tipo = "com-falha", "cenario"
	if err := s.GravarItem(it); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(raiz, it.Arquivo), []byte("adulterado"), 0o644); err != nil {
		t.Fatal(err)
	}
	falhas, err := s.VerificarEstilo("com-falha")
	if err != nil || len(falhas) != 1 || !errors.Is(falhas[0], ErrHashDivergente) {
		t.Fatalf("falhas=%v err=%v, queria 1 ErrHashDivergente", falhas, err)
	}
}

// TestVerificarEstiloInexistente cobre o erro de LerEstilo propagado por VerificarEstilo.
func TestVerificarEstiloInexistente(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if _, err := s.VerificarEstilo("nao-existe"); !errors.Is(err, ErrEstiloInexistente) {
		t.Fatalf("err = %v, queria ErrEstiloInexistente", err)
	}
}

// TestLerEstiloIDInvalido cobre o kebab invalido em LerEstilo.
func TestLerEstiloIDInvalido(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	for _, id := range []string{"", "-x", "x-", "Maiuscula", "a--b"} {
		if _, err := s.LerEstilo(id); !errors.Is(err, ErrEstiloInexistente) {
			t.Errorf("LerEstilo(%q) err=%v, queria ErrEstiloInexistente", id, err)
		}
	}
}

// TestLerItensTipoVazioListaTodos cobre lerItens com tipo vazio (via Buscar sem filtro de tipo).
func TestLerItensTipoVazioListaTodos(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	if err := s.NovoEstilo(estiloRascunho("misto")); err != nil {
		t.Fatal(err)
	}
	if err := s.GravarItem(Item{ID: "p1", Tipo: "pose", Estilo: "misto"}); err != nil {
		t.Fatal(err)
	}
	if err := s.GravarItem(Item{ID: "h1", Tipo: "humor", Estilo: "misto"}); err != nil {
		t.Fatal(err)
	}
	itens, err := s.lerItens("misto", "")
	if err != nil || len(itens) != 2 {
		t.Fatalf("itens=%v err=%v, queria 2", itens, err)
	}
}

// TestHashArquivoErroLeitura cobre o erro de io.Copy (abrir um diretorio como se fosse arquivo).
func TestHashArquivoErroLeitura(t *testing.T) {
	if _, err := HashArquivo(t.TempDir()); err == nil {
		t.Fatal("esperava erro ao tentar ler um diretorio como arquivo")
	}
}

// TestLerEstiloJSONCorrompido cobre o erro de decodificacao em LerEstilo.
func TestLerEstiloJSONCorrompido(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	caminho := filepath.Join(raiz, "inventario", "video", "estilos", "ruim", "estilo.json")
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte("{nao e json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.LerEstilo("ruim"); err == nil {
		t.Fatal("esperava erro de decodificacao")
	}
}

// TestGravarJSONErroDeCodificacao cobre o erro de json.MarshalIndent em gravarJSON.
func TestGravarJSONErroDeCodificacao(t *testing.T) {
	if err := gravarJSON(filepath.Join(t.TempDir(), "x.json"), make(chan int)); err == nil {
		t.Fatal("esperava erro ao codificar um canal")
	}
}

// TestLerAprovacaoCaminhoVazio cobre o erro de caminho vazio em LerAprovacao.
func TestLerAprovacaoCaminhoVazio(t *testing.T) {
	if _, err := LerAprovacao("", "x"); !errors.Is(err, ErrAprovacaoAusente) {
		t.Fatalf("err = %v, queria ErrAprovacaoAusente", err)
	}
}

// TestLerAprovacaoJSONInvalido cobre o erro de decodificacao em LerAprovacao.
func TestLerAprovacaoJSONInvalido(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "aprovacao.json")
	if err := os.WriteFile(caminho, []byte("{nao e json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LerAprovacao(caminho, "x"); !errors.Is(err, ErrAprovacaoInvalida) {
		t.Fatalf("err = %v, queria ErrAprovacaoInvalida", err)
	}
}

// TestLerAprovacaoAprovadoPorVazio cobre o branch aprovado_por vazio.
func TestLerAprovacaoAprovadoPorVazio(t *testing.T) {
	caminho := gravarAprovacao(t, t.TempDir(), Aprovacao{Estilo: "x", PersonagensTeste: []string{"p1"}})
	if _, err := LerAprovacao(caminho, "x"); !errors.Is(err, ErrAprovacaoInvalida) {
		t.Fatalf("err = %v, queria ErrAprovacaoInvalida (aprovado_por vazio)", err)
	}
}

// TestValidarEstiloSemReferenciasPropaga ja existe para a leitura real; aqui cobrimos LerValidacao.
func TestLerValidacaoAusenteEInvalida(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	if _, err := s.LerValidacao("sem-validacao"); err == nil {
		t.Fatal("esperava erro: validacao.json ausente")
	}
	caminho := filepath.Join(raiz, "inventario", "video", "estilos", "ruim2", "validacao.json")
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte("{nao e json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.LerValidacao("ruim2"); err == nil {
		t.Fatal("esperava erro de decodificacao de validacao.json")
	}
}

// TestBuscarPersonagensEstiloInexistentePropagaErro cobre o branch de erro de Buscar em BuscarPersonagens.
func TestBuscarPersonagensEstiloInexistentePropagaErro(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if _, err := s.BuscarPersonagens("nao-existe", "", ""); err == nil {
		t.Fatal("esperava erro: estilo inexistente")
	}
}

// TestValidarEstiloSemReferencias cobre o branch "sem referencias" de ValidarEstilo (estilo.json
// editado a mao depois de criado, sem passar por NovoEstilo).
func TestValidarEstiloSemReferencias(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	e := estiloRascunho("sem-ref-val")
	if err := s.NovoEstilo(e); err != nil {
		t.Fatal(err)
	}
	e.Referencias = nil
	if err := gravarJSON(s.caminhoEstilo("sem-ref-val"), e); err != nil {
		t.Fatal(err)
	}
	arq := gravarAprovacao(t, t.TempDir(), Aprovacao{Estilo: "sem-ref-val", AprovadoPor: "x", PersonagensTeste: []string{"p1"}})
	if _, err := s.ValidarEstilo("sem-ref-val", arq, time.Now()); !errors.Is(err, ErrEstiloInvalido) {
		t.Fatalf("err = %v, queria ErrEstiloInvalido", err)
	}
}

// TestVerificarEstiloItemJSONCorrompido cobre o erro de decodificacao dentro de lerItens, via
// VerificarEstilo (que chama lerItens com tipo vazio).
func TestVerificarEstiloItemJSONCorrompido(t *testing.T) {
	raiz := t.TempDir()
	s := Store{Raiz: raiz}
	if err := s.NovoEstilo(estiloRascunho("item-ruim")); err != nil {
		t.Fatal(err)
	}
	caminho := filepath.Join(raiz, "inventario", "video", "estilos", "item-ruim", "itens", "cenario", "x.json")
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte("{nao e json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.VerificarEstilo("item-ruim"); err == nil {
		t.Fatal("esperava erro de decodificacao de item")
	}
}
