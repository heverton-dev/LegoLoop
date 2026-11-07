# PARTE VI — Apêndices

## Apêndice A: Glossário

**Bounded pass (verificação delimitada).** Rodada de revisão com fronteira fixada antes de começar: uma captura em lote, um bloco de correções, no máximo mais uma checagem.

**Briefing.** Conjunto de diretrizes de quem pediu a interface: estética, paleta, época, tom. Tem prioridade sobre o gosto do agente.

**Cartão fantasma.** Cartão com borda de 1px sob uma sombra larga e suave; declara elevação duas vezes.

**craft-floor.md.** Arquivo de referência do upstream que define o piso de qualidade: o que verificar e o que recusar.

**DESIGN.md.** Especificação de design gerada a partir do código entregue.

**Gate binário.** Ponto de parada com apenas duas saídas: passa ou não passa.

**Modo de superfície.** Classificação de uma tela por objetivo: Operate, Persuade, Read ou Experience.

**Out-of-distribution.** Resultado que sai da média estatística do que um modelo produziria por padrão.

**PRODUCT.md.** Contexto durável da aplicação: público, problema e tom.

**Refinamento.** Melhoria que preserva identidade, comportamento e texto.

**Redesenho.** Novo mundo visual em que a interface antiga serve apenas de evidência.

**Sub-loop.** Grupo de comandos com um objetivo único, reutilizável de forma isolada.

## Apêndice B: Relatório de conformidade

Este apêndice compara o manual de 01/11/2026 com o estado do repositório em 07/11/2026. O objetivo é que o leitor saiba exatamente onde o texto do manual e o código coincidem.

| Ponto verificado       | Manual                        | Código                       | Situação        |
| :--------------------- | :---------------------------- | :--------------------------- | :-------------- |
| Arquivo do loop        | `impeccable-design.loop.json` | Existe em `loops/`           | Conforme        |
| Número de nós          | 6 sub-loops                   | 6 nós lineares               | Conforme        |
| Sexto nó               | Iterar                        | `subloop-iterar-documentar`  | Diverge no nome |
| Peça executora         | Não especificada              | `skills/impeccable` em todos | Informativo     |
| Template do Studio     | Ícone `palette`               | Ícone `palette`              | Conforme        |
| Nome do template       | Não especificado              | "...(6 Sub-Loops)"           | Informativo     |
| Estágios e gates P1-P4 | 5 estágios, 4 gates           | Gates `file_exists` por nó   | Diverge         |
| Nome do revisor        | `impeccabe-finish-reviewer`   | `impeccable-finish-reviewer` | Erro de grafia  |

**Leitura do relatório.** O DAG declara apenas gates de existência de arquivo por nó; os gates de julgamento P1 a P4 são protocolo da skill, conduzido no chat. No sexto nó, o código agrupa iteração e documentação, enquanto o manual trata `document` como comando do sub-loop Construir. Para o revisor, o livro usa a grafia do agente instalado.

## Apêndice C: Referência rápida de comandos

| Comando     | Sub-loop  | Função resumida                          |
| :---------- | :-------- | :--------------------------------------- |
| `critique`  | Avaliar   | Revisão heurística de UX                 |
| `audit`     | Avaliar   | Acessibilidade, performance, viewports   |
| `init`      | Construir | Cria PRODUCT.md                          |
| `shape`     | Construir | Planeja fluxos sem código                |
| `document`  | Construir | Gera DESIGN.md do código                 |
| `extract`   | Construir | Isola tokens e componentes               |
| `layout`    | Melhorar  | Espaçamento, grid, ritmo                 |
| `typeset`   | Melhorar  | Escala e pesos tipográficos              |
| `colorize`  | Melhorar  | Paleta com significado                   |
| `animate`   | Melhorar  | Micro-interações funcionais              |
| `delight`   | Melhorar  | Detalhes memoráveis                      |
| `overdrive` | Melhorar  | Expressividade; só Persuade e Experience |
| `polish`    | Refinar   | Micropolimento final                     |
| `harden`    | Refinar   | Erros, loading, i18n, limites            |
| `onboard`   | Refinar   | Primeira execução e estados vazios       |
| `distill`   | Refinar   | Remove ruído visual                      |
| `bolder`    | Refinar   | Aumenta contraste e presença             |
| `quieter`   | Refinar   | Reduz estimulação                        |
| `clarify`   | Corrigir  | UX copy e mensagens                      |
| `adapt`     | Corrigir  | Telas extremas                           |
| `optimize`  | Corrigir  | Gargalos de renderização                 |
| `live`      | Iterar    | Variantes ao vivo no navegador           |
| `generate`  | Iterar    | Variações simultâneas                    |
