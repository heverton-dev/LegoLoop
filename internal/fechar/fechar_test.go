package fechar_test

import (
	"testing"

	"legoloop/internal/fechar"
)

func TestValidateAndNormalizeSlug(t *testing.T) {
	tests := []struct {
		name      string
		req       fechar.Request
		wantSlug  string
		wantCanon string
		wantErr   bool
	}{
		{
			name:      "slug simples sem prefixo",
			req:       fechar.Request{Slug: "video-motion"},
			wantSlug:  "video-motion",
			wantCanon: "legoloop-video-motion",
			wantErr:   false,
		},
		{
			name:      "slug já com prefixo legoloop-",
			req:       fechar.Request{Slug: "legoloop-video-rea"},
			wantSlug:  "video-rea",
			wantCanon: "legoloop-video-rea",
			wantErr:   false,
		},
		{
			name:      "slug com espacos ao redor",
			req:       fechar.Request{Slug: "  video-motion  "},
			wantSlug:  "video-motion",
			wantCanon: "legoloop-video-motion",
			wantErr:   false,
		},
		{
			name:    "slug vazio deve falhar",
			req:     fechar.Request{Slug: ""},
			wantErr: true,
		},
		{
			name:    "slug com espaco interno deve falhar",
			req:     fechar.Request{Slug: "video motion"},
			wantErr: true,
		},
		{
			name:    "slug com path traversal deve falhar",
			req:     fechar.Request{Slug: "../malicioso"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, canon, err := fechar.NormalizeSlug(tt.req.Slug)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeSlug(%q) err = %v, wantErr = %v", tt.req.Slug, err, tt.wantErr)
			}
			if !tt.wantErr {
				if slug != tt.wantSlug {
					t.Errorf("slug = %q, want %q", slug, tt.wantSlug)
				}
				if canon != tt.wantCanon {
					t.Errorf("canon = %q, want %q", canon, tt.wantCanon)
				}
				if err := tt.req.Validate(); err != nil {
					t.Errorf("req.Validate() falhou inesperadamente: %v", err)
				}
			}
		})
	}
}

func TestExecutarOitoEtapas(t *testing.T) {
	tmp := t.TempDir()
	deps := fechar.Deps{
		RepoDir:   tmp,
		VaultDir:  tmp,
		Home:      tmp,
		SkipBuild: true,
	}

	req := fechar.Request{
		Slug:      "video-motion",
		DryRun:    false,
		SkipBuild: true,
	}

	res, err := fechar.Executar(nil, req, deps)
	if err != nil {
		t.Fatalf("Executar falhou: %v", err)
	}
	if !res.Sucesso {
		t.Fatalf("esperava sucesso no fechamento, obteve false")
	}
	if len(res.Etapas) != 8 {
		t.Fatalf("esperava 8 etapas executadas, obteve %d", len(res.Etapas))
	}
}
