# Relatório: gate `cites_exist` e gates por nó no `swarm run`

Data: 2026-11-16. Continua o relatório `2026-11-16_RELATORIO_SWARM_RUN_AGENTES_LLM.md`, que mostrou 14 de 19 citações `arquivo:linha` inexistentes aprovadas por gates estruturais.

## 1. O que foi entregue

1. **Builtin `cites_exist(artefato, base)`** em `internal/gate/cites.go`.
   - Verdadeiro só se o artefato traz ao menos uma citação `arquivo:linha` e todas existem sob a pasta `base`.
   - Uma citação existe quando o arquivo é regular, fica dentro de `base` (sem `..`, sem caminho absoluto) e o intervalo citado cabe no arquivo.
   - Formatos aceitos: `a.go:12`, `a.go:12-15`, `a.go (linha 12)`, `a.go (função X, linha 12)`, `a.go:line 12`. O número pertence ao caminho mais próximo à esquerda na mesma linha.
   - Sem nenhuma citação, o gate reprova. Acima de 500 citações, reprova.
   - Prova existência, não sentido: a linha pode existir e não dizer o que o texto afirma.
   - Registrado em `builtins`, no catálogo e em `contentCalls` (conta como prova de conteúdo).
2. **Avaliação de gate por nó no `swarm run`** com `-gates` (CLI) e `gates` (MCP `swarm_run`).
   - Depois que o worker sai com código 0, o gate do nó (`output.gate`) é avaliado com raiz no worktree do worker.
   - Gate falso reprova o nível, encerra a execução e limpa workers e worktrees. O nível seguinte não começa.
   - A saída traz `gate_ok` por worker. Sem `-gates`, o comportamento antigo vale (gate do nó ignorado).

Antes desta entrega, o `swarm run` não avaliava gate nenhum.

## 2. Testes

- `internal/gate/cites_test.go`: 17 casos (citação válida, fora do arquivo, arquivo inexistente, forma "linha N", intervalo, duas citações na mesma linha, travessia `../`, sem citação, limite de 500, `base`, erros de compilação). O catálogo e a aridade continuam verdes.
- `internal/cli/swarm_run_e2e_test.go`: gate aprova com todos os artefatos, gate reprova um nível (worker com código 0 sem o artefato), e sem `-gates` o gate é ignorado.
- Mutação: desligando a avaliação de gate, o teste de reprovação falha (`code=0`, esperado 1). Com o código restaurado passa.
- Suíte completa (`go test ./... -count=1`) sem falhas.

## 3. Prova 1: os 3 artefatos do relatório anterior

Avaliação direta de `cites_exist(<artefato>, '.')` com raiz no repositório:

| Artefato | Resultado |
| :--- | :--- |
| `a_riscos.md` | falso |
| `b_riscos.md` | falso |
| `c_consolidado.md` | falso |
| controle positivo (citações verdadeiras: `swarm_run.go:51-52`, `:100`, `swarm_dag.go (linha 70)`) | verdadeiro |

## 4. Prova 2: nova execução com agentes de LLM e `-gates`

Mesmo DAG de auditoria (`scripts/fixtures/swarm_llm_auditoria.dag.json`), agora com `cites_exist` no gate de cada nó, `claude -p` via 9Router (`t-agente-128k`, ferramentas só de leitura).

- Comando: `legoloop swarm run <dag> -dir <quadro> -repo <repo> -timeout 300s -gates -- bash scripts/swarm_agente_llm.sh`.
- Resultado: **saída 1**, `nivel 0: gate reprovou em 2 no(s)`: `a` e `b`.
- Nível 0: `a` e `b` com código 0 e `gate_ok=false`, 15 s (17:30:39 a 17:30:54 UTC). O nível 1 (`c`) não iniciou.
- Limpeza: 1 entrada em `git worktree list` e `git status` idêntico ao de antes.

Causa da reprovação (os outros termos do gate eram verdadeiros: `a` com 226 palavras e `b` com 201, ambos com "Risco 1"):

- `a`: 3 citações no formato `arquivo.go:line N`. `swarm_runtime.go:387` (arquivo com 209 linhas) e `swarm_run.go:210` (148 linhas) não existem. `swarm_run.go:124` existe.
- `b`: 4 citações, 3 inexistentes: `swarm_dag.go:352` (70 linhas), `cli/swarm_run.go:124` e `:187` (56 linhas).

O modelo continuou inventando linhas nesta execução, e agora o gate impede que o texto passe adiante.

## 5. Limites

1. `cites_exist` prova que a linha existe, não que o texto sobre ela é verdadeiro. Uma citação em faixa e errada passa.
2. Só valida citação com caminho e número de linha. Menção a arquivo sem número não é citação.
3. A extração é por heurística de texto (extensões `go`, `md`, `json`, `sh`, `py`, `js`, `ts`, `tsx`, `css`, `html`, `yaml`, `toml`, `mod`, `sql`). Formato de citação fora dos listados passa despercebido.
4. A avaliação de gate usa a raiz do worktree do worker. Artefato escrito fora dele precisa ser copiado para lá (o script faz `cp` para `out/`).

## 6. Evidências brutas

- `docs/reports/data/2026-11-16_swarm_llm_gate/`: `a_riscos.md`, `b_riscos.md`, `resultado.json`, `stderr.txt` da execução com `-gates`.
- `docs/reports/data/2026-11-16_swarm_llm/`: saídas da execução anterior, sem gate de citação.
- O conteúdo dos agentes é dado de teste e contém afirmações falsas.
