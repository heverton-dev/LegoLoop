# 📐 Diretório de Templates Canônicos do LegoLoop (SSOT)

Este diretório centraliza os modelos oficiais (**Single Source of Truth**) para todos os artefatos de governança, esteira e peças atômicas do ecossistema LegoLoop.

Todos os templates aqui presentes foram desenhados com base nos **padrões abertos da indústria (Anthropic MCP, Google DeepMind Agentic Coding, OpenAI Swarm e BDD/TDD)**, sendo 100% universais e agnósticos a qualquer harness (Antigravity, Claude Code, Cursor, Windsurf ou Terminal puro).

---

## 1. A Cascata de Governança em 4 Estágios

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                       A CASCATA CANÔNICA EM 4 ESTÁGIOS                     │
├──────────────┬──────────────┬──────────────┬────────────────────────────────┤
│ 1. INTENT    │ 2. SPEC      │ 3. PLAN      │ 4. ROUTE                       │
│ "O Porquê"   │ "O Quê"      │ "O Como"     │ "A Máquina"                    │
│ TEMPLATE_    │ TEMPLATE_    │ TEMPLATE_    │ TEMPLATE_                      │
│ INTENT.md    │ SPEC.md      │ PLAN.md      │ ROUTE.json                     │
└──────────────┴──────────────┴──────────────┴────────────────────────────────┘
```

1. **[`TEMPLATE_INTENT.md`](TEMPLATE_INTENT.md):** Alinhamento de intenção com a Entrevista Socrática do Addy Osmani adaptada em PT-BR (funil até >= 95% de confiança).
2. **[`TEMPLATE_SPEC.md`](TEMPLATE_SPEC.md):** Fusão Addy Osmani (Premissas + Boundaries) + Matt Pocock (Costuras de Teste) + AIDD-Spec (Critérios Binários BDD `AC-XX`).
3. **[`TEMPLATE_PLAN.md`](TEMPLATE_PLAN.md):** Fatiamento Vertical do Addy Osmani com Governança em 4 Níveis de Supervisão Humana (`nano`, `micro`, `epic`, `autonomous`).
4. **[`TEMPLATE_ROUTE.json`](TEMPLATE_ROUTE.json):** Matriz determinística isomórfica executável em JSON com isolamento de duas camadas e contingência dupla (harness + modelo).

---

## 2. A Sextilha Canônica de Peças Padrão Ouro (Score >= 0.85)

Estes templates definem o **Padrão Ouro da Indústria** para cada um dos 6 tipos de peças que compõem a caixa de ferramentas do LegoLoop:

| Tipo de Peça | Arquivo de Template | Especificação / Autoridade de Referência |
|:---|:---|:---|
| **1. Skill** | [`TEMPLATE_PIECE_SKILL.md`](TEMPLATE_PIECE_SKILL.md) | Google DeepMind & Anthropic (YAML Frontmatter + Guardrails Negativos + I/O + Testes) |
| **2. Script** | [`TEMPLATE_PIECE_SCRIPT.ps1`](TEMPLATE_PIECE_SCRIPT.ps1) | Engenharia de Sistemas (Validação de Parâmetros + Idempotência `-DryRun` + Exit Codes Semânticos) |
| **3. Hook** | [`TEMPLATE_PIECE_HOOK.json`](TEMPLATE_PIECE_HOOK.json) | OpenAI / Interceptadores Seguros (Fail-Safe + Sanitização de Segredos + Timeout) |
| **4. MCP** | [`TEMPLATE_PIECE_MCP.json`](TEMPLATE_PIECE_MCP.json) | Anthropic MCP Standard 2024-11-05 (JSON-RPC com Schemas Tipados estritos) |
| **5. Gate** | [`TEMPLATE_PIECE_GATE.json`](TEMPLATE_PIECE_GATE.json) | Quality Gates Determinísticos (Expressões Booleanas com verificação mecânica de terminal) |
| **6. Loop** | [`TEMPLATE_PIECE_LOOP.json`](TEMPLATE_PIECE_LOOP.json) | OpenWorkflows / n8n-style (DAG Acíclico Componível com isolamento estéril por nó) |

---

## 3. Régua Matricial de Maturação (Badges & Cores)

Toda peça avaliada pelo motor `internal/reputation` recebe uma nota matemática de $0.00$ a $1.00$ e é classificada visualmente no Studio:

* 🥇 **`ouro`** (`#F59E0B` - Dourado): Score $\ge 0.85$ (Atinge integralmente os requisitos dos templates Padrão Ouro).
* 🥈 **`prata`** (`#94A3B8` - Prata): Score $0.60 - 0.84$ (Funcional, mas falta rigor de testes ou contratos formais).
* 🥉 **`bronze`** (`#F97316` - Bronze/Laranja): Score $0.35 - 0.59$ (Peça anêmica, texto corrido sem verificação).
* ☣️ **`toxica`** (`#EF4444` - Vermelho Alerta): Score $< 0.35$ (Prolixa, sem comandos, risco de alucinação).

---

## 4. Sub-Loop de Elevação de Peças (Upgrade)

Qualquer peça em nível Bronze ou Prata pode ser elevada ao Padrão Ouro acionando a esteira canônica:
`Diagnóstico de Gaps` ➔ `Entrevista Socrática` ➔ `SPEC` ➔ `PLAN` ➔ `ROUTE` ➔ `BUILD TDD` ➔ `Peça Padrão Ouro`.
