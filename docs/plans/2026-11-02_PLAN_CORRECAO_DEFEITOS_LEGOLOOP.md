# 🗺️ PLAN: Correção Canônica de Defeitos e Usabilidade do LegoLoop

> **Status:** Proposta para Revisão (Tríade de Planejamento)  
> **Data:** 2026-11-02  
> **Especificação Base:** [2026-11-02_SPEC_CORRECAO_DEFEITOS_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-02_SPEC_CORRECAO_DEFEITOS_LEGOLOOP.md)

---

## 🏛️ ÉPICO 01: Correções Críticas P0 e P1 (Servidor, Runtime e CSP)
> **Objetivo:** Desbloquear a execução da esteira no Studio, servir arquivos essenciais de ícones, remover todas as violações de segurança CSP e tornar o botão de exportação funcional.  
> **Modelo Recomendado:** Claude 3.5 Sonnet / Gemini 2.0 Flash (--effort medium)

### 🔹 Micro-ticket FIX-01.1: Rotas Estáticas em `server.go` (`icons.js` e `/index.html`)
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/server.go`, `internal/studio/static_test.go`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **TDD:** Adicionar testes em `static_test.go` esperando 200 OK para `/icons.js` e `/index.html` (Red), em seguida atualizar `server.go` (Green).

### 🔹 Micro-ticket FIX-01.2: Retorno Booleano de `validate()` em `app.js`
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/static/app.js`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Ajustar `validate()` para `return Boolean(data.ok)`, permitindo que `initRunModal()` abra `#modal-run`.

### 🔹 Micro-ticket FIX-01.3: Eliminação de Estilos Inline & Classes em `style.css`
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/static/style.css`, `internal/studio/static/index.html`, `internal/studio/static/components/lego-item-modal.js`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Criar classes `.is-hidden`, `.btn--compact`, etc., substituindo todos os 15 atributos `style="..."` inline para cumprir a política CSP `style-src 'self'`.

### 🔹 Micro-ticket FIX-01.4: Implementação da Exportação PNG (`#btn-export-png`)
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/static/app.js`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Anexar listener em `#btn-export-png` para capturar a renderização visual dos nós do Canvas e disponibilizar o download do PNG.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 01
* [ ] Testes do Épico 01 verdes comprovados.
* [ ] Validação do usuário para avançar ao Épico 02.

---

## 🏛️ ÉPICO 02: Polimento P2/P3, Homogeneização REST e Recompilação
> **Objetivo:** Sanar duplicações visuais, padronizar códigos de retorno HTTP e recompilar o binário oficial.  
> **Modelo Recomendado:** Gemini 2.0 Flash (--effort medium)

### 🔹 Micro-ticket FIX-02.1: Fallback de Templates & Favicon Inline
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/static/app.js`, `internal/studio/static/index.html`, `internal/studio/static/inventario.html`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Evitar duplicação do nome do template em fallback e adicionar link rel="icon" com data URI SVG.

### 🔹 Micro-ticket FIX-02.2: Padronização de JSON de Erro REST
* **Tipo:** Micro-ticket
* **Alvo:** `internal/studio/assets_api.go`, `internal/studio/loops_api.go`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** Substituir `http.Error` por `writeError` para manter formato JSON uniforme.

### 🔹 Micro-ticket FIX-02.3: Revalidação Geral & Recompilação Obrigatória
* **Tipo:** Micro-ticket
* **Alvo:** `bin/legoloop.exe`
* **Rota:** Gemini 2.0 Flash (--effort medium)
* **Entregáveis:** `go test ./...` verde e execução de `go build -o bin/legoloop.exe ./cmd/legoloop`.

### 🛑 STOP-GATE FINAL: Entrega & Fechamento
* [ ] Demonstração de que todas as correções funcionam na prática.
* [ ] Aprovação do usuário.
