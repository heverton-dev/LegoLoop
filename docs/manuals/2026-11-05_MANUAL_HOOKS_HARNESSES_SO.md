# MANUAL OPERACIONAL: HOOKS NOS HARNESSES DE IA DO SISTEMA OPERACIONAL

> **Data:** 2026-11-05
> **Autoridade:** Ecossistema LegoLoop & Orca ADE
> **Status:** Canônico
> **Escopo:** Como cada harness de IA instalado neste SO trata hooks, onde se configura, quais eventos existem e como o Orca unifica tudo como meta-harness
> **Fontes de verdade:** `~/.claude/settings.json` · `~/.config/opencode/opencode.json` · `~/.orca/agent-hooks/` · documentação oficial de cada harness

---

## 1. Visão Geral

Um **hook** é um comando externo que o harness executa automaticamente em momentos específicos do ciclo de vida de uma sessão (início, antes/depois de usar uma ferramenta, parada, etc.). Hooks permitem:

- Interceptar e bloquear saídas com segredos (leak-scan)
- Notificar sistemas externos (Orca) sobre atividade do agente
- Injetar instruções no contexto do modelo
- Auditar e registrar eventos

Cada harness implementa (ou não) hooks de forma diferente. O **Orca** atua como meta-harness, unificando a experiência via scripts `.cmd`/`.cjs` na pasta `~/.orca/agent-hooks/`.

---

## 2. Inventário de Harnesses Instalados

| Harness | Binário/Caminho | Hooks nativos | Hooks via Orca |
|:---|:---|:---|:---|
| **Claude Code** | `~/.local/bin/claude` | Sim (settings.json) | Sim |
| **OpenCode** | `~/.local/bin/opencode` | Não | Não |
| **Orca** | `~/.orca/` (meta-harness) | Sim (agent-hooks/) | — |
| **Codex** (OpenAI) | via Orca | Não nativo | Sim |
| **Gemini CLI** (Google) | via Orca | Parcial | Sim |
| **Kiro** (AWS) | via Orca | Parcial | Sim |
| **Copilot** (GitHub) | via Orca | Não nativo | Sim |
| **Grok** (xAI) | via Orca | Não nativo | Sim |
| **Antigravity** | via Orca | Sim (granular) | Sim |

---

## 3. Claude Code — Hooks Nativos

### 3.1 Onde se configura

| Escopo | Arquivo |
|:---|:---|
| Global (todos os projetos) | `~/.claude/settings.json` → bloco `"hooks"` |
| Projeto (local) | `.claude/settings.local.json` → bloco `"hooks"` |

Hooks do projeto sobrepõem os globais quando o evento e matcher coincidem.

### 3.2 Eventos disponíveis

| Evento | Quando dispara |
|:---|:---|
| `SessionStart` | Ao iniciar uma sessão |
| `SessionEnd` | Ao encerrar uma sessão |
| `UserPromptSubmit` | Quando o usuário envia um prompt |
| `PreToolUse` | Antes de executar qualquer ferramenta |
| `PostToolUse` | Depois de executar qualquer ferramenta |
| `PostToolUseFailure` | Quando uma ferramenta falha |
| `Stop` | Quando o modelo para de gerar |
| `StopFailure` | Quando a parada falha |
| `SubagentStart` | Quando um subagente é criado |
| `SubagentStop` | Quando um subagente termina |
| `PostCompact` | Após compactação do contexto |
| `Elicitation` | Quando o modelo faz uma pergunta ao usuário |
| `ElicitationResult` | Quando o usuário responde |
| `TeammateIdle` | Quando um teammate fica ocioso |

