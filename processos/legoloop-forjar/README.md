# Processo `legoloop-forjar`

> **Slug:** `legoloop-forjar`  
> **Motor:** `internal/forja`  
> **Versão:** `1.0.0`  
> **Padrão:** Processo canônico 5/5 (DAG, README, Triagem, Anticorpos, Skill)

## Objetivo

Preparar e verificar o ambiente do usuário em qualquer sistema operacional de forma determinística, reversível e idempotente. A Forja configura os harnesses detectados no `HOME`, aplica o perfil de contexto escolhido (`minimo`, `padrao` ou `gateway`), preserva conteúdo próprio com Peças gerenciadas ligadas ao `AGENTS.md`, grava cópia de segurança com `MANIFESTO.tsv` e emite o par `RELATORIO.md` e `RELATORIO.json`.

## Quatro Gatilhos, Um Motor

Todos os gatilhos convergem para `forja.Executar`:

1. **Instalação:** `legoloop instalar --apply`
2. **CLI:** `legoloop forjar --dry-run | --apply | --desfazer <data> | --verificar`
3. **Agente / MCP:** `call_mcp_tool` invocando a ferramenta de instalação/forja
4. **Skill:** `/legoloop-forjar` (`skills/legoloop-forjar/SKILL.md`)

## Etapas do DAG (`processo.dag.json`)

| Nó | Nome | Papel |
| :--- | :--- | :--- |
| `e01` | `detectar-harnesses` | Detecta os harnesses instalados no `HOME` e valida opções e Path Jail |
| `e02` | `simular-forja` | Simula alterações com `--dry-run` sem gravar no disco |
| `e03` | `aplicar-forja` | Aplica o perfil com backup prévio em `~/.aidd/backups/forja_<data>/<harness>/MANIFESTO.tsv` |
| `e04` | `verificar-orcamento` | Valida o orçamento de tokens por harness contra `docs/governanca/orcamento_tokens.json` |
| `e05` | `emitir-relatorio` | Consolida e valida `RELATORIO.md` e `RELATORIO.json` livres de segredos |

## Artefatos do Processo

- `processos/legoloop-forjar/processo.dag.json`
- `processos/legoloop-forjar/README.md`
- `processos/legoloop-forjar/triagem.json`
- `processos/legoloop-forjar/ANTICORPOS.md`
- `skills/legoloop-forjar/SKILL.md`
