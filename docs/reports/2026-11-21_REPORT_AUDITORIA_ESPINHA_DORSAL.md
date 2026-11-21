# Auditoria da espinha dorsal do LegoLoop (insumo da I04)

Data: 2026-11-21. Escopo: verbos, Processos, motores e requisitos. Método: só leitura, sem subagentes, cada achado com o comando que o provou.
Limite honesto: é uma auditoria de amplitude com medições reais, não uma leitura linha a linha dos 51 pacotes. O que não foi medido está na seção 5.

## 1. O que está sólido (medido)

| Prova | Resultado |
| :--- | :--- |
| `GOOS=linux` e `GOOS=darwin`: `go build ./...` e `go vet ./...` | limpos |
| `go test -count=2 -shuffle=on ./...` | sem falha (ordem dos testes não importa) |
| `legoloop validate` nos 33 Processos | 33 de 33 passam |
| Nós de Processo sem Gate (fora `stop_gate`) | 0 de 198 |
| Saída de 6 verbos de leitura rodados 2 vezes (`validate`, `listar`, `matriz status`, `doctor`, `skills-custo`, `anticorpos status`) | hashes idênticos |
| Stubs, `TODO` e funções vazias em código de produção | 0 |
| Cobertura: `dag` 96,4%, `roteamento` 99,3%, `matriz` 92,7%, `plangen` 91,9%, `compor` 91,5%, `gate` 90,9%, `interview` 89,4% | alta |

## 2. Lacunas, com severidade e prova

| ID | Sev. | Lacuna | Prova |
| :--- | :--- | :--- | :--- |
| G01 | Crítica | A Matriz de peças, única fonte do `fabricar`, está defasada e nada a atualiza sozinha: 505 novas, 2 alteradas, 6.146 removidas, 21 esperando o cérebro. | `legoloop matriz status` |
| G02 | Crítica | A premissa do produto (usar as peças do SO e dos discos) não é cumprida pelo `fabricar`. Ele ignora `TargetDrives` e avisa "O fabricar não varre drives. As peças vêm da matriz". O comentário do campo ("Vazio = todos os drives do SO") é falso. A descoberta só acontece por `canonize --root` explícito e pelo Cofre. | `internal/fabricar/engine.go:198`, `request.go:27`, `fonte_matriz.go:19` |
| G03 | Alta | `isolation.executor` é decorativo. Nenhum código o lê. Os rótulos não acompanham o tipo do nó: 69 nós `script` estão como `orchestrix_hermetic`; `stop_gate` tem 3 rótulos (6 `local_process`, 5 `orchestrix_hermetic`, 1 ausente); 20 de 198 nós não têm o campo. | contagem tipo x executor em `processos/*/pipeline.dag.json` |
| G04 | Alta | Vocabulário de tipo de nó duplicado: o DAG usa `unit_no`, a ROUTE e o README usam `unit_task`. `TypeSubloop` e `TypeSubprocesso` são a mesma string. O `runner` usa o literal `dag.NodeType("unit_no")` em vez da constante. | `internal/dag/types.go:26-32`, `internal/runner/runner.go:195` |
| G05 | Alta | **(corrigido na SPEC da I04)** `internal/paths` é uma fachada: constantes de nome, 4 funções que só delegam a `caminhos` (2 pontos de chamada) e um jail de caminhos hostis que nenhum arquivo de produção chama. Os importadores se dividem entre `caminhos` (35) e `paths` (15). | `grep` de símbolos de `paths` em `internal`, `cmd` e `tools` |
| G06 | Média | `internal/vault` (Cofre de Segredos) não tem importador em produção; só o próprio teste o usa. | `grep -rl 'internal/vault"'` |
| G07 | Média | Universalidade só compilada, nunca executada em Linux ou macOS. Literais do Windows no código comum: `C:\Windows\Fonts` (`videorea/cartoon_render.go:193`), `python` fixo (`assets/argv.go:119`; em Linux comum é `python3`), `powershell` e `cmd /c` (`assets/argv.go:115,123`). | `grep` de literais |
| G08 | Média | 40 verbos no mesmo nível, sem hierarquia: os da espinha (`fabricar`, `compor`, `run`, `exec`, `validate`, `canonize`, `fechar`) se misturam com os de domínio (`video-motion`, `video-rea`, `canal`, `cartoon-shorts`, `web-qa`, `legopress`). `matrix` está descontinuado e ainda aparece no help. | `legoloop` (help) |
| G09 | Média | Nascimento 5 de 5 incompleto: 3 Processos sem `ANTICORPOS.md` (`legoloop-agent-desktop`, `legoloop-matriz`, `listar`). 10 de 33 Processos fogem do prefixo `legoloop-` (`listar`, `fechar`, `addy-osmani`...). | listagem de `processos/` |
| G10 | Média | Cobertura baixa justamente no motor do Swarm e na triagem: `board` 65,3%, `triagem` 74,5%, `cli` 77,8%. | `go test -cover` |
| G11 | Média | Escritas de estado não atômicas: 6 em `fechar`, 6 em `cli`, 4 em `brain`, 4 em `orchestrix` com `os.WriteFile`; só 13 arquivos usam `os.Rename`. Não verificado quais gravam estado crítico. | `grep` por pacote |
| G12 | Baixa | Erros ignorados (`_ =`): `cli` 36, `studio` 31, `board` 30, `runner` 18, `orchestrix` 15. Exemplo: `saveRunTelemetry` ignora falha ao gravar. | `grep` por pacote |
| G13 | Baixa | Não existe `.github/workflows` (sem CI). A I01 (Épico 02) já cobre. | `ls .github/workflows` |
| G14 | Baixa | 27 funções `init()` e 51 variáveis globais de tipo mapa, fatia ou ponteiro. Não verificado se são mutáveis. | `grep` |

