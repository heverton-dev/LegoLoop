---
name: legoloop-impeccable-design
description: Use quando precisar executar o processo impecable-design / legoloop-impeccable-design de design system, UI/UX de alta fidelidade e acessibilidade.
consumes: Artifact:Idea
produces: Artifact:Design
artifact: DESIGN.md
gate: file_exists('DESIGN.md')
suggested_model: sonnet
---

# 🧩 legoloop-impeccable-design — Design de Alta Fidelidade Impeccable

> **Processo:** `impeccable-design`  
> **Versão:** 1.0.0  
> **Referência:** `processos/impeccable-design/processo.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for necessário auditar, construir, refinar e documentar interfaces de usuário (UI/UX), design tokens, sistemas de cor e acessibilidade no padrão Impeccable.  
**Quando NÃO usar:** Para Tickets puramente de backend ou sem interface visual.

* **Consumes:** Especificações em `docs/specs/` ou telas existentes.
* **Produces:** Relatórios de auditoria, direção visual, tokens CSS, validação responsiva e documentação do design system em `DESIGN.md`.

---

## 2. Estrutura de Execução do Grafo DAG

O processo engloba 6 subprocessos modulares e sequenciais:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `subprocesso-avaliar` | Avaliação de UI/UX e Diagnóstico | `skills/impeccable` | `docs/reports/AUDITORIA_DESIGN.md` |
| `subprocesso-construir` | Construção da Direção Visual | `skills/impeccable` | `docs/specs/DIRECAO_VISUAL.md` |
| `subprocesso-melhorar` | Design Tokens & Sistema de Cores | `skills/impeccable` | `src/styles/theme.css` |
| `subprocesso-refinar` | Refinamento & Hardening UI | `skills/impeccable` | `docs/reports/HARDENING_UI.md` |
| `subprocesso-corrigir` | Correções de Acessibilidade & Responsividade | `skills/impeccable` | `docs/reports/RESPONSIVE_VERIFICATION.md` |
| `subprocesso-iterar-documentar` | Documentação do Design System | `skills/impeccable` | `DESIGN.md` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Avaliação de UI/UX e Diagnóstico (`subprocesso-avaliar`):** Avaliar heurísticas de design, contraste e hierarquia visual, registrando o diagnóstico em `docs/reports/AUDITORIA_DESIGN.md`.
2. **Construção da Direção Visual (`subprocesso-construir`):** Definir pilares estéticos, tipografia e espaçamentos em `docs/specs/DIRECAO_VISUAL.md`.
3. **Design Tokens & Sistema de Cores (`subprocesso-melhorar`):** Estabelecer paleta OKLCH/Hex e variáveis CSS em `src/styles/theme.css`.
4. **Refinamento & Hardening UI (`subprocesso-refinar`):** Refinar micro-interações, consistência de Peças e limites visuais em `docs/reports/HARDENING_UI.md`.
5. **Correções de Acessibilidade & Responsividade (`subprocesso-corrigir`):** Validar WCAG 2.1 AA/AAA e layouts multi-viewport em `docs/reports/RESPONSIVE_VERIFICATION.md`.
6. **Documentação do Design System (`subprocesso-iterar-documentar`):** Consolidar guidelines e especificações canônicas no manifesto final `DESIGN.md`.

---

## 4. Invariantes & Regras Negativas

- **Zero HTML em Visualizações:** Respeitar regras de marcação limpa sem injeções incorretas de HTML.
- **Tokens Canônicos:** Não utilizar cores e medidas mágicas (hardcoded); usar sempre variáveis e tokens.
- **Portão de Saída:** Validação mecânica com verificação obrigatória de integridade de arquivos.
