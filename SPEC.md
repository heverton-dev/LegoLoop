# 🧩 LegoLoop — Especificação Arquitetural Unificada (Master Spec v1.0)

> **Projeto:** LegoLoop (The Composable Looping & Harness Fabric)  
> **Status:** Aprovado para Implementação  
> **Autor / Arquitetura:** Ecossistema AIDD + Symbiont Engine + 9Router Gateway  
> **Data:** 30/10/2026  

---

## 1. Visão Executiva & Proposta de Valor

O **LegoLoop** é a primeira plataforma unificada de **Harness & Looping Engineering** que transforma skills, MCPs, scripts e modelos de IA em blocos composíveis ("peças Lego") no estilo visual **N8N / Make**, eliminando o inchaço de contexto e o desperdício de tokens através de **Orchestrix Cirúrgicas Estéreis**.

Ele unifica três pilares fundamentais já construídos no ambiente:
1. **Ecossistema AIDD:** As regras invioláveis (Zero Stubs, Economia Severa, TDD Estrito) e o baú de skills canônicas.
2. **SYMBIONT:** O motor monólito em Go de altíssima performance, o scanner paralelo de 32 threads (que mapeia 22k+ ativos no SO) e o canvas visual de esteiras.
3. **9Router:** O gateway cognitivo blindado contra incompatibilidades de schema (Gemini v1beta) e estouro de contexto em failover (>100k tokens).

---

## 2. Os Quatro Pilares de Engenharia do LegoLoop

```
                       ┌──────────────────────────────────────────────┐
                       │           LEGOLOOP: A FORJA TOTAL            │
                       │   (Ecossistema AIDD + Symbiont + 9Router)    │
                       └──────────────────────┬───────────────────────┘
                                              │
         ┌───────────────────┬────────────────┴──────────────────┬───────────────────┐
         ▼                   ▼                                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐                 ┌─────────────────┐ ┌─────────────────┐
│     HARNESS     │ │     LOOPING     │                 │      GRAPH      │ │     ROUTER      │
│   ENGINEERING   │ │   ENGINEERING   │                 │   ENGINEERING   │ │   ENGINEERING   │
├─────────────────┤ ├─────────────────┤                 ├─────────────────┤ ├─────────────────┤
│ • Blindagem CLIs│ │ • N8N Studio    │                 │ • Grafo DAG     │ │ • Orchestrix por│
│ • Claude Code,  │ │ • Blocos Lego   │                 │ • Matriz Nó/Port│ │   Bloco Estéril │
│   AGY, OpenCode │ │ • Handoff Fixo  │                 │ • Rota ótima    │ │ • Provider-Aware│
│ • Canonização SO│ │ • Gates Pass/   │                 │ • Compatibilida-│ │ • 9Router Combos│
│   via Symlinks  │ │   Fail estritos │                 │   de de schemas │ │ • Anti-Subagente│
└─────────────────┘ └─────────────────┘                 └─────────────────┘ └─────────────────┘
```

---

### Pilar 1: Harness Engineering & A Forja do SO
* **Canonização Transparente (Symlink Federation):** Desduplica os 22.000 ativos espalhados nos discos (C:, D:, AppData) para um cofre central (`~/.aidd/canonical/`), substituindo duplicatas por links simbólicos (`mklink`). Zero quebra de projetos legados e 80% de economia de disco/RAM.
* **Universalidade de Harnesses:** Suporte nativo a 13 ambientes locais (Claude Code, Antigravity, OpenCode, Cursor, Codex, terminal CLI puro).
* **Comando de Forja (`legoloop forge`):** Inicializa o ambiente em 3 segundos, checa dependências essenciais e aplica as regras do harness.

### Pilar 2: Looping Engineering & Canvas N8N
* **Anatomia do Bloco Lego:**
  * **Entrada Tipada:** O bloco só aceita entradas que respeitem o schema exigido.
  * **Handoff Imutável (Saída Rígida):** O artefato gerado é fixo (ex: `BRIEFING.md`, `SPEC.md`, `ESTRUTURA_REVERSA.json`).
  * **Gates de Qualidade (Pass/Fail):** A esteira só transiciona para o próximo nó se a condição de saída for 100% satisfeita.
* **Composabilidade Livre:** O usuário pode alternar entre `interview` ou `aidd-grill` na concepção; pode conectar um sub-loop de Engenharia Reversa (`Obscura + REA + Playwright`) ou pular etapas.

### Pilar 3: Graph Engineering (A Malha Semântica)
* **Vértices e Arestas:** Cada skill, MCP e script é um nó do Grafo de Dependências com tipos de consumo e produção.
* **Cálculo de Rota Ótima:** Se o desenvolvedor define a entrada `URL Web` e o objetivo `Código Limpo em Go`, o Graph Engine traça automaticamente o menor caminho de blocos válidos.
* **Validação Semântica de Conexão:** Alerta em tempo real no canvas caso um bloco incompatível seja conectado a outro.

### Pilar 4: Router Engineering & Orchestrix Universal
* **Orchestrix Cirúrgica por Bloco (Regra Inviolável):** Nenhum bloco compartilha memória de chat com outro. Cada nó é executado em um subprocesso estéril com contexto zero (2k a 5k tokens), lê apenas o artefato do disco, gera seu output e é finalizado. Subagentes pesados e gastadores de tokens são terminantemente proibidos.
* **Provider-Awareness Passivo:** O Symbiont audita o que o usuário realmente possui (Ollama local `:11434`, chaves de API ativas, sessões ativas do Claude Code) e só sugere/permite modelos reais.
* **Gateway Cognitivo (9Router):** Se o 9Router estiver rodando, expõe combos otimizados (`code-pro`, `code-fast`), sanitiza schemas recursivamente para o Gemini v1beta e desvia contextos >100k tokens para modelos de 1M, evitando erro fatal no Kiro.

---

## 3. Formato do Contrato de Execução (`pipeline.dag.json`)

```json
{
  "$schema": "https://legoloop.dev/schemas/dag.v1.json",
  "name": "greenfield-clean-craft",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "node_01_interview",
      "type": "skill",
      "ref": "skills/interview",
      "routing": {
        "mode": "suggested",
        "model": "gemini-3.5-flash-lite",
        "fallback": "haiku"
      },
      "inputs": { "prompt": "var.idea" },
      "output": {
        "artifact": "BRIEFING.md",
        "gate": "file_exists('BRIEFING.md') && word_count('BRIEFING.md') > 100"
      }
    },
    {
      "id": "node_02_spec",
      "type": "skill",
      "ref": "skills/aidd-spec",
      "depends_on": ["node_01_interview"],
      "routing": {
        "mode": "auto",
        "combo": "code-pro"
      },
      "inputs": { "briefing": "file('BRIEFING.md')" },
      "output": {
        "artifact": "SPEC.md",
        "gate": "has_binary_criteria('SPEC.md')"
      }
    }
  ]
}
```

---

## 4. Matriz de Compatibilidade e Requisitos

| Componente | Requisito Mínimo | Requisito Recomendado |
|---|---|---|
| **OS** | Windows 10/11 x64, Linux x64, macOS ARM | Windows 11 x64 com PowerShell 7 |
| **Runtime Core** | Binário compilado `legoloop.exe` (Go ~20MB) | Zero dependências externas para rodar |
| **Navegador Headless** | Playwright Chromium | Obscura Rust Engine (~30MB RAM) |
| **Roteador Cognitivo** | Chamada direta via CLI (Claude/AGY) | 9Router Gateway local (:20128) |
