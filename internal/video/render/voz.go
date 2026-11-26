package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

const VozPadrao = "pt-BR-ThalitaMultilingualNeural"

// OpcoesVoz carrega a escolha de voz e velocidade vindas de marca.json.
type OpcoesVoz struct {
	Voz        string
	Velocidade string
}

// Sintetizador e o contrato de qualquer provedor de voz; devolve WAV PCM de 16 bits.
type Sintetizador interface {
	Nome() string
	Pago() bool
	Sintetizar(ctx context.Context, texto string, op OpcoesVoz) ([]byte, error)
}

// GateCreditos autoriza (nil) ou nega (erro) uma chamada paga.
type GateCreditos func(ctx context.Context, provedor string) error

// Registro guarda os provedores disponiveis por nome.
type Registro struct {
	provedores map[string]Sintetizador
}

func NovoRegistro() *Registro {
	return &Registro{provedores: map[string]Sintetizador{}}
}

// Registrar rejeita nome vazio e nome repetido.
func (r *Registro) Registrar(s Sintetizador) error {
	nome := strings.TrimSpace(s.Nome())
	if nome == "" {
		return invalido("provedor sem nome")
	}
	if _, ok := r.provedores[nome]; ok {
		return invalido("provedor %q ja registrado", nome)
	}
	r.provedores[nome] = s
	return nil
}

func (r *Registro) Existe(nome string) bool {
	_, ok := r.provedores[nome]
	return ok
}

func (r *Registro) Obter(nome string) (Sintetizador, error) {
	if s, ok := r.provedores[nome]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrProvedorDesconhecido, nome)
}

// ResultadoVoz descreve o que a Etapa 6 gravou em audio/.
type ResultadoVoz struct {
	Provedor   string
	Voz        string
	ArquivoVoz string
	Linhas     []string
	DuracoesS  []float64
}

// Locutor sintetiza o roteiro inteiro por um provedor registrado.
type Locutor struct {
	Registro *Registro
	Gate     GateCreditos
}

// preparo agrupa o que a sintese precisa depois de validar provedor, gate e pronuncia.
type preparo struct {
	s    Sintetizador
	op   OpcoesVoz
	voz  string
	pron map[string]string
}

func (l Locutor) preparar(ctx context.Context, projDir string, marca roteiro.Marca, script roteiro.Script) (preparo, error) {
	var p preparo
	if len(script.Lines) == 0 {
		return p, invalido("roteiro sem linhas")
	}
	if l.Registro == nil {
		return p, fmt.Errorf("%w: %q", ErrProvedorDesconhecido, marca.Voz.Motor)
	}
	s, err := l.Registro.Obter(marca.Voz.Motor)
	if err != nil {
		return p, err
	}
	if err := l.autorizar(ctx, s); err != nil {
		return p, err
	}
	p.s, p.voz = s, strings.TrimSpace(marca.Voz.Voz)
	if p.voz == "" {
		p.voz = VozPadrao
	}
	p.op = OpcoesVoz{Voz: p.voz, Velocidade: marca.Voz.Velocidade}
	p.pron, err = carregarPronuncia(projDir)
	return p, err
}

// SintetizarRoteiro grava audio/linha_NN.wav e audio/voice.wav somente se todas as linhas passarem.
func (l Locutor) SintetizarRoteiro(ctx context.Context, projDir string, marca roteiro.Marca, script roteiro.Script) (ResultadoVoz, error) {
	res := ResultadoVoz{Provedor: marca.Voz.Motor, Voz: marca.Voz.Voz}
	p, err := l.preparar(ctx, projDir, marca, script)
	if err != nil {
		return res, err
	}
	res.Voz = p.voz
	trechos, err := sintetizarLinhas(ctx, p.s, script, p.pron, p.op)
	if err != nil {
		return res, err
	}
	return gravarVoz(projDir, res, trechos)
}

func (l Locutor) autorizar(ctx context.Context, s Sintetizador) error {
	if !s.Pago() {
		return nil
	}
	if l.Gate == nil {
		return fmt.Errorf("%w: provedor pago %q sem gate configurado", ErrCreditosNegados, s.Nome())
	}
	if err := l.Gate(ctx, s.Nome()); err != nil {
		return fmt.Errorf("%w: %v", ErrCreditosNegados, err)
	}
	return nil
}

func sintetizarLinhas(ctx context.Context, s Sintetizador, script roteiro.Script, pron map[string]string, op OpcoesVoz) ([][]byte, error) {
	trechos := make([][]byte, 0, len(script.Lines))
	for i, linha := range script.Lines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		b, err := s.Sintetizar(ctx, AplicarPronuncia(linha.Text, pron), op)
		if err != nil {
			return nil, fmt.Errorf("diorama-render: linha %d: %w", i, err)
		}
		if _, _, err := lerWAV(b); err != nil {
			return nil, fmt.Errorf("diorama-render: linha %d: %w", i, err)
		}
		trechos = append(trechos, b)
	}
	return trechos, nil
}

func gravarVoz(projDir string, res ResultadoVoz, trechos [][]byte) (ResultadoVoz, error) {
	voz, err := juntarWAV(trechos)
	if err != nil {
		return res, err
	}
	for i, t := range trechos {
		f, dados, _ := lerWAV(t)
		rel := fmt.Sprintf("audio/linha_%02d.wav", i)
		if err := gravarSeguro(projDir, rel, t); err != nil {
			return res, err
		}
		res.Linhas = append(res.Linhas, rel)
		res.DuracoesS = append(res.DuracoesS, duracaoWAV(f, len(dados)))
	}
	res.ArquivoVoz = "audio/voice.wav"
	return res, gravarSeguro(projDir, res.ArquivoVoz, voz)
}

func gravarSeguro(projDir, rel string, dados []byte) error {
	caminho, err := resolverSeguro(projDir, rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	return os.WriteFile(caminho, dados, 0o644)
}

// carregarPronuncia le pronuncia.json opcional (mapa palavra para grafia fonetica).
func carregarPronuncia(projDir string) (map[string]string, error) {
	caminho, err := resolverSeguro(projDir, "pronuncia.json")
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(caminho); os.IsNotExist(err) {
		return nil, nil
	}
	var mapa map[string]string
	if err := lerJSON(projDir, "pronuncia.json", &mapa, true); err != nil {
		return nil, fmt.Errorf("%w: pronuncia.json: %v", ErrSchemaInvalido, err)
	}
	return mapa, nil
}

// AplicarPronuncia troca apenas palavras inteiras, sem tocar em trechos maiores.
func AplicarPronuncia(texto string, mapa map[string]string) string {
	if len(mapa) == 0 {
		return texto
	}
	var saida, palavra strings.Builder
	descarga := func() {
		p := palavra.String()
		if troca, ok := mapa[p]; ok {
			p = troca
		}
		saida.WriteString(p)
		palavra.Reset()
	}
	for _, r := range texto {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			palavra.WriteRune(r)
			continue
		}
		descarga()
		saida.WriteRune(r)
	}
	descarga()
	return saida.String()
}
