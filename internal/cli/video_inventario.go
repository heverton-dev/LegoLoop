package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

const usoVideoInventario = `Uso: legoloop video-inventario <subcomando> [opções]

Subcomandos:
  validar   --estilo ID --aprovacao ARQ [--raiz DIR] [--json]  muda o estilo de rascunho para validado com a aprovação humana gravada em arquivo
  verificar --estilo ID [--raiz DIR] [--json]                  confere o estilo e o hash dos arquivos pesados dos itens
`

func cmdVideoInventario(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprint(d.Stdout, usoVideoInventario)
		return 0
	}
	switch args[0] {
	case "validar":
		return cmdVideoInventarioValidar(args[1:], d)
	case "verificar":
		return cmdVideoInventarioVerificar(args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s", args[0], usoVideoInventario)
		return 2
	}
}

type flagsVideoInventario struct {
	raiz, estilo, aprovacao string
	emJSON                  bool
}

func parseVideoInventario(nome string, args []string, d Deps, comAprovacao bool) (flagsVideoInventario, bool) {
	var f flagsVideoInventario
	fs := flag.NewFlagSet("video-inventario "+nome, flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	fs.StringVar(&f.raiz, "raiz", ".", "raiz do repositório que contém inventario/video")
	fs.StringVar(&f.estilo, "estilo", "", "id do estilo (kebab-case)")
	fs.BoolVar(&f.emJSON, "json", false, "saída em JSON")
	if comAprovacao {
		fs.StringVar(&f.aprovacao, "aprovacao", "", "arquivo JSON da aprovação humana")
	}
	if err := fs.Parse(args); err != nil {
		return f, false
	}
	if f.estilo == "" {
		fmt.Fprintf(d.Stderr, "video-inventario %s: --estilo é obrigatório\n", nome)
		return f, false
	}
	return f, true
}

func cmdVideoInventarioValidar(args []string, d Deps) int {
	f, ok := parseVideoInventario("validar", args, d, true)
	if !ok {
		return 2
	}
	s := inventario.Store{Raiz: f.raiz}
	e, err := s.ValidarEstilo(f.estilo, f.aprovacao, time.Now())
	if err != nil {
		return responderVideoInventario(d, f.emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	return responderVideoInventario(d, f.emJSON, map[string]any{"ok": true, "estilo": e.ID, "status": e.Status}, 0)
}

func cmdVideoInventarioVerificar(args []string, d Deps) int {
	f, ok := parseVideoInventario("verificar", args, d, false)
	if !ok {
		return 2
	}
	falhas, err := inventario.Store{Raiz: f.raiz}.VerificarEstilo(f.estilo)
	if err != nil {
		return responderVideoInventario(d, f.emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	if len(falhas) > 0 {
		erros := make([]string, 0, len(falhas))
		for _, e := range falhas {
			erros = append(erros, e.Error())
		}
		return responderVideoInventario(d, f.emJSON, map[string]any{"ok": false, "erros": erros}, 1)
	}
	return responderVideoInventario(d, f.emJSON, map[string]any{"ok": true, "estilo": f.estilo}, 0)
}

func responderVideoInventario(d Deps, emJSON bool, payload map[string]any, code int) int {
	if !emJSON {
		fmt.Fprintf(d.Stdout, "%v\n", payload)
		return code
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-inventario: codificando saida: %v\n", err)
		return 1
	}
	fmt.Fprintln(d.Stdout, string(b))
	return code
}
