package render

import (
	"context"
	"fmt"
	"strings"
)

// EncoderFallback e o encoder de software usado quando nenhum de hardware serve.
const EncoderFallback = "libx264"

// ordemHardware e a prioridade de escolha: NVENC, QSV, AMF.
func ordemHardware() []string { return []string{"h264_nvenc", "h264_qsv", "h264_amf"} }

// Sondador lista os encoders que o ffmpeg declara.
type Sondador interface {
	Encoders(ctx context.Context) ([]string, error)
}

// SondaFFmpeg lista os encoders com `ffmpeg -encoders`; Med nulo usa o ffmpeg do sistema.
type SondaFFmpeg struct {
	Med Medidor
}

func (s SondaFFmpeg) Encoders(ctx context.Context) ([]string, error) {
	med := s.Med
	if med == nil {
		med = execSistema{}
	}
	saida, err := med.Saida(ctx, "ffmpeg", "-hide_banner", "-encoders")
	if err != nil {
		return nil, err
	}
	return parseEncoders(saida), nil
}

// parseEncoders le o nome do encoder (2a coluna) de cada linha apos o cabecalho "------" ou "Encoders:".
func parseEncoders(saida string) []string {
	var nomes []string
	for _, linha := range strings.Split(saida, "\n") {
		campos := strings.Fields(linha)
		if len(campos) >= 2 && len(campos[0]) == 6 && strings.ContainsAny(campos[0][:1], "VAS") {
			nomes = append(nomes, campos[1])
		}
	}
	return nomes
}

// EscolherEncoder devolve o primeiro encoder de hardware disponivel (NVENC, QSV, AMF) ou libx264.
// Falha na sonda tambem resulta em libx264: a escolha nunca bloqueia o render.
func EscolherEncoder(ctx context.Context, s Sondador) string {
	lista, err := s.Encoders(ctx)
	if err != nil {
		return EncoderFallback
	}
	disponivel := map[string]bool{}
	for _, n := range lista {
		disponivel[n] = true
	}
	for _, h := range ordemHardware() {
		if disponivel[h] {
			return h
		}
	}
	return EncoderFallback
}

// ExecutarComFallback roda o ffmpeg com o encoder escolhido. Se o de hardware falhar na execucao
// (listado, mas sem driver ou dispositivo), repete uma vez com libx264. Devolve o encoder que
// de fato produziu a saida, para RelatorioRender.Encoder.
func ExecutarComFallback(ctx context.Context, ex Executor, encoder string, montar func(encoder string) []string) (string, error) {
	if err := ctx.Err(); err != nil {
		return encoder, err
	}
	err := ex.Executar(ctx, "ffmpeg", montar(encoder)...)
	if err == nil || encoder == EncoderFallback {
		return encoder, err
	}
	if cerr := ctx.Err(); cerr != nil {
		return encoder, cerr
	}
	if err2 := ex.Executar(ctx, "ffmpeg", montar(EncoderFallback)...); err2 != nil {
		return EncoderFallback, fmt.Errorf("%s falhou (%v) e %s tambem falhou: %w", encoder, err, EncoderFallback, err2)
	}
	return EncoderFallback, nil
}
