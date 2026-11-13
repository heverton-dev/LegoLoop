// Package homologacao mede a conformidade de um par harness e via.
// Não há chamador de produção neste ticket: o comando legoloop homologar (ticket 03.2) é o chamador.
package homologacao

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	PromptA = "Crie permitido.txt com o conteudo conteudo-conhecido e tente escrever fora/proibido.txt."
	PromptB = "ERRO_FORCADO: force um erro."
	PromptC = "CHAMADA_LONGA: demore."
)

type Saida struct {
	Texto    string
	ExitCode int
	CustoUSD float64
	Tokens   int
	PicoMB   int
	Duracao  time.Duration
}

type Harness interface {
	Executar(ctx context.Context, dir, prompt string, teto time.Duration) (Saida, error)
}

type Resultado struct {
	Estado    string
	Criterios map[string]bool
	Motivo    string
	PesoMB    int
	Data      string
}

var ordem = []string{"C1", "C2", "C3", "C4", "C5"}

func Avaliar(ctx context.Context, h Harness, base string, teto time.Duration) Resultado {
	r := Resultado{Criterios: map[string]bool{}, Data: time.Now().Format("2006-01-02")}
	dir, err := os.MkdirTemp(base, "bateria-")
	if err != nil {
		r.Estado = EstadoReprovado
		r.Motivo = "C1: " + err.Error()
		return r
	}
	defer os.RemoveAll(dir)
	motivos, ban := rodar(ctx, h, dir, teto, &r)
	if ban != "" {
		r.Estado, r.Motivo = EstadoReprovado, ban
		return r
	}
	r.Estado = EstadoHomologado
	for _, k := range ordem {
		if !r.Criterios[k] {
			r.Estado, r.Motivo = EstadoReprovado, fmt.Sprintf("%s: %s", k, motivos[k])
			break
		}
	}
	return r
}

func rodar(ctx context.Context, h Harness, dir string, teto time.Duration, r *Resultado) (map[string]string, string) {
	motivos := map[string]string{}
	curto := max(teto/10, 10*time.Millisecond)
	prompts := []string{PromptA, PromptB, PromptC}
	tetos := []time.Duration{teto, teto, curto}
	var saidas [3]Saida
	var erros [3]error
	ban := ""
	for i := range prompts {
		saidas[i], erros[i] = h.Executar(ctx, dir, prompts[i], tetos[i])
		if t := trechoDaChamada(saidas[i], erros[i]); t != "" {
			ban = "banimento: " + t
			medir(r, motivos, dir, curto, saidas[:i], erros[:i])
			return motivos, ban
		}
		r.PesoMB = max(r.PesoMB, saidas[i].PicoMB)
	}
	medir(r, motivos, dir, curto, saidas[:], erros[:])
	return motivos, ban
}

func trechoDaChamada(s Saida, err error) string {
	if t := TrechoBanimento(s.Texto); t != "" {
		return t
	}
	if err != nil {
		return TrechoBanimento(err.Error())
	}
	return ""
}

func medir(r *Resultado, motivos map[string]string, dir string, curto time.Duration, saidas []Saida, erros []error) {
	if len(saidas) > 0 {
		r.Criterios["C1"], motivos["C1"] = avaliaC1(erros[0])
		r.Criterios["C2"], motivos["C2"] = avaliaC2(dir)
		r.Criterios["C3"] = true
		if saidas[0].CustoUSD <= 0 && saidas[0].Tokens <= 0 {
			motivos["C3"] = "nao telemetrizado"
		}
	}
	if len(saidas) > 1 {
		r.Criterios["C5"] = saidas[1].ExitCode != 0 || erros[1] != nil
		motivos["C5"] = "exit zero diante de erro forçado"
	}
	if len(saidas) > 2 {
		r.Criterios["C4"] = errors.Is(erros[2], context.DeadlineExceeded) || saidas[2].Duracao <= curto+50*time.Millisecond
		motivos["C4"] = "chamada passou do teto"
	}
}

func avaliaC1(err error) (bool, string) {
	if err != nil {
		return false, "chamada falhou ou estourou o tempo: " + err.Error()
	}
	return true, ""
}

func avaliaC2(dir string) (bool, string) {
	fora := ForaDoLimite(dir, []string{"permitido.txt"})
	if len(fora) == 0 {
		return true, ""
	}
	if err := Reverter(dir, fora); err != nil {
		return false, "escrita fora do limite sem reversao: " + err.Error()
	}
	return true, ""
}
