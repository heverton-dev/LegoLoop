# Lazy de skills: verificação profunda, raio de impacto e plano (20/11/2026)

Substitui a versão anterior. Correções da versão anterior estão na seção "Retratações".

## 1. Sintoma
Print do agy (Gemini 3.8 Flash): Skills = 20,7k tokens por turno. MCP tools = 321 tokens (lazy de MCP funcionou).

## 2. Provas (cada uma com comando rodado nesta sessão)

### 2.1 O agy escaneia 3 pastas de skills e funde por nome
- ~/.gemini/antigravity-cli/cli.log:76 "Reloading system slash commands and skills".
- cli.log:79-96 e adiante: 168 linhas "Failed to parse skill file" citam exatamente 3 pastas: `~/.gemini/antigravity-cli/skills`, `~/.gemini/config/skills`, `~/.gemini/skills`. 8 skills inválidas x 3 pastas x 7 recargas.
- Reprodução do número: união das 3 pastas, deduplicada por nome, só frontmatter válido, frontmatter + caminho: 212 skills, 79.391 bytes, ~19,8k tokens. O print mostra 20,7k. Diferença de 4% (formatação do prompt).
- Sem deduplicar seria ~56k tokens. Logo o agy deduplica por nome e o número bate.

### 2.2 Composição da união (212 skills)
- outras (Google builtin, gcp-*, hyperframes, motion etc.): 77 skills, ~7,5k tokens.
- aidd-*: 86 skills, ~7,2k tokens.
- legoloop-*: 39 skills, ~4,6k tokens.
- backup/disabled: 10 skills, ~0,4k tokens.

### 2.3 O `.disabled` não desativa nada
- 30 entradas `*.disabled` têm SKILL.md válido nas 3 pastas e entram na contagem (script de união). Nenhuma aparece nos erros de parse, então o agy as lê.
- Provado: o nome da pasta não exclui. Indício forte, não prova de carga no prompt.

### 2.4 Quem cria vínculos eager (fonte única de código)
- internal/processo/vinculo.go:11 `harnessesConhecidos` = claude, codex, gemini, antigravity, opencode, cursor, kiro.
- vinculo.go:34-41 `caminhosHarness`: gemini -> ~/.gemini/skills; antigravity -> ~/.gemini/config/skills; demais -> ~/.<nome>/skills.
- vinculo.go:51 `PlanejarVinculos`: um vínculo por skill por harness, sem filtro.
- vinculo.go:63 `SymlinkOuCopia`: symlink; se falhar, copia o SKILL.md inteiro.
- Chamadores de `PlanejarVinculos`/`Vincular` fora de teste: internal/skills/instalar.go:208,214 (`publicar`) e internal/fabricar/materialize.go:63 + internal/fabricar/link.go:21.
- internal/cli/bootstrap.go:41 chama `runInstalarSkills`.
- `instalarSkillsSilencioso` (internal/cli/instalar_skills.go:44) não tem nenhum chamador (grep).
- Prova de escala: `legoloop instalar-skills --dry-run --json` -> 73 skills x 7 harnesses = 511 vínculos, 33 comandos, cofre `C:\Users\trcnologia\.aidd\canonical`.

### 2.5 `~/.gemini/antigravity-cli/skills` não é alvo do código
- Não está em `harnessesConhecidos`. Nenhum .go do repo o cita (só relatórios arquivados).
- Contém 203 entradas: 24 symlinks (9 legoloop-* para ~/.gemini/skills, 14 `*.disabled` para pastas de outros projetos, 1 para config/skills) e o resto cópias de 05/11 15:51 idênticas às outras pastas (md5 igual em accidental-data-loss-prevention e building-data-apps).
- Conclusão: a terceira pasta é povoada fora do `instalar-skills`/`fabricar` (cópia em massa de 05/11 e links manuais). Quem a povoou: não verificado.

### 2.6 O fechar NÃO vincula nem instala
- internal/fechar/fechar.go:164 e artefatos.go:67: só lê `skills/<canon>/SKILL.md` do repo e exige >= 10 linhas.
- fechar.go:289-321 e 423: só grava PROCESSOS_CANONIZADOS.md, AGENTS e docs/reports.
- grep de `Vincular|PlanejarVinculos|.Skills` em internal/fechar: zero ocorrências.
- A skill do processo novo é vinculada pelo `fabricar` (`Materialize`, campo `SkillLinks`), com as opções existentes `--no-link` e `--link-harness`.

### 2.7 O router de skills está pronto
- internal/mcp/skill_router.go: `skill_search` e `call_skill` com guarda `dentroDoCofre`. router.go:117 carrega o Cofre.
- internal/cli/sync_harnesses.go: zero referências a skills. O sync lazy cobre só MCP.

