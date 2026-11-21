# Parte II — O Arsenal de 8 Peças Complementares

## Por que exatamente estas oito ferramentas

Cada ferramenta do processo `legoloop-web-qa` cobre uma classe distinta de defeito sem duplicar esforço nem depender de serviços pagos fechados:

| Camada              | Ferramenta             | O que detecta                                  | Modo de Execução                       |
| :------------------ | :--------------------- | :--------------------------------------------- | :------------------------------------- |
| Integridade de Rede | `lychee`                 | Links, âncoras e assets com HTTP >= 400 ou timeout | Binário Rust CLI (`--format json`)       |
| Render Rápido Local | `obscura`                | Falhas de carregamento, bloqueio de rede local e DOM inicial | Binário Rust (`--allow-private-network`) |
| Runtime e Sessão    | `playwright`             | Erros de console, `pageerror`, falhas de rota e sessões com `storageState` | Chromium headless (`1440x900` e `390x844`) |
| UI Comportamental   | `renderlint`             | Botões mortos (`dead-control`), conteúdo encoberto (`obscured-content`) e `text-overflow` | CLI Node sobre navegador real          |
| Consistência Visual | `aidd-linter-design`     | Desvios de design tokens, espaçamento e contraste estrutural | Skill canônica de design tokens        |
| Acessibilidade      | `axe-core`               | Violações WCAG (`critical` e `serious`), ARIA e foco por teclado | `@axe-core/cli` (`--exit`)                 |
| Saúde Multi-Rota    | `unlighthouse`           | Queda de performance, SEO e boas práticas abaixo de 90 em todas as rotas | `@unlighthouse/cli` (`--no-open`)          |
| Regressão Visual    | `lost-pixel`             | Diferença pixel a pixel contra baseline aprovada | CLI open-source em `.lostpixel/baseline` |

## Separação entre `obscura` e `playwright`

O motor `obscura` executa verificação ultrarrápida em V8 nativo (~30MB de RAM) com permissão explícita de rede local (`--allow-private-network`), validando prontidão da página antes de acionar sessões completas de Chromium no `playwright`. O `playwright`, por sua vez, garante fidelidade total de eventos JavaScript, captura de `console.error`, injeção opcional de `storageState` para fluxos autenticados e teste nos dois viewports padrão.
