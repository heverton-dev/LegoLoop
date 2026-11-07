# 📘 MANUAL OPERACIONAL: FUNÇÃO LOOP MAP NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Ingestão, dissecação e mapeamento de frameworks, autores e repositórios em DAGs e templates do LegoLoop  
> **Fontes de verdade:** `skills/legoloop-loop-map/SKILL.md` · `internal/studio/templates.go` · `docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md`

---

## 1. Visão Geral

A função **`legoloop-loop-map`** automatiza a transformação de metodologias consagradas de agentes e desenvolvimento com IA (ex: Matt Pocock, Addy Osmani, Impeccable) em loops modulares de primeira classe no LegoLoop, gerando manual operacional, arquivo DAG `.loop.json` e template embutido no Canvas.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop-map <autor_ou_framework> <caminho_ou_url>`
* `/loop-map <autor_ou_framework> <caminho_ou_url>`

### 2.2 Entradas Aceitas
* Nome do Autor ou Framework (ex: `Addy Osmani`, `Matt Pocock`).
* Repositório local ou URL Git contendo as diretrizes ou código.
* Prefixo de domínio canônico (ex: `addy-osmani`, `matt-pocock`).

---

## 3. Fluxo de Execução com 3 Stop-Gates

1. **Fase 1 (Dissecação Semântica):** Análise do grafo semântico e separação entre fluxo mestre e sub-loops modulares ➔ **Stop-Gate G1**.
2. **Fase 2 (Redação do Manual ISO):** Síntese da filosofia e contratos operacionais no manual `docs/manuals/` ➔ **Stop-Gate G2**.
3. **Fase 3 (Materialização DAG JSON):** Geração do `loops/<slug>.loop.json` validado pelo schema `dag.v1.json`.
4. **Fase 4 (Injeção TDD em Go & Build):** Inserção do template no Looping Studio Canvas (`internal/studio/templates.go`), testes e recompilação do binário ➔ **Stop-Gate G3**.
