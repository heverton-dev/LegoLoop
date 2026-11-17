# 📋 SPEC: Auditoria Profunda Full-Stack do LegoLoop

> **Status:** Proposta para Revisão (Tríade de Planejamento)  
> **Data:** 2026-11-01  
> **Referência:** [2026-11-01_INTENT_AUDITORIA_PROFUNDA_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_AUDITORIA_PROFUNDA_LEGOLOOP.md)

---

## 1. Visão Geral e Arquitetura de Diagnóstico
A auditoria cobre as três camadas fundamentais do ecossistema LegoLoop:
1. **Camada de Apresentação (Studio Web UI):** `internal/studio/static/` (`index.html`, `app.js`, `style.css`, `inspector.js`, `palette.js`, `slots.js`, `components/lego-item-modal.js`, etc.).
2. **Camada de Servidor & APIs HTTP Locais:** `internal/studio/` (`server.go`, `assets_api.go`, `loops_api.go`, `run_api.go`, `pieces.go`, `templates.go`, `security.go`).
3. **Camada de CLI & Motores de Grafo:** `cmd/legoloop/`, `internal/dag/`, `internal/runner/`, `internal/canon/`, `internal/doctor/`.

---

## 2. Contratos de Dados & I/O de Auditoria

### A. Matriz de Severidade de Bugs:
* **Crítico (P0):** Crash do processo, quebra total de renderização (tela branca), perda de dados ou deadlock de execução.
* **Alto (P1):** Botão ou funcionalidade essencial inoperante, rota de API retornando 500 ou quebrando contrato JSON, nós do Canvas impossíveis de conectar ou editar.
* **Médio (P2):** Erros/warnings não fatais no console, falha de layout/sobreposição de elementos, campos de formulário sem validação adequada, inconsistências de estilo.
* **Baixo (P3):** Polimento estético, labels truncadas, falta de feedback tátil/cursor incorreto, melhorias menores de usabilidade.

### B. I/O dos Scripts de Diagnóstico:
* **Entrada:** Código-fonte estático e instância ativa do Studio em `http://127.0.0.1:7878` (ou porta isolada de teste).
* **Saída:** Logs estruturados de console, mapa de rotas testadas com códigos HTTP e artefato final consolidado em `docs/reports/2026-11-01_REPORT_AUDITORIA_PROFUNDA.md`.

---

## 3. Critérios Binários de Aceite (Gates Mecânicos)
* [ ] **G1 (Varredura Estática de Frontend):** Inspeção completa de todos os arquivos JS/CSS/HTML em `internal/studio/static/` detectando referências nulas a DOM IDs, listeners sem elemento alvo, handlers quebrados e classes CSS órfãs.
* [ ] **G2 (Auditoria de Contratos de API):** Cobertura de 100% dos endpoints REST do Studio (`/api/pieces`, `/api/loops`, `/api/assets`, `/api/run`, `/api/templates`, etc.) verificando status HTTP, formato de erro padronizado e headers CORS/segurança.
* [ ] **G3 (Auditoria Dinâmica E2E via Playwright):** Execução autônoma de suíte de testes de UI acionando: Canvas (zoom, pan, adicionar nó, conectar nós, mover nós), Paleta de Componentes, Modal de Detalhes, Inspector Lateral, Barra Superior e Ações de Salvar/Exportar.
* [ ] **G4 (Auditoria de Logs & Runtime UI):** Zero exceções JavaScript não tratadas (`Uncaught TypeError`, `UnhandledPromiseRejection`) durante o fluxo padrão de operação do Studio.
* [ ] **G5 (Auditoria de Comandos e Motores CLI):** Teste determinístico dos comandos `legoloop doctor`, `validate`, `run` e `canonize` com verificação de saídas normais e códigos de retorno de erro.
* [ ] **G6 (Relatório Factual Consolidado):** Emissão de relatório técnico estruturado com catálogo detalhado de cada defeito, reprodução exata, impacto e plano de correção sugerido.
