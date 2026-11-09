package videomotion

import (
	"errors"
)

var (
	// ErrUnsupportedAspectRatio indica formato de tela não cadastrado.
	ErrUnsupportedAspectRatio = errors.New("videomotion: proporção de aspecto não suportada")
)

// RenderJobSpec parametriza a codificação de vídeo.
type RenderJobSpec struct {
	InputVideo  string
	AudioTrack  string
	AspectRatio AspectRatio
	TargetLUFS  float64
	OutputPath  string
}

// RenderProfile define resoluções e comandos determinísticos de saída.
type RenderProfile struct {
	Width       int
	Height      int
	AspectRatio AspectRatio
	FilterGraph string
	Spec        RenderJobSpec
}

// BuildRenderProfile valida a especificação e calcula a geometria de saída.
func BuildRenderProfile(spec RenderJobSpec) (*RenderProfile, error) {
	w, h, filter, err := resolveAspectGeometry(spec.AspectRatio)
	if err != nil {
		return nil, err
	}
	if spec.TargetLUFS == 0 {
		spec.TargetLUFS = -14.0
	}
	return &RenderProfile{
		Width:       w,
		Height:      h,
		AspectRatio: spec.AspectRatio,
		FilterGraph: filter,
		Spec:        spec,
	}, nil
}

func resolveAspectGeometry(aspect AspectRatio) (int, int, string, error) {
	switch aspect {
	case AspectVertical:
		return 1080, 1920, "scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2", nil
	case AspectHorizontal:
		return 1920, 1080, "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2", nil
	case AspectSquare:
		return 1080, 1080, "scale=1080:1080:force_original_aspect_ratio=decrease,pad=1080:1080:(ow-iw)/2:(oh-ih)/2", nil
	default:
		return 0, 0, "", ErrUnsupportedAspectRatio
	}
}

// BuildFFmpegCommand constrói a lista de argumentos para encode broadcast H.264 / AAC.
func (p *RenderProfile) BuildFFmpegCommand() []string {
	loudnormFilter := BuildLoudnormFilter(p.Spec.TargetLUFS)
	args := []string{
		"-y",
		"-i", p.Spec.InputVideo,
	}
	if p.Spec.AudioTrack != "" {
		args = append(args, "-i", p.Spec.AudioTrack)
	}
	args = append(args,
		"-vf", p.FilterGraph,
		"-af", loudnormFilter,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "18",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "320k",
		p.Spec.OutputPath,
	)
	return args
}
