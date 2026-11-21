# Varredura do contexto do Claude Code: lazy em todos os elementos (21/11)

Fonte das provas: `/context all` de sessão nova (41,6k) e strings do binário `~/.local/share/claude/versions/2.1.330` (grep -a).

## Mapa por categoria (tokens do /context)
- System prompt 4,1k: fixo. Sem alavanca local.
- System tools 22k: fixo; ferramentas MCP já são lazy via ToolSearch (MCP tools = 0,68k).
- MCP server instructions 1,7k: vêm dos conectores claude.ai (Docs, Data, Meetings). Alavanca: env `ENABLE_CLAUDEAI_MCP_SERVERS=false` (binário) desliga o carregamento automático. Custa os conectores; reavaliar uso.
- Custom agents ~0,3k (4 impeccable do usuário): baixo. Alavanca: mover para o Cofre.
- Memory files 3,4k: MEMORY.md 2,2k (índice, cresce a cada ticket), stept 0,6k, RTK 0,37k, CLAUDE.md 0,23k. Alavanca: enxugar o índice e mover RTK.md para skill lazy.
- Skills 5,6k: ver abaixo. É onde há alavanca de configuração comprovada.

## Alavancas de skills comprovadas no binário (2.1.330)
- `skillOverrides` {nome: on | name-only | user-invocable-only | off}: "name-only" lista a skill sem descrição (lazy real: o modelo sabe que existe e carrega sob demanda); "user-invocable-only" esconde do modelo e mantém /nome.
- `skillListingBudgetFraction` (padrão 0,01 da janela em caracteres x 4): teto duro da lista. Janela 1m x 4 x 0,01 = 40k caracteres ≈ 10k tokens, o que explica o teto de ~9,9k visto. Com 0,001 o teto cai para ~1k tokens.
- `skillListingMaxDescChars` (padrão 1536): corta cada descrição.
- `syncClaudeAiSkills: false`: remove as 9 skills do claude.ai (~2,1k).
- `disableBundledSkills`: remove built-ins por inteiro (perde recurso; preferir name-only).

## Plano (nada aplicado ainda)
1. `syncClaudeAiSkills=false` (-2,1k).
2. `skillOverrides` name-only para built-ins e fable (-3,2k, mantém os recursos sob demanda).
3. `skillListingBudgetFraction=0,001` como teto de segurança (alvo <= 1k).
4. Memória: enxugar MEMORY.md e levar RTK.md para o Cofre.
5. Repetir no agy com o mesmo método (migrar agy-ide, opencode).

## Alertas
- O hook leak-scan acusou o valor de CONTABO_SECRET_USER na saída da leitura do settings.json do Claude (chaves `autoMode` ou `hooks`). Rotacionar a chave e tirar o valor do arquivo.
- Não verificado: se name-only reduz de fato o custo medido no /context; medir após aplicar.

## Resultado medido (21/11)
- Skills no /context do Claude: 5,6k -> 827 tokens (meta 100-1k atingida). Sessão nova, após `syncClaudeAiSkills=false`, `skillOverrides` name-only (22 skills) e `skillListingBudgetFraction=0,001`.
- Total da frente: 20,7k -> 0,83k (-96%).
- Backups: settings.json.bak-lazyctx-20261121 (nos dois settings: ~/.claude e perfil Orca).
- Pendentes: Memory 3,4k, agy-ide (6,9k), opencode, conectores claude.ai (1,7k de instruções).

