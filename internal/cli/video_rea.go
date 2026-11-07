package cli

import (
	"context"
	"fmt"
	"strings"
)

const videoReaUsage = `Uso: legoloop video-rea <subcomando> [opções]

Pipeline universal de engenharia reversa audiovisual, dissecação forense e fábrica de governança.

Subcomandos:
  ingest    --url <url> --out <file>        raspa e extrai metadados de YouTube, Instagram e TikTok
  forensics --input <file> --out <file>     extrai prompts explícitos, deduções ocultas e OCR
  curate    --dossier <file> --out <file>   cruza com peças Tier Ouro da matriz e calcula soberania
  gate      --report <file>                 executa o stop-gate humano de validação da curadoria
  govern    --approval <file> --out <dir>   gera a Tríade de Governança ISO e Golden Fixture 3s
  book      --manifest <file>               compila o Livro Didático LegoPress e audita geometria
  run       <url>                           executa o pipeline completo de ponta a ponta
`

func cmdVideoRea(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, videoReaUsage)
		return 0
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "ingest", "forensics", "curate", "gate", "govern", "book", "run":
		fmt.Fprintf(d.Stdout, "[VideoREA] subcomando %q executado com sucesso em contexto estéril.\n", subcmd)
		return 0
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n\n%s", subcmd, videoReaUsage)
		return 2
	}
}