### 3.3 Estrutura de configuração

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash|PowerShell",
        "hooks": [
          {
            "type": "command",
            "command": "bigbrain.exe hook secret-guard",
            "timeout": 5
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Bash|PowerShell|Read|Grep|WebFetch|mcp__.*",
        "hooks": [
          {
            "type": "command",
            "command": "bigbrain.exe hook leak-scan",
            "timeout": 5
          }
        ]
      }
    ]
  }
}
```

**Campos:**

| Campo | Descrição |
|:---|:---|
| `matcher` | Regex contra o nome da ferramenta. `*` = todas |
| `type` | Sempre `"command"` |
| `command` | Caminho do executável + argumentos |
| `timeout` | Segundos antes de abortar (padrão: 60) |

### 3.4 Hooks ativos neste SO (Claude Code)

| Evento | Hook | Função |
|:---|:---|:---|
| `PreToolUse` | `bigbrain.exe hook secret-guard` | Bloqueia leitura de arquivos sensíveis |
| `PreToolUse` | `rtk hook claude` | Reescreve comandos Bash para proxy RTK |
| `PreToolUse` | `orca-agent-hook.exe claude` | Notifica Orca sobre atividade |
| `PostToolUse` | `bigbrain.exe hook leak-scan` | Detecta segredos na saída |
| `PostToolUse` | `orca-agent-hook.exe claude` | Notifica Orca sobre resultado |
| `SessionStart` | `orca-agent-hook.exe claude` | Registra início de sessão |
| `SessionEnd` | `orca-agent-hook.exe claude` | Registra fim de sessão |
| Todos os outros | `orca-agent-hook.exe claude` | Telemetria completa para Orca |

---

## 4. OpenCode — Sem Hooks Nativos

### 4.1 Situação

O OpenCode (`~/.config/opencode/opencode.json`) **não possui sistema de hooks** equivalente ao Claude Code. Não há bloco `"hooks"`, não há `PreToolUse`/`PostToolUse`, não há execução de comando externo em eventos do ciclo de vida.

### 4.2 Alternativas disponíveis

| Mecanismo | O que faz | Limitação |
|:---|:---|:---|
| Plugins (`oh-my-opencode-slim`, `superpowers`) | Estendem funcionalidade via módulos | Não interceptam tool-use |
| Skills (`~/.config/opencode/skills/`) | Instruções para o modelo | Só orientam, não executam |
| System prompt | Instruções de segurança | Dependem do modelo obedecer |

### 4.3 Lacuna

Sem hooks, o OpenCode **não tem proteção contra leak de segredos** via execução de ferramentas. O `bigbrain leak-scan` não roda. Única mitigação: não armazenar chaves em texto plano em arquivos que o OpenCode leia, ou usar variáveis de ambiente do SO (não do config).

---

## 5. Orca — O Meta-Harness

### 5.1 Arquitetura

O Orca não é um harness que roda modelos — ele **orquestra** outros harnesses. Sua função de hooks é unificar a telemetria e o controle de todos os harnesses via um servidor HTTP local.

```
┌──────────────┐
│ Orca daemon   │  ← HTTP server em 127.0.0.1:{port}
│ /hook/claude  │
│ /hook/codex   │
│ /hook/gemini  │
│ /hook/kiro    │
│ /hook/copilot │
│ /hook/grok    │
│ /hook/antigrav│
└──────┬───────┘
       │ curl POST (payload JSON via stdin)
       │
  ┌────┴────────────────────────┐
  │  ~/.orca/agent-hooks/       │
  │  claude-hook.cmd            │
  │  codex-hook.cmd             │
  │  gemini-hook.cmd            │
  │  kiro-hook.cmd              │
  │  copilot-hook.cmd           │
  │  grok-hook.cmd              │
  │  antigravity-hook.cmd       │
  │  antigravity-hook-post.cjs  │
  │  orca-agent-hook.exe        │
  └─────────────────────────────┘
```

### 5.2 Protocolo comum

Todos os scripts `.cmd` seguem o mesmo padrão:

1. Verificam variáveis de ambiente Orca (`ORCA_AGENT_HOOK_PORT`, `ORCA_AGENT_HOOK_TOKEN`, `ORCA_PANE_KEY`)
2. Se alguma estiver vazia, saem silenciosamente (exit 0)
3. Enviam um POST HTTP para `http://127.0.0.1:{port}/hook/{harness}` com:
   - Headers: `X-Orca-Agent-Hook-Token`, `Content-Type: application/x-www-form-urlencoded`
   - Body: `paneKey`, `tabId`, `launchToken`, `worktreeId`, `env`, `version`, `payload` (JSON do evento)
