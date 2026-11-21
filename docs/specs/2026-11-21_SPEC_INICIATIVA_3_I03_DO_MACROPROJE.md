---
id: SPEC-intv_1795274711307895500
status: APPROVED
date: 2026-11-21
schema_version: 2
upstream: [INTENT-intv_1795274711307895500]
covers: [R1, R2, R3, R4, R5, R6]
---

# 📋 SPEC: Despachante de executores no runner (I03)

> **Status:** Aprovado pelo Usuário (Stop-Gate 02, 2026-11-21)  
> **Data:** 2026-11-21  
> **Versão:** 1.0.0  
> **Autor:** Claude (Fusão Addy Osmani + Matt Pocock + AIDD-Spec)  
> **Iniciativa:** I03 — Despachante de executores no runner  
> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md](../intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md)  
> **Alvo Primário:** internal/dag, internal/runner, internal/board, internal/plangen, internal/routedag, internal/cli  
> **Consome:** INTENT aprovada · **Produz:** PLAN, ROUTE e DAG  
> **Template:** TEMPLATE_SPEC.md v2.0.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)

> [!IMPORTANT]
> Alerte no Stop-Gate 02 caso alguma premissa esteja incorreta.
1. Runtime: Go conforme o `go.mod`, só biblioteca padrão, zero dependências novas.
2. Fatos conferidos no código antes desta SPEC:
   * `go list` mostra que `internal/runner` importa `dag`, `gate` e `triagem`, e `internal/board` importa `dag`, `gate` e `taxonomy`. O `runner` importar o `board` não cria ciclo.
   * `Orchestrix.Execute` (`internal/orchestrix/orchestrix.go:104`) é o único executor e já despacha pelo tipo do nó: `script` roda local, `skill` roda pelo LLM. Não existe um segundo executor "local".
   * `isolation.executor` é gravado pelo `plangen` e lido por nenhum código de produção.
   * O Swarm executa por nível (`board.LevelsFromDAG` e `SwarmBoard.RunLevels`). Um nó só vira worker se não for `stop_gate`, tiver papel (`role` ou `papel`) e tiver `command`.
   * O vocabulário que existe nos documentos é `orchestrix_hermetic` (1.063 ocorrências), `local_process` (144), `native_swarm` (5) e `human_in_the_loop` (4). Os dois últimos só aparecem em ROUTEs arquivadas (`docs/archive/quinteto/routes`) e em nenhum código.
3. A I03 só executa código depois que a I02 fechar (mexe em `internal/routedag` e na ROUTE).
4. Nenhum LLM entra no despacho nem na validação.

### Decisões do Stop-Gate 02 (pergunta, sugestão e motivo)
1. **Nome do valor do Swarm (corrige a INTENT).** A INTENT aprovada cria o valor `swarm_worker`. Fato verificado: o valor `native_swarm` já existe nos documentos. **Sugestão:** usar `native_swarm` e não criar `swarm_worker`. **Motivo:** dois nomes para o mesmo destino é vocabulário duplicado, e a lei 15 manda um termo por significado. O conjunto aceito passa a ser `local_process`, `orchestrix_hermetic` e `native_swarm`. O valor `human_in_the_loop`, que só existe em ROUTE arquivada, passa a ser recusado.
2. **`local_process` e `orchestrix_hermetic` (corrige a INTENT).** A INTENT diz que o primeiro roda "no runner local" e o segundo "pelo Orchestrix". Fato verificado: há um só executor, e ele já despacha pelo tipo do nó. **Sugestão:** os dois valores continuam no mesmo executor das `Options`, sem criar um segundo executor. **Motivo:** um segundo executor seria código novo sem comportamento diferente a provar, e separar os dois é refatoração fora do escopo. A validação aceita os dois, e a diferença futura fica registrada como pendência.
3. **Executor ausente.** Muitos DAGs antigos não têm `isolation.executor`. **Sugestão:** ausente vale `orchestrix_hermetic`, que é o comportamento de hoje. Valor presente e desconhecido continua recusado, sem fallback. **Motivo:** preserva a paridade sem violar a regra "sem fallback" da INTENT, porque ausente não é desconhecido.
4. **Pré-condições do `native_swarm`.** **Sugestão:** a validação prévia recusa `native_swarm` em nó `stop_gate`, sem papel ou sem `command`, com o mesmo erro tipado. **Motivo:** o `board` ignoraria o nó em silêncio ou falharia no meio do nível.
5. **Ticket do plangen.** **Sugestão:** o `Ticket` do `plangen` ganha o campo opcional `executor`, validado pela mesma regra, e é assim que o planejamento escolhe o Swarm de propósito. **Motivo:** hoje o `plangen` não tem como emitir `native_swarm`.

