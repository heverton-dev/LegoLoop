# 🎯 INTENT: Auditoria Profunda Full-Stack do LegoLoop (UI/UX, Usabilidade, Aviabilidade e Motores)

> **Status:** Proposta para Alinhamento Socrático (STOP-GATE HUMANO 01)  
> **Data:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** Full-Stack LegoLoop (`web/`, `cmd/`, `internal/`, rotas HTTP, Studio Canvas, CLI e Motores DAG)

---

## 1. Contexto e Motivação
O ecossistema LegoLoop apresenta acúmulo de inconsistências, falhas de usabilidade e bugs de diversas ordens percebidos no uso diário. É imperativo executar uma auditoria técnica profunda e autônoma, mapeando desde quebras de contratos de API, erros silenciosos de runtime/console e falhas de interação visual (cliques, drag-and-drop, modais e formulários no Studio) até fragilidades nos motores de execução e CLI.

---

## 2. Objetivos Principais
1. **Engenharia Reversa Estática & Análise de Rotas:** Varredura completa da base de código do frontend (`web/` e componentes de UI) e do backend HTTP/CLI em Go, mapeando contratos e potenciais fontes de falhas.
2. **Auditoria Dinâmica E2E Autônoma via Playwright/DevTools:** Testar interações reais do Studio UI (canvas, nós, menus, botões, modais, submissões) com captura de logs de console, erros de rede e comportamentos inesperados.
3. **Preservação de Armazenamento:** Proibição de dump indiscriminado de screenshots e capturas de tela redundantes, gravando apenas evidências estritamente necessárias para correções internas.
4. **Inventário Abrangente de Bugs & Relatório Diagnóstico:** Catalogar todos os defeitos encontrados por severidade (Crítico, Alto, Médio, Baixo), acompanhados da causa-raiz exata e recomendação técnica.

---

## 3. Escopo Detalhado

### O que entra:
* **Studio Web UI:**
  * Usabilidade e aviabilidade de botões, barras de ferramentas, painéis laterais e modais.
  * Drag-and-drop, conexão e manipulação de nós no Canvas n8n-style.
  * Validação de formulários, inputs de configuração e tratamento de estados vazios/erro.
  * Erros de JavaScript/TypeScript no console e warnings do navegador.
* **Backend & Contratos de API:**
  * Integridade dos endpoints HTTP locais consumidos pelo frontend.
  * Tratamento de payloads incorretos e códigos de status HTTP.
* **CLI & Motores de Execução:**
  * Validador de DAG (`legoloop validate`).
  * Executor autônomo e paralelismo (`legoloop run`).
  * Diagnóstico de ambiente (`legoloop doctor`).
  * Desduplicação canônica (`legoloop canonize`).
* **Relatório Final Consolidado:**
  * Artefato `docs/reports/2026-11-01_REPORT_AUDITORIA_PROFUNDA.md`.

### O que NÃO entra:
* Modificações diretas ou refatorações de código antes da aprovação da Tríade de Planejamento (esta fase é estritamente diagnóstica e de especificação).
* Gravação excessiva de mídia/vídeo ou screenshots de apresentação que consumam espaço em disco.
* Subagentes desnecessários ou chamadas fora da esteira canônica.

---

## 4. Metodologia de Execução (Modo 100% Autônomo)
* **Fase 1 (Diagnóstico Estático):** Inspeção com `codebase-memory-mcp` e linters nos componentes de UI e rotas Go.
* **Fase 2 (Diagnóstico Dinâmico E2E):** Subida do Studio local em porta de teste e automação via Playwright/DevTools acionando cada botão, modal e fluxo do Canvas.
* **Fase 3 (Validação CLI & Motores):** Execução de testes de estresse e suites automatizadas no CLI (`doctor`, `validate`, `run`, `canonize`).
* **Fase 4 (Consolidação):** Emissão da Matriz de Defeitos e Relatório Geral.

---

## 5. Critérios de Conclusão do Loop 1
* [x] Entrevista socrática realizada e respostas compiladas.
* [x] Documento de Intenção formalizado com padrão ISO.
* [ ] **Aprovação explícita do usuário no Stop-Gate 01.**
