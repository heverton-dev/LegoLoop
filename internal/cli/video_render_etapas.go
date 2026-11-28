package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/motor"
	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

// formatoMotorPadrao e o formato padrao quando nao especificado.
const formatoMotorPadrao = "16:9"

func cmdVideoRenderPrevoo(ctx context.Context, args []string, d Deps) int {
	projeto, ok := parseProjetoRender("prevoo", args, d.Stderr)
	if !ok {
		return 2
	}
	return rodarPrevoo(ctx, projeto, render.DepsPreVoo{}, motor.Dependencias{}, d)
}

// rodarPrevoo valida contratos, voz, ferramentas e pastas; depois confere o motor local.
func rodarPrevoo(ctx context.Context, projeto string, dr render.DepsPreVoo, dm motor.Dependencias, d Deps) int {
	rel, _, err := render.ExecutarPreVoo(ctx, projeto, dr)
	for _, it := range rel.Itens {
		fmt.Fprintf(d.Stdout, "%s: %s %s\n", it.Nivel, it.Nome, it.Detalhe)
	}
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render prevoo: %v\n", err)
		return 1
	}
	if err := motor.PreVoo(dm); err != nil {
		fmt.Fprintf(d.Stderr, "video-render prevoo: %v\n", err)
		return 1
	}
	fmt.Fprintln(d.Stdout, "pre-voo aprovado")
	return 0
}

func cmdVideoRenderVoz(ctx context.Context, args []string, d Deps) int {
	projeto, ok := parseProjetoRender("voz", args, d.Stderr)
	if !ok {
		return 2
	}
	return rodarVoz(ctx, projeto, render.EtapaVozSistema, d)
}

// rodarVoz sintetiza o roteiro, confere com o Whisper e grava audio/ e data/timings.json.
func rodarVoz(ctx context.Context, projeto string, montar func(reprovas int) render.EtapaVoz, d Deps) int {
	c, err := render.CarregarContratos(projeto)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render voz: %v\n", err)
		return 1
	}
	tm, err := montar(c.Config.ReprovasWERMax).Executar(ctx, projeto, c.Marca, c.Script)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render voz: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "voz concluida: %s/%s, %d linhas, WER %.3f\n", tm.Provedor, tm.Voz, len(tm.Linhas), tm.WERGlobal)
	return 0
}

func cmdVideoRenderMix(ctx context.Context, args []string, d Deps) int {
	projeto, ok := parseProjetoRender("mix", args, d.Stderr)
	if !ok {
		return 2
	}
	return rodarMix(ctx, projeto, render.MixDeps{}, d)
}

// rodarMix mixa voz, SFX e trilha em audio/soundtrack.wav com loudness medido.
func rodarMix(ctx context.Context, projeto string, deps render.MixDeps, d Deps) int {
	rel, err := render.ExecutarMix(ctx, projeto, deps)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render mix: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "mix concluido: %.1f LUFS, %.1f dBTP, %d eventos SFX\n", rel.LUFS, rel.TruePeakDBTP, rel.Eventos)
	return 0
}

func cmdVideoRenderMotor(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("video-render motor", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	projeto := fs.String("projeto", "", "diretorio do projeto de video")
	raiz := fs.String("raiz", ".", "raiz do LegoLoop (inventario/video/estilos)")
	motorDir := fs.String("motor-dir", "motor-local", "pasta do motor portado")
	formato := fs.String("formato", "", "formato do video (16:9, 9:16, 1:1)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *projeto == "" {
		fmt.Fprintln(d.Stderr, "video-render motor: --projeto e obrigatorio")
		return 2
	}
	fmtFinal := resolverFormatoMotor(*projeto, *formato)
	ch := motor.ChecadorPython{Script: filepath.Join(*motorDir, "check_timeline.py")}
	mr := motor.RenderMotor{
		Formato: fmtFinal,
		Ex:      motor.ExecutorPython{Script: filepath.Join(*motorDir, "run_render.py"), Args: []string{"--formato", fmtFinal}},
		Sonda:   motor.SondaFFmpeg{},
	}
	return rodarMotor(ctx, *projeto, inventario.Store{Raiz: *raiz}, ch, mr, d)
}

func resolverFormatoMotor(projeto, formato string) string {
	if formato != "" {
		return formato
	}
	if p := motor.FormatoPrimarioProjeto(projeto); p != "" {
		return p
	}
	return formatoMotorPadrao
}

// rodarMotor gera a timeline do motor (CHECK OK) e renderiza o MP4 com a voz.
func rodarMotor(ctx context.Context, projeto string, st inventario.Store, ch motor.Checador, mr render.MotorRender, d Deps) int {
	rel, err := motor.Executar(ctx, projeto, st, ch, mr)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render motor: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "motor concluido: %d cortes, voz %.2f s, %.1f LUFS\n", rel.Cortes, rel.DuracaoVozS, rel.LUFS)
	return 0
}