---

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)

* **Sempre Fazer:**
  * TDD estrito (teste falhando antes do código de produção).
  * Funções com menos de 35 linhas (`go run ./tools/funclen internal`).
  * Validar todos os nós antes de executar o primeiro.
  * Recompilar o binário ao fim de todo ticket com Go (AGENTS.md, regra 5).
  * Executar o Swarm só sob as seis condições da lei 4 do `AGENTS.md`.
* **Perguntar Antes:**
  * Adicionar dependência ao `go.mod`.
  * Alterar contratos públicos de CLI, do DAG ou de `ROUTE.json` além do que esta SPEC descreve.
  * Alterar regras globais do `AGENTS.md`.
* **Nunca Fazer:**
  * Alterar a INTENT da I01 ou a da I02.
  * Usar LLM ou heurística de tamanho ou carga na escolha do executor.
  * Cair para um executor padrão quando o valor é desconhecido.
  * Mudar o comportamento dos nós `orchestrix_hermetic` e `local_process` (paridade).
  * Criar função vazia, `TODO` ou mock em código de produção (Lei #1).
  * Commitar ou publicar sem o usuário pedir.

---

## 3. Visão Geral, Arquitetura & Costuras de Teste

O `runner` hoje executa todo nó pelo mesmo `Executor` das `Options`. A I03 acrescenta duas coisas: a validação dos executores antes do nível 0 e o desvio dos nós `native_swarm` de cada nível para o Swarm.

### Diagrama de Fronteiras (Mermaid)
```mermaid
flowchart LR
  D[DAG] --> V[dag.ValidarExecutores]
  V -->|erro tipado, todos os nós| X[recusa antes do nível 0]
  V -->|ok| L[para cada nível]
  L --> A[nós local_process, orchestrix_hermetic ou ausente]
  L --> S[nós native_swarm]
  A --> E[Options.Executor]
  S --> W[Options.Swarm]
  W --> B[board.RunLevels]
  E --> R[Report]
  B --> R
```

### Costuras de Teste (Test Seams — Pocock)
* `dag.ValidarExecutores(d)`: função pura sobre um `*dag.DAG` montado no teste.
* `runner.Options.Swarm`: interface `SwarmRunner` com o método `RunLevels(ctx, [][]board.RunNode, board.RunOptions) ([]board.RunLevelResult, error)`. O teste injeta um falso e conta as chamadas.
* `runner.Options.Executor`: o teste injeta um falso que conta as chamadas por nó.
* `board.RunNodeOf(n, global, usados)`: função extraída de `LevelsFromDAG`, testada isolada.
* `plangen.Render` e `routedag.Projetar`: fixtures em `t.TempDir()`.

---

## 4. Contratos de Dados & Schemas Formais

**Valores aceitos de `isolation.executor`:** `local_process`, `orchestrix_hermetic`, `native_swarm`. Ausente vale `orchestrix_hermetic`.

**Erro tipado (`internal/dag`):** `ErrExecutorDesconhecido`, usado com `errors.Is`. A mensagem lista uma linha por nó inválido:
```
nó TICK-03.1: isolation.executor "xyz" desconhecido; aceitos: local_process, orchestrix_hermetic, native_swarm
nó TICK-04.2: native_swarm exige papel e command; nó stop_gate não vai ao Swarm
```
**Regras de `subagents_allowed`:** `false` vale para qualquer executor; `"ephemeral"` só vale com `native_swarm`.

**Interface do `runner`:**
```go
type SwarmRunner interface {
    RunLevels(ctx context.Context, levels [][]board.RunNode, opt board.RunOptions) ([]board.RunLevelResult, error)
}
```
Campo novo: `Options.Swarm SwarmRunner`. Com `Swarm` nulo e nenhum nó `native_swarm`, o comportamento é idêntico ao de hoje.

**Ticket do `plangen`:** campo opcional `"executor"` no JSON de Tickets.

---

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

> [!NOTE]
> Cada `AC-XX` é verificável por teste ou comando com `exit code == 0`. A ordem dos Épicos é a ordem de implementação.

### Critérios de Aceitação por Épico

### Épico 01 — Vocabulário e validação (`internal/dag`)

* [ ] `AC-E01.1`: **Conjunto aceito e ausente**
  * **Dado que:** um DAG com nós de executor `local_process`, `orchestrix_hermetic`, `native_swarm` e um nó sem `isolation.executor`.
  * **Quando:** rodar `go test ./internal/dag/ -run TestExecutorAceito`.
  * **Então:** o teste sai com 0 e prova que os quatro casos passam e que `xyz` e `human_in_the_loop` são recusados.
* [ ] `AC-E01.2`: **Todos os nós inválidos de uma vez**
  * **Dado que:** um DAG com três nós de executor inválido em níveis diferentes.
  * **Quando:** rodar `go test ./internal/dag/ -run TestValidarExecutoresListaTodos`.
  * **Então:** o erro satisfaz `errors.Is(err, ErrExecutorDesconhecido)` e a mensagem traz uma linha por nó, com o nó, o valor recebido e os valores aceitos.
* [ ] `AC-E01.3`: **Regra do `ephemeral`**
  * **Dado que:** nós com `subagents_allowed` igual a `false` e a `"ephemeral"`, em executores diferentes.
  * **Quando:** rodar `go test ./internal/dag/ -run TestSubagentesEfemerosSoNoSwarm`.
  * **Então:** `"ephemeral"` com `native_swarm` passa e `"ephemeral"` com os outros dois executores é recusado com o mesmo erro tipado.
* [ ] `AC-E01.4`: **Pré-condições do `native_swarm`**
  * **Dado que:** nós `native_swarm` do tipo `stop_gate`, sem papel e sem `command`.
  * **Quando:** rodar `go test ./internal/dag/ -run TestNativeSwarmExigePapelECommand`.
  * **Então:** os três são recusados com o erro tipado e um nó com papel e `command` passa.
* [ ] `AC-E01.5`: **`legoloop validate` recusa**
  * **Dado que:** o binário recompilado e um `pipeline.dag.json` com executor `xyz`.
  * **Quando:** rodar `legoloop validate` nesse arquivo.
  * **Então:** o comando sai com código diferente de 0 e imprime a linha do nó inválido.

### Épico 02 — Contrato no `plangen` e no `routedag`

* [ ] `AC-E02.1`: **Ticket escolhe o Swarm**
  * **Dado que:** um `tickets.json` com um Ticket de `"executor": "native_swarm"` e os demais sem o campo.
  * **Quando:** rodar `go test ./internal/plangen/ -run TestTicketExecutor`.
  * **Então:** a ROUTE gerada tem `isolation.executor` igual a `native_swarm` nesse nó e igual ao valor de hoje nos outros.
* [ ] `AC-E02.2`: **Ticket com executor desconhecido é recusado**
  * **Dado que:** um Ticket com `"executor": "xyz"`.
  * **Quando:** rodar `go test ./internal/plangen/ -run TestTicketExecutorDesconhecido`.
  * **Então:** `plangen.Render` devolve o erro tipado e não grava PLAN nem ROUTE.
* [ ] `AC-E02.3`: **A projeção do DAG preserva e valida**
  * **Dado que:** uma ROUTE com um nó `native_swarm` e outra com um nó `xyz`.
  * **Quando:** rodar `go test ./internal/routedag/...`.
  * **Então:** a primeira projeta um DAG que valida com `native_swarm` preservado, e a segunda é recusada com o erro tipado.

### Épico 03 — Despachante no `runner`

* [ ] `AC-E03.1`: **Recusa antes de executar qualquer nó**
  * **Dado que:** um DAG de três níveis em que só um nó do último nível tem executor inválido, e um `Executor` falso que conta chamadas.
  * **Quando:** rodar `go test ./internal/runner/ -run TestRunRecusaAntesDoNivelZero`.
  * **Então:** `runner.Run` devolve o erro tipado com o nó listado e o `Executor` falso registra zero chamadas.
* [ ] `AC-E03.2`: **Despacho por nível**
  * **Dado que:** um DAG de dois níveis, cada um com nós `local_process`, `orchestrix_hermetic` e `native_swarm`, um `Executor` falso e um `SwarmRunner` falso.
  * **Quando:** rodar `go test ./internal/runner/ -run TestDespachoPorNivel`.
  * **Então:** o `SwarmRunner` recebe em cada chamada só os nós `native_swarm` do nível, o `Executor` recebe os demais, e os níveis rodam na ordem topológica.
* [ ] `AC-E03.3`: **Nível misto termina com os dois lados**
  * **Dado que:** um nível com um nó local e um `native_swarm`, o do Swarm mais lento.
  * **Quando:** rodar `go test ./internal/runner/ -run TestNivelMistoEsperaOsDoisLados`.
  * **Então:** os dois começam sem esperar um pelo outro e o nível seguinte só começa depois que ambos terminam.
* [ ] `AC-E03.4`: **Falha no Swarm interrompe os níveis seguintes**
  * **Dado que:** um `SwarmRunner` falso que reprova o worker de um nó.
  * **Quando:** rodar `go test ./internal/runner/ -run TestFalhaNoSwarmMarcaSkipped`.
  * **Então:** o `Report.OK` é falso, o nó está `error` ou `gate_failed`, e os nós dos níveis seguintes estão `skipped`.
* [ ] `AC-E03.5`: **DAG sem `native_swarm` não toca no Swarm**
  * **Dado que:** um DAG só com `local_process`, `orchestrix_hermetic` e executor ausente, e `Options.Swarm` nulo.
  * **Quando:** rodar `go test ./internal/runner/ -run TestSemNativeSwarmNaoUsaSwarm`.
  * **Então:** a execução termina sem pânico, o `Report` é o de hoje e nenhuma chamada ao Swarm acontece.
* [ ] `AC-E03.6`: **`legoloop run` recusa sem criar `worktree`**
  * **Dado que:** o binário recompilado, um repositório git de teste e um DAG com executor `xyz`.
  * **Quando:** rodar `legoloop run` nesse DAG e depois `git worktree list`.
  * **Então:** o `run` sai com 1, imprime todos os nós inválidos e a lista de `worktree` não ganha entrada.

### Épico 04 — Ligação real com o `board`

* [ ] `AC-E04.1`: **`RunNodeOf` extraído sem mudar o `board`**
  * **Dado que:** o corpo do laço de `board.LevelsFromDAG`.
  * **Quando:** rodar `go test ./internal/board/...`.
  * **Então:** os testes atuais passam, `LevelsFromDAG` chama `RunNodeOf`, e um teste novo cobre o nó sem papel e o nó sem `command`.
* [ ] `AC-E04.2`: **`legoloop run` entrega ao Swarm real**
  * **Dado que:** o binário recompilado e um DAG de um nível com um nó `native_swarm`, papel e `command` de um processo falso que sai com 0.
  * **Quando:** rodar `go test ./internal/cli/ -run TestRunNativeSwarmBinario`.
  * **Então:** o `run` sai com 0, o relatório mostra o nó `passed` e o `SwarmBoard` registra o worker.
* [ ] `AC-E04.3`: **Sem `native_swarm`, nenhum diretório do Swarm**
  * **Dado que:** o binário recompilado e um DAG só com executores locais.
  * **Quando:** rodar `legoloop run` e listar o diretório padrão do Swarm.
  * **Então:** o diretório não é criado.

### Épico 05 — Paridade e governança viva

* [ ] `AC-E05.1`: **Todos os DAGs do repositório continuam válidos**
  * **Dado que:** os DAGs de `docs/dags` e os `pipeline.dag.json` de `processos`.
  * **Quando:** rodar `go test ./internal/dag/ -run TestDAGsDoRepositorioValidam`.
  * **Então:** o teste sai com 0 e nenhum DAG é recusado pela nova validação.
* [ ] `AC-E05.2`: **Relatório idêntico para DAG sem `native_swarm`**
  * **Dado que:** um DAG de referência com os dois executores atuais e um `Executor` falso determinístico.
  * **Quando:** rodar `go test ./internal/runner/ -run TestParidadeDoRelatorio`.
  * **Então:** o `Report` é igual ao registrado antes da mudança (golden gravado no Red).
* [ ] `AC-E05.3`: **Documentos vivos**
  * **Dado que:** o fim da I03.
  * **Quando:** rodar `wc -c AGENTS.md` e `grep -c "native_swarm" AGENTS.md README.md MEMORY.md`.
  * **Então:** o `AGENTS.md` tem no máximo 24000 bytes e cada arquivo cita `native_swarm`.
* [ ] `AC-E05.4`: **Anticorpos e prompt de execução**
  * **Dado que:** cada achado corrigido na I03.
  * **Quando:** rodar `legoloop anticorpos status legoloop-concepcao` e `ls docs/prompts`.
  * **Então:** o comando sai com 0 e existe `docs/prompts/*_PROMPT_EXECUCAO_I03.md`.

### Requisitos Não Funcionais

* [ ] `AC-NF.1`: **Zero LLM no despacho**
  * **Dado que:** `internal/dag` e o código novo do `internal/runner`.
  * **Quando:** rodar `grep -rln "LLMClient\|ResolvedorLLM" internal/dag internal/runner`.
  * **Então:** o `grep` não acha nenhum arquivo.
* [ ] `AC-NF.2`: **Portão do projeto íntegro**
  * **Dado que:** todos os Épicos entregues.
  * **Quando:** rodar os comandos da seção 6.
  * **Então:** todos saem com 0 e `gofmt -l .` não imprime nada.

### Rastreabilidade Requisito para Critério

| R | Requisito | AC |
| :--- | :--- | :--- |
| `R1` | Vocabulário e validação de executores | `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4`, `AC-E01.5` |
| `R2` | Contrato no `plangen` e no `routedag` | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` |
| `R3` | Validação prévia e despacho por nível no `runner` | `AC-E03.1` a `AC-E03.6`, `AC-NF.1` |
| `R4` | Ligação real com o `board` | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` |
| `R5` | Paridade dos executores atuais | `AC-E05.1`, `AC-E05.2` |
| `R6` | Governança viva e portão | `AC-E05.3`, `AC-E05.4`, `AC-NF.2` |

> Cada `R` tem ao menos um `AC`. O `covers` do front matter lista todos os `R`.

---

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)

