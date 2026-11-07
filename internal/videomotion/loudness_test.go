package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
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
