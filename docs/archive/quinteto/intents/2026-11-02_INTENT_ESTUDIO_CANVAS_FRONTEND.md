# INTENT — Dissecação, Auditoria e Correção do Studio & Canvas (Front-End / UI / UX)

**Data:** 2026-11-02  
**Autor:** Antigravity + Usuário  
**Status:** Aprovado em Interview-Me  

---

## 1. Contexto e Motivação
Após uma bateria preliminar de engenharia reversa no LegoLoop, o acesso real à plataforma demonstrou a permanência de falhas críticas de Front-End, Interface de Usuário (UI) e Experiência do Usuário (UX) no LegoLoop Studio.
Problemas observados diretamente no uso:
1. **Canvas Viewport e Scroll:** Erros de rolagem e colisão entre scroll nativo e pan/zoom virtual, descalibrando ponteiros e nós.
2. **Ciclo de Vida do Card (Ligar / Desligar):** Ausência ou inoperância do toggle de ativação/desativação (bypass/enabled) nos cards de peças no canvas.
3. **Modais e Diálogos:** Falha de abertura, fechamento e visualização de itens em modais críticos (Execução, Templates, Cofre).
4. **Violações de CSP (Content Security Policy):** Bloqueio estrito no navegador de estilos inline necessários para transformações 2D do canvas.
5. **Transparência e Factualidade:** Necessidade de conduzir os testes de forma visível e factual, eliminando suposições e relatórios falsos.

---

## 2. Escopo Delimitado
- **Foco Estrito:** Front-End do LegoLoop Studio (`internal/studio/static/`) e contratos de segurança imediatos (`internal/studio/security.go`).
- **Ondas Posteriores Desacopladas:** Motores de execução backend (Onda 2 e Onda 3) só serão abordados após a estabilização completa do Studio & Canvas.

---

## 3. Resultados Esperados
1. Documentação canônica completa: `INTENT`, `SPEC`, `PLAN` e `ROUTE`.
2. Correção definitiva do scroll/pan do canvas com eliminação de desvios de coordenadas.
3. Implementação e funcionamento real do botão/toggle de Ligar/Desligar (Ativar/Desativar) nos cards do canvas.
4. Ajuste do Content Security Policy (CSP) para não travar transformações CSS e animações.
5. Funcionamento 100% verificado de todos os modais da interface.
6. Validação automatizada e visual com Playwright visível.