### Comandos de Verificação Mecânica
* **Build:** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Test:** `go test -count=1 ./...`
* **Estática e formatação:** `go vet ./...` e `gofmt -l .` (saída vazia).
* **Funções:** `go run ./tools/funclen internal`
* **Taxonomia:** `go run ./tools/taxlinter internal cmd loops skills docs/templates`
* **Dependências intactas:** `git diff --exit-code go.mod`

### Execução pelo Swarm
Cada Ticket vira um worker do Swarm sob as seis condições da lei 4 do `AGENTS.md`: contexto só com os artefatos do Ticket, ferramentas listadas, um worker por Ticket, saída em disco com Gate mecânico, sem sub-workers e `write_boundary` `strict`.

### Critério de Done Global
A especificação é **100% CONCLUÍDA** quando:
- [ ] 100% dos critérios `AC-E{N}.{M}` e `AC-NF.{K}` estiverem atendidos.
- [ ] O comando de **Test** passar com exit code 0.
- [ ] `go vet`, `gofmt -l`, `funclen` e `taxlinter` não reportarem nada.
- [ ] O comando de **Build** recompilar o binário.
- [ ] Zero arquivos temporários ou órfãos residuais.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**  
*As premissas, as fronteiras, os contratos, os critérios binários e as cinco decisões acima atendem à sua expectativa?*  
*Somente após a aprovação explícita do usuário, gerar o PLAN (`TEMPLATE_PLAN.md`) e a ROUTE (`TEMPLATE_ROUTE.json`).*