## 3. Efeito nas iniciativas planejadas

| Iniciativa | Lacunas que a tocam | Efeito |
| :--- | :--- | :--- |
| I01 (Forja) | G01, G02, G05, G07, G13 | A Forja instala e impõe o SO, mas o `fabricar` não descobre peças. Sem resolver G01 e G02, o produto instalado nasce com a Matriz defasada. A fachada de caminhos (G05) precisa de uma só dona antes do renomear do módulo. |
| I02 (especificação e planejamento) | G04, G09 | O vocabulário `unit_task` e `unit_no` precisa ser unificado antes de a triagem e os Processos novos o citarem. Os Processos novos nascem 5 de 5. |
| I03 (despachante) | G03, G04, G10 | A SPEC da I03 assume que o executor é um rótulo com significado. Não é. A validação por valor sem semântica por tipo de nó reproduz o problema. O `board` tem a menor cobertura e é o alvo da I03. |

## 4. O que isso diz sobre "relógio suíço"

- A parte **determinística de leitura e de geração** está boa e medida (seção 1).
- A parte **de descoberta** (G01, G02) e a parte **de vocabulário e contrato** (G03, G04, G05) não são determinísticas nem universais hoje. É onde a espinha dorsal falha.

## 4A. Correções e verificações feitas depois da SPEC da I04 (2026-11-21)

**Correções ao que está acima**
- **Correção de G05:** ver a linha da tabela. `paths` não era uma segunda fonte independente.
- **Correção de G01 e G02:** a detecção de Matriz defasada por hash `sha256` já existe (`matriz.Staleness`) e o varredor `assets.Scanner` também. A lacuna é que o `fabricar` nunca atualiza a Matriz e ignora `TargetDrives`.
- **Correção de G14:** `dag`, `runner`, `plangen`, `triagem`, `upstream` e `board` não têm `init()` nem variável global mutável. `gate`, `orchestrix` e `interview` têm só tabelas de consulta, e os `init()` só registram mensagens. A lacuna caiu para "sem achado".

