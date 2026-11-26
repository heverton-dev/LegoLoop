package videomotion_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/videomotion"
)

func TestMeasureAndNormalizeLUFS(t *testing.T) {
	// Saída simulada do ffmpeg -filter:a loudnorm (summary JSON ou key-value)
	rawOutput := `
[Parsed_loudnorm_0 @ 000001] 
Input Integrated:    -18.4 LUFS
Input True Peak:      -3.1 dBFS
Input LRA:             6.2 LU
Output Integrated:   -14.1 LUFS
Output True Peak:     -1.5 dBFS
Output LRA:            5.8 LU
Output Threshold:    -24.3 LUFS
Normalization Type:   Dynamic
`
	metrics, err := videomotion.ParseLoudnormOutput(rawOutput)
	if err != nil {
		t.Fatalf("ParseLoudnormOutput falhou: %v", err)
	}

	if metrics.OutputIntegratedLUFS != -14.1 {
		t.Errorf("Integrated esperado -14.1, obteve %f", metrics.OutputIntegratedLUFS)
	}

	if metrics.OutputTruePeakDBFS != -1.5 {
		t.Errorf("True Peak esperado -1.5, obteve %f", metrics.OutputTruePeakDBFS)
	}

	// Validação de conformidade: target -14.0 com tolerância de 1.0 LUFS
	err = videomotion.ValidateLoudnessConformity(*metrics, -14.0, 1.0)
	if err != nil {
		t.Errorf("ValidateLoudnessConformity falhou: %v", err)
	}
}

func TestMeasureAndNormalizeLUFS_OutOfRange(t *testing.T) {
	metrics := videomotion.LoudnessMetrics{
		OutputIntegratedLUFS: -20.5, // muito baixo (-20.5 vs alvo -14.0)
		OutputTruePeakDBFS:   -1.5,
	}
	err := videomotion.ValidateLoudnessConformity(metrics, -14.0, 1.0)
	if err == nil {
		t.Fatal("esperava erro para áudio fora da conformidade de -14 LUFS")
	}
}

func TestBuildLoudnormFilter(t *testing.T) {
	filter := videomotion.BuildLoudnormFilter(-14.0)
	expected := "loudnorm=I=-14:TP=-1.5:LRA=11"
	if filter != expected {
		t.Errorf("filtro esperado '%s', obteve '%s'", expected, filter)
	}
}

func TestValidateLoudnessConformity_TetoDeTruePeak(t *testing.T) {
	casos := []struct {
		nome       string
		pico       float64
		esperaErro bool
	}{
		{"-1.0 excede o teto", -1.0, true},
		{"-1.2 excede o teto", -1.2, true},
		{"-1.4 excede o teto", -1.4, true},
		{"-1.5 esta no teto", -1.5, false},
		{"-1.6 abaixo do teto", -1.6, false},
		{"-3.0 abaixo do teto", -3.0, false},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			m := videomotion.LoudnessMetrics{OutputIntegratedLUFS: -14.0, OutputTruePeakDBFS: c.pico}
			err := videomotion.ValidateLoudnessConformity(m, -14.0, 0.5)
			if c.esperaErro {
				if !errors.Is(err, videomotion.ErrLoudnessConformityFailed) || !strings.Contains(err.Error(), "-1.5") {
					t.Fatalf("esperava ErrLoudnessConformityFailed citando o teto -1.5: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("nao deveria reprovar: %v", err)
			}
		})
	}
}

func TestTruePeakTetoTemFonteUnica(t *testing.T) {
	if videomotion.TruePeakTetoDBTP != -1.5 {
		t.Fatalf("teto de True Peak deve ser -1.5 dBTP: %v", videomotion.TruePeakTetoDBTP)
	}
	filtro := videomotion.BuildLoudnormFilter(-14.0)
	quer := "TP=" + strconv.FormatFloat(videomotion.TruePeakTetoDBTP, 'f', 1, 64)
	if !strings.Contains(filtro, quer) {
		t.Fatalf("o filtro (%s) e a validacao devem usar o mesmo teto (%s)", filtro, quer)
	}
}

func TestValidateLoudnessConformity_LUFSContinuaValendo(t *testing.T) {
	m := videomotion.LoudnessMetrics{OutputIntegratedLUFS: -14.6, OutputTruePeakDBFS: -2.0}
	if err := videomotion.ValidateLoudnessConformity(m, -14.0, 0.5); !errors.Is(err, videomotion.ErrLoudnessConformityFailed) {
		t.Fatalf("LUFS fora de -14 +- 0,5 deve reprovar mesmo com pico bom: %v", err)
	}
}
