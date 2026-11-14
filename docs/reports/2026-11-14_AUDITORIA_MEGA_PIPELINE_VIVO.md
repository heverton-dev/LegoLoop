# Auditoria Bottom-Up do Mega-Pipeline (Documento Vivo)

> **Sessão:** `95d7dbd7-0074-4021-91ef-bcd376890df0`
> **Data:** 2026-11-14
> **Alvo:** `loops/mega-pipeline.loop.json`
> **Objetivo:** Auditar e homologar cada camada no nível nanométrico (do nano ao macro), garantindo uso de peças Padrão Ouro e contratos fechados.

---

## Método Repetível de Auditoria por Camada (4 Passos)

1. **Auditar a peça/etapa atual:** passo a passo, subcomandos, funções, entradas, saídas, gates, hooks e dependências.
2. **Investigar no ecossistema interno:** buscar em `loops/`, `MEMORY.md` e `docs/` se a fase já foi lapidada ou transformada em motor nativo.
3. **Consultar a Matriz de Peças:** filtrar `legoloop matriz ver --json` para identificar a melhor peça de encaixe (Ouro/Prata/Bronze/Tóxica).
4. **Consolidar e decidir:** registrar o diagnóstico neste documento antes de editar contratos ou avançar para a próxima camada.

---

## Camada 1 — `l1-concepcao` (Ficha Técnica Completa)

- **Referência atual no `mega-pipeline.loop.json`:** `skills/aidd-grill` (reprovada: referencia scripts externos inexistentes em `componentes/`).
- **Peça Ouro Substituta:** Motor nativo `internal/interview` via CLI `legoloop fabricar entrevista`.
- **Regras da Baia:**
  - Não inventar respostas não fornecidas pelo usuário (`"Não informado na entrevista"`).
  - Só liberar `INTENT.md` quando `Confidence >= 95` (ou `--aceitar-rascunho` explícito).
  - Isolamento estéril por sessão (`<Vault>/sessions/<id>.json`).
- **Ferramentas:** CLI `legoloop fabricar entrevista`, `internal/interview` (`Engine`, `Store`, `RenderIntent`), `internal/brain` (`brain.Pick`).
- **Cérebro & Esforço (`routing`):** `mode: auto` via `brain.Pick` (9Router -> Ollama -> CLI harness), esforço `medium`.
- **Subcamadas (Processamento, Input, Output e Gate):**
  1. **Subcamada 1.1 — `iniciar "<ideia>"` (`runIniciar`):**
     - **Input (`context_in`):** string `idea` + arquivos do projeto via `LoadProjectContext(docs)`.
     - **Processamento:** `eng.Start(ctx, idea)` cria sessão `intv_<nano>` com `Confidence: 10` e gera a 1ª rodada de perguntas.
     - **Output (`context_out`):** `<Vault>/sessions/<id>.json` + perguntas no stdout.
     - **Gate:** `exit_code == 0` e arquivo `<id>.json` existente.
  2. **Subcamada 1.2 — `responder <sessao> "<texto>"` (`runResponder`):**
     - **Input (`context_in`):** `<Vault>/sessions/<id>.json` + resposta do usuário (`text`).
     - **Processamento:** `eng.Step(ctx, id, text)` atualiza espelho socrático (`parseMirror`) e recalcula `Confidence`.
     - **Output (`context_out`):** `<Vault>/sessions/<id>.json` atualizado + próxima rodada (ou conclusão se `>= 95`).
     - **Gate:** `Confidence >= 95`.
  3. **Subcamada 1.3 — `intent <sessao>` (`runIntent`):**
     - **Input (`context_in`):** `<Vault>/sessions/<id>.json` (com `Mirror` preenchido).
     - **Processamento:** `eng.MaterializeIntent` preenche `TEMPLATE_INTENT` via `RenderIntent` sem placeholders `{{...}}`.
     - **Output (`context_out`):** `docs/intents/YYYY-MM-DD_INTENT_<NOME>.md`.
     - **Gate:** `file_exists('docs/intents/YYYY-MM-DD_INTENT_<NOME>.md')` e ausência de `{{`.
- **Decisão:** Camada 1 homologada com `legoloop fabricar entrevista` substituindo `aidd-grill`.

---

## Decisões Arquiteturais Transversais

- **Ficha Técnica Visual ("Carta de Trunfo"):** Deve ser derivada em tempo real no Studio (ou regerada deterministicamente em `fabricar`, `refinar` e `fechar`) direto do `pipeline.dag.json` + Matriz de Peças (SSOT), nunca mantida à mão.
- **Canonização Universal (SO + Harnesses):** Todo processo do LegoLoop deve residir no cofre (`~/.aidd/canonical`) e estar vinculado dinamicamente a todos os harnesses detectados na máquina do usuário (`legoloop doctor` / `instalar-skills`) para uso imediato no SO e em qualquer agente.

## Camada 2 — `l2-design` (Ficha Técnica Completa)

- **Referência atual no `mega-pipeline.loop.json`:** `skills/impeccable` (`tier=bronze`).
- **Input (`context_in`):** `docs/intents/INTENT.md`.
- **Processamento:** `impeccable` v4.5.0 (`shape`, `document`, `critique`, `audit`, `polish`) via `npx impeccable`.
- **Output (`context_out`):** `DESIGN.md`.
- **Gate:** `file_exists('DESIGN.md')`.
- **Diagnóstico:** `impeccable` é específico para UI/UX e frontend; não deve bloquear projetos sem interface visual.
- **Decisão:** A execução de `l2-design` (`impeccable`) deve ser condicional à existência de superfície visual, decidida por uma **Peça de Triagem / Roteador de Entrada** antes da esteira linear.

---

## Pausa Estratégica — Peça de Triagem / Roteador Universal

- **Decisão:** Pausar a auditoria do `mega-pipeline` na Camada 2 para projetar e implementar a **Peça de Triagem / Roteador de Entrada** como componente obrigatório de todo processo no LegoLoop (gerada automaticamente ao final de `legoloop fabricar` e operável via CLI, Chat/Skill e Studio).
- **Retomada:** Após implementar e homologar a Peça de Triagem Universal, retomar este documento vivo a partir da Camada 3 (`l3-triade` -> Sexteto).

---

## Próximas Camadas (Pendentes)

- [ ] Camada 2: `l2-design` (`skills/impeccable`)
- [ ] Camada 3: `l3-triade` (`skills/aidd-spec` -> Sexteto)
- [ ] Camada 4: `l3b-swarm` (`skills/legoloop-swarm`)
- [ ] Camada 5: `l4-tdd` (`skills/aidd-tdd`)
- [ ] Camada 6: `l5-homologue` (`skills/aidd-homologue`)
- [ ] Camada 7: `l6-fechar-canonico` (`fechar`)