**Lacunas novas, achadas ao executar em Linux (WSL Ubuntu, go1.27.0)**

| ID | Sev. | Lacuna | Prova |
| :--- | :--- | :--- | :--- |
| G15 | Alta | `go build` e `go vet` passam, mas `go test ./...` falha em 9 testes de 5 pacotes. Três são defeitos reais: o jail aceita `..\fora.txt`, `C:\...` e `\\server\...`; o lock da memória da triagem não é exclusivo; a guarda de deleção aceita `docs\vocab\dicionario.json`. Seis dependem de ferramenta instalada (`pandoc`, `pymupdf`, harness). | `go test -count=1 ./...` no WSL |
| G16 | Alta | 14 nós de Processo chamam `powershell` (em `legoloop-triagem`, `legoloop-video-master`, `legoloop-canal-redc`) e 14 chamam `python`: esses Processos só rodam no Windows. | contagem do primeiro argumento de `command` nos DAGs |
| G17 | Média | `processos/fechar` (v1.0.0, 6 nós) duplica `processos/legoloop-fechar` (v1.1.0, 8 nós) e guarda um `triagem.json.legoloop-bak`. Por isso o prefixo não pode ser inserido em `fechar`. | `ls` e leitura dos dois DAGs |
| G18 | Baixa | O relato do Andon diz "Nó executado com exit code 0" quando o nó foi recusado antes de executar (sem `--allow-scripts`). | `internal/runner/sbar.go:37` e execução real |
| G19 | Baixa | Arquivo temporário rastreado: `scratch/i6_inventario_legados.json`. | `git ls-files scratch` |

**Verificado e sem defeito**
- Execução de ponta a ponta duas vezes (DAG de 2 nós `script`): relatório idêntico, descontadas as durações, e artefatos com o mesmo hash. Só `flight_record.json` e `memoria_execucao.json` mudam, pois guardam horário e sessão.
- `go test -race` em Linux, 12 pacotes da espinha: nenhuma corrida de dados. A única falha é o lock da triagem, que é funcional.
- O Gate de aceite do projeto aprova as SPEC da I02, da I03 e da I04.
- Gravações de estado crítico não atômicas, por ponto de chamada: `brain/antibodies.go:295,300`, `fechar/fechar.go:295,308,321` (documentos do usuário), `fechar/quarteto.go:95`, `fechar/destilar.go:113`, `cli/sync_harnesses.go:308` (configuração de harness do usuário), `cli/validate_studio.go:57` e `cli/rota.go:251`.
- Custo do mutador: cerca de 1.960 pontos de mutação nos sete pacotes da espinha; estimativa de 1 a 3 horas por rodada completa em um núcleo (estimativa, não medição).
- Contratos `consumes`/`produces`: as skills `aidd-spec`, `aidd-planner` e `aidd-tickets` não declaram esses campos no front matter.

## 5. Não verificado (precisa de rodada própria)

Resolvido na seção 4A: execução de ponta a ponta duas vezes, `-race` em Linux, estado crítico das gravações, globais e `init()`, contratos das skills do Cofre. Continua não verificado:
- Qualquer execução em macOS (não há máquina; o CI da I04 é o juiz).
- Unidades de rede e removíveis na descoberta.
- Desempenho e memória sob DAG grande.

Lista original, mantida para o histórico:

- Execução real de ponta a ponta (`run` e `exec`) com um DAG só de scripts, rodada duas vezes, comparando o relatório.
- `go test -race` e concorrência.
- Qualquer comportamento em Linux e macOS.
- Quais `os.WriteFile` gravam estado crítico (G11) e quais `init()`/globais são mutáveis (G14).
- Contratos `consumes`/`produces` das skills do Cofre global.
- Desempenho e consumo de memória sob DAG grande.
