package board

import (
	"context"
	"fmt"
)

// runLevelNotificado avisa o início de cada nó, roda o nível e avisa o desfecho de cada nó do nível.
func (sb *SwarmBoard) runLevelNotificado(ctx context.Context, li int, lvl []RunNode, opt RunOptions) (*RunLevelResult, []string, error) {
	if opt.OnNodeStart != nil {
		for _, n := range lvl {
			if n.NodeID != `` {
				opt.OnNodeStart(n.NodeID)
			}
		}
	}
	res, names, err := sb.runLevel(ctx, li, lvl, opt)
	if opt.OnNodeEnd != nil {
		notificarFimDoNivel(lvl, res, err, opt.OnNodeEnd)
	}
	return res, names, err
}

// notificarFimDoNivel reporta cada nó do nível; nó sem resultado de worker conta como reprovado quando o nível falhou.
func notificarFimDoNivel(lvl []RunNode, res *RunLevelResult, errNivel error, fim func(string, bool, int, string)) {
	vistos := map[string]bool{}
	ids := map[string]string{}
	for _, n := range lvl {
		ids[n.Name] = n.NodeID
	}
	if res != nil {
		for _, w := range res.Workers {
			vistos[w.Name] = true
			passou := w.ExitCode == 0 && (w.GateOK == nil || *w.GateOK)
			if id := ids[w.Name]; id != `` {
				fim(id, passou, w.ExitCode, fmt.Sprintf(`worker %s: exit %d, status %s`, w.Name, w.ExitCode, w.Status))
			}
		}
	}
	if errNivel == nil {
		return
	}
	for _, n := range lvl {
		if !vistos[n.Name] && n.NodeID != `` {
			fim(n.NodeID, false, -1, errNivel.Error())
		}
	}
}
