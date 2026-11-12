// Package dag define o contrato do arquivo DAG (SPEC.md §3), seu parser estrito
// e a validação estrutural (IDs, dependências, ciclos, handoffs e gates).
package dag

import "legoloop/internal/gate"

// Mode é a modalidade de roteamento do nó (SPEC_ROTEAMENTO_CANONIZACAO §1).
type Mode string

const (
	ModeSuggested Mode = "suggested"
	ModeManual    Mode = "manual"
	ModeAuto      Mode = "auto"
)

// NodeType define como a Orchestrix executa o nó.
type NodeType string

const (
	// TypeSkill despacha um prompt de skill para um modelo (Orchestrix isolada).
	TypeSkill NodeType = "skill"
	// TypeScript executa um comando local sem shell.
	TypeScript NodeType = "script"
	// Nó subprocesso no dialeto antigo: referencia outro Processo do acervo (Ref = nome do
	// Processo) e roda isolado e estéril. É a representação interna depois do Parse.
	TypeSubloop NodeType = "subloop"
	// TypeSubprocesso é o nome novo do tipo acima; o leitor converte no Parse.
	TypeSubprocesso NodeType = "subprocesso"
	// TypeStopGate representa uma barreira de auditoria humana no Processo.
	TypeStopGate NodeType = "stop_gate"
	// Nó de Ticket: unidade cognitiva de engenharia.
	TypeUnitTask NodeType = "unit_task"
)

// ExecMode define onde a peça roda: na pasta do projeto ou em workspace isolado.
type ExecMode string

const (
	ExecInline   ExecMode = "inline"
	ExecIsolated ExecMode = "isolated"
)

// Execution é a configuração de execução da peça. Ausente = inline.
type Execution struct {
	Mode ExecMode `json:"mode"`
}

// MaxSlotEntries limita cada lista de slot (gates, scripts, hooks pre/post) de uma peça.
const MaxSlotEntries = 16

// Hooks são comandos executados antes (pre) e depois (post) da peça.
type Hooks struct {
	Pre  [][]string `json:"pre,omitempty"`
	Post [][]string `json:"post,omitempty"`
}

// Slots são os encaixes da peça: gates extras, scripts e hooks. Cada comando é um argv
// (sem shell) e passa pela mesma validação dos nós script.
type Slots struct {
	Gates   []string   `json:"gates,omitempty"`
	Scripts [][]string `json:"scripts,omitempty"`
	Hooks   Hooks      `json:"hooks,omitzero"`
}

// TriggerType é o gatilho que dispara o Processo.
type TriggerType string

const (
	TriggerManual     TriggerType = "manual"
	TriggerCron       TriggerType = "cron"
	TriggerFileChange TriggerType = "file_change"
	TriggerLoopEnd    TriggerType = "loop_end"
	// TriggerProcessoFim é o nome novo do gatilho acima; o leitor converte no Parse.
	TriggerProcessoFim TriggerType = "processo_fim"
)

// Representação interna depois do Parse: o leitor converte o dialeto novo nestes valores, que
// os pacotes consumidores comparam com as constantes exportadas do dialeto antigo.
const (
	tipoComposto      = TypeSubloop
	gatilhoFimInterno = TriggerLoopEnd
)

// campoAlvo é a chave JSON do Processo alvo no dialeto antigo do trigger.
const campoAlvo = "loop"

// Trigger é o gatilho único de um Processo (DAG v2). Cada tipo aceita só o seu campo.
type Trigger struct {
	Type TriggerType `json:"type"`
	Cron string      `json:"cron,omitempty"`
	Path string      `json:"path,omitempty"`
	Loop string      `json:"loop,omitempty"`
	// Processo é o campo novo do alvo; o leitor move o valor para o campo antigo no Parse.
	Processo string `json:"processo,omitempty"`
	// Conflito marca trigger com alvos diferentes nos dois dialetos; a validação recusa.
	Conflito bool `json:"-"`
}

// alvo devolve o endereço do Processo alvo na representação interna (campo do dialeto antigo).
func (t *Trigger) alvo() *string { return &t.Loop }

// EffortLevel define o teto de raciocínio cognitivo do nó (Harness Engineering).
type EffortLevel string

const (
	EffortLow    EffortLevel = "low"
	EffortMedium EffortLevel = "medium"
	EffortHigh   EffortLevel = "high"
)

// RoutingConfig descreve qual cérebro processa o nó.
type RoutingConfig struct {
	Mode     Mode        `json:"mode"`
	Harness  string      `json:"harness,omitempty"`
	Model    string      `json:"model,omitempty"`
	Effort   EffortLevel `json:"effort,omitempty"`
	Fallback string      `json:"fallback,omitempty"`
	Combo    string      `json:"combo,omitempty"`
}

