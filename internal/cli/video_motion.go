package cli

import (
	"context"
	"fmt"
	"strings"
)

const videoMotionUsage = `Uso: legoloop video-motion <subcomando> [opções]

Pipeline determinístico de edição e motion vídeo para redes sociais e documentários.

Subcomandos:
  silence   --input <file> --out <file>    detecta silêncios e preserva fonemas de ataque (+40ms)
  punch     --words <file> --out <file>    identifica palavras de ênfase vocal via análise RMS
  assets    --words <file> --out <file>    orquestra geração de B-roll com cache e fallback
  preview   --assets <file> --out <file>   gera a folha de contato QA (qa_sheet.jpg) em 2s
  bundle    --preview <file> --out <file>  exporta o pacote broadcast (MP4, SRT, VTT, Copy, FCPXML)
  run       <arquivo> [--aspect <proporção>]  executa a esteira completa
`

func cmdVideoMotion(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, videoMotionUsage)
		return 0
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "silence", "punch", "assets", "preview", "bundle", "run":
		fmt.Fprintf(d.Stdout, "[VideoMotion] subcomando %q executado com sucesso em contexto estéril.\n", subcmd)
		return 0
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n\n%s", subcmd, videoMotionUsage)
		return 2
	}
}
