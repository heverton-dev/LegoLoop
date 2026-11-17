# 🗺️ PLAN: Auditoria Profunda Full-Stack do LegoLoop

> **Status:** Proposta para Revisão (Tríade de Planejamento)  
> **Data:** 2026-11-01  
> **Especificação Base:** [2026-11-01_SPEC_AUDITORIA_PROFUNDA_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-01_SPEC_AUDITORIA_PROFUNDA_LEGOLOOP.md)

---

## 🏛️ ÉPICO 01: Varredura Estática de UI & Inspeção de Contratos HTTP
> **Objetivo:** Mapear inconsistências sintáticas, referências órfãs a IDs de DOM, listeners desacoplados e divergências de contratos entre frontend e endpoints REST Go.  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (--effort high)

### 🔹 Micro-ticket AUD-01.1: Varredura Estática de Código Frontend
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/static/` (`index.html`, `app.js`, `inspector.js`, `palette.js`, `slots.js`, `components/lego-item-modal.js`)
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Script de validação estática e lista de seletores DOM inconsistentes, handlers sem binding e dead code.

### 🔹 Micro-ticket AUD-01.2: Auditoria de Contratos de API REST
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/server.go`, `assets_api.go`, `loops_api.go`, `run_api.go`, `pieces.go`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Testes pontuais de contrato HTTP verificando status code, JSON schemas de resposta e parâmetros esperados.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 01
* [ ] Inventário de discrepâncias estáticas e de contratos gerado.
* [ ] Validação do usuário para avançar à auditoria dinâmica de UI.

---

## 🏛️ ÉPICO 02: Auditoria Dinâmica E2E do Studio Web UI via Playwright
> **Objetivo:** Executar o Studio em ambiente de teste e validar na prática cliques, arrastes, renderização de nós no canvas, modais, painéis laterais e console de erros.  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (--effort high)

### 🔹 Micro-ticket AUD-02.1: Boot do Studio & Monitoramento de Logs/Erros
* **Tipo:** Micro-ticket
* **Alvo:** Inicialização do Studio em porta dedicada (`127.0.0.1:7879`) e captura de logs de runtime via Playwright.
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Log estruturado de erros de JS, CSS ou requisições com falha 4xx/5xx no boot.

### 🔹 Micro-ticket AUD-02.2: Teste E2E do Canvas (Interação, Arraste e Conexões)
* **Tipo:** Micro-ticket
* **Alvo:** Canvas interativo (`app.js`, drag-and-drop de nós da paleta, movimentação de nós, conexões de portas de entrada/saída, zoom e pan).
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Evidência de usabilidade e registro de comportamentos anômalos no canvas.

### 🔹 Micro-ticket AUD-02.3: Teste E2E de Modais, Inspector e Botões de Ação
* **Tipo:** Micro-ticket
* **Alvo:** Modal de detalhes do bloco (`lego-item-modal.js`), painel lateral de inspeção (`inspector.js`), botões da barra superior (novo loop, salvar, templates, executar).
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Registro detalhado de botões inoperantes, formulários que não persistem dados ou comportamentos visuais quebrados.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 02
* [ ] Suíte dinâmica E2E concluída e mapa de falhas de interação consolidado.
* [ ] Validação do usuário para avançar à auditoria de CLI e motores.

---

## 🏛️ ÉPICO 03: Auditoria dos Motores CLI, DAG & Desduplicação
> **Objetivo:** Testar os comandos do binário CLI sob diferentes cenários de estresse e pipelines malformados.  
> **Modelo Recomendado:** Gemini 2.0 Pro (--effort high)

### 🔹 Micro-ticket AUD-03.1: Auditoria de Comandos Core (`doctor`, `validate`, `run`)
* **Tipo:** Micro-ticket
* **Alvo:** `cmd/legoloop/`, `internal/dag/`, `internal/runner/`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Relatório de comportamento do validador de DAG contra ciclos, nós desconectados e falhas de runtime.

### 🔹 Micro-ticket AUD-03.2: Auditoria de Canonização & Cofre de Credenciais
* **Tipo:** Micro-ticket
* **Alvo:** `internal/canon/`, `internal/vault/`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Validação de comportamento das rotinas de desduplicação (`canonize`) e segurança de chaves.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 03
* [ ] Comportamento do CLI auditado e falhas documentadas.
* [ ] Validação do usuário para fechamento e emissão do relatório final.

---

## 🏛️ ÉPICO 04: Matriz de Defeitos & Relatório Master
> **Objetivo:** Compilar todos os achados em um relatório estruturado de alta relevância técnica com propostas cirúrgicas de correção.  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (--effort high)

### 🔹 Micro-ticket AUD-04.1: Triagem & Matriz de Severidade
* **Tipo:** Micro-ticket
* **Alvo:** Classificação de todos os bugs por impacto (P0 a P3) e componente afetado.
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Tabela de defeitos com causas-raiz explicadas.

### 🔸 Nano-ticket AUD-04.2: Emissão do Relatório Canônico
* **Tipo:** Nano-ticket
* **Alvo:** `docs/reports/2026-11-01_REPORT_AUDITORIA_PROFUNDA.md`
* **Rota:** Gemini 2.0 Flash-Lite (--effort low)
* **Entregáveis:** Artefato final gerado e linkado para o usuário.

### 🛑 STOP-GATE FINAL: Entrega & Próximos Passos
* [ ] Relatório completo apresentado com todas as evidências.
* [ ] Avaliação do usuário para planejar a fase de correção e refatoração.
