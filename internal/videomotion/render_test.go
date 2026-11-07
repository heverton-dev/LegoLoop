package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestRenderMultiAspect(t *testing.T) {
	cases := []struct {
		aspect         videomotion.AspectRatio
		expectedW      int
		expectedH      int
		expectedFilter string
	}{
		{
			aspect:         videomotion.AspectVertical,
			expectedW:      1080,
			expectedH:      1920,
			expectedFilter: "scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2",
		},
		{
			aspect:         videomotion.AspectHorizontal,
			expectedW:      1920,
			expectedH:      1080,
			expectedFilter: "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2",
		},
		{
			aspect:         videomotion.AspectSquare,
			expectedW:      1080,
			expectedH:      1080,
			expectedFilter: "scale=1080:1080:force_original_aspect_ratio=decrease,pad=1080:1080:(ow-iw)/2:(oh-ih)/2",
		},
	}

	for _, tc := range cases {
		spec := videomotion.RenderJobSpec{
			InputVideo:  "input.mp4",
			AudioTrack:  "audio.wav",
			AspectRatio: tc.aspect,
			TargetLUFS:  -14.0,
			OutputPath:  "output.mp4",
		}

		profile, err := videomotion.BuildRenderProfile(spec)
		if err != nil {
			t.Fatalf("BuildRenderProfile falhou para %s: %v", tc.aspect, err)
		}

		if profile.Width != tc.expectedW || profile.Height != tc.expectedH {
			t.Errorf("dimensão para %s: esperada %dx%d, obteve %dx%d", tc.aspect, tc.expectedW, tc.expectedH, profile.Width, profile.Height)
		}

		cmdArgs := profile.BuildFFmpegCommand()
		cmdStr := joinArgs(cmdArgs)

		if !containsStr(cmdStr, tc.expectedFilter) {
			t.Errorf("filtro de escala para %s ausente nos argumentos:\n%s", tc.aspect, cmdStr)
		}

		if !containsStr(cmdStr, "loudnorm=I=-14:TP=-1.5:LRA=11") {
			t.Errorf("filtro loudnorm ausente nos argumentos para %s", tc.aspect)
		}

		if !containsStr(cmdStr, "-c:v libx264") || !containsStr(cmdStr, "-pix_fmt yuv420p") {
			t.Errorf("codecs obrigatórios ausentes para %s", tc.aspect)
		}
	}
}

func TestRenderMultiAspect_InvalidAspect(t *testing.T) {
	spec := videomotion.RenderJobSpec{
		InputVideo:  "in.mp4",
		AspectRatio: "21:9", // aspect não homologado
		OutputPath:  "out.mp4",
	}
	_, err := videomotion.BuildRenderProfile(spec)
	if err == nil {
		t.Fatal("esperava erro para aspect ratio não suportado")
	}
}
