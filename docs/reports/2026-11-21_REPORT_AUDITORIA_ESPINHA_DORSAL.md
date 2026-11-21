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
| G05 | Alta | Duas fachadas "únicas" de caminhos: `internal/caminhos` (35 importadores) e `internal/paths` (15). Ambas se declaram a fonte única. | `go list` e `grep` de importadores |
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

## 5. Não verificado (precisa de rodada própria)

- Execução real de ponta a ponta (`run` e `exec`) com um DAG só de scripts, rodada duas vezes, comparando o relatório.
- `go test -race` e concorrência.
- Qualquer comportamento em Linux e macOS.
- Quais `os.WriteFile` gravam estado crítico (G11) e quais `init()`/globais são mutáveis (G14).
- Contratos `consumes`/`produces` das skills do Cofre global.
- Desempenho e consumo de memória sob DAG grande.
