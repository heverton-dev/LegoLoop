package fabricar

import (
	"fmt"
	"path/filepath"
	"strings"

	"legoloop/internal/processo"
)

// planLinks monta a prévia dos vínculos: todos os harnesses detectados, salvo --no-link. Com --link-harness a
// lista é restrita aos pedidos; pedido sem harness detectado vira aviso. Não toca o disco.
func planLinks(req Request, hs []processo.Harness, skillPath, name string) ([]processo.Vinculo, []string) {
	if req.NoLink || skillPath == "" {
		return nil, nil
	}
	var warns []string
	if len(req.LinkHarnesses) > 0 {
		hs, warns = restrictHarnesses(hs, req.LinkHarnesses)
	}
	return processo.PlanejarVinculos(hs, name, filepath.Dir(skillPath)), warns
}

// restrictHarnesses mantém só os harnesses pedidos e avisa dos pedidos que não foram detectados.
func restrictHarnesses(detected []processo.Harness, wanted []string) ([]processo.Harness, []string) {
	var out []processo.Harness
	var warns []string
	for _, w := range wanted {
		found := false
		for _, h := range detected {
			if strings.EqualFold(h.Nome, strings.TrimSpace(w)) {
				out, found = append(out, h), true
			}
		}
		if !found {
			warns = append(warns, fmt.Sprintf("harness %q não detectado: nenhum vínculo criado para ele", w))
		}
	}
	return out, warns
}

// linkDests devolve os destinos da prévia, para o usuário ver o que será ligado antes de gravar.
func linkDests(vs []processo.Vinculo) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = v.Destino
	}
	return out
}
