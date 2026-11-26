package inventario

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSchemas(t *testing.T) {
	t.Run("erros tipados distintos e com prefixo", testErrosTipados)
	t.Run("AC-E07.2 tipo de estilo valido e invalido", testTipoDeEstilo)
	t.Run("AC-E07.6 area comum so aceita itens sem estilo", testTipoComum)
	t.Run("Estilo serializa com as chaves da SPEC", testEstiloJSON)
	t.Run("Item serializa com as chaves da SPEC", testItemJSON)
	t.Run("Referencia serializa com as chaves da SPEC", testReferenciaJSON)
	t.Run("mensagens dos erros iguais a SPEC", testMensagensErros)
}

func testReferenciaJSON(t *testing.T) {
	r := Referencia{Tipo: "video", Origem: "https://exemplo/ref.mp4"}
	exigirChaves(t, chavesJSON(t, r), []string{"tipo", "origem"}, []string{"nota"})
	r.Nota = "luz de recorte"
	exigirChaves(t, chavesJSON(t, r), []string{"tipo", "origem", "nota"}, nil)
}

func testMensagensErros(t *testing.T) {
	esperado := map[error]string{
		ErrEstiloDuplicado:   "video-inventario: estilo ja existe",
		ErrEstiloInexistente: "video-inventario: estilo inexistente",
		ErrEstiloNaoValidado: "video-inventario: estilo nao validado nao aceita item aprovado",
		ErrTipoInvalido:      "video-inventario: tipo de item invalido",
		ErrHashDivergente:    "video-inventario: hash do arquivo diverge",
		ErrHerancaCiclica:    "video-inventario: heranca de estilos em ciclo",
		ErrParteNaoSuportada: "video-inventario: parte nao suportada pelo renderizador",
	}
	for err, msg := range esperado {
		if err.Error() != msg {
			t.Fatalf("mensagem %q, esperado %q", err.Error(), msg)
		}
	}
}

func testErrosTipados(t *testing.T) {
	erros := []error{
		ErrEstiloDuplicado, ErrEstiloInexistente, ErrEstiloNaoValidado,
		ErrTipoInvalido, ErrHashDivergente, ErrHerancaCiclica, ErrParteNaoSuportada,
	}
	for i, a := range erros {
		if a == nil {
			t.Fatalf("erro %d nulo", i)
		}
		for j, b := range erros {
			if i != j && errors.Is(a, b) {
				t.Fatalf("erros %d e %d se confundem", i, j)
			}
		}
	}
	if ErrTipoInvalido.Error() != "video-inventario: tipo de item invalido" {
		t.Fatalf("mensagem inesperada: %q", ErrTipoInvalido.Error())
	}
}

func testTipoDeEstilo(t *testing.T) {
	validos := []string{"pele", "traje", "penteado", "adorno", "aderecos", "pose", "humor", "personagem", "cenario", "objeto"}
	for _, tipo := range validos {
		if err := ValidarTipo(tipo, false); err != nil {
			t.Fatalf("tipo %q deveria ser aceito: %v", tipo, err)
		}
	}
	for _, tipo := range []string{"", "chapeu", "sfx", "PENTEADO"} {
		if err := ValidarTipo(tipo, false); !errors.Is(err, ErrTipoInvalido) {
			t.Fatalf("tipo %q deveria retornar ErrTipoInvalido, veio %v", tipo, err)
		}
	}
}

func testTipoComum(t *testing.T) {
	for _, tipo := range []string{"sfx", "trilha", "voz"} {
		if err := ValidarTipo(tipo, true); err != nil {
			t.Fatalf("tipo comum %q deveria ser aceito: %v", tipo, err)
		}
	}
	for _, tipo := range []string{"penteado", "personagem", ""} {
		if err := ValidarTipo(tipo, true); !errors.Is(err, ErrTipoInvalido) {
			t.Fatalf("tipo %q em comum/ deveria retornar ErrTipoInvalido, veio %v", tipo, err)
		}
	}
}

func chavesJSON(t *testing.T, v any) map[string]any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func exigirChaves(t *testing.T, m map[string]any, presentes, ausentes []string) {
	t.Helper()
	for _, k := range presentes {
		if _, ok := m[k]; !ok {
			t.Fatalf("chave %q ausente", k)
		}
	}
	for _, k := range ausentes {
		if _, ok := m[k]; ok {
			t.Fatalf("chave %q vazia deveria ser omitida", k)
		}
	}
}

func testEstiloJSON(t *testing.T) {
	e := Estilo{
		SchemaVersion: 1, ID: "molde-pilula", Nome: "Molde", Status: "rascunho", Versao: "1.0.0",
		Renderizador: "molde-pilula", PartesSuportadas: []string{"penteado"},
		PromptMestre: "p", Paleta: []string{"#fff"},
		Luz:         &Luz{Rim: true, Cor: "#ffd", Quantidade: 0.5},
		Referencias: []Referencia{{Tipo: "imagem", Origem: "ref.png"}},
	}
	presentes := []string{"schema_version", "id", "nome", "status", "versao", "renderizador", "partes_suportadas", "prompt_mestre", "paleta", "luz", "referencias"}
	exigirChaves(t, chavesJSON(t, e), presentes, []string{"herda"})
}

func testItemJSON(t *testing.T) {
	it := Item{SchemaVersion: 1, ID: "coque", Tipo: "penteado", Estilo: "molde-pilula", EstiloVersao: "1.0.0", Origem: "p/2026-11-26"}
	presentes := []string{"schema_version", "id", "tipo", "estilo", "estilo_versao", "tags", "parametros", "origem", "aprovado"}
	exigirChaves(t, chavesJSON(t, it), presentes, []string{"arquivo", "hash", "herdavel", "obsoleto"})
}
