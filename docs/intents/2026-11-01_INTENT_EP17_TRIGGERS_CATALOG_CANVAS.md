# INTENT_EP17: Triggers no Catálogo Lateral e Inicialização de Pipelines no Canvas N8N

> **Status:** Aprovado para Especificação & Fatiamento  
> **Data:** 01/11/2026  
> **Caso de Teste do Mega-Pipeline:** Épico 17 — Triggers de Primeira Classe no Studio LegoLoop

---

### 1. 📌 Problemas Diagnosticados & Foco da Intenção

1. **Ausência de Triggers no Catálogo Lateral:**
   * O contrato DAG v2 do LegoLoop (`internal/dag/types.go`) exige obrigatoriamente um gatilho formal (`Trigger`: manual, cron, file_change, loop_end), porém a Sidebar Esquerda do Studio lista apenas 7 categorias (`Controle`, `Skills`, `MCPs`, `Loops`, `Scripts`, `Hooks`, `Gates`). Não há categoria nem peças visíveis para Triggers.
2. **Nó `github-board-trigger` Perdido em Outra Categoria:**
   * O bloco canônico `github-board-trigger` foi criado em `ControlBlocks` com `Kind: "trigger_board"`. Como a paleta não possui pasta para gatilhos, ele foi alocado como "Skill", gerando inconsistência ontológica.
3. **Ponto de Partida Desconexo do Padrão N8N:**
   * Em ferramentas padrão de mercado (n8n, Zapier), todo pipeline inicia obrigatoriamente por um nó de gatilho visual. No Studio, o Canvas inicia vazio e aceita qualquer bloco solto, sem ancoragem no gatilho que disparará a esteira.
4. **Falta de Configuração Paramétrica do Gatilho:**
   * Ao compor uma esteira, se o usuário não configurar explicitamente um gatilho, o backend assume `manual` em silêncio, sem permitir que o usuário configure expressões `cron` (ex: `*/15 * * * *`), caminhos observados (`path` em `file_change`) ou loops predecessores (`loop` em `loop_end`) pelo Inspetor.

---

### 2. 🏛️ Escopo da Solução

#### 2.1 Categoria Canônica "Triggers" no Menu Lateral
* A Sidebar ganha a categoria **Triggers** no topo da lista sanfona (antes de *Controle*).
* Ícone distintivo em SVG (raio/disparador N8N em âmbar).
* Peças nativas disponibilizadas:
  1. `⚡ Disparo Manual` (`trigger_manual`) — Disparo sob demanda via UI ou CLI.
  2. `⏰ Agendador Cron` (`trigger_cron`) — Execução agendada periódica (5 campos).
  3. `📁 Observador de Arquivos` (`trigger_file_change`) — Vigilância de caminhos no workspace.
  4. `🔁 Encadeamento de Loop` (`trigger_loop_end`) — Disparado ao término de outro loop.
  5. `📋 GitHub Board Kanban` (`trigger_board`) — Gatilho por movimentação de cards no GitHub Project.

#### 2.2 Nó de Gatilho no Canvas N8N
* O nó de Trigger atua como ponto de partida visual da esteira (estilo n8n).
* Não possui portas de entrada (`consumes`), possuindo apenas porta de saída (`produces: Idea/URL/Artifact`).
* Regra estrita: O Canvas aceita **no máximo 1 nó de Trigger ativo** (ou substitui o trigger existente ao soltar um novo).

#### 2.3 Inspetor de Gatilho Lateral Direito
* Ao clicar no nó Trigger no Canvas, o painel Inspetor exibe campos dinâmicos:
  * Se `cron`: campo de expressão cron validada com preview explicativo.
  * Se `file_change`: seletor de caminho do workspace validado contra loops infinitos.
  * Se `loop_end`: seletor de loops salvos no acervo (`internal/loop/store`).
* Ao serializar a esteira (`/api/compose`), o trigger é injetado diretamente no cabeçalho do DAG v2.

---

### 3. 🎯 Critérios de Sucesso do Épico

* [ ] Categoria `Triggers` visível no menu lateral com contador e itens corretos.
* [ ] Arrastar ou clicar em um Trigger adiciona o nó raiz no Canvas.
* [ ] Inspetor permite parametrizar `cron`, `path` e `loop`.
* [ ] Validação do Studio e serialização DAG v2 envia o `Trigger` correto para `/api/compose`.
* [ ] Suíte de testes (Go e JS) 100% verde com zero regressões.