## agy-ide e opencode (21/11)
- agy-ide: 76 itens movidos para ~/.gemini/_skills_backup_20261121_ide (RESTAURAR.sh); custo 59 skills (~6,9k) -> 1 (~61 tokens, estimativa).
- opencode: 80 itens movidos para ~/.opencode/_skills_backup_20261121_lazy (RESTAURAR.sh); 87 -> 1 skill (~57 tokens, estimativa). Todas já estavam no Cofre.
- 7 skills com frontmatter inválido (Conexão Elite: contextual-architect, ds-sculptor, master-assembler, motion-master, quality-gatekeeper, tech-lead, technical-polyglot) ganharam frontmatter derivado do próprio H1/"Função:" e entraram no Cofre. Skills aninhadas em repositórios clonados (caveman, cavecrew, lean-ctx) copiadas ao Cofre. Cofre: 299.
- Diretórios sem SKILL.md: 24 estão vazios (0 bytes). 8 têm equivalente no Cofre (aidd-*); 16 não têm fonte (fluxo*, bridge, generate, orca-plan-orchestrator, planos-auditoria-runner). Não se inventa SKILL.md sem conteúdo-fonte.
- Não verificado: busca das novas skills (o MCP em execução tem índice antigo; reiniciar mcp-router).

## Medição autônoma e A/B (21/11)
- Instrumento: `MSYS_NO_PATHCONV=1 claude -p "/context" --model claude-sonnet-5-5` (headless, reproduzível; sem a variável o MSYS converte /context em caminho).
- Headless, sem skillOverrides: total 23,6k; Skills 2,1k; System tools 14,3k (+16,5k adiadas, fora do total).
- Headless, com skillOverrides (22 name-only): total 23,6k; Skills 0,8k; System tools 15,7k.
- Conclusão: name-only move ~1,3k de Skills para System tools (descrição da ferramenta Skill); ganho líquido ~0. O que derruba Skills de 9,9k é o teto 0,001 + syncClaudeAiSkills=false + remoção das fontes.
- O salto 22k -> 30,4k de System tools no /context interativo NÃO é reproduzível no headless e não é explicado pelo skillOverrides (+1,4k no máximo). Causa em aberto.
- Estado final mantido: skillOverrides (22) + teto 0,001 + sync=false, nos dois settings.json.

## Meta de 97%: medição real e platô (21/11)
Instrumento verdadeiro: `claude -p "ok" --output-format json` (usage.input + cache_creation + cache_read). O /context de categorias é estimativa e reatribui tokens entre System tools e Skills.
- Custos reais medidos por A/B: conectores claude.ai +3,46k; plugin fable +0,73k; vercel+postman+agent-skills +2,27k; skillOverrides e teto de skills já não mudam nada.
- Estado original reconstruído (plugins, conectores, fable, sync e teto padrão ligados): 37,4k. Hoje: 23,2k.
- Piso fixo do Claude Code (prompt do sistema + ferramentas nativas): ~21,7k; não controlável pelo usuário.
- Controlável real: ~22k (original, estimado) -> ~1,5k (hoje) = ~93% de redução.
- Feito nesta rodada: conectores desligados por padrão + lançador `claude-conectores` (religa via --settings); fable movido ao Cofre (4 skills) e plugin desabilitado; agentes impeccable ao Cofre (agente-*) com backup; MEMORY.md em dois níveis (índice completo em indice-completo.md); CLAUDE.md global compactado (704 -> 337 B).
- Restante controlável (~1,5k): stept-rules 0,6k (contrato do usuário, travado por testes), MEMORY.md+CLAUDE.md ~0,25k, stubs MCP legoloop/mobbin ~0,35k, skills ~0,1-0,8k.
- Para 97% faltariam ~0,8k: exigiria comprimir ou tornar lazy o stept-rules (decisão do usuário) e esconder o swarm_* do nível superior do router.
- Achado de armadilha: `user-invocable-only` aumenta System tools em ~6,3k; não usar. Chaves `fable:*` em skillOverrides não casam.

