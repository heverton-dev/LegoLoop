# Lazy de skills e engenharia de contexto: causa, solução, resultado e processo

## A causa dos problemas

O problema nasceu de uma suposição que parecia correta e não era. Aplicamos ao conjunto de skills o mesmo conceito que funcionou no MCP: deixar os recursos fora do contexto e buscá-los sob demanda. Mas o processo `instalar-skills` continuava gravando cada skill do Cofre como vínculo direto nas pastas de cada harness. Assim, o agy via centenas de skills em todo turno e mostrava 20,7k tokens só de lista de skills. O mesmo código de vínculo "eager" existia também no `fabricar`, que espalha a skill que cria.

Há uma segunda causa, descoberta lendo o prompt real que o agy envia (nos arquivos WAL das conversas). O agy tem um orçamento fixo de cerca de 60 KB, ou ~20,5k tokens, para o bloco de skills. Remover algumas skills só abria espaço para outras candidatas entrarem. O custo só cai quando o total de candidatas fica abaixo do teto. Várias fontes alimentavam esse bloco: as pastas do agy, plugins (inclusive um com sufixo ".disabled", que carrega mesmo assim) e skills nativas.

No Claude Code o quadro era diferente e maior do que eu havia estimado. Os 77 comandos em `~/.claude/commands`, as skills em `~/.agents/skills`, três plugins (Vercel, Postman e agent-skills), as skills sincronizadas da conta no claude.ai, os conectores claude.ai (3,5k tokens por turno) e oito ferramentas nativas pesadas somavam o contexto. O `/context` interativo também enganou várias vezes: ele reatribui tokens entre categorias, então a categoria "Skills" pode cair sem que o total caia.

## Como resolvemos

Primeiro corrigimos o código. O vínculo passou a ser lazy por padrão: só a skill-índice `legoloop-skills` vai para os harnesses, e as outras vivem no Cofre, achadas por `skill_search` e lidas por `call_skill`. A flag `--eager` restaura o comportamento antigo, e um arquivo de exceções permite manter skills específicas visíveis. Corrigimos também o ranking da busca (nome pesa mais que descrição, palavras vazias não pontuam, empates são determinísticos), para que a skill certa apareça em primeiro lugar. Tudo isso foi feito em TDD, com testes escritos antes do código e gates limpos (`go vet`, `gofmt`, limite de 35 linhas por função). Criamos o comando `skills-custo` para medir o custo por harness.

Depois migramos os dados, sempre com backup e script `RESTAURAR.sh`: as skills do agy, do Claude (`~/.agents/skills`), do agy-ide e do opencode foram para o Cofre, e os 77 comandos e os agentes impeccable também. Para os plugins, as skills do `fable` foram copiadas ao Cofre e o plugin foi desabilitado. As skills com frontmatter quebrado foram consertadas.

Em seguida fizemos a engenharia de contexto do Claude Code, medindo cada passo pelo consumo real da API (`usage` de uma chamada `claude -p`), e não pelo `/context`. Isso revelou as alavancas certas: o `permissions.deny` remove uma ferramenta do prompt (PowerShell, Task, Workflow, ScheduleWakeup e outras somavam cerca de 12k tokens), o `syncClaudeAiSkills=false` e o `ENABLE_CLAUDEAI_MCP_SERVERS=false` removem skills e conectores da conta, e o teto `skillListingBudgetFraction` limita a lista. Algumas ideias não funcionaram e foram descartadas: `user-invocable-only` aumentou o custo, e o `skillOverrides` só movia tokens de categoria sem ganho líquido.

Por fim, construímos o modo gateway. O `legoloop mcp-router --gateway` expõe uma única ferramenta, `call_mcp_tool`. Pelo servidor embutido `ws` ela faz shell e arquivos, e pelo servidor `legoloop` ela faz busca de skills, leitura de skills, busca de MCPs e o swarm. Nenhum recurso se perde: tudo continua acessível, só deixa de ocupar o contexto até ser usado.

## O resultado

Medido pela API: o ambiente original custava cerca de 41,4k tokens por turno (37,1k com a configuração, 39,6k com as skills e os comandos, e mais a memória antiga). Hoje o padrão enxuto custa 11,3k (cerca de −73%). O perfil `claude-gateway` custa 2,6k com as suas regras de memória (−93,7%) e 1,45k sem memória (−96,5%). A categoria Skills, que originou o problema, caiu de 9,9k para 0,11k (−98,9%); no agy a lista caiu de 20,7k para ~1,6k.

A meta de 97% do contexto total não foi atingida, e isso está medido: o piso do próprio Claude Code, com uma ferramenta MCP vazia, é 1,33k (−96,8%). O harness cobra cerca de 0,41k só por existir uma ferramenta. A prova funcional do gateway passou: uma sessão gravou um arquivo, leu de volta, rodou um comando e achou uma skill, tudo por uma ferramenta.

## O processo para sempre obter o resultado

1. Instale skills sempre com `legoloop instalar-skills` (lazy por padrão) e crie com `fabricar` sem `--eager`. Não coloque skills diretamente nas pastas dos harnesses nem instale plugins pesados sem medir. Ferramentas de terceiros que instalam skills por conta própria desfazem a economia.
2. Depois de adicionar skills ao Cofre, reinicie os harnesses e o `mcp-router`: o índice é carregado na partida.
3. No Claude Code mantenha as configurações em `~/.claude/settings.json` e também em `settings.json` do perfil Orca. São arquivos diferentes, e o Orca usa o do perfil. Mantenha o `deny` das ferramentas pesadas, os conectores desligados por padrão e o teto de skills.
4. Escolha o lançador pela necessidade: `claude` (padrão enxuto, 11,3k), `claude-gateway` para sessões longas (2,6k, ou 1,45k com `--sem-memoria`), `claude-minimo` (~6,3k), `claude-conectores` quando precisar dos conectores claude.ai e `claude-completo` quando precisar de todas as ferramentas.
5. Novas memórias entram em `indice-completo.md`. O `MEMORY.md` deve continuar curto, com o ponteiro e a linha de regras fixas.
6. Para medir, use `claude -p "ok" --output-format json` com `MSYS_NO_PATHCONV=1`, somando `input_tokens`, `cache_creation_input_tokens` e `cache_read_input_tokens`. Use `legoloop skills-custo` para o custo por harness. Não confie só no `/context`.
7. Se algo faltar, os backups com `RESTAURAR.sh` desfazem cada migração.

## Pendências e ressalvas

Não verifiquei o `/context` interativo do perfil gateway. O custo real do agy-ide e do opencode ainda não foi medido abrindo cada um. O agy também depende de fontes fora do nosso controle (plugins e skills nativas). Por último, o hook de segurança acusou o valor de `CONTABO_SECRET_USER` na saída da leitura do `settings.json`: essa chave deve ser rotacionada e o valor retirado do arquivo.
