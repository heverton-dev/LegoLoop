# 📋 SPEC: Correção Canônica de Defeitos e Usabilidade do LegoLoop

> **Status:** Proposta para Revisão (Tríade de Planejamento)  
> **Data:** 2026-11-02  
> **Referência:** [2026-11-02_INTENT_CORRECAO_DEFEITOS_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_CORRECAO_DEFEITOS_LEGOLOOP.md)

---

## 1. Visão Geral e Arquitetura de Correção

A implementação das correções atua cirurgicamente nas seguintes camadas:
1. **Camada de Servidor Go (`internal/studio/server.go`, `internal/studio/assets_api.go`, `internal/studio/loops_api.go`):**
   - Roteador estático do ServeMux: inclusão de `"GET /icons.js"` e `"GET /index.html"`.
   - Homogeneização de status de erro com `writeError(w, code, msg)`.
2. **Camada de Script & Lógica de UI (`internal/studio/static/app.js`):**
   - Função `validate()`: retornar explicitamente `Boolean(data.ok)` para destravar `initRunModal()`.
   - Botão `#btn-export-png`: implementação de exportação do canvas para SVG/PNG via canvas nativo do navegador.
   - Fallback de templates: uso de placeholder neutro em vez de duplicação de texto.
3. **Camada de Estilo & Conformidade CSP (`internal/studio/static/style.css`, `index.html`, `lego-item-modal.js`):**
   - Criação de classes CSS utilitárias para substituir todos os atributos `style="..."` inline, eliminando as 15 violações de CSP.
   - Inclusão de Favicon SVG embutido em `index.html` e `inventario.html`.

---

## 2. Contratos & Critérios Binários de Aceite (Gates Mecânicos)

* [ ] **G1 (Desbloqueio de Execução):** A função `validate()` retorna `true` quando `data.ok === true`. Um teste unitário em `app.js` ou Playwright confirma que o modal `#modal-run` abre ao clicar em `#btn-run`.
* [ ] **G2 (Servir `icons.js` & `/index.html`):** Testes unitários em `server_test.go` confirmam que `GET /icons.js` responde `200 OK` com `Content-Type: text/javascript` e `GET /index.html` responde `200 OK` com `Content-Type: text/html`.
* [ ] **G3 (Zero Violações de CSP):** Varredura automatizada com Playwright na porta de teste confirma que o boot do Studio emite zero erros de CSP (`style-src`). Nenhum atributo `style="..."` permanece em `index.html` ou `lego-item-modal.js`.
* [ ] **G4 (Exportação PNG Ativa):** O elemento `#btn-export-png` possui manipulador de evento anexado; o clique gera um arquivo de imagem para download.
* [ ] **G5 (Padronização JSON):** Endpoints de erro em `assets_api.go` e `loops_api.go` respondem com `Content-Type: application/json` e payload contendo a chave `"error"`.
* [ ] **G6 (Suíte Verde & Recompilação):** `go test ./...` passa com 100% de sucesso e o binário `bin/legoloop.exe` é recompilado com timestamp atualizado.
