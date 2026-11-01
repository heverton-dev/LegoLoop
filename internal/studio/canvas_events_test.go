package studio

import (
	"strings"
	"testing"
)

// CA-18.1 a CA-18.6: Verificação rigorosa de fiação do Canvas N8N (Pan, Zoom, Drag e Atalhos).
func TestCanvasEventWiring(t *testing.T) {
	appJS := readStatic(t, "app.js")

	// CA-18.1 & CA-18.2: Pan confiável: detecção canPanCanvas, suporte a deltaX/deltaY e pointer capture
	if !strings.Contains(appJS, "canPanCanvas") {
		t.Errorf("app.js deve conter a função canPanCanvas para controle de pan")
	}
	if !strings.Contains(appJS, "deltaX") || !strings.Contains(appJS, "deltaY") {
		t.Errorf("app.js deve suportar rolagem 2D via deltaX e deltaY")
	}
	if !strings.Contains(appJS, "setPointerCapture") {
		t.Errorf("app.js deve usar setPointerCapture para prevenir perda de ponteiro durante o arrasto")
	}

	// CA-18.3: Zoom contínuo e clamping com zoomBy
	if !strings.Contains(appJS, "zoomBy") {
		t.Errorf("app.js deve conter a função zoomBy")
	}
	if !strings.Contains(appJS, "Math.min(3.0, Math.max(0.2") {
		t.Errorf("app.js deve aplicar clamping estrito de zoom entre 0.2 e 3.0")
	}

	// CA-18.4 & CA-18.5: Drag confiável de nós com listeners no window
	if !strings.Contains(appJS, `window.addEventListener("pointermove"`) ||
		!strings.Contains(appJS, `window.addEventListener("pointerup"`) {
		t.Errorf("app.js deve registrar pointermove e pointerup no window para drag estável")
	}

	// CA-18.6: Atalhos de teclado universais (Delete, Backspace, Ctrl+0 para fitToScreen)
	if !strings.Contains(appJS, "Delete") || !strings.Contains(appJS, "Backspace") {
		t.Errorf("app.js deve tratar teclas Delete e Backspace para remoção")
	}
	if !strings.Contains(appJS, "fitToScreen()") || !strings.Contains(appJS, "Digit0") {
		t.Errorf("app.js deve tratar atalho Ctrl+0 (Digit0) para fitToScreen")
	}
}
