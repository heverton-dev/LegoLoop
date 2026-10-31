# 🧩 Espec v3 — Studio, Inventário, Triggers, Sub-loops e Execução

> **Status:** Decisões aprovadas em 31/10/2026 · **Substitui/estende:** `SPEC_STUDIO_GRAPH_CANONIZACAO.md`
> **Regras:** `AGENTS.md` (Zero Stubs, TDD estrito, PT-BR, Orchestrix estéril, sem subagentes pesados)

---

## 1. Decisões fechadas

| ID | Decisão |
|---|---|
| D1 | Triggers do MVP: `manual`, `cron`, `file_change`, `loop_end`. `webhook` e `git_event` = **fase 2** (registrado). |
| D2 | Hook = comando `pre` ou `post` de uma peça (≠ gate, que decide passa/falha). Inventário lista todos os hooks (global, projeto, scaffold, por harness), origem visível. |
| D3 | Sub-loop = peça isolada, contexto estéril, devolve só o artefato final. |
| D4 | MCP: inventário lista todos. Sem cliente MCP real ⇒ **não executa** (sem stub). Falta de MCP ⇒ sugere instalar via **catálogo embutido** (offline). |
| D5 | Materialização: **Executar no Studio** + **Skill orquestradora** no MVP. Scaffold = depois. Contrato único: `pipeline.dag.json`. |
| D6 | Modo de execução por peça: `inline` (padrão) ou `isolated`. Isolado **obrigatório** em peças paralelas que escrevem. |
| D7 | Sugestões 100% **determinísticas** (sem LLM), calculadas só do inventário. Sempre com **personalização** (o usuário escreve o próprio). |
| D8 | Canonizar mantém os limites atuais: só duplicatas exatas (SHA-256), hardlink no mesmo volume se symlink exigir admin, ignora repos git sem `--include-git`. UI: **Simular** antes de **Aplicar**, explicando limites. |

---

## 2. Requisitos (critérios binários)

### 2.1 Interface
- **R-UI-1** Paleta virtualizada: renderiza só itens visíveis; com 5.000 peças, nós no DOM da paleta ≤ 100.
- **R-UI-6** Categorias da paleta: `Todas`, `Skills`, `MCPs`, `Loops`, `Scripts`, `Hooks`, `Gates`. O usuário escolhe a categoria e só ela é listada; cada categoria mostra o contador de itens do inventário. `Todas` agrupa por categoria, com selo de tipo em cada item. Categoria vazia mostra estado vazio explícito (nunca inventa item). `Loops` lista o acervo do usuário (arrastar cria peça `subloop`). `Hooks` e `Gates` são arrastados **sobre um nó** e entram no slot correspondente (hook pergunta `pre` ou `post`); soltar no canvas vazio é recusado com Toast explicando.
- **R-UI-7** Busca da paleta: campo de texto acima da lista, resultado a cada tecla (≤ 100 ms para 5.000 itens), sem diferenciar maiúsculas/minúsculas nem acentos, casando por nome, id e descrição. Atua dentro da categoria escolhida; em `Todas` busca em tudo. Sem resultado: estado vazio com o termo digitado. `Esc` limpa o campo; `/` foca o campo.
- **R-UI-8** Laterais fixas: as colunas esquerda (peças e inventário) e direita (inspetor) ficam fixas na viewport e **só a lista interna rola na vertical**, com a mesma scrollbar de 4px. Cabeçalho da paleta (busca e categorias) permanece visível; a página nunca ganha scroll vertical próprio.
- **R-UI-2** Toda scrollbar da aplicação: largura ≤ 4px (vertical e horizontal), thumb com opacidade 0,5, trilho transparente. Verificável por regra CSS global e teste de snapshot do `style.css`.
- **R-UI-3** Zero diálogo nativo: nenhuma ocorrência de `alert(`, `confirm(`, `prompt(` em `internal/studio/static`. Teste Go falha o build se encontrar.
- **R-UI-4** Componentes próprios `Dialog` (confirmação, consentimento de scripts) e `Toast` (feedback), com foco preso, `Esc` fecha, `aria-live`.
- **R-UI-5** Galeria de templates na tela inicial; escolher um carrega o canvas validado (≥ 3 templates: greenfield, engenharia reversa, auditoria).

