package videomotion

// AspectRatio representa a proporção de tela do vídeo final.
type AspectRatio string

const (
	// AspectVertical proporção vertical 9:16 (1080x1920).
	AspectVertical AspectRatio = "9:16"
	// AspectHorizontal proporção widescreen 16:9 (1920x1080).
	AspectHorizontal AspectRatio = "16:9"
	// AspectSquare proporção quadrada 1:1 (1080x1080).
	AspectSquare AspectRatio = "1:1"
)