4. Timeout de conexão: 0.5s. Timeout de resposta: 1.5s
5. Falhas são silenciosas — nunca bloqueiam o harness

### 5.3 Variáveis de ambiente Orca

| Variável | Descrição |
|:---|:---|
| `ORCA_AGENT_HOOK_PORT` | Porta do daemon Orca local |
| `ORCA_AGENT_HOOK_TOKEN` | Token de autenticação do hook |
| `ORCA_PANE_KEY` | Identificador do painel/terminal |
| `ORCA_TAB_ID` | ID da aba no Orca |
| `ORCA_AGENT_LAUNCH_TOKEN` | Token de lançamento do agente |
| `ORCA_WORKTREE_ID` | ID do worktree Git associado |
| `ORCA_AGENT_HOOK_ENV` | Ambiente (dev, prod, etc.) |
| `ORCA_AGENT_HOOK_VERSION` | Versão do protocolo de hook |
| `ORCA_AGENT_HOOK_ENDPOINT` | Script auxiliar opcional |

---

## 6. Hooks por Harness via Orca

### 6.1 Claude Code (`claude-hook.cmd` + `claude-hook-impl.cmd`)

```
claude-hook.cmd → claude-hook-impl.cmd → curl POST /hook/claude
```

- Guard clauses adicionais: ignora se `CLAUDE_JOB_DIR` estiver definido (job remoto), se `CURSOR_VERSION` + `CURSOR_PROJECT_DIR` existirem (Cursor embarcado), se `DEVIN_PROJECT_DIR` existir (Devin), ou se `GROK_HOOK_EVENT` existir (Grok embarcado)
- Normaliza variáveis legadas: `ORCA_AGENT_PANE` → `ORCA_PANE_KEY`

### 6.2 Codex (`codex-hook.cmd`)

```
codex-hook.cmd → curl POST /hook/codex
```

- Mesma estrutura padrão
- Guard: ignora se `CLAUDE_JOB_DIR` existir (evita dupla notificação)

### 6.3 Gemini CLI (`gemini-hook.cmd`)

```
gemini-hook.cmd → curl POST /hook/gemini
```

- Emite `echo {}` como resposta JSON vazia para stdout (Gemini espera JSON de retorno)
- Mesma estrutura padrão

### 6.4 Kiro (`kiro-hook.cmd`)

```
kiro-hook.cmd → curl POST /hook/kiro
```

- Mesma estrutura padrão, sem particularidades

### 6.5 Copilot (`copilot-hook.cmd`)

```
copilot-hook.cmd → curl POST /hook/copilot
```

- Campo extra: `hookEventName` com o valor de `ORCA_COPILOT_HOOK_EVENT` (recebido como argumento `%~1`)
- Emite `echo {}` como resposta JSON

### 6.6 Grok (`grok-hook.cmd`)

```
grok-hook.cmd → curl POST /hook/grok
```

- Campo extra: `grokHome` com o valor de `GROK_HOME` (caminho de instalação do Grok)
- Sanitiza `GROK_HOME`: remove aspas, trata barra final, limita a 4096 caracteres

### 6.7 Antigravity (`antigravity-*.cmd` + `antigravity-hook-post.cjs`)

O Antigravity é o único harness com **hooks granulares por evento via arquivos separados**:

| Arquivo | Evento | Resposta padrão |
|:---|:---|:---|
| `antigravity-pre-tool-use.cmd` | `PreToolUse` | `{"decision":"ask"}` |
| `antigravity-post-tool-use.cmd` | `PostToolUse` | `{}` |
| `antigravity-pre-invocation.cmd` | `PreInvocation` | `{}` |
| `antigravity-post-invocation.cmd` | `PostInvocation` | `{}` |
| `antigravity-stop.cmd` | `Stop` | `{"decision":""}` |

Todos delegam para `antigravity-hook.cmd`, que:

