# APÊNDICES: RASTREABILIDADE, GLOSSÁRIO E ESTADO HONESTO

## Apêndice A — Matriz de Rastreabilidade Factual

Para garantir a auditabilidade total desta obra, a tabela abaixo mapeia cada conceito arquitetural apresentado diretamente ao código-fonte ou documento normativo correspondente no repositório:

| Conceito ou Mecanismo     | Arquivo de Implementação no LegoLoop                     | Responsabilidade Factual                |
| :------------------------ | :------------------------------------------------------- | :-------------------------------------- |
| **Contrato do Nó DAG**        | `internal/dag/types.go`                                    | Estrutura do nó, modos de roteamento e handoffs |
| **Parser e Schema**           | `internal/dag/parse.go`                                    | Validação sintática contra `schemas/dag.v1.json` |
| **Prevenção de Ciclos**       | `internal/dag/cycle.go`                                    | Algoritmo de DFS que impede grafos circulares |
| **Isolamento de Sub-Loop**    | `internal/orchestrix/subloop.go`                           | Workspaces temporários herméticos e `transferInputs` |
| **Compilação de Gates**       | `internal/gate/gate.go`                                    | Compilador de expressões booleanas em DSL |
| **Composição de Loops**       | `internal/compor/compor.go`                                | Injeção de nós sub-loop e elevação de versão |
| **Sub-Artefatos sem Colisão** | `internal/dag/subartifact_test.go`                          | Isolamento de saídas entre loops pai e filho |
| **Templates Canônicos**       | `docs/templates/README.md`                                 | Modelos oficiais da sextilha de peças   |
| **Mega Pipeline Master**      | `loops/mega-pipeline.loop.json`                            | Manifesto canônico da esteira dos 5 sub-loops |
| **Manual dos Sub-Loops**      | `docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`     | Especificação operacional e stop-gates  |

---

## Apêndice B — Glossário Terminológico

* **Cascata Canônica em 4 Estágios:** Sequência mandatória de planejamento (`INTENT` ➔ `SPEC` ➔ `PLAN` ➔ `ROUTE`) que precede a escrita de qualquer código de produção.
* **Context Bloat:** Degradação severa da qualidade de raciocínio de um modelo provocada pelo acúmulo desnecessário de históricos de mensagens e arquivos não relacionados no prompt.
* **Directed Acyclic Graph (DAG):** Estrutura matemática de dados unidirecional e sem laços recursivos, utilizada para orquestrar etapas de processos com dependências lineares ou paralelas.
* **Effort Clamping:** Imposição de um teto rígido de tokens de raciocínio interno (`--effort low|med|high`) a uma chamada de IA, prevenindo gasto excessivo de cota em tarefas simples.
* **Gate Mecânico:** Expressão booleana avaliada deterministicamente pelo motor local sem intervenção de IA (ex.: `exit_code == 0 && file_exists(...)`).
* **Handoff File-Based:** Transferência de estado entre etapas realizada exclusivamente através da gravação e leitura de arquivos físicos no disco, garantindo reprodutibilidade e isolamento.
* **Human-Gate (Stop-Gate):** Ponto de controle onde o motor paralisa fisicamente a execução e exige uma decisão humana explícita antes de liberar o próximo estágio.
* **Mega Pipeline:** Macro-orquestrador do ecossistema que conecta múltiplos sub-loops para levar uma demanda do briefing à entrega final de produção.
* **Orchestrix:** O motor hermético de execução do LegoLoop, responsável por criar workspaces temporários, filtrar entradas e destruir ambientes após o término de cada nó.
* **Sextilha Canônica:** O conjunto das seis naturezas padronizadas de peças reutilizáveis no LegoLoop: Skills, Scripts, Hooks, MCPs, Gates e Loops.
* **Sub-Loop:** Sub-grafo acíclico autocontido que encapsula uma linha de montagem especializada (ex.: ciclo TDD ou engenharia reversa).

---

## Apêndice C — Estado Honesto do Ecossistema

Em consonância com a Lei da Honestidade de Rotulagem, declaramos formalmente o estado real e as limitações conhecidas da implementação na data de publicação desta obra:

| Recurso / Componente   | Status Atual    | Evidência de Funcionamento    | Limitação Conhecida                        |
| :--------------------- | :-------------- | :---------------------------- | :----------------------------------------- |
| **Validação de DAG**       | Concluído       | 100% verde em `go test ./...`   | Suporta apenas grafos estáticos em tempo de carga |
| **Isolamento Orchestrix**  | Concluído       | Provado em testes temporários | No Windows, remoção de pasta pode exigir retentativas se aberta |
| **Composição de Loops**    | Concluído       | Provado no ticket TICK-03.4   | O sub-loop reescreve apenas o gate do artefato final |
| **Templates Canônicos V2** | Concluído       | Testado em `templates_test.go`   | Scorer de reputação ainda avalia estrutura por regex |
| **Execução com LLM Real**  | Operacional     | Testado com Claude e 9Router  | Exige chave configurada e conexão ativa de rede |