## Rodada 2: ferramentas nativas e piso (21/11)
Medidas reais (usage da API), todas por A/B com `--settings`:
- Piso absoluto: prompt curto, sem ferramentas, sem MCP, sem settings = 894 tokens. Prompt padrão sem ferramentas = 4,3k. Ferramentas nativas padrão = +18,5k.
- `permissions.deny` REMOVE a ferramenta do prompt. Custos individuais: Task/Agent 2,9k; PowerShell 3,7k; Workflow 1,9k; ScheduleWakeup 1,7k; ReportFindings 0,8k; ShareOnboardingGuide 0,4k; ListAgents 0,4k. Demais ferramentas deferidas custam ~6 tokens.
- Essenciais mantidas (custo real): Bash 1,7k, Grep 1,3k, Skill 0,84k, Read 0,7k, Edit 0,35k, Glob 0,28k, Write 0,24k. ToolSearch é obrigatória: negá-la faz o MCP carregar por inteiro (+31k).
- Aplicado nos dois settings.json (backup .bak-deny-20261121): deny das 8 ferramentas pesadas. Nenhum recurso se perde: `claude-completo` (lançador novo) ignora as configurações de usuário e monta as suas em arquivo temporário, sem deny e com conectores ligados (medido: 26,2k).
- Resultado real: 37,4k (original reconstruído) -> 23,2k -> 11,3k. Redução de ~70% do contexto total e ~75% contando memória e skills de usuário que não consegui reconstruir.
- 97% do TOTAL não é atingível: custaria 1,1k, abaixo do prompt padrão sem ferramentas (4,3k). O piso depende do Claude Code, não do usuário. Só se alcançaria trocando o prompt do sistema, o que muda o comportamento do produto.
- Erro corrigido: o `real.sh` falhava de forma intermitente (saída vazia); agora tenta 3 vezes.

## Rodada 3: perfil mínimo e prova do limite (21/11)
- `claude-minimo` (novo): prompt enxuto + Bash, Read, Edit, Write, ToolSearch; skills e MCPs por demanda via legoloop. Provado: executou Bash e respondeu "funciona". Custo por turno ~6,3k reais (12,9k em 2 chamadas).
- Variantes medidas: Bash+ToolSearch 5,0k; +Read/Edit/Write 6,3k; +Grep/Glob/Skill 8,7k.
- Limite provado: prompt curto + só Bash = 4,0k; zero ferramentas = 0,9k. Uma sessão com ferramentas não cai de ~2,6k (Bash 1,7k + piso 0,9k). 97% do contexto original (~45k) = ~1,4k: inatingível por qualquer configuração útil.
- Melhor resultado real: 45k -> 6,3k (-86%) no perfil mínimo; 45k -> 11,3k (-75%) no perfil padrão enxuto. Todos os recursos permanecem acessíveis (claude-completo religa tudo).

## Rodada 4: modo gateway de uma ferramenta (21/11)
Código (TDD, Red->Green): `Router.EnableGateway`, `internal/mcp/workspace_tools.go` (servidores embutidos `ws` e `legoloop`), flag `legoloop mcp-router --gateway`. Testes: `internal/mcp/gateway_test.go` (7) e `TestMcpRouterGatewayExpoeUmaFerramenta`. Gates: go vet, gofmt e funclen limpos. Falha alheia pré-existente: `TestLegoPressCLI_Auditar`.
- tools/list em modo gateway = só `call_mcp_tool`. Skills (`skill_search`, `call_skill`), `mcp_search`, tools de Processo e `swarm_*` passam a ser tools do servidor embutido `legoloop`; shell e arquivos pelo servidor `ws` (bash, read_file, write_file, edit_file; edit recusa trecho ambíguo ou ausente).
- Medido por turno (usage da API): piso sem ferramentas 0,92k; gateway de 4 ferramentas 3,05k com memória e 1,96k sem; gateway de 1 ferramenta 2,60k com memória e 1,48k sem memória.
- Provado em sessão real (`claude-gateway --sem-memoria`): `echo` via ws.bash, gravação de arquivo via ws.write_file e `skill_search` via legoloop retornaram `shell-ok | arquivo-ok | legoloop-fechar`.
- Redução: 1,48k contra 37,4k (original reconstruído e medido) = -96,0%; contra ~46k (estimativa com memória e skills de usuário) = -96,8%. Com memória (2,6k): -93%.
- Falta ~0,36k para 97% estrito sobre 37,4k. O restante é overhead fixo do harness por ferramenta MCP (~0,55k para 1 ferramenta, descrição e schema já mínimos: comprimir mais rendeu só 33 tokens).
- Armadilhas: sem as settings do usuário (`--setting-sources ""`) a ferramenta fica sem permissão em -p, por isso o lançador passa `--allowedTools`; caminhos do lançador usam `cygpath -m` (MSYS_NO_PATHCONV quebra `$HOME`).

