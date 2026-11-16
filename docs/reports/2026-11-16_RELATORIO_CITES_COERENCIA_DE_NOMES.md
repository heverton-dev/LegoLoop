# Relatório: coerência de nomes no gate `cites_exist`

Data: 2026-11-16. Continua `2026-11-16_RELATORIO_GATE_CITES_EXIST.md`.

## 1. Veredito

- **Implementado e correto:** `cites_exist` agora também exige que cada identificador entre crases na mesma linha da citação apareça no arquivo citado.
- **Quase inerte no formato real dos agentes:** das 15 citações em faixa coletadas, a checagem reprovou 0. Só 2 das 36 linhas de citação traziam identificador, e os dois existem no arquivo citado.
- **Efeito prático:** quem barra os agentes continua sendo a checagem de intervalo (23 das 38 citações caem fora do arquivo). O identificador sozinho não ancora o sentido.

## 2. O que mudou

`internal/gate/cites.go`:

- `CiteVerdicts(root, artefato, base)` (exportada): devolve, por citação, `InRange` (o arquivo existe e o intervalo cabe nele) e `MissingIdent` (identificadores da mesma linha do artefato ausentes dos arquivos citados nela). Serve a gates e a relatórios.
- Identificador válido: trecho entre crases com nome de código, com ao menos 3 caracteres, que não seja caminho nem número de linha. `Func()` vale `Func`; `pkg.Func` vale `Func`.
- A busca é por palavra inteira (`\bnome\b`) em qualquer ponto do arquivo citado. Se a linha cita 2 arquivos, vale o identificador estar em um deles.
- `cites_exist` passa só se todas as citações têm `OK()`.

## 3. Testes

- `internal/gate/cites_test.go`: 11 casos novos de identificador (presente, ausente, em outra linha do artefato, com `()`, pontuado, curto, caminho entre crases, dois arquivos, pedaço de palavra, texto que não é identificador) e um teste de `CiteVerdicts` que aponta qual citação e qual identificador falharam.
- Mutação: desligando a checagem de identificador, os dois testes novos falham. Com o código restaurado passam.
- Suíte completa (`go test ./... -count=1`) sem falhas.

## 4. Medição

Execução nova de LLM com `-gates` (mesmo DAG e modelo `t-agente-128k`): saída 1, gate reprovou `a` e `b`, `c` não iniciou; os 2 workers ficaram `killed` e o repositório ficou idêntico. Uma entrada extra em `git worktree list` pertence a outra sessão (`.gemini/antigravity-cli/...`), não aos workers.

`CiteVerdicts` sobre todos os artefatos de LLM coletados (3 execuções, 7 arquivos):

| Medida | Valor |
| :--- | :--- |
| Citações | 38 |
| Fora do arquivo (intervalo inexistente) | 23 |
| Em faixa | 15 |
| Em faixa e reprovadas por identificador | 0 |

Por que 0: das 36 linhas do artefato com citação e número de linha, só 2 traziam identificador entre crases na mesma linha (`swarmRunDAG`, que existe em `swarm_tools.go`). O formato dos agentes põe o identificador em outra linha (`Cenário`, `Correção`) e descreve código que deveria existir, não código que existe.

## 5. Conclusão e limite

1. A coerência por identificador só funciona se o agente escrever o nome do código na mesma linha da citação. O prompt atual não exige isso.
2. Incluir as linhas seguintes do bloco no teste reprovaria correções legítimas (`Use WaitGroup` descreve algo ausente por definição). Por isso não foi feito.
3. A âncora de sentido confiável é a citação com trecho literal: o agente escreve `arquivo:linha` e o texto exato da linha entre crases, e o gate confere que o trecho aparece perto da linha citada. Isso ainda não existe.

## 6. Evidências brutas

- `docs/reports/data/2026-11-16_swarm_llm_semantica/`: `a_riscos.md`, `b_riscos.md`, `resultado.json` e `stderr.txt` da execução com o gate de identificador.
- `docs/reports/data/2026-11-16_swarm_llm_gate/` e `docs/reports/data/2026-11-16_swarm_llm/`: execuções anteriores.
- O conteúdo dos agentes é dado de teste e contém afirmações falsas.
