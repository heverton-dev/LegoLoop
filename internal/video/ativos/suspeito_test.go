package ativos

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

func pngSuspeito(t *testing.T, w, h int, c color.RGBA) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func padraoSuspeito() PadraoEstilo {
	return PadraoEstilo{Largura: 64, Altura: 36, CorMedia: [3]uint8{40, 60, 120}, ToleranciaCor: 30}
}

func entregaSuspeito(confirmada bool) EntregaFlow {
	it := inventario.Item{ID: "praca-dia", Tipo: "cenario", Aprovado: true}
	return EntregaFlow{Item: it, RenomeacaoConfirmada: confirmada}
}

func hashSuspeito(d []byte) string {
	s := sha256.Sum256(d)
	return hex.EncodeToString(s[:])
}

func TestSuspeitoEntregaLimpaNaoMarca(t *testing.T) {
	dados := pngSuspeito(t, 64, 36, color.RGBA{40, 60, 120, 255})
	av := AvaliarEntregaFlow(dados, entregaSuspeito(true), nil, padraoSuspeito())
	if av.Suspeito || len(av.Motivos) != 0 || av.ExigeFolhaContato {
		t.Fatalf("entrega limpa marcada: %+v", av)
	}
	if av.Item.Hash != hashSuspeito(dados) {
		t.Fatalf("hash nao registrado: %q", av.Item.Hash)
	}
}

func TestSuspeitoSemConfirmacaoDeRenomeacao(t *testing.T) {
	dados := pngSuspeito(t, 64, 36, color.RGBA{40, 60, 120, 255})
	av := AvaliarEntregaFlow(dados, entregaSuspeito(false), nil, padraoSuspeito())
	exigirSuspeito(t, av, MotivoSemRenomeacao)
}

func TestSuspeitoHashCruzadoComOutroTipo(t *testing.T) {
	dados := pngSuspeito(t, 64, 36, color.RGBA{40, 60, 120, 255})
	outro := inventario.Item{ID: "ficha-ana", Tipo: "personagem", Hash: hashSuspeito(dados)}
	mesmoTipo := inventario.Item{ID: "praca-noite", Tipo: "cenario", Hash: hashSuspeito(dados)}
	av := AvaliarEntregaFlow(dados, entregaSuspeito(true), []inventario.Item{mesmoTipo}, padraoSuspeito())
	if av.Suspeito {
		t.Fatalf("hash igual de mesmo tipo marcado: %+v", av.Motivos)
	}
	av = AvaliarEntregaFlow(dados, entregaSuspeito(true), []inventario.Item{mesmoTipo, outro}, padraoSuspeito())
	exigirSuspeito(t, av, MotivoHashCruzado)
}

func TestSuspeitoForaDoPadraoDeEstilo(t *testing.T) {
	cor := AvaliarEntregaFlow(pngSuspeito(t, 64, 36, color.RGBA{250, 250, 10, 255}), entregaSuspeito(true), nil, padraoSuspeito())
	exigirSuspeito(t, cor, MotivoCorMedia)
	dim := AvaliarEntregaFlow(pngSuspeito(t, 36, 64, color.RGBA{40, 60, 120, 255}), entregaSuspeito(true), nil, padraoSuspeito())
	exigirSuspeito(t, dim, MotivoDimensao)
	ilegivel := AvaliarEntregaFlow([]byte("nao e imagem"), entregaSuspeito(true), nil, padraoSuspeito())
	exigirSuspeito(t, ilegivel, MotivoImagemIlegivel)
}

func TestSuspeitoConfirmacaoExigeFolhaDeContato(t *testing.T) {
	av := AvaliarEntregaFlow([]byte("x"), entregaSuspeito(false), nil, padraoSuspeito())
	if _, err := ConfirmarSuspeito(av, " ", "heverton"); !errors.Is(err, ErrFolhaContatoAusente) {
		t.Fatalf("sem folha deveria falhar: %v", err)
	}
	if _, err := ConfirmarSuspeito(av, "folhas/praca.png", ""); !errors.Is(err, ErrFolhaContatoAusente) {
		t.Fatalf("sem aprovador deveria falhar: %v", err)
	}
	it, err := ConfirmarSuspeito(av, "folhas/praca.png", "heverton")
	if err != nil {
		t.Fatal(err)
	}
	if !it.Aprovado || it.Parametros["status"] != "confirmado" || it.Parametros["folha_contato"] != "folhas/praca.png" {
		t.Fatalf("confirmacao nao registrada: %+v", it)
	}
}

func exigirSuspeito(t *testing.T, av AvaliacaoEntrega, motivo string) {
	t.Helper()
	if !av.Suspeito || av.Item.Aprovado || !av.ExigeFolhaContato {
		t.Fatalf("esperado suspeito nao aprovado com folha: %+v", av)
	}
	if av.Item.Parametros["status"] != "suspeito" {
		t.Fatalf("status no item: %v", av.Item.Parametros["status"])
	}
	for _, m := range av.Motivos {
		if m == motivo {
			return
		}
	}
	t.Fatalf("motivo %q ausente em %v", motivo, av.Motivos)
}