## Rodada 5: linha de base medida e piso do harness (21/11)
Linha de base MEDIDA (usage da API, `--setting-sources project,local` + settings sem deny, plugins/conectores/fable/sync ligados, teto padrão):
- A) config original: 37,1k. B) + 98 skills, 77 comandos e 4 agentes restaurados via `--add-dir`: 39,6k (a lista de skills respeita o teto de 1% da janela, por isso só +2,5k). C) + memória antiga (MEMORY.md 4,7 KB + RTK.md + CLAUDE.md): 42,4k. Descontando ~1,0k de memória atual que A já continha: ~41,4k.
- Custos isolados (anexados a um prompt): 97 skills 11,8k; 77 comandos 4,7k; memória antiga 2,8k; agentes 0,4k. Os 2 primeiros somam mais que o teto, por isso não se somam.
Piso do harness (medido):
- Sem ferramentas, prompt curto: 0,92k. Uma ferramenta MCP vazia (descrição "x") custa +0,41k: piso de qualquer sessão capaz de agir = 1,33k. Nomes de servidor e de ferramenta não mudam o custo (1.320 a 1.330).
- `claude-gateway --sem-memoria` final: 1,45k (prompt vazio 1,41k). Schema sem propriedades economiza só 41 tokens e QUEBRA a chamada (o modelo erra os argumentos); descrição de 150 B com 3 propriedades é o ponto estável.
Resultado final: 41,4k -> 1,45k = -96,5% (-96,6% sobre 42,4k). O piso físico (1,33k) daria -96,8%: 97% estrito é inatingível por qualquer sessão que consiga agir.
Com as suas regras de memória (stept): 2,6k = -93,7%.

## Rodada 6: lançadores no Orca (21/11)
- Criados `claude-gateway.cmd`, `claude-minimo.cmd` e `claude-completo.cmd` (bash.exe do Git chamando o script); antes só `claude-conectores` tinha `.cmd`. Testado por `cmd.exe /c`: gateway = 1.448 tokens.
- Argumentos do Orca (`--dangerously-skip-permissions --model sonnet --effort medium`) passam pelo lançador (`"$@"`) sem custo.
- `--chrome` custa: gateway 1.448 -> 13.519 (+12,1k, porque sem ToolSearch as ferramentas do claude-in-chrome carregam por inteiro); claude padrão 11.423 -> 13.713 (+2,3k). No gateway, não usar `--chrome`.
- O Orca define CLAUDE_CONFIG_DIR para o perfil; os lançadores leem `settings.json` de lá.

## Rodada 7: lançadores fora do Orca (21/11)
- Falha real num PowerShell comum: o PATH do Windows não tem os utilitários do Git (`dirname`, `cygpath`, `mktemp`); os lançadores quebravam com "dirname: command not found".
- Correção: os 4 scripts (`claude-gateway`, `claude-minimo`, `claude-conectores`, `claude-completo`) agora acham o `gw.exe` pelo próprio caminho (`${0//\//}`), acrescentam `/usr/bin` do Git ao PATH e usam `USERPROFILE`; os `.cmd` também prefixam o PATH. `gateway.mcp.json` usa o caminho absoluto do `legoloop.exe`.
- Teste com PATH mínimo do Windows (System32, .local\bin, go\bin, Python314) e config do perfil Orca: padrão 11,6k; gateway 1,48k (2,57k com memória); mínimo 6,4k; conectores 14,6k; completo 22,5k.
- Fora do Orca o login de `~/.claude` está expirado (credenciais de 18/11) e o `claude` comum falha igual aos lançadores ("login expired ... credit"). O perfil Orca tem login de 21/11. Solução do usuário: `claude /login` no terminal ou apontar `CLAUDE_CONFIG_DIR` para o perfil Orca.