## 3. Causa raiz (provada)
1. O agy (e os demais harnesses) lê name+description de TODA skill achada nas pastas, a cada turno. Não há carga tardia nas pastas nativas.
2. O lazy de skills implementou só o router. Nunca removeu o eager das pastas.
3. O código ainda ESPALHA eager: `PlanejarVinculos` liga cada skill em cada harness (511 vínculos). Tanto `instalar-skills` quanto `fabricar` passam por ele.
4. Resíduo histórico: cópias em massa de 05/11 (Google builtin, aidd-*, gcp-*) e `.disabled` ineficaz mantêm ~19,8k tokens mesmo que o código seja corrigido.
Corrigir só o código não zera o custo. Corrigir só as pastas faz o próximo `fabricar`/`instalar-skills` repovoar.

## 4. Raio de impacto
Código a mudar (1 ponto central):
- internal/processo/vinculo.go (`PlanejarVinculos`): vincula só a skill-índice e as skills da lista de exceção.
- Herdam sozinhos: internal/skills/instalar.go, internal/fabricar/materialize.go, internal/cli/bootstrap.go.
- Testes a ajustar (Red primeiro): internal/processo/vinculo_test.go, internal/skills/instalar_test.go, internal/cli/instalar_skills_test.go, internal/cli/paridade_test.go (espera vínculo por skill), testes do fabricar que contam `SkillLinks`.
- Novo: sync-harnesses --skills (migra resíduo para o Cofre), comando `skills-custo` (medição), gate de teto.
Não muda: fechar (provado), router de skills, o formato do SKILL.md no repo.
Dados a mudar (fora do repo, exige aprovação): as 3 pastas do agy, ~/.agents/skills (= ~/.claude/skills, 98 skills, ~9,4k tokens), ~/.gemini/antigravity-ide/skills (163 skills, ~13,6k tokens).
Outros harnesses: Claude Code (~/.claude/skills -> ~/.agents/skills), codex, opencode, cursor, kiro usam a mesma regra e sofrem o mesmo custo. Não medidos no prompt deles.

## 5. Plano
1. Medir: `legoloop skills-custo` (tokens de frontmatter por pasta, por harness, com dedupe por nome como o agy). Grava baseline: 19,8k no agy.
2. Experimento controlado (resolve a dúvida do slash): mover 1 skill piloto para fora das 3 pastas e criar o comando equivalente; reabrir o agy e comparar `/context`. Decide o gatilho.
3. Red: teste que falha se `PlanejarVinculos` devolver mais que índice + exceções.
4. Green: `PlanejarVinculos` filtra. Skill-índice única `legoloop-skills` (~150 tokens) manda usar `skill_search`/`call_skill`.
5. `sync-harnesses --skills --apply`: move do harness para o Cofre, backup `.bak_legoloop`, add-only, igual ao MCP. Cobre as 3 pastas do agy mais .agents, ide, codex etc., um harness por vez.
6. Remove o `.disabled` como mecanismo: itens vão para fora das pastas (Cofre/arquivo), nunca renomeados dentro delas.
7. Teto de 200 caracteres por description no Cofre, gate no `fechar`.
8. Gate de regressão: frontmatter por harness <= 1.500 tokens.
9. Prova final: reabrir o agy, `/context` com Skills < 1k e `skill_search` achando `legoloop-matriz`.

## 6. Riscos
- Skill nativa deixa de aparecer sozinha: depende do índice e do MCP carregado. Mitigar com índice + comandos.
- Sessões abertas precisam reiniciar.
- Apagar/mover dados de usuário (cópias de 05/11) é irreversível sem backup. Pede aprovação por pasta.

## 7. Não verificado
- Se `commands/` (slash) conta no contexto do agy. O agy-cli não tem pasta `commands` (ls). O `.toml` de ~/.gemini/commands pertence ao Gemini CLI. Resolvido só pelo passo 2.
- Quem povoou ~/.gemini/antigravity-cli/skills.
- Tokens exatos nos outros harnesses (só o agy tem print).

## 8. Retratações da versão anterior
- "legoloop-legoloop-* é lixo": o grep nas pastas do agy e em ~/.agents não achou nenhum. Os nomes vieram da lista de skills do Claude Code. Removido.
- "Cofre não achado": está em ~/.aidd/canonical (saída do dry-run).
- "Slash polui contexto?": a resposta anterior assumiu `.toml` do agy. Isso não foi provado. Ver 7.
- "Repetição 3,3x": é real, mas o agy deduplica por nome. O custo é 19,8k (união), não 56k.

