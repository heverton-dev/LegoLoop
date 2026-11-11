package fabricar_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/ste"
)

// captura guarda o plano que o motor mostrou e responde sempre a mesma decisão.
type captura struct {
	planos []fabricar.Plano
	ok     bool
	antes  int
}

func (c *captura) aprovar(r *rig, t *testing.T) fabricar.ApproveFunc {
	return func(p fabricar.Plano, _ func([]string) fabricar.Plano) (bool, error) {
		c.planos = append(c.planos, p)
		c.antes = r.written(t)
		return c.ok, nil
	}
}

func TestCartaoPlanoUnico(t *testing.T) {
	t.Run("o plano lista arquivos e links e nada está gravado", func(t *testing.T) {
		r := newRig(t, goodPieces())
		comHarnesses(t, r, "claude", "gemini")
		c := &captura{ok: true}
		r.deps.Aprovar = c.aprovar(r, t)
		r.run(t)
		p := c.planos[0]
		if len(p.Arquivos) != 10 || len(p.Vinculos) != 2 || len(p.Detectados) != 2 || len(p.Ligados) != 2 {
			t.Fatalf("plano incompleto: %+v", p)
		}
		if c.antes != 0 {
			t.Errorf("nada pode estar gravado ao mostrar o plano: %d", c.antes)
		}
		if len(p.Sobrescreve) != 0 || len(p.Removivel) != 0 {
			t.Errorf("sem Force e sem mídia removível não há pausa extra: %+v", p)
		}
	})

	t.Run("recusar não grava nada e devolve ErrAborted", func(t *testing.T) {
		r := newRig(t, goodPieces())
		comHarnesses(t, r, "claude")
		c := &captura{ok: false}
		r.deps.Aprovar = c.aprovar(r, t)
		if _, err := fabricar.Fabricar(context.Background(), r.req, r.deps); !errors.Is(err, fabricar.ErrAborted) {
			t.Fatalf("esperava ErrAborted: %v", err)
		}
		if n := r.written(t); n != 0 {
			t.Errorf("recusou mas gravou %d", n)
		}
	})

	t.Run("sobrescrever aparece no plano", func(t *testing.T) {
		r := newRig(t, goodPieces())
		r.run(t)
		r.req.Force = true
		c := &captura{ok: true}
		r.deps.Aprovar = c.aprovar(r, t)
		r.run(t)
		if len(c.planos[0].Sobrescreve) != 10 {
			t.Errorf("os 10 arquivos existentes devem ser listados: %v", c.planos[0].Sobrescreve)
		}
	})

	t.Run("mídia removível aparece no plano", func(t *testing.T) {
		usb := skillApta("usb", "/media/usb/s/SKILL.md", 0.95, "construir")
		usb.Perfil.CasosDeUso = []string{"tdd"} // apta como alfa; o Score maior a escolhe
		pcs := append(goodPieces(), usb)
		r := newRig(t, pcs)
		r.req.TrustRemovable = true
		c := &captura{ok: true}
		r.deps.Aprovar = c.aprovar(r, t)
		r.run(t)
		if len(c.planos[0].Removivel) != 1 || c.planos[0].Removivel[0] != "/media/usb/s/SKILL.md" {
			t.Errorf("a peça removível deve ser listada: %v", c.planos[0].Removivel)
		}
	})

	t.Run("personalizar troca os harnesses ligados", func(t *testing.T) {
		r := newRig(t, goodPieces())
		comHarnesses(t, r, "claude", "gemini", "codex")
		r.deps.Aprovar = func(p fabricar.Plano, replan func([]string) fabricar.Plano) (bool, error) {
			novo := replan([]string{"gemini"})
			if len(novo.Vinculos) != 1 || len(novo.Ligados) != 1 || len(novo.Detectados) != 3 {
				t.Errorf("replanejar deve manter os detectados e trocar os ligados: %+v", novo)
			}
			return true, nil
		}
		res := r.run(t)
		if len(res.Artifacts.SkillLinks) != 1 || !strings.Contains(res.Artifacts.SkillLinks[0], ".gemini") {
			t.Fatalf("só o gemini deve ser ligado: %v", res.Artifacts.SkillLinks)
		}
	})

	t.Run("replanejar com lista vazia não liga nada", func(t *testing.T) {
		r := newRig(t, goodPieces())
		comHarnesses(t, r, "claude")
		r.deps.Aprovar = func(p fabricar.Plano, replan func([]string) fabricar.Plano) (bool, error) {
			replan(nil)
			return true, nil
		}
		if res := r.run(t); len(res.Artifacts.SkillLinks) != 0 {
			t.Fatalf("lista vazia = nenhum link: %v", res.Artifacts.SkillLinks)
		}
	})

	t.Run("o cartão do plano segue o perfil STE-PT", func(t *testing.T) {
		r := newRig(t, goodPieces())
		comHarnesses(t, r, "claude", "gemini")
		r.req.Force = true
		c := &captura{ok: true}
		r.deps.Aprovar = c.aprovar(r, t)
		r.run(t)
		ficha := c.planos[0].Cartao()
		if v := ste.LintCartao(ficha); len(v) != 0 {
			t.Errorf("violações STE-PT: %v\n%s", v, ficha.Render())
		}
		if len(ficha.Parametros) == 0 || ficha.Parametros[0].Nome != "harnesses" {
			t.Errorf("o cartão deve oferecer os harnesses para personalizar: %+v", ficha.Parametros)
		}
	})
}