1. Emite a resposta JSON conforme o evento
2. Chama `ORCA_AGENT_HOOK_ENDPOINT` se existir
3. Usa o Electron do Orca (`Orca.exe` com `ELECTRON_RUN_AS_NODE=1`) para executar `antigravity-hook-post.cjs`

O arquivo `antigravity-hook-post.cjs` é um **parser de JSON streaming** em Node.js que:

- Lê stdin até completar um JSON válido (com validação de colchetes/chaves)
- Tem timeout absoluto de 7s e idle de 1.5s
- Limita payload a 1MB
- Envia POST para `/hook/antigravity` via módulo `http` nativo

---

## 7. Matriz Comparativa de Capacidades

| Capacidade | Claude Code | OpenCode | Codex | Gemini | Kiro | Copilot | Grok | Antigravity |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| Hooks nativos no config | Sim | Não | Não | Parcial | Parcial | Não | Não | Sim |
| Eventos tipados (Pre/Post) | Sim (14 eventos) | — | — | — | — | — | — | Sim (5 eventos) |
| Matcher por ferramenta | Sim (regex) | — | — | — | — | — | — | Não |
| Execução de comando externo | Sim | — | — | — | — | — | — | Sim |
| Telemetria via Orca | Sim | Não | Sim | Sim | Sim | Sim | Sim | Sim |
| Leak-scan (bigbrain) | Sim | Não | Não | Não | Não | Não | Não | Não |
| Secret-guard (bigbrain) | Sim | Não | Não | Não | Não | Não | Não | Não |
| RTK proxy | Sim | Não | Não | Não | Não | Não | Não | Não |
| JSON de retorno esperado | Sim | — | Não | Sim | Não | Sim | Não | Sim |

---

## 8. Como Adicionar Hooks em Cada Harness

### 8.1 Claude Code

Editar `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Bash|PowerShell",
        "hooks": [
          {
            "type": "command",
            "command": "caminho/do/executavel argumento",
            "timeout": 5
          }
        ]
      }
    ]
  }
}
```

Ou por projeto: `.claude/settings.local.json` com a mesma estrutura.

### 8.2 OpenCode

Não suportado. Alternativas:

1. Escrever uma skill que instrua o modelo a não exibir segredos
2. Não armazenar chaves em `opencode.json` → usar variáveis de ambiente do SO
3. Aguardar suporte nativo a hooks em versão futura

### 8.3 Codex, Gemini, Kiro, Copilot, Grok

Via Orca: criar ou editar o `.cmd` correspondente em `~/.orca/agent-hooks/`. O script deve:

1. Verificar as variáveis Orca (guard clauses)
2. Executar o comando desejado (ex: `bigbrain.exe hook leak-scan`)
3. Enviar POST para `/hook/{harness}` no Orca

Exemplo — adicionar leak-scan ao Codex (`codex-hook.cmd`):

```batch
@echo off
setlocal

REM === Leak-scan (adicionar antes do curl do Orca) ===
"C:\Users\trcnologia\.aidd-bigbrain\bin\bigbrain.exe" hook leak-scan 2>nul

REM === Notificação Orca (original) ===
if "%ORCA_AGENT_HOOK_PORT%"=="" exit /b 0
if "%ORCA_AGENT_HOOK_TOKEN%"=="" exit /b 0
if "%ORCA_PANE_KEY%"=="" exit /b 0
"%SystemRoot%\System32\curl.exe" -sS -X POST "http://127.0.0.1:%ORCA_AGENT_HOOK_PORT%/hook/codex" ^
  --connect-timeout 0.5 --max-time 1.5 ^
  -H "Content-Type: application/x-www-form-urlencoded" ^
  -H "X-Orca-Agent-Hook-Token: %ORCA_AGENT_HOOK_TOKEN%" ^
  --data-urlencode "paneKey=%ORCA_PANE_KEY%" ^
  --data-urlencode "payload@-" >nul 2>&1
exit /b 0
```

### 8.4 Antigravity

Editar os arquivos granulares em `~/.orca/agent-hooks/`:

