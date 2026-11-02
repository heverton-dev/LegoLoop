# 🎯 INTENT: Correção Canônica dos Defeitos Críticos e de Usabilidade do LegoLoop

> **Status:** Proposta para Alinhamento (STOP-GATE HUMANO 01)  
> **Data:** 2026-11-02  
> **Autor:** Antigravity AI + Usuário  
> **Referência:** [2026-11-01_REPORT_AUDITORIA_PROFUNDA.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-01_REPORT_AUDITORIA_PROFUNDA.md)  
> **Alvo:** `internal/studio/` (`server.go`, `static/app.js`, `static/style.css`, `static/index.html`, `static/components/lego-item-modal.js`)

---

## 1. Contexto e Motivação
A auditoria profunda comprovou 10 defeitos técnicos no LegoLoop. O objetivo desta esteira é executar a correção cirúrgica com TDD estrito e zero stubs, priorizando o desbloqueio da execução da esteira (P0), eliminação dos 15 erros de CSP (P1), serving do script de ícones (P1) e ativação do botão de exportação PNG (P1).

---

## 2. Objetivos Principais
1. **Desbloqueio do Botão "Executar" (BUG-RUN-01 - P0):** Garantir que `validate()` retorne `data.ok === true` para que `initRunModal()` abra o modal de execução e dispare a esteira.
2. **Registro de Rota de Ícones & Alias Index (BUG-SRV-01 / BUG-SRV-02 - P1/P2):** Mapear `GET /icons.js` e `GET /index.html` em `server.go`, eliminando erros 404 e restauração de ícones SVG.
3. **Eliminação de 15 Violações de CSP (BUG-CSP-01 - P1):** Migrar estilos inline de `index.html` e `lego-item-modal.js` para classes dedicadas em `style.css`.
4. **Ativação da Exportação de Imagem PNG (BUG-UI-01 - P1):** Implementar handler no botão `#btn-export-png` para exportar o Canvas visual.
5. **Correção de Fallbacks e Endpoints de Erro (BUG-TMPL-01 / BUG-API-01 - P2):** Ajustar template fallback e homogeneizar respostas REST em JSON.
6. **Recompilação Fresca Obrigatória:** Recompilar `bin/legoloop.exe` ao final do ciclo.

---

## 3. Escopo Detalhado

### O que entra:
* **ÉPICO 01 (Correções Críticas P0 & P1):**
  * `app.js`: Retorno booleano em `validate()`.
  * `server.go`: Inclusão de `GET /icons.js` e `GET /index.html`.
  * `style.css` + `index.html` + `lego-item-modal.js`: Eliminação de `style="..."` inline para conformidade CSP total.
  * `app.js`: Implementação de exportação PNG do canvas em `#btn-export-png`.
* **ÉPICO 02 (Polimento P2 & P3):**
  * Correção de duplicidade em nomes de templates (`app.js`).
  * Padronização de JSON de erro (`writeError`) em `assets_api.go` e `loops_api.go`.
  * Adição de favicon inline data URI.
* **Testes Automatizados:** Testes em Go e suites Node/Playwright validando cada correção antes e depois.

### O que NÃO entra:
* Refatoração arquitetural de outros subsistemas de motor (DAG/Runner) que foram validados como 100% íntegros.
* Uso de subagentes ou polling.

---

## 4. Critérios de Conclusão do Loop 1
* [x] Documento de Intenção formalizado com prefixo ISO.
* [ ] **Aprovação explícita do usuário no Stop-Gate 01.**
