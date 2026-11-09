package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestDetectSilenceAndCut(t *testing.T) {
	// Amostra simulada de energia de áudio em intervalos de 50ms (0 a 2000ms = 2s)
	// Fala: 100ms a 600ms (alta energia), Silêncio: 600ms a 1400ms (800ms de pausa), Fala: 1400ms a 1900ms
	frames := make([]videomotion.AudioEnergyFrame, 40)
	for i := range frames {
		timeMs := int64(i * 50)
		energy := 0.01 // ruído de fundo
		if (timeMs >= 150 && timeMs <= 600) || (timeMs >= 1450 && timeMs <= 1900) {
			energy = 0.8 // fala
		}
		frames[i] = videomotion.AudioEnergyFrame{
			TimestampMs: timeMs,
			Energy:      energy,
		}
	}

	cfg := videomotion.SilenceDetectorConfig{
		SilenceThreshold: 0.1,
		MinPauseMs:       400,
		PreAttackMs:      40,
	}

	segments, err := videomotion.DetectSilenceAndCut(frames, cfg)
	if err != nil {
		t.Fatalf("DetectSilenceAndCut falhou: %v", err)
	}

	if len(segments) != 2 {
		t.Fatalf("esperava 2 segmentos de fala, obteve %d", len(segments))
	}

	// Primeiro segmento: fala começa em 150ms -> com pre-attack de 40ms, deve começar em 110ms
	if segments[0].StartMs != 110 {
		t.Errorf("segmento 0: StartMs esperado 110, obteve %d", segments[0].StartMs)
	}
	if segments[0].EndMs < 600 {
		t.Errorf("segmento 0: EndMs esperado >= 600, obteve %d", segments[0].EndMs)
	}

	// Segundo segmento: fala começa em 1450ms -> com pre-attack de 40ms, deve começar em 1410ms
	if segments[1].StartMs != 1410 {
		t.Errorf("segmento 1: StartMs esperado 1410, obteve %d", segments[1].StartMs)
	}
}

func TestDetectSilenceAndCut_EmptyFrames(t *testing.T) {
	cfg := videomotion.SilenceDetectorConfig{
		SilenceThreshold: 0.1,
		MinPauseMs:       400,
		PreAttackMs:      40,
	}
	_, err := videomotion.DetectSilenceAndCut(nil, cfg)
	if err == nil {
		t.Fatal("esperava erro para frames vazios, obteve nil")
	}
}
