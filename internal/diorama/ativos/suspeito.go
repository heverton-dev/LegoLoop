package ativos

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	_ "image/jpeg" // registra o decodificador JPEG
	_ "image/png"  // registra o decodificador PNG
	"math"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

// Motivos de entrega suspeita do Flow (AC-E09.5, achado A11).
const (
	MotivoSemRenomeacao  = "sem_confirmacao_renomeacao"
	MotivoHashCruzado    = "hash_cruzado_outro_tipo"
	MotivoCorMedia       = "cor_media_fora_do_padrao"
	MotivoDimensao       = "dimensao_fora_do_padrao"
	MotivoImagemIlegivel = "imagem_ilegivel"
)

// ErrFolhaContatoAusente bloqueia a confirmacao de suspeito sem folha de contato e aprovador.
var ErrFolhaContatoAusente = errors.New("diorama-ativos: confirmacao de suspeito exige folha de contato e aprovador")

// EntregaFlow descreve o download registrado de um job de imagem do Flow.
type EntregaFlow struct {
	Item                 inventario.Item
	RenomeacaoConfirmada bool
}

// PadraoEstilo fixa a dimensao e a cor media esperadas do estilo.
type PadraoEstilo struct {
	Largura       int
	Altura        int
	CorMedia      [3]uint8
	ToleranciaCor float64
}

// AvaliacaoEntrega traz o item avaliado e os motivos de suspeita.
type AvaliacaoEntrega struct {
	Item              inventario.Item
	Suspeito          bool
	Motivos           []string
	ExigeFolhaContato bool
}

// AvaliarEntregaFlow marca o item como suspeito quando falta a renomeacao, o hash cruza
// com item de outro tipo ou a imagem foge do padrao do estilo; suspeito nunca sai aprovado.
func AvaliarEntregaFlow(dados []byte, e EntregaFlow, conhecidos []inventario.Item, p PadraoEstilo) AvaliacaoEntrega {
	it := e.Item
	it.Parametros = copiarParametros(it.Parametros)
	it.Hash = hashBytes(dados)
	var motivos []string
	if !e.RenomeacaoConfirmada {
		motivos = append(motivos, MotivoSemRenomeacao)
	}
	if hashCruzado(it, conhecidos) {
		motivos = append(motivos, MotivoHashCruzado)
	}
	motivos = append(motivos, motivosEstilo(dados, p)...)
	av := AvaliacaoEntrega{Item: it, Motivos: motivos}
	if len(motivos) == 0 {
		return av
	}
	av.Suspeito = true
	av.ExigeFolhaContato = true
	av.Item.Aprovado = false
	av.Item.Parametros["status"] = "suspeito"
	av.Item.Parametros["suspeito_motivos"] = append([]string(nil), motivos...)
	return av
}

// ConfirmarSuspeito libera o item suspeito so com folha de contato e aprovador humano.
func ConfirmarSuspeito(av AvaliacaoEntrega, folhaContato, aprovadoPor string) (inventario.Item, error) {
	folha := strings.TrimSpace(folhaContato)
	quem := strings.TrimSpace(aprovadoPor)
	if folha == "" || quem == "" {
		return inventario.Item{}, ErrFolhaContatoAusente
	}
	it := av.Item
	it.Parametros = copiarParametros(it.Parametros)
	it.Parametros["status"] = "confirmado"
	it.Parametros["folha_contato"] = folha
	it.Parametros["confirmado_por"] = quem
	it.Aprovado = true
	return it, nil
}

func hashCruzado(it inventario.Item, conhecidos []inventario.Item) bool {
	for _, k := range conhecidos {
		if k.Hash == it.Hash && k.Tipo != it.Tipo && k.ID != it.ID {
			return true
		}
	}
	return false
}

// motivosEstilo checa dimensao e cor media da imagem contra o padrao do estilo.
func motivosEstilo(dados []byte, p PadraoEstilo) []string {
	img, _, err := image.Decode(bytes.NewReader(dados))
	if err != nil {
		return []string{MotivoImagemIlegivel}
	}
	var motivos []string
	b := img.Bounds()
	if p.Largura > 0 && p.Altura > 0 && (b.Dx() != p.Largura || b.Dy() != p.Altura) {
		motivos = append(motivos, MotivoDimensao)
	}
	if p.ToleranciaCor > 0 && distanciaCor(corMedia(img), p.CorMedia) > p.ToleranciaCor {
		motivos = append(motivos, MotivoCorMedia)
	}
	return motivos
}

func corMedia(img image.Image) [3]float64 {
	b := img.Bounds()
	var soma [3]float64
	n := float64(b.Dx() * b.Dy())
	if n == 0 {
		return soma
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			soma[0] += float64(r >> 8)
			soma[1] += float64(g >> 8)
			soma[2] += float64(bl >> 8)
		}
	}
	return [3]float64{soma[0] / n, soma[1] / n, soma[2] / n}
}

func distanciaCor(m [3]float64, alvo [3]uint8) float64 {
	var s float64
	for i := range m {
		d := m[i] - float64(alvo[i])
		s += d * d
	}
	return math.Sqrt(s)
}

func hashBytes(dados []byte) string {
	s := sha256.Sum256(dados)
	return hex.EncodeToString(s[:])
}

func copiarParametros(orig map[string]any) map[string]any {
	out := make(map[string]any, len(orig)+2)
	for k, v := range orig {
		out[k] = v
	}
	return out
}
