package fabricar

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// planLinks monta a prévia dos vínculos: com --eager, todos os harnesses detectados; sem ele, nenhum (a skill fica no Cofre). --no-link zera. Com --link-harness a
// lista é restrita aos pedidos; pedido sem harness detectado vira aviso. Não toca o disco.
func planLinks(req Request, hs []processo.Harness, skillPath, name string) ([]processo.Vinculo, []string) {
	if req.NoLink || skillPath == "" {
		return nil, nil
	}
	// lazy por padrão (prova de 20/11: cada skill vinculada custa contexto a todo turno do harness).
	// Pedir --eager ou --link-harness é pedir o vínculo.
	if !(processo.Politica{Eager: req.Eager || len(req.LinkHarnesses) > 0}).Vincula(name) {
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