func cmdVideoRenderMarca(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("video-render marca", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	proj := fs.String("projeto", "", "diretorio do projeto de video")
	motorDir := fs.String("motor-dir", "motor-local", "pasta do motor portado")
	versao := fs.Int("versao", 1, "versao do render (gera renders/vN.mp4)")
	if err := fs.Parse(args); err != nil || *proj == "" {
		if *proj == "" {
			fmt.Fprintln(d.Stderr, "video-render marca: --projeto e obrigatorio")
		}
		return 2
	}
	entrada := filepath.Join(*proj, "render.mp4")
	saida := filepath.Join(*proj, "renders", fmt.Sprintf("v%d.mp4", *versao))
	renderPy := filepath.Join(*motorDir, "run_render.py")
	rel, err := motor.AplicarMarca(ctx, *proj, entrada, saida, renderPy, motor.ExecutorSistemaComando{})
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render marca: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "marca aplicada: %s (%s %dx%d)\n", rel.Saida, rel.Formato, rel.Largura, rel.Altura)
	return 0
}

func cmdVideoRenderCortes(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("video-render cortes", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	proj := fs.String("projeto", "", "diretorio do projeto de video")
	motorDir := fs.String("motor-dir", "motor-local", "pasta do motor portado")
	minS := fs.Float64("min", 15.0, "duracao minima em segundos")
	maxS := fs.Float64("max", 59.0, "duracao maxima em segundos")
	if err := fs.Parse(args); err != nil || *proj == "" {
		if *proj == "" {
			fmt.Fprintln(d.Stderr, "video-render cortes: --projeto e obrigatorio")
		}
		return 2
	}
	return executarCortes(ctx, *proj, *motorDir, *minS, *maxS, d)
}

func carregarEspecificacaoCortes(proj string) (motor.EspecificacaoCortes, error) {
	var espec motor.EspecificacaoCortes
	caminho := filepath.Join(proj, "cortes.json")
	if dados, err := os.ReadFile(caminho); err == nil {
		if err := json.Unmarshal(dados, &espec); err == nil && len(espec.Cortes) > 0 {
			return espec, nil
		}
	}
	var tl motor.Timeline
	tlDados, err := os.ReadFile(filepath.Join(proj, "timeline.json"))
	if err != nil {
		return espec, err
	}
	if err := json.Unmarshal(tlDados, &tl); err != nil {
		return espec, err
	}
	sug := motor.SugerirCortes(tl)
	return motor.EspecificacaoCortes{Cortes: sug, CartaoFinal: "Caso completo no canal", DuracaoCartao: 3.2}, nil
}

func executarCortes(ctx context.Context, proj, motorDir string, minS, maxS float64, d Deps) int {
	espec, err := carregarEspecificacaoCortes(proj)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render cortes: %v\n", err)
		return 1
	}
	cmd := motor.ExecutorSistemaComando{}
	resumos, err := motor.PrepararCortes(ctx, proj, espec, cmd)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render cortes: %v\n", err)
		return 1
	}
	renderPy := filepath.Join(motorDir, "run_render.py")
	res, err := motor.RenderizarCortes(ctx, proj, resumos, renderPy, cmd)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render cortes: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "cortes concluidos: %d shorts 9:16 gerados em %s\n", len(res), filepath.Join(proj, "renders", "shorts"))
	return 0
}
