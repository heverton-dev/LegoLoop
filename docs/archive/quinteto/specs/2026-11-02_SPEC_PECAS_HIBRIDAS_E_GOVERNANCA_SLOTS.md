# 📋 SPEC: Taxonomia de Peças Híbridas, Triggers e Governança de Slots

> **Status:** Proposta para Revisão (Stop-Gate 02)  
> **Data:** 2026-11-02  
> **Referência:** [docs/intents/2026-11-02_INTENT_PECAS_HIBRIDAS_E_GOVERNANCA_SLOTS.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_PECAS_HIBRIDAS_E_GOVERNANCA_SLOTS.md)  

---

## 1. Visão Geral e Arquitetura
Esta especificação estabelece:
1. **O Modelo Híbrido (*Dual-Role*):** Expansão do motor de DAG para tratar `gate`, `hook`, `mcp` e `binary` como nós de primeira classe, além do seu uso em slots.
2. **Nós de Gate Dual-Mode:** Capacidade do nó de gate atuar como *Stop-Gate binário* (pausa para intervenção/aprovação) e como *bifurcação condicional de rotas* (`branch` com `on_pass` e `on_fail`).
3. **Triggers como Cabeça Universal:** Validação mandatória de que todo DAG possui um gatilho canônico raiz (`manual`, `cron`, `file_change`, `loop_end`).
4. **Governança de Slots Canônicos com Override Seguro:** Padrão onde peças canônicas são imutáveis e consumidas intactas por padrão, permitindo `slots_override` declarativo por nó ou Fork total da peça para o projeto.

---

## 2. Contratos de Dados & I/O

### 2.1 Expansão de `NodeType` (`internal/dag/types.go`)
```go
type NodeType string

const (
    TypeSkill    NodeType = "skill"
    TypeScript   NodeType = "script"
    TypeSubloop  NodeType = "subloop"
    TypeGate     NodeType = "gate"     // Novo: Stop-gate ou bifurcação condicional
    TypeHook     NodeType = "hook"     // Novo: Event listener de ciclo de vida do DAG
    TypeMCP      NodeType = "mcp"      // Novo: Execução mecânica direta de tool MCP
    TypeBinary   NodeType = "binary"   // Novo: Execução direta de executável local
)
```

### 2.2 Estrutura de `SlotsOverride` (`internal/dag/types.go`)
```go
type SlotsOverride struct {
    AddGates      []string   `json:"add_gates,omitempty"`
    RemoveGates   []string   `json:"remove_gates,omitempty"`
    AddScripts    [][]string `json:"add_scripts,omitempty"`
    RemoveScripts [][]string `json:"remove_scripts,omitempty"`
    AddHooks      *Hooks     `json:"add_hooks,omitempty"`
}
```

### 2.3 Configuração de Gate Dual-Mode (`Node.GateConfig`)
```go
type GateConfig struct {
    Mode    string `json:"mode"`              // "stop_gate" ou "branch"
    Expr    string `json:"expr"`              // Expressão booleana avaliada sobre a entrada/artefato
    OnPass  string `json:"on_pass,omitempty"` // ID do nó de destino em caso de sucesso (modo branch)
    OnFail  string `json:"on_fail,omitempty"` // ID do nó de destino em caso de reprovação (modo branch)
    Message string `json:"message,omitempty"` // Mensagem de pausa para o operador humano
}
```

---

## 3. Critérios Binários de Aceite por Node

### ÉPICO 01 — Taxonomia & Expansão de Tipos no DAG
* [ ] `AC-E01.1`: `internal/dag/types.go` declara formalmente as constantes `TypeGate`, `TypeHook`, `TypeMCP` e `TypeBinary` no enum `NodeType`.
* [ ] `AC-E01.2`: `internal/dag/validate.go` aceita nós de tipo `gate`, exigindo declaração de `GateConfig.Expr` válida e modo (`stop_gate` ou `branch`).
* [ ] `AC-E01.3`: `internal/dag/validate.go` rejeita nós `mcp` que não especifiquem `server` e `tool` no campo de comando/configuração.
* [ ] `AC-E01.4`: Suíte de testes `internal/dag/validate_test.go` cobre validação positiva e negativa de cada um dos 4 novos tipos de nós.

### ÉPICO 02 — Padronização Universal de Triggers
* [ ] `AC-E02.1`: `internal/dag/validate.go` valida que `dag.Trigger != nil` e que seu tipo pertence ao conjunto canônico (`manual`, `cron`, `file_change`, `loop_end`).
* [ ] `AC-E02.2`: DAGs herdados (v1) sem campo trigger são normalizados em memória com trigger default `manual` e recebem warning suave na validação sem bloquear a execução.
* [ ] `AC-E02.3`: Testes unitários em `internal/dag/trigger_test.go` comprovam a normalização do nó raiz virtual `#trigger`.

### ÉPICO 03 — Governança de Slots & Mecanismo de Override
* [ ] `AC-E03.1`: Struct `Node` recebe o campo opcional `SlotsOverride *SlotsOverride `json:"slots_override,omitempty"``.
* [ ] `AC-E03.2`: Função `ResolveEffectiveSlots(n Node, canonical Slots) Slots` em `internal/dag/slots.go` retorna `canonical` intacto quando `SlotsOverride == nil`.
* [ ] `AC-E03.3`: Quando `SlotsOverride` estiver presente, `ResolveEffectiveSlots` adiciona os itens de `Add*` e filtra/exclui os itens de `Remove*` com equivalência exata de tokens/nomes.
* [ ] `AC-E03.4`: Testes unitários em `internal/dag/slots_test.go` cobrem: herança pura (sem override), adição de gates, remoção cirúrgica de scripts e conflito nulo.

### ÉPICO 04 — Adaptação Visual no Looping Studio
* [ ] `AC-E04.1`: Canvas visual renderiza o nó de Trigger com ícone `⚡` e cor temática de gatilho, conectado aos nós de entrada do DAG.
* [ ] `AC-E04.2`: Nós de Gate exibem visualmente o badge `🛑 Stop-Gate` ou bifurcação visual com saídas `pass` e `fail`.
* [ ] `AC-E04.3`: O painel Inspetor (`internal/studio/static/inspector.js`) exibe a seção "Slots & Overrides" permitindo desativar ou adicionar slots da peça selecionada.

---

## 4. Critério de Done Global
> A demanda será considerada **CONCLUÍDA** quando 100% dos critérios binários `AC-E01.1` a `AC-E04.3` estiverem marcados como `[x]`, a suíte `go test ./...` estiver 100% verde e o binário `bin/legoloop.exe` recompilado com sucesso.

---

🛑 **Stop-Gate 02:** Aguardando aprovação do usuário para a Tríade de Planejamento (SPEC + PLAN + ROUTE).
