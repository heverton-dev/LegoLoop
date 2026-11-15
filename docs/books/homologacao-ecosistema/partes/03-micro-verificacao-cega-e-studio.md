# PARTE III — Verificação Cega e Homologação Ao Vivo

## 1. A Primeira Rodada Cega: Por que Reprovar é Obrigatório

Na Etapa 5 (`TICK-07.1`), auditores independentes em contexto limpo avaliaram o repositório sem confiar em relatórios prévios dos autores.

A Rodada 1 reprovou a entrega com nota global **9,07 / 10**, abaixo da barra mínima de **9,5 / 10** na Frente 1 (`8,6`) e na Frente 3 (`9,3`):

| Frente na Rodada 1            | Nota           | Pontos Apontados pelo Red-Team                  |
| :---------------------------- | :------------: | :---------------------------------------------- |
| Arquitetura, Go e Taxonomia   | 8,6            | Stub `nil, nil` em `internal/motion/sync.go` e ausência de Path Jail em `/api/assets/translate` |
| Canais CLI, Harnesses e Swarm | 9,3            | Ausência de `--help` com exit `0` em `cmdSwarm` e sinônimo léxico em artefatos do Swarm |
| UI/UX Studio Web e Swarm UI   | 9,3            | Combos legados em `templates.json`, emojis residuais no modal e erro 404 de `/favicon.ico` |

## 2. Remediação 5.1 e Homologação Real no Navegador

Corrigidos os seis desvios apontados na Rodada 1, o binário `./bin/legoloop.exe studio --addr 127.0.0.1:18799` foi iniciado com perfil de usuário isolado e validado com Playwright, Chrome DevTools MCP, REA MCP e Obscura nas resoluções `1440x900`, `1024x768` e `390x844`, nos modos claro e escuro.

![Tela principal do Studio Web em 1440x900 após correções](/img/2026-11-14_HOMOLOG_STUDIO.png){width=90%}

Todos os caminhos críticos foram confirmados sem erros ou avisos no console JavaScript:

- **Navegação sem recorte (`overflowX: false`):** barras e cartões responderam sem sobreposição nas três rotas.
- **Inspetor recuperável (`#btn-toggle-inspector`):** alternância confirmada entre `280px` e `0px`.
- **Guarda de matriz em `+ Novo Processo`:** sem matriz prévia no cofre, o endpoint devolveu HTTP `422` com orientação explícita para `legoloop matriz gerar`.
- **Rota `/favicon.ico`:** resposta HTTP `200` com `image/svg+xml` em todas as páginas (`index.html`, `inventario.html` e `swarm.html`).

![Tela de Inventário de Peças e Harnesses no Studio Web](/img/2026-11-14_HOMOLOG_INVENTARIO.png){width=90%}

![Tela do Swarm UI com painel de workers e eventos](/img/2026-11-14_HOMOLOG_SWARM.png){width=90%}

## 3. Veredito da Reavaliação Cega (Rodada 5.2)

Submetida a base corrigida a nova verificação cega (`TICK-07.1`), todas as frentes superaram o corte de `9,5 / 10`, alcançando média global **9,87 / 10**:

| Frente Auditada                       | Linha de Base  | Rodada 1       | Rodada 2 (Final) | Status         |
| :------------------------------------ | :------------: | :------------: | :--------------: | :------------- |
| Frente 1 — Arquitetura, Go e Taxonomia | 4,8            | 8,6            | 9,9              | Aprovado       |
| Frente 2 — Canais CLI, Harnesses, MCP e Swarm | 5,2            | 9,3            | 9,8              | Aprovado       |
| Frente 3 — UI/UX Studio Web e Swarm UI | 3,8            | 9,3            | 9,9              | Aprovado       |
| **Média Global Ponderada**                | **4,60**           | **9,07**           | **9,87**             | **Liberado**       |
