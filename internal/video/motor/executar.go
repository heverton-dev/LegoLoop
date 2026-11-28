package motor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

// ErrEstiloAusente indica projeto.json sem estilo_visual.
var ErrEstiloAusente = errors.New("motor-executar: projeto.json sem estilo_visual")

// arquivoSoundtrack e a trilha que o mix grava e o motor usa como audio do render.
const arquivoSoundtrack = "audio/soundtrack.wav"

// lerJSONProjeto le um JSON do projeto e nomeia o arquivo no erro.
func lerJSONProjeto(projDir, rel string, alvo any) error {
	dados, err := os.ReadFile(filepath.Join(projDir, filepath.FromSlash(rel)))
	if err != nil {
		return fmt.Errorf("motor-executar: %s: %w", rel, err)
	}
	if err := json.Unmarshal(dados, alvo); err != nil {
		return fmt.Errorf("motor-executar: %s invalido: %w", rel, err)
	}
	return nil
}

// CarregarEntrada le os contratos do projeto e o inventario do estilo e monta a Entrada da ponte.
func CarregarEntrada(projDir string, st inventario.Store) (Entrada, error) {
	var proj struct {
		EstiloVisual string `json:"estilo_visual"`
	}
	var sb ativos.Storyboard
	var el ElencoPonte
	var tm render.Timings
	for _, c := range []struct {
		rel  string
		alvo any
	}{{"projeto.json", &proj}, {"cenas.json", &sb}, {"elenco.json", &el}, {"data/timings.json", &tm}} {
		if err := lerJSONProjeto(projDir, c.rel, c.alvo); err != nil {
			return Entrada{}, err
		}
	}
	if proj.EstiloVisual == "" {
		return Entrada{}, ErrEstiloAusente
	}
	est, err := st.LerEstilo(proj.EstiloVisual)
	if err != nil {
		return Entrada{}, err
	}
	fundos, err := st.Buscar(proj.EstiloVisual, "cenario", nil)
	if err != nil {
		return Entrada{}, err
	}
	voz, err := render.DuracaoVoz(projDir)
	if err != nil {
		return Entrada{}, err
	}
	return Entrada{Cenas: sb, Elenco: el, Timings: tm, Fundos: fundos, Estilo: est, DuracaoVozS: voz, FPS: fpsPadrao, Audio: arquivoSoundtrack}, nil
}

// Executar grava timeline.json do motor (exige CHECK OK) e so entao roda o render do motor.
func Executar(ctx context.Context, projDir string, st inventario.Store, ch Checador, mr render.MotorRender) (render.RelatorioMotor, error) {
	e, err := CarregarEntrada(projDir, st)
	if err != nil {
		return render.RelatorioMotor{}, err
	}
	if err := Gerar(e, projDir, ch); err != nil {
		return render.RelatorioMotor{}, err
	}
	return render.ExecutarRenderMotor(ctx, projDir, mr)
}

// RenderMotor liga Renderizar ao contrato render.MotorRender.
type RenderMotor struct {
	Formato string
	Ex      Executor
	Sonda   Sonda
}

// Renderizar roda o motor e confere o MP4 contra a voz.
func (r RenderMotor) Renderizar(_ context.Context, projDir string, vozS float64, fps int) error {
	formato := r.Formato
	if formato == "" {
		formato = FormatoPrimarioProjeto(projDir)
	}
	_, err := Renderizar(PedidoRender{Dir: projDir, Formato: formato, DuracaoVozS: vozS, FPS: fps}, r.Ex, r.Sonda)
	return err
}

// FormatoPrimarioProjeto devolve o primeiro formato de projeto.json ou vazio.
func FormatoPrimarioProjeto(projDir string) string {
	var proj struct {
		Formatos []string `json:"formatos"`
		Formato  string   `json:"formato"`
	}
	if err := lerJSONProjeto(projDir, "projeto.json", &proj); err != nil {
		return ""
	}
	if len(proj.Formatos) > 0 && proj.Formatos[0] != "" {
		return proj.Formatos[0]
	}
	return proj.Formato
}