### 2.2 Peça (nó)
- **R-NODE-1** Cada peça tem 3 slots: `gates[]`, `scripts[]`, `hooks{pre[],post[]}`. **MCP não é slot: é peça** (nó próprio, ver R-MCP).
- **R-NODE-10** Peça genérica (decisão do usuário em 31/10/2026): skill, script ou MCP do inventário que não declara contrato vira nó com entrada e saída `Qualquer` (`Artifact:Any`). O usuário define o nome do artefato de saída e o gate (padrão `file_exists`); liga com qualquer peça. Os 8 blocos do catálogo seguem com validação de tipos. O navegador só envia o `id` opaco do item: nunca caminho nem comando de hook.
- **R-NODE-11** Resolução no servidor: script do inventário vira argv por extensão (`.ps1`, `.sh`, `.py`, `.js`, `.bat`); hook vira argv por tokenização sem shell. Comando que exige shell (pipe, redirecionamento, `$VAR`, `%VAR%`) não é sugerido nem aceito. Hook com credencial no comando é recusado (a credencial iria para o `pipeline.dag.json`). A resposta da API mascara caminhos (`~`, `<projeto>`); o arquivo gravado guarda os reais.
- **R-NODE-9** Slots são **opcionais**: ausente, `{}` ou listas vazias são escolhas válidas do usuário, sem erro e sem aviso. O sistema **sugere** mas nunca preenche sozinho; "deixar em branco" é uma ação explícita e permitida na interface.
- **R-NODE-2** Painel "Sugeridos" por slot: itens do inventário compatíveis por tipo de artefato; lista vazia mostra estado vazio explícito, nunca inventa item.
- **R-NODE-3** Personalização: o usuário informa gate/script/hook próprio; passa pela mesma validação de segurança (SEC-01/03).
- **R-NODE-4** Desconectar: `Delete` na aresta selecionada ou arrastar a ponta para fora.
- **R-NODE-5** Trocar peça: preserva arestas compatíveis; incompatíveis ficam marcadas e bloqueiam a materialização até resolvidas. Nunca apaga silenciosamente.
- **R-NODE-6** Vizinhos sugeridos: skills, MCPs e sub-loops do acervo; instalados primeiro, ausentes rotulados "não instalado".
- **R-NODE-7** `execution.mode`: `inline` | `isolated`. Validador recusa `inline` em peça de camada paralela com `output.artifact` compartilhado/sobreposto.
- **R-NODE-8** Isolamento: worktree git se o projeto for repo; senão cópia temporária do que a peça lê. Descartado ao fim; só o artefato declarado volta.

### 2.3 Trigger
- **R-TRG-1** Todo DAG exige exatamente 1 trigger; validação recusa 0 ou >1.
- **R-TRG-2** `cron`: expressão de 5 campos validada; `file_change`: caminho relativo à raiz (sem `..`, sem link); `loop_end`: referencia loop do acervo, sem ciclo.
- **R-TRG-3** Gatilhos executam pelo mesmo runner; `webhook`/`git_event` são rejeitados com mensagem "fase 2".

### 2.4 Inventário
- **R-INV-1** Rota `/inventario` com abas: Skills, MCPs, Scripts, Hooks, Gates, Modelos, Harnesses, Binários/CLIs, Loops.
- **R-INV-2** Varredura reutiliza o hasher paralelo; cache por (caminho, mtime, tamanho, hash); re-varredura incremental.
- **R-INV-3** Atualização automática: vigia de arquivos nas raízes conhecidas + re-varredura periódica; novo ativo aparece sem ação do usuário.
- **R-INV-4** Harnesses/CLIs: varre `PATH` e diretórios conhecidos (claude, agy, opencode, mimo, omp, cursor, codex…); mostra versão quando obtida sem executar código não confiável.
- **R-INV-5** Hooks: todas as origens (global/projeto/scaffold/harness) com coluna de escopo; somente leitura, nunca executados pelo inventário.
- **R-INV-6** Chaves de API: só presença, nunca valor (regra existente).
- **R-INV-7** Botão Canonizar conforme D8: Simular → relatório → Aplicar (com Dialog próprio).