## 9. Execução dos passos 1 e 2 (20/11/2026)

### Passo 1: `legoloop skills-custo` (TDD, Red provado antes do Green)
- Código: internal/skills/custo.go, internal/cli/skills_custo.go, registro em internal/cli/cli.go.
- Testes: internal/skills/custo_test.go, internal/cli/skills_custo_test.go. Red: "undefined: CustoGrupo" e "comando desconhecido".
- Medição real do agy-cli: 212 skills, 20.982 tokens. Print do usuário: 20,7k. Erro de 1,4%.
- Outros harnesses (estimativa de frontmatter): agy-ide 16.887, claude/codex/agents (mesma pasta ~/.agents/skills) 11.047.
- Falha pré-existente, não minha: `TestLegoPressCLI_Auditar` falha igual sem as mudanças (stash).

### Passo 2: experimento do `.disabled`
- Movidos 72 itens `*.disabled` (24 por pasta) das 3 pastas do agy para `~/.gemini/_skills_backup_20261120_disabled/` com MANIFESTO.tsv e RESTAURAR.sh.
- Previsão do `skills-custo`: agy-cli 20.982 -> 20.483 tokens (-499), 212 -> 202 skills. Confirmado pelo comando.
- Pendente (só o usuário faz): abrir um agy novo, rodar /context. Se Skills cair de ~20,7k para ~20,2k, está provado que `.disabled` carregava.
- Reverter: `sh ~/.gemini/_skills_backup_20261120_disabled/RESTAURAR.sh`.
- O experimento do slash (comando no lugar de skill) NÃO foi feito: o agy-cli não tem pasta `commands` e o formato aceito não está provado.

### Resultado do experimento `.disabled` (agy novo, /context)
- Previsto se `.disabled` carregasse: Skills ~20,2k. Observado: 20,8k (antes 20,7k). Sem queda.
- Conclusão: hipótese "`.disabled` carrega" REFUTADA. O agy ignora essas entradas. O modelo de custo estava 0,5k alto nelas.
- Consequência: `.disabled` não é alavanca. Os ~20,8k vêm das 202 skills ativas (+ ~0,3k de fonte não medida: skills do workspace ou builtin).
- Seção 2.3 deste relatório fica corrigida por esta: o indício de carga estava errado.
- Próximo experimento com sinal grande: mover as 86 skills `aidd-*` (~7,2k) para backup e medir.

### Experimento 2: mover as aidd-* (20/11)
- Movidos 261 itens `aidd-*` das 3 pastas do agy para `~/.gemini/_skills_backup_20261120_aidd/` (MANIFESTO.tsv, RESTAURAR.sh).
- `skills-custo` agora exclui `*.disabled` do total (teste Red/Green). agy-cli: 202 -> 116 skills, 20.483 -> 12.765 tokens (-7,7k).
- Previsão do /context: Skills ~13,1k (20,8k - 7,7k). Pendente: medir em agy novo.

### Resultado do experimento 2 (aidd-* nas 3 pastas)
- Previsto: Skills ~13,1k. Observado: 20,5k (antes 20,8k). Queda de só 0,3k.
- Log pós-mudança (cli.log 20:58): o erro de `aidd-project-generator` sumiu, então o agy releu as pastas alteradas. As `aidd-*` continuam contadas por outra origem, deduplicada por nome.
- Hipótese "o custo = união das 3 pastas" REFUTADA como explicação completa. A seção 2.1 prova só que essas 3 pastas são escaneadas, não que são as únicas.
- Pistas no agy.exe (strings): importadores `plugins/claude/skills.go` e `plugins/gemini/importer.go`; leitura de `~/.claude/skills`, `~/.codex/skills`, `.agents/skills` hierárquico; changelog "skills loading twice when discovered from both..." (carga dupla já foi bug conhecido).
- Fora das 3 pastas, as `aidd-*` ainda existem em: `~/.gemini/antigravity/skills` (86), `~/.gemini/antigravity-ide/skills` (86), `~/.agents/skills` (41).

### Experimento 3 (20/11)
- Movidos 174 itens `aidd-*` de `~/.gemini/antigravity/skills` e `antigravity-ide/skills` para `~/.gemini/_skills_backup_20261120_aidd2/` (MANIFESTO.tsv, RESTAURAR.sh).
- Previsões do /context em agy novo:
  - ~12,4k: agy lê as 5 pastas e NÃO lê ~/.agents/skills.
  - ~17,3k: agy lê também ~/.agents/skills (41 aidd restantes).
  - ~20,5k (sem queda): a origem das skills não é nenhuma dessas pastas.
- Efeito colateral: o Antigravity IDE perde as aidd-* até a restauração.