// HandoffContract é a saída imutável do nó e sua condição de saída.
type HandoffContract struct {
	Artifact string `json:"artifact"`
	// SubArtifact só vale no nó subprocesso: o caminho do artefato final dentro do diretório do
	// subprocesso. O orchestrix lê esse caminho no subprocesso e grava em Artifact no pai.
	// Vazio = o mesmo caminho de Artifact.
	SubArtifact string `json:"sub_artifact,omitempty"`
	Schema      string `json:"schema,omitempty"`
	Gate        string `json:"gate"`
}

// Node é uma peça Lego do Processo.
type Node struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name,omitempty"`
	Type                NodeType          `json:"type"`
	Nature              string            `json:"nature,omitempty"`
	CognitiveComplexity string            `json:"cognitive_complexity,omitempty"`
	Ref                 string            `json:"ref,omitempty"`
	DependsOn           []string          `json:"depends_on,omitempty"`
	Dependencies        []string          `json:"dependencies,omitempty"`
	Command             []string          `json:"command,omitempty"`
	TimeoutSec          int               `json:"timeout_sec,omitempty"`
	Routing             RoutingConfig     `json:"routing"`
	Inputs              map[string]string `json:"inputs,omitempty"`
	Output              HandoffContract   `json:"output"`
	Execution           *Execution        `json:"execution,omitempty"`
	Slots               *Slots            `json:"slots,omitempty"`
	Isolation           map[string]any    `json:"isolation,omitempty"`
	TokenBudget         map[string]any    `json:"token_budget,omitempty"`
	Tooling             map[string]any    `json:"tooling,omitempty"`
	ContextIn           []string          `json:"context_in,omitempty"`
	ContextOut          []string          `json:"context_out,omitempty"`
	TestSeam            map[string]string `json:"test_seam,omitempty"`
	Verification        map[string]any    `json:"verification,omitempty"`
	Limits              map[string]any    `json:"limits,omitempty"`
	RollbackStrategy    string            `json:"rollback_strategy,omitempty"`
	SpecRef             string            `json:"spec_ref,omitempty"`
	StopGate            any               `json:"stop_gate,omitempty"`
	Status              string            `json:"status,omitempty"`
}

// EffectiveMode devolve o modo de execução, assumindo inline quando não declarado.
func (n Node) EffectiveMode() ExecMode {
	if n.Execution == nil || n.Execution.Mode == "" {
		return ExecInline
	}
	return n.Execution.Mode
}

// HasSlots indica se a peça tem algum gate, script ou hook de slot configurado.
func (n Node) HasSlots() bool {
	s := n.Slots
	return s != nil && len(s.Gates)+len(s.Scripts)+len(s.Hooks.Pre)+len(s.Hooks.Post) > 0
}

// DAG é o Processo validado. Só é obtido via Parse, garantindo invariantes.
type DAG struct {
	Schema             string         `json:"$schema,omitempty"`
	ID                 string         `json:"id,omitempty"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	Version            string         `json:"version,omitempty"`
	SchemaVersion      int            `json:"schema_version,omitempty"`
	Date               string         `json:"date,omitempty"`
	UpstreamSignatures map[string]any `json:"upstream_signatures,omitempty"`
	ExecutionStrategy  map[string]any `json:"execution_strategy,omitempty"`
	Trigger            *Trigger       `json:"trigger,omitempty"`
	Nodes              []Node         `json:"nodes"`

	gates    map[string]*gate.Expr
	layers   [][]string
	byID     map[string]int
	warnings []string
}

// EffectiveTrigger devolve o gatilho do Processo; DAG v1 sem trigger é tratado como manual.
func (d *DAG) EffectiveTrigger() Trigger {
	if d.Trigger == nil {
		return Trigger{Type: TriggerManual}
	}
	return *d.Trigger
}

// Warnings devolve avisos não fatais da validação (ex.: DAG v1 sem trigger).
func (d *DAG) Warnings() []string { return append([]string(nil), d.warnings...) }

// Layers devolve as camadas topológicas: nós de uma mesma camada podem rodar em paralelo.
func (d *DAG) Layers() [][]string {
	out := make([][]string, len(d.layers))
	for i, l := range d.layers {
		out[i] = append([]string(nil), l...)
	}
	return out
}

// Node devolve o nó pelo ID, ou nil.
func (d *DAG) Node(id string) *Node {
	i, ok := d.byID[id]
	if !ok {
		return nil
	}
	return &d.Nodes[i]
}

// Gate devolve o gate compilado do nó, ou nil.
func (d *DAG) Gate(id string) *gate.Expr { return d.gates[id] }