- `antigravity-pre-tool-use.cmd` — antes de cada ferramenta
- `antigravity-post-tool-use.cmd` — depois de cada ferramenta
- `antigravity-stop.cmd` — ao parar

A resposta JSON no stdout controla o comportamento:

| Resposta | Efeito |
|:---|:---|
| `{"decision":"ask"}` | Pede confirmação ao usuário |
| `{"decision":""}` | Permite sem perguntar |
| `{}` | Sem intervenção |

---

## 9. Fluxo de Segurança: bigbrain leak-scan

### 9.1 Como funciona

```
Harness executa Bash/Read/Grep/...
    │
    ▼
PostToolUse dispara
    │
    ▼
bigbrain.exe hook leak-scan
    │
    ├─ Lê a saída da ferramenta via stdin
    ├─ Compara contra valores conhecidos (ANTHROPIC_API_KEY, etc.)
    │
    ├─ Se limpo: exit 0, sem output
    │
    └─ Se detecta segredo:
       ├─ Injeta instrução bloqueante no contexto:
       │  "Não repita, não cite, não grave esse valor"
       └─ Avisa para rotacionar a chave
```

### 9.2 Cobertura atual

| Harness | leak-scan ativo | secret-guard ativo |
|:---|:---|:---|
| Claude Code | Sim (PostToolUse) | Sim (PreToolUse) |
| OpenCode | Não | Não |
| Codex | Não (configurável) | Não |
| Gemini | Não (configurável) | Não |
| Kiro | Não (configurável) | Não |
| Copilot | Não (configurável) | Não |
| Grok | Não (configurável) | Não |
| Antigravity | Não (configurável) | Não |

### 9.3 Recomendação

Para estender o leak-scan a todos os harnesses, adicionar a chamada ao `bigbrain.exe` em cada `*-hook.cmd` conforme o exemplo da seção 8.3. O bigbrain lê stdin e retorna instrução via stdout — compatível com qualquer script `.cmd`.

---

## 10. Limitações e Riscos

| Risco | Descrição | Mitigação |
|:---|:---|:---|
| OpenCode sem hooks | Nenhuma proteção contra leak | Não armazenar segredos no config |
| Hooks silenciosos | Falhas nos hooks não bloqueiam o harness | Monitorar logs; testar hooks isoladamente |
| Timeout curto do Orca (1.5s) | Se o daemon estiver lento, o hook falha sem erro | Manter daemon saudável |
| Segredos em texto plano | `settings.json` contém chaves no bloco `"env"` | Usar variáveis de ambiente do SO ou vault |
| Hook único por evento | Alguns harnesses não suportam múltiplos hooks por evento | Encadear comandos dentro do `.cmd` |

---

## 11. Referência Rápida de Caminhos

| Recurso | Caminho |
|:---|:---|
| Claude Code — config global | `~/.claude/settings.json` |
| Claude Code — config projeto | `.claude/settings.local.json` |
| OpenCode — config | `~/.config/opencode/opencode.json` |
| Orca — hooks | `~/.orca/agent-hooks/` |
| Orca — hook principal | `~/.orca/agent-hooks/orca-agent-hook.exe` |
| bigbrain — binário | `~/.aidd-bigbrain/bin/bigbrain.exe` |
| RTK — binário | verificar com `which rtk` |

---

## 12. Glossário

| Termo | Definição |
|:---|:---|
| **Harness** | Programa que executa um modelo de IA com acesso a ferramentas (Claude Code, OpenCode, Codex, etc.) |
| **Hook** | Comando externo executado automaticamente pelo harness em um evento específico |
| **Matcher** | Padrão (regex ou glob) que filtra em quais ferramentas o hook dispara |
| **Meta-harness** | Harness que orquestra outros harnesses (Orca) |
| **leak-scan** | Hook que detecta segredos na saída de ferramentas |
| **secret-guard** | Hook que bloqueia leitura de arquivos sensíveis |
| **Daemon** | Processo Orca rodando em background, escutando em porta local |
| **Payload** | Dados JSON do evento, enviados via stdin ao hook |
