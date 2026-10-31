package dag

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"legoloop/internal/gate"
)

// majorVersion lê o componente maior de version; vazio equivale a 1 (contrato legado).
func majorVersion(v string) (int, error) {
	if v == "" {
		return 1, nil
	}
	head, _, _ := strings.Cut(v, ".")
	n, err := strconv.Atoi(head)
	if err != nil || n < 1 || n > 2 {
		return 0, fmt.Errorf("dag: versão %q não suportada (use 1.x ou 2.x)", v)
	}
	return n, nil
}

// checkTrigger aplica a regra de versão: v2 exige trigger; v1 sem trigger vira manual com aviso.
func (d *DAG) checkTrigger() error {
	major, err := majorVersion(d.Version)
	if err != nil {
		return err
	}
	if d.Trigger == nil {
		if major >= 2 {
			return errors.New("dag: versão 2 exige o campo trigger")
		}
		d.warnings = append(d.warnings, "DAG v1 sem trigger: assumido gatilho manual")
		return nil
	}
	if err := validateTrigger(*d.Trigger); err != nil {
		return err
	}
	if d.Trigger.Type == TriggerFileChange {
		return d.checkWatchedPath(d.Trigger.Path)
	}
	return nil
}

// checkWatchedPath impede o loop infinito: nenhuma peça pode gravar dentro do que o trigger vigia.
func (d *DAG) checkWatchedPath(watched string) error {
	for _, n := range d.Nodes {
		if overlaps(watched, n.Output.Artifact) {
			return fmt.Errorf("dag: nó %q grava %q dentro do caminho vigiado %q (loop infinito)", n.ID, n.Output.Artifact, watched)
		}
	}
	return nil
}

func validateTrigger(t Trigger) error {
	switch t.Type {
	case TriggerManual:
		return t.stray("")
	case TriggerCron:
		if err := t.stray("cron"); err != nil {
			return err
		}
		return validateCron(t.Cron)
	case TriggerFileChange:
		if err := t.stray("path"); err != nil {
			return err
		}
		return validatePath("trigger file_change: path", t.Path)
	case TriggerLoopEnd:
		if err := t.stray("loop"); err != nil {
			return err
		}
		if !idRe.MatchString(t.Loop) {
			return fmt.Errorf("trigger loop_end exige loop válido, recebeu %q", t.Loop)
		}
		return nil
	case "webhook", "git_event":
		return fmt.Errorf("trigger %q é da fase 2 e ainda não está disponível", t.Type)
	}
	return fmt.Errorf("trigger desconhecido %q", t.Type)
}

// stray recusa campos que não pertencem ao tipo do trigger.
func (t Trigger) stray(allowed string) error {
	fields := []struct{ name, value string }{{"cron", t.Cron}, {"path", t.Path}, {"loop", t.Loop}}
	for _, f := range fields {
		if f.value != "" && f.name != allowed {
			return fmt.Errorf("trigger %s não aceita o campo %q", t.Type, f.name)
		}
	}
	return nil
}

func validateExecution(e *Execution) error {
	if e == nil {
		return nil
	}
	if e.Mode != ExecInline && e.Mode != ExecIsolated {
		return fmt.Errorf("execution.mode inválido %q (use inline ou isolated)", e.Mode)
	}
	return nil
}

// validateSlots aplica aos slots as mesmas regras de segurança dos nós (SEC-01/03).
func validateSlots(s *Slots) error {
	if s == nil {
		return nil
	}
	if len(s.Gates) > MaxSlotEntries {
		return fmt.Errorf("slot gates excede %d entradas", MaxSlotEntries)
	}
	for _, expr := range s.Gates {
		if _, err := gate.Compile(expr); err != nil {
			return fmt.Errorf("slot gates: %w", err)
		}
	}
	lists := []struct {
		name string
		cmds [][]string
	}{{"scripts", s.Scripts}, {"hooks.pre", s.Hooks.Pre}, {"hooks.post", s.Hooks.Post}}
	for _, l := range lists {
		if err := validateSlotCommands(l.name, l.cmds); err != nil {
			return err
		}
	}
	return nil
}

func validateSlotCommands(name string, cmds [][]string) error {
	if len(cmds) > MaxSlotEntries {
		return fmt.Errorf("slot %s excede %d entradas", name, MaxSlotEntries)
	}
	for _, c := range cmds {
		if err := validateCommand(c); err != nil {
			return fmt.Errorf("slot %s: %w", name, err)
		}
	}
	return nil
}

// checkExecution proíbe duas peças inline na mesma camada paralela com artefatos sobrepostos.
func (d *DAG) checkExecution() error {
	for _, layer := range d.layers {
		if err := d.checkLayerOverlap(layer); err != nil {
			return err
		}
	}
	return nil
}

func (d *DAG) checkLayerOverlap(layer []string) error {
	for i, a := range layer {
		for _, b := range layer[i+1:] {
			na, nb := d.Node(a), d.Node(b)
			if na.EffectiveMode() != ExecInline || nb.EffectiveMode() != ExecInline {
				continue
			}
			if overlaps(na.Output.Artifact, nb.Output.Artifact) {
				return fmt.Errorf("dag: nós %q e %q rodam em paralelo no modo inline com artefatos sobrepostos (%q e %q); use execution.mode=isolated em um deles",
					a, b, na.Output.Artifact, nb.Output.Artifact)
			}
		}
	}
	return nil
}

func normPath(p string) string {
	return strings.ToLower(filepath.ToSlash(filepath.Clean(filepath.FromSlash(p))))
}

// overlaps diz se um caminho é igual ou está dentro do outro (a raiz "." contém tudo).
func overlaps(a, b string) bool {
	na, nb := normPath(a), normPath(b)
	if na == "." || nb == "." || na == nb {
		return true
	}
	return strings.HasPrefix(na, nb+"/") || strings.HasPrefix(nb, na+"/")
}
