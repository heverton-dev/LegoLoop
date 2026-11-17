# Relatório: `swarm run` com agentes de LLM reais como workers

Data: 2026-11-16. Janela da execução bem-sucedida: 17:17:12Z a 17:19:00Z (108 s).

## 1. Veredito

- **Orquestração: aprovada.** Níveis em ordem, 2 agentes em paralelo real, worktrees isolados, repositório intacto e limpeza completa.
- **Qualidade dos agentes: reprovada.** 14 das 19 citações `arquivo:linha` distintas apontam para linhas que não existem. Os 3 gates passaram mesmo assim, porque medem só estrutura.

## 2. O que foi executado

DAG: `scripts/fixtures/swarm_llm_auditoria.dag.json`, 3 nós, 2 níveis.

- Nível 0: `a` (auditor-qualidade) e `b` (engenheiro-tdd) auditam em paralelo o `swarm run` (somente leitura).
- Nível 1: `c` (engenheiro-backend) consolida os riscos de `a` e `b`.

Comando:

```bash
legoloop swarm run scripts/fixtures/swarm_llm_auditoria.dag.json \
  -dir <quadro-descartavel> -repo <repositorio-legoloop> -timeout 300s \
  -- bash <repo>/scripts/swarm_agente_llm.sh
```

Cada worker é `claude -p` (Claude Code 2.1.316) roteado pelo 9Router, modelo `t-agente-128k`, com `--allowedTools "Read Grep Glob"`, `--disallowedTools "Bash Edit Write NotebookEdit WebFetch WebSearch"`, sem MCP e sem persistência de sessão. O cwd de cada worker é um worktree destacado do repositório.

Antes, um teste de fumaça com `code-fast` e `t-agente-128k` respondeu `ok` em 5 s cada.

## 3. Resultado por nível

| Nível | Workers | Início (UTC) | Fim (UTC) | Códigos |
| :--- | :--- | :--- | :--- | :--- |
| 0 | a, b | 17:17:13.035 | 17:17:30.673 | 0, 0 |
| 1 | c | 17:17:30.673 | 17:18:58.229 | 0 |

Todos os workers terminaram na primeira tentativa (`tentativa=1`).

## 4. Provas da orquestração

- **Paralelismo real:** `a` rodou 11,5 s e `b` 9,6 s, com janelas sobrepostas (início de cada um antes do fim do outro).
- **Ordem por nível:** `c` iniciou depois do fim de `a` e de `b` e durou 85,1 s.
- **Gates (motor `internal/gate`):** os 3 do DAG deram verdadeiro (`a_riscos.md`, `b_riscos.md`, `c_consolidado.md`).
- **Limpeza:** os 3 workers ficaram `killed` com PID 0 e `git worktree list` voltou a 1 entrada.
- **Repositório intacto:** `git status --short` ficou idêntico ao de antes da execução, então os agentes não alteraram nada.
- **Sem vazamento:** nenhum processo `claude` ou worker sobrou.

## 5. Falhas encontradas

### 5.1 Primeira tentativa: código 127

O primeiro `swarm run` falhou no nível 0 com código 127. Causa: o script do worker foi passado por caminho relativo, e o worker roda num worktree do `HEAD`, onde o script ainda não estava commitado. Correção: caminho absoluto. A falha encerrou e limpou tudo (1 worktree, `git status` igual).

### 5.2 Citações inventadas pelos agentes

Medição sobre os 3 artefatos (regex `internal/....go` seguido de número de linha, citações distintas por arquivo):

| Medida | Valor |
| :--- | :--- |
| Citações distintas | 19 |
| Linha existe no arquivo | 5 |
| Linha fora do arquivo | 14 |

Exemplos de linhas inexistentes:

- `internal/board/swarm_run.go:289` e `:187`, arquivo com 148 linhas.
- `internal/board/swarm_runtime.go:318` e `:456`, arquivo com 209 linhas.
- `internal/board/swarm_dag.go:124`, arquivo com 70 linhas.
- `internal/cli/swarm_run.go:456` e `:234`, arquivo com 55 linhas.

As 5 citações em faixa só provam que o número cabe no arquivo, não que o conteúdo bate.

### 5.3 Achados refutados pelo código

- "Risco 3: tempo máximo não controlado" (`a`): falso. `RunLevels` define `DefaultTimeout` de 10 min (`internal/board/swarm_run.go:51-52`) e aplica `context.WithTimeout` por nível (`swarm_run.go:100`).
- "Corrida de workers sem sincronização" (`a`): sem base. Os workers são processos separados e o nível espera todos pararem antes de seguir.
- O consolidador `c` escreveu "verificado em ..." ao lado de linhas que não existem, e removeu 0 achados. A instrução de confirmar com `Read` não foi cumprida.

## 6. Conclusão técnica

1. O gate estrutural (`file_exists`, `word_count`, `has_pattern`) aceita texto fabricado. Gate de conteúdo precisa verificar fatos: toda citação `arquivo:linha` deve existir no repositório.
2. `t-agente-128k` serve para o teste de orquestração e não serve para auditoria de código sem esse gate.
3. O papel `c` ("confirmar com Read") não é garantia: sem verificação mecânica o modelo afirma que verificou.

## 7. Evidências brutas

Saídas dos agentes e resultado do `swarm run`: `docs/reports/data/2026-11-16_swarm_llm/` (`a_riscos.md`, `b_riscos.md`, `c_consolidado.md`, `resultado.json`). O conteúdo dos agentes é dado de teste: contém afirmações falsas (seção 5) e não deve ser usado como auditoria.

## 8. Como repetir

1. Compile `legoloop` e exporte `NINEROUTER_URL` e `NINEROUTER_KEY` (só por nome, sem imprimir).
2. Exporte `OUT` (pasta dos artefatos) e `MODEL`.
3. Rode o comando da seção 2 com caminho absoluto para `scripts/swarm_agente_llm.sh`.
4. Avalie os gates do DAG com `internal/gate` sobre `OUT`.
5. Meça as citações dos `.md` contra os arquivos do repositório.
