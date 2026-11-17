# 📋 SPEC: Blindagem da Rota Formal no `aidd-legoloop-mega-pipeline`

> **Status:** Aprovado pelo Usuário  
> **Data de Aprovação:** 2026-11-01  
> **Referência de Intenção:** [docs/intents/2026-11-01_INTENT_BLINDAGEM_ROTA_SKILL.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_BLINDAGEM_ROTA_SKILL.md)  
> **Alvo:** `C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`

---

## 1. Visão Geral e Contrato

O objetivo desta especificação é eliminar a lacuna comportamental onde o agente pula a geração do mapa formal de rota em `docs/routes/`. A skill deve instituir a **Tríade de Planejamento do Loop 3**:
* `SPEC.md` ➔ Define o QUE e os CRITÉRIOS BINÁRIOS.
* `PLAN.md` ➔ Define a HIERARQUIA e os ÉPICOS/TICKETS.
* `ROUTE.json` ➔ Define o GRAFO EXECUTÁVEL DE EXECUÇÃO ESTÉRIL (dependências, entradas, saídas, modelos e gates por nó).

O avanço para o Loop 4 (Build) é fisicamente bloqueado enquanto a Tríade não estiver completa no disco e aprovada pelo usuário.

---

## 2. Requisitos Técnicos de Modificação no `SKILL.md`

### 2.1. Reforço no Bloco de Alerta e nas Leis
* Atualizar o alerta inicial `> [!CAUTION]` para declarar que tentar gerar código sem `ROUTE.json` é quebra crítica de protocolo.
* Na Lei 3 ("Plano Único e Mapa Formal de Rota"), amarrar que todo plano é acompanhado obrigatoriamente do seu arquivo `docs/routes/YYYY-MM-DD_ROUTE_<NOME>.json`.

### 2.2. Atualização do Diagrama Mermaid e Descrição do Loop 3
* No diagrama Mermaid, o Loop 3 deve indicar explicitamente a emissão dos 3 artefatos: `SPEC.md`, `PLAN.md` e `ROUTE.json`.
* O Gate do Loop 3 (STOP-GATE HUMANO 02) deve exigir a aprovação expressa do trio antes da liberação do Loop 4.

### 2.3. Adição do Template Canônico de Rota JSON na Seção 4
* Incluir o esqueleto canônico:
  ```json
  {
    "name": "route-<nome-da-demanda>",
    "version": "1.0.0",
    "date": "YYYY-MM-DD",
    "intent_ref": "docs/intents/YYYY-MM-DD_INTENT_*.md",
    "spec_ref": "docs/specs/2026-11-01_SPEC_*.md",
    "plan_ref": "docs/plans/2026-11-01_PLAN_*.md",
    "execution_strategy": {
      "engine": "orchestrix_hermetic",
      "subagents_allowed": false,
      "max_parallel_workers": 1,
      "tdd_enforced": true
    },
    "nodes": [
      {
        "id": "TICK-01.1",
        "nature": "micro",
        "execution_mode": "synchronous",
        "depends_on": [],
        "routing": {
          "mode": "suggested",
          "model": "flash",
          "effort": "medium"
        },
        "inputs": {
          "spec_source": "docs/specs/..."
        },
        "output": {
          "artifacts": [
            "caminho/do/arquivo"
          ],
          "gate": "has_pattern('...', '...')"
        }
      }
    ]
  }
  ```

---

## 3. Critérios Binários de Aceite (Gates Mecânicos)

* [ ] **G1 (Cláusula de Rota Inviolável):** O `SKILL.md` declara que o mapa formal de rota em `docs/routes/` é pré-requisito indispensável para qualquer execução de código.
* [ ] **G2 (Diagrama com Tríade de Planejamento):** O Mermaid do fluxo exibe claramente a saída conjunta de `SPEC`, `PLAN` e `ROUTE`.
* [ ] **G3 (Template Canônico JSON de Rota):** A Seção de Templates possui o template JSON completo de rota com campos `intent_ref`, `spec_ref`, `plan_ref`, `nodes` e `gate`.
* [ ] **G4 (Stop-Gate 02 Travado):** O texto do Stop-Gate 02 exige aprovação humana da Tríade (`SPEC.md` + `PLAN.md` + `ROUTE.json`).
