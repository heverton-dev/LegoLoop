package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/gate"
)

// executeSlotCommands executa uma lista de comandos (argv) em sequência no workdir.
func executeSlotCommands(ctx context.Context, cmds [][]string, workdir string) error {
	for _, argv := range cmds {
		if len(argv) == 0 {
			continue
		}
		cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(cmdCtx, argv[0], argv[1:]...)
		cmd.Dir = workdir
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "SYSTEMROOT=" + os.Getenv("SYSTEMROOT")}

		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("comando do slot falhou %v: %s (%w)", argv, string(out), err)
		}
	}
	return nil
}

// evalSlotGates avalia todas as expressões em n.Slots.Gates no ambiente do nó.
// Todos os gates de slot precisam passar para o nó ter sucesso.
func evalSlotGates(slots *dag.Slots, root *os.Root, exitCode int) (bool, string, error) {
	if slots == nil || len(slots.Gates) == 0 {
		return true, "", nil
	}

	for _, expr := range slots.Gates {
		g, err := gate.Compile(expr)
		if err != nil {
			return false, "", fmt.Errorf("gate de slot inválido %q: %w", expr, err)
		}

		ok, err := g.Eval(gate.Env{Root: root, ExitCode: exitCode})
		if err != nil {
			return false, "", err
		}
		if !ok {
			return false, "gate de slot não satisfeito: " + expr, nil
		}
	}

	return true, "", nil
}
