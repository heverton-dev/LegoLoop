package render

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executor roda um programa externo sem shell; injetavel nos testes.
type Executor interface {
	Executar(ctx context.Context, nome string, args ...string) error
}

type execSistema struct{}

func (execSistema) Executar(ctx context.Context, nome string, args ...string) error {
	saida, err := exec.CommandContext(ctx, nome, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", nome, err, strings.TrimSpace(string(saida)))
	}
	return nil
}

// EdgeTTS e o adaptador gratuito padrao: edge-tts gera MP3 e o ffmpeg converte para WAV PCM 24 kHz mono.
type EdgeTTS struct {
	Exec Executor
}

func (EdgeTTS) Nome() string { return "edge-tts" }
func (EdgeTTS) Pago() bool   { return false }

// Sintetizar usa o formato --opcao=valor para que o texto nunca seja lido como opcao.
func (e EdgeTTS) Sintetizar(ctx context.Context, texto string, op OpcoesVoz) ([]byte, error) {
	dir, err := os.MkdirTemp("", "edgetts-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	mp3, wav := filepath.Join(dir, "fala.mp3"), filepath.Join(dir, "fala.wav")
	args := []string{"--voice=" + op.Voz, "--text=" + texto, "--write-media=" + mp3}
	if op.Velocidade != "" {
		args = append(args, "--rate="+op.Velocidade)
	}
	if err := e.Exec.Executar(ctx, "edge-tts", args...); err != nil {
		return nil, err
	}
	conv := []string{"-y", "-i", mp3, "-ar", "24000", "-ac", "1", "-c:a", "pcm_s16le", wav}
	if err := e.Exec.Executar(ctx, "ffmpeg", conv...); err != nil {
		return nil, err
	}
	return os.ReadFile(wav)
}

// RegistroPadrao devolve o registro com os provedores gratuitos embutidos.
func RegistroPadrao(exec Executor) *Registro {
	if exec == nil {
		exec = execSistema{}
	}
	r := NovoRegistro()
	edge := EdgeTTS{Exec: exec}
	r.provedores[edge.Nome()] = edge
	return r
}
