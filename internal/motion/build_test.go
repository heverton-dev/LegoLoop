package motion

import (
	"strings"
	"testing"
)

func TestBuildHTML(t *testing.T) {
	tests := []struct {
		name       string
		repos      []string
		prompt     string
		wantCanvas bool
		wantThree  bool
		wantErr    bool
	}{
		{
			name:       "threejs generates canvas and three import",
			repos:      []string{"threejs"},
			prompt:     "grafico de barras 3D de receita anual",
			wantCanvas: true,
			wantThree:  true,
			wantErr:    false,
		},
		{
			name:       "pixelfx generates canvas",
			repos:      []string{"pixelfx"},
			prompt:     "tempestade pixel art",
			wantCanvas: true,
			wantThree:  false,
			wantErr:    false,
		},
		{
			name:    "empty repos returns error",
			repos:   []string{},
			prompt:  "test",
			wantErr: true,
		},
		{
			name:    "nil repos returns error",
			repos:   nil,
			prompt:  "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := BuildHTML(tt.repos, tt.prompt, "1920x1080")
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildHTML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err != ErrNoRepoSelected {
					t.Errorf("BuildHTML() error = %v, want ErrNoRepoSelected", err)
				}
				return
			}
			if !strings.Contains(html, "<canvas") {
				t.Error("BuildHTML() output missing <canvas>")
			}
			if tt.wantThree && !strings.Contains(html, "three") {
				t.Error("BuildHTML() output missing three.js reference")
			}
			if len(html) > 50*1024 {
				t.Errorf("BuildHTML() output is %d bytes, want < 50KB", len(html))
			}
			if !strings.Contains(html, tt.prompt) {
				t.Error("BuildHTML() output missing prompt text")
			}
		})
	}
}