## Rodada 8: todos os harnesses lazy (21/11)
Regra do usuário: todo harness do SO usa as técnicas desenvolvidas.
- Inventário (pastas de skills com SKILL.md): ainda eager antes desta rodada: opencode (`~/.config/opencode/skills` 66 + `skill` 26; a pasta que eu migrara antes, `~/.opencode/skills`, era outra), cline 41, hermes 38, windsurf 45, mimocode 26, crush/devin/goose/augment/junie 18 cada, pi 8, grok 5. Já lazy: agents, claude, codex, cursor, copilot, trae, roo, continue, qwen, agy.
- Migração com backup centralizado em `~/.aidd/backups/skills_20261121/<harness>/` (MANIFESTO.tsv + RESTAURAR.sh): Cofre 250 -> 440 skills (pacotes de skills aninhados foram abertos e copiados ao Cofre).
- Código (TDD, Red provado antes do Green): `processo.harnessesConhecidos` de 7 para 23 harnesses com `layoutsEspeciais` (opencode, windsurf, crush, devin, goose, pi, antigravity); `skills-custo` cobre todos (24 linhas, todas com 1 skill); `sync-harnesses` ganhou pi e OMP (o MCP `rea` direto, 207 tools, deu lugar ao roteador lazy); novo comando `skills-migrar` (`internal/skills/migrar.go`).
- Proteção descoberta no teste real: pacote sem SKILL.md no topo que é repositório de ferramenta (.git, package.json, Cargo.toml, go.mod, pyproject.toml, setup.py, Makefile) NÃO é movido (caveman, lean-ctx, rtk, agent-desktop, headroom ficaram no lugar).
- Busca com 440 skills continua precisa (legoloop-fechar, tech-lead, caveman, aidd-tdd no topo).
- Resíduos fora do controle: skills nativas do agy (15-16), plugins do grok (`installed-plugins/vercel-plugin` 30, marketplace-cache, bundled 22) e comandos do opencode (46); não verificado o custo real de cada harness no prompt.
- Falhas pré-existentes e alheias: `TestLegoPressCLI_Auditar`, `TestProcessosCincoPorCinco`.


## Rodada 9: o grafo semântico (`codebase-memory-mcp`) e a economia de tokens

Método: três perguntas de código no próprio repositório. Para cada uma, comparei o texto que o agente leria com `git grep` e leitura de arquivo contra a resposta do grafo (`codebase-memory-mcp cli ...`). Os tokens são reais: soma de `input_tokens`, `cache_creation_input_tokens` e `cache_read_input_tokens` de uma chamada `claude -p --model haiku`, menos a linha de base (12.688 tokens do prompt vazio).

| Pergunta | Sem o grafo | Com o grafo | Diferença |
| :--- | :--- | :--- | :--- |
| Quem chama `garantirIndice` | `git grep`: 181; `git grep` e ler os arquivos achados: 8.587 | `trace_path`: 148 | −18% contra o `grep`; −98% contra ler os arquivos |
| Ler a função `aplicarMigracao` | arquivo inteiro: 3.382 | `get_code_snippet`: 586 | −83% |
| Estrutura de `internal/mcp/router.go` | arquivo inteiro: 7.595 | `get_file_outline`: 4.215 | −44% |

Leitura honesta: o ganho é grande contra a leitura de arquivos inteiros e pequeno contra um `git grep` simples. O grafo não substitui o `Grep` para texto exato, e o ganho de uma pergunta só vale se o agente leria o arquivo todo. Amostra de 3 perguntas, 1 execução cada: o número mostra a ordem de grandeza e não uma média.

Correções da rodada: o `mcp_search` exigia a frase inteira colada no texto e devolvia ordem aleatória (`internal/mcp/busca_tokens.go`); o binário agora responde pelo PATH (`~/.local/bin/codebase-memory-mcp` e `.cmd`, que apontam para `~/tools/codebase-memory-mcp/`); o grafo foi reindexado (48.234 nós e 106.059 arestas).
