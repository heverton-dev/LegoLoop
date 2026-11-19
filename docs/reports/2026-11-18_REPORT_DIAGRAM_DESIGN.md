# Relatório de Design do Diagrama: LegoLoop Diagram Design

**Data:** 2026-11-18
**Slug:** `legoloop-diagram-design`
**Artefato fonte:** `docs/diagrams/legoloop-diagram-design/export/legoloop-diagram-design.svg`

## 1. Escopo

Este relatório registra o estado do diagrama de seis nós do processo `legoloop-diagram-design`. Cobre o SVG exportado, os anticorpos registrados e os critérios verificados nesta sessão.

## 2. Estrutura do diagrama

| Nó | Rótulo | Papel visual |
|---|---|---|
| e01 | Inicializar slug | Etapa |
| e02 | Selecionar tipo | Etapa |
| e03 | Gerar SVG e HTML | Foco editorial |
| e04 | Animar em CSS | Etapa |
| e05 | Exportar e validar | Foco editorial |
| e06 | Fechar e relatar | Etapa |

O fluxo segue um circuito de cinco setas. Não há setas de retorno.

## 3. Verificação do orçamento visual

| Critério | Limite | Medido | Resultado |
|---|---|---|---|
| Nós | 9 | 6 | Aprovado |
| Setas | 12 | 5 | Aprovado |
| Elementos coral | 2 | 2 (e03 e e05) | Aprovado |
| Callouts de anotação | 2 | 0 | Aprovado |

## 4. Verificação de acessibilidade e tokens

| Critério | Resultado | Prova |
|---|---|---|
| `role="img"` e `aria-labelledby` no `<svg>` | Presente | Linha 2 do SVG |
| `<title>` e `<desc>` preenchidos | Presentes | IDs com prefixo `legoloop-diagram-design-` |
| `<title>` antes de `<defs>` | Correto | Ordem dos elementos no SVG |
| Fonte serifada | Ausente | Só Geist e Geist Mono |
| Emojis ou símbolos decorativos | Ausentes | Nenhum caractere no SVG |
| Conectores ortogonais | Aprovado por leitura | Todas as linhas são horizontais ou verticais |
| Legenda em faixa inferior | Presente | Linha 360 do SVG |

## 5. Critérios não verificados nesta sessão

- `scripts/verify_geometry.py`: não executado.
- `scripts/self_check.py`: não executado.
- Gate `file_exists` do relatório: não verificado.
- Gate do nó e05 (exit code 0): não executado.

## 6. Desvios conhecidos

1. **Animação ausente.** O SVG não contém `@keyframes`, `stroke-dasharray` nem `stroke-dashoffset`. O nó e04 descreve uma animação que o arquivo não tem.
2. **Cores em hex literal.** O SVG usa `#f5f5f5`, `#2d3142` e `#bf4520` em vez de `var(--paper)`, `var(--ink)` e `var(--accent)`. Os valores coincidem com os tokens, mas a regra de tokens canônicos não é cumprida.
3. **Legenda sem formas geométricas.** Todos os nós são retângulos. O anticorpo `3d4e5f6a` cita cinco formas, e o SVG não as usa.
4. **Anticorpos de outra versão.** Os erros citam "Camada 2", "fabricar" e um card "fechar" com nove etapas. Nenhum desses elementos existe neste SVG. Por isso, o campo `antidoto_aplicado: true` não pode ser confirmado para este artefato.
5. **Entradas duplicadas.** O mesmo conjunto de dez anticorpos aparece em dois arquivos de contexto.
6. **Dependência de rede para fontes.** O `@import` do Google Fonts carrega Geist. Sem rede, o texto cai para `system-ui`, o que muda a métrica de largura.
7. **Classes sem uso.** `.lead` e `.diagram-wrap` estão declaradas, mas não há elementos correspondentes no SVG.

## 7. Pendências

1. Adicionar `@keyframes flow-dash` com `stroke-dasharray` e `stroke-dashoffset` nas setas, ou remover a animação do nó e04.
2. Trocar os hex literais por `var(--paper)`, `var(--ink)` e `var(--accent)`.
3. Regenerar os anticorpos a partir do SVG atual, ou marcar as entradas atuais como de versão anterior.
4. Executar `scripts/verify_geometry.py` e `scripts/self_check.py` sobre o SVG.
5. Remover a menção a cinco formas geométricas, ou desenhar as formas na legenda.

## 8. Status

Entrega parcial. O diagrama tem estrutura válida, acessibilidade e orçamento visual dentro dos limites. Há desvios de animação, tokens e anticorpos, listados nas seções 6 e 7.