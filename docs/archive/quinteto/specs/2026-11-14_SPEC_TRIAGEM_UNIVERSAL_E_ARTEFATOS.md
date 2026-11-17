# SPEC — Peça de Triagem Universal e Artefatos do Fabricar

- **Data:** 2026-11-14
- **INTENT:** `docs/intents/2026-11-14_INTENT_TRIAGEM_UNIVERSAL_E_ARTEFATOS.md`

## 1. Contrato da Peça de Triagem (`triagem.json`)
Cada processo possui uma Peça de Triagem (`loops/<slug>/triagem.json` e projeção no cofre) com:
- `schema_version`: `"triagem.v1"`
- `processo`: slug do processo
- `regras_deterministicas`: lista de `{padrao_slug, comandos_cli, artefatos_ativos, prioridade}`
- `criterio_ambiguidade`: descrição e catálogo resumido dos nós/etapas para classificação probabilística
- `pontos_entrada`: `["cli", "chat", "studio"]`

## 2. Motor Híbrido (`internal/triagem`)
1. **Camada Determinística (custo zero de LLM):**
   - Match exato por comando/slash/slug na entrada.
   - Match por artefato/estado ativo inequívoco no workspace.
   - Match único por regra cadastrada em `triagem.json`.
2. **Camada Probabilística (LLM apenas na ambiguidade):**
   - Acionada quando a entrada é linguagem natural aberta sem match determinístico único ou em caso de empate.
   - Retorna `{processo, etapa_inicial, confianca, justificativa, precisa_confirmacao}`.

## 3. Artefatos do `legoloop fabricar` (13 Artefatos)
Remover `template` de `plannedFiles` (`internal/fabricar/camada4.go`) e garantir geração de:
1-9: `skill`, `loop`, `dag`, `processo`, `intent`, `spec`, `plan`, `route`, `dag_iso`
10: `triagem` (`triagem.json`)
11: `manual` (`docs/manuals/YYYY-MM-DD_MANUAL_<SLUG>.md`)
12: `livro` (`docs/books/<slug>/livro.json`)
13: `fichas_tecnicas` (`docs/fichas/<slug>/<no>.md` para cada nano-processo/nó do DAG)

## 4. Pontos de Entrada
- **CLI:** `legoloop triagem "<entrada>" [--workspace DIR] [--json]` e `legoloop migrar-loops [--dry-run] [--apply]`
- **Chat/Skill:** seção de triagem padronizada em `SKILL.md` instruindo o agente a invocar `legoloop triagem` antes de adivinhar rota.
- **Studio:** endpoint `POST /api/triagem` e barra de entrada rápida no Studio consumindo `internal/triagem`.

## 5. Saneamento de `loops/`
Auditar os 10 diretórios em `loops/`, remover processos/peças legadas inúteis ou quebradas e gerar `triagem.json` + artefatos faltantes via `legoloop migrar-loops`.