### 2.5 Acervo e Sub-loops
- **R-LOOP-1** Acervo: salvar, listar, abrir, duplicar, remover loops (arquivo `.loop.json` no cofre).
- **R-LOOP-2** Peça `subloop` referencia um loop do acervo; portas derivadas do artefato final (entrada = trigger/inputs declarados, saída = artefato final tipado).
- **R-LOOP-3** Detecção de ciclo entre loops (A→B→A) e profundidade máxima 8; validador recusa.
- **R-LOOP-4** Execução do sub-loop é isolada e estéril; relatório agrega tokens e modelos de cada peça interna.

### 2.6 Materialização e Execução
- **R-RUN-1** Botão **Executar** no Studio: endpoint de execução com progresso por peça (SSE), tokens, modelo efetivo, notas de failover, resultado de gates.
- **R-RUN-2** Scripts só executam após consentimento no `Dialog` próprio (equivale a `--allow-scripts`); sem consentimento, peça `script` falha com motivo claro.
- **R-RUN-3** **Materializar** grava `pipeline.dag.json` + **skill orquestradora** (< 20 linhas) no cofre canônico, que chama `legoloop run`; registra o loop no acervo.
- **R-RUN-4** Segurança do Studio preservada: loopback, Host allowlist, mesma origem, CSP `'self'`, corpo ≤ 1 MiB.
- **R-RUN-5** Scaffold de projeto: **fora do MVP** (registrado na fase 2).

### 2.7 MCP
- **R-MCP-1** Inventário lista todos os MCPs encontrados.
- **R-MCP-2** Catálogo embutido (`embed`) com nome, descrição, origem e comando de instalação; peça que exige MCP ausente mostra "Instalar" com origem. Não instala sozinho sem consentimento.
- **R-MCP-3** Nó tipo `mcp` continua **não executável** até existir cliente MCP real; o validador recusa com mensagem explícita.

---

## 3. Mudanças de contrato (`pipeline.dag.json` v2)

```json
{
  "name": "...",
  "version": "2.0.0",
  "trigger": { "type": "cron", "cron": "0 9 * * 1-5" },
  "nodes": [{
    "id": "n1", "type": "skill|script|subloop",
    "ref": "skills/interview",
    "execution": { "mode": "inline" },
    "slots": {
      "gates": ["file_exists('BRIEFING.md')"],
      "hooks": { "pre": [], "post": [] }
    },
    "routing": { "mode": "suggested", "model": "...", "fallback": "..." },
    "output": { "artifact": "BRIEFING.md", "gate": "..." }
  }]
}
```

Retrocompatibilidade: DAG v1 (sem `trigger`) abre no Studio como `manual` com aviso; o runner CLI aceita v1.

---

## 4. Jornada

`baixar` → `legoloop forge` (ambiente, harnesses, 9Router) → primeira varredura → (opcional) Canonizar simulado → template ou canvas vazio → definir trigger → validar → **Executar** ou **Materializar** → loop no acervo → reuso como sub-loop.

---

## 5. Fora do escopo (fase 2)

Triggers `webhook` e `git_event`; scaffold de projeto; cliente MCP real; instalação automática de MCP; registry público de MCPs.

## 6. Riscos

| Risco | Mitigação |
|---|---|
| Vigia de arquivos pesado em 22k ativos | Cache incremental + vigia só nas raízes conhecidas; meta RAM do Studio ≤ 45 MB de baseline Go + 10 MB |
| Worktree sem git | Cópia temporária (R-NODE-8) |
| Trigger `file_change` em loop infinito (peça grava o que vigia) | Validador recusa artefato dentro do caminho vigiado |
| Hook executando código arbitrário | Mesma validação de scripts (SEC-01/03) + consentimento |
| Sub-loop recursivo | R-LOOP-3 |
