# 📘 MANUAL OPERACIONAL: FUNÇÃO MEGA PIPELINE NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Orquestração determinística de Épicos em 5 sub-loops modulares com contexto estéril e gates de qualidade  
> **Fontes de verdade:** `skills/legoloop-mega-pipeline/SKILL.md` · `docs/templates/` · `docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md`

---

## 1. Visão Geral

A função **`legoloop-mega-pipeline`** é o orquestrador de engenharia do LegoLoop. Ela decompõe demandas complexas em uma esteira sequencial inviolável de 5 sub-loops (`Concepção` ➔ `Engenharia Reversa / Design` ➔ `Especificação & Tríade` ➔ `Construção TDD por Épico` ➔ `Homologação E2E & Fechamento`), executados em contexto estéril por bloco (*Orchestrix cirúrgica*) e supervisionados por Stop-Gates humanos.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop`
* `/mega-pipeline`
* `/legoloop-orchestrator`
* `/canonize`

### 2.2 Pela CLI Mecânica
```bash
# Auditar ambiente e provedores disponíveis
legoloop doctor --json

# Validar DAG, contratos e gates de um pipeline
legoloop validate pipeline.dag.json

# Executar pipeline em modo autônomo com workers paralelos
legoloop run pipeline.dag.json --parallel 4 --timeout 30m --json

# Abrir o Studio Canvas
legoloop studio --addr 127.0.0.1:7878
```

---

## 3. As 7 Leis Invioláveis do Mega Pipeline

1. **Zero Context Bloat (Orchestrix Cirúrgica por Bloco):** Cada bloco nasce e morre em isolamento estéril.
2. **Stop-Gates de Auditoria Humana Obrigatórios:** Parada mandatória entre cada fase e entre cada Épico.
3. **Plano Único e Mapa Formal de Rota:** A Tríade `SPEC.md` + `PLAN.md` + `ROUTE.json` é pré-requisito absoluto do código.
4. **Provider-Awareness & Roteamento Cognitivo com Clamping:** Roteamento por complexidade com clamping de esforço (`--effort low|medium|high`).
5. **TDD Estrito & Zero Stubs:** Red-Green obrigatório, zero `TODO` ou mocks simulados.
6. **Nomenclatura Canônica com Prefixo ISO:** Todos os artefatos com `YYYY-MM-DD_<NOME>.<ext>`.
7. **Idioma PT-BR & Executor Silencioso:** Comunicação técnica direta, telegráfica e sem repetições.
