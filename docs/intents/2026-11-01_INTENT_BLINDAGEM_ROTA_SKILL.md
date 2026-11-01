# 🎯 INTENT: Blindagem da Rota Formal no `aidd-legoloop-mega-pipeline`

> **Status:** Aprovado pelo Usuário  
> **Data de Aprovação:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** `C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`

---

## 1. Contexto e Fato Diagnosticado

Em teste prático realizado com um modelo na sessão seguinte, foi constatada uma falha comportamental imediata:
* **Pulo da Etapa de Rota:** O modelo gerou `INTENT.md`, `SPEC.md` e `PLAN.md`, mas **pulou completamente a geração do mapa formal de rota** (`docs/routes/YYYY-MM-DD_ROUTE_*.json`) e tentou disparar o `BUILD` diretamente.
* **Causa Raiz:** A regra #4 do [`AGENTS.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/AGENTS.md) exige que *"Toda execução pós-plano deve seguir o mapa formal em `docs/routes/` com execução estéril por bloco"*. No entanto, na skill `aidd-legoloop-mega-pipeline`, nós omitimos a Rota da lista de artefatos obrigatórios do Loop 3 e não incluímos o seu template canônico em JSON nem a trava de gate.

---

## 2. Objetivos Principais

1. **Instituir a Rota Formal (`docs/routes/YYYY-MM-DD_ROUTE_<TEMA>.json`) como Artefato Mandatório:**
   * Nenhum projeto pode transicionar do planejamento para o Build sem um mapa formal de rota em JSON.
2. **Adicionar o Template Canônico de Rota (`ROUTE.json`) no `SKILL.md`:**
   * Padronizar a estrutura com `intent_ref`, `spec_ref`, `plan_ref`, `execution_strategy` e array de `nodes` (com `nature`, `depends_on`, `routing`, `inputs`, `output`, `gate`).
3. **Reforçar o Stop-Gate Humano do Loop 3:**
   * A condição binária para liberação do Loop 4 (Build) passa a ser: **`SPEC.md` + `PLAN.md` + `ROUTE.json` 100% aprovados pelo usuário**.
   * Se o agente tentar codificar sem o `ROUTE.json` no disco, ele deve ser imediatamente interrompido por sanção contratual.

---

## 3. Escopo

* **O que entra:**
  * Atualização da Seção 1 (Leis Fundamentais) para incluir o Mapa Formal de Rota.
  * Atualização do Diagrama Mermaid e Seção 3 (Loop 3 passa a entregar SPEC, PLAN e ROUTE).
  * Inclusão do Template 4 (`docs/routes/YYYY-MM-DD_ROUTE_<TEMA>.json`) na Seção 4 (Templates Canônicos).
  * Atualização do Stop-Gate 02 para exigir expressamente a aprovação da Rota.
* **O que NÃO entra:**
  * Alterações no código Go (permanece para a fase subsequente).

---

## 4. Alinhamento com o Usuário

* Você concorda que o mapa de rota deve ser gerado no formato **JSON estrito** (conforme o padrão já existente no LegoLoop em `docs/routes/`) e ser um requisito obrigatório para liberar o Build?
