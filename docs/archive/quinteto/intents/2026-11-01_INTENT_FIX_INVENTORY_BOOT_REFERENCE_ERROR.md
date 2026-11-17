# 🎯 INTENT: Correção do ReferenceError no Boot do Inventário do Studio

> **Status:** Proposta para Alinhamento Socrático e Validação Humana (Stop-Gate 01)  
> **Data:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** Frontend do LegoLoop Studio ([`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js))

---

## 1. Contexto e Motivação

Ao carregar a página de inventário (`/inventario`) do LegoLoop Studio, a interface falha prematuramente exibindo um toast de erro vermelho:
`"Não foi possível iniciar o inventário: translateActiveItem is not defined"`.

### Causa Raiz Factual (Dissecação REA)
* O modal de inspeção de peças foi refatorado para o Web Component canônico [`<lego-item-modal>`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/components/lego-item-modal.js), que já encapsula internamente seus botões de fechar, abas WAI-ARIA e a ação de tradução via método de instância `this.translateActiveItem()`.
* No entanto, na função `boot()` de [`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js#L351-L357), permaneceu um bloco residual órfão de listeners que tenta referenciar `translateActiveItem` e `switchModalTab` no escopo global.
* Por não existirem tais identificadores no escopo léxico do arquivo, o interpretador JavaScript lança um `ReferenceError`, capturado pelo `.catch()` de inicialização que interrompe toda a execução do inventário e bloqueia a exibição das peças.

---

## 2. Objetivos Principais

1. **Eliminação do ReferenceError no Boot:**
   * Remover os listeners legados e órfãos do modal em `inventory.js`, assegurando que a inicialização `boot()` ocorra de forma limpa, sem exceções.
2. **Preservação Integral das Capacidades do Web Component:**
   * Assegurar que o componente `<lego-item-modal>` continue funcionando 100% de forma autônoma, respondendo à abertura via `openItemModal(item)`, com todas as abas, tradução e canonização operacionais.
3. **Validação e Conformidade:**
   * Criar teste automatizado para validar que os scripts do Studio inicializam sem referências órfãs e que a rota `/api/assets` é consumida sem erros no frontend.
4. **Governança Canônica:**
   * Conduzir o fluxo completo de ponta a ponta: `INTENT` ➔ `SPEC` ➔ `PLAN` ➔ `ROUTE` ➔ `BUILD` (em worktree isolada) ➔ `REPORT` de homologação.

---

## 3. Escopo

### O que entra:
* Limpeza e saneamento de listeners legados em [`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js).
* Adição/atualização de testes automatizados em [`internal/studio/`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/) para garantir que o bundle estático não possua referências indefinidas.
* Execução em worktree hermética git.

### O que NÃO entra:
* Modificações no motor backend de ativos (`internal/assets/`) ou na API de tradução (`internal/studio/translate.go`), que já estão funcionais e verdes nos testes.
* Alteração visual ou de layout nas telas do Studio.

---

## 4. Questões para Validação Socrática (Stop-Gate 01)

1. **Concorda com o escopo:** Apenas sanear o bloco residual em `inventory.js` e reforçar a suíte de testes de UI estática, mantendo todo o comportamento do Web Component intacto?
2. **Critério de Aceite:** O teste de aceitação deve validar que nenhum símbolo não declarado é referenciado no boot dos scripts do Studio?
