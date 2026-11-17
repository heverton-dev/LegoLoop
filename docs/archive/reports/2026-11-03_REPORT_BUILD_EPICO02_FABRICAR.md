# 📊 REPORT: BUILD do Épico 02 (CLI `legoloop fabricar` e skill do Chat) — `SG-E02`

> **Data:** 2026-11-03  
> **Escopo:** `TICK-02.1` a `TICK-02.4` do PLAN `2026-11-03_PLAN_LEGOLOOP_FABRICAR.md`  
> **Estado:** 4 tickets `PASSED`; aguardando a autorização humana do `SG-E02` para iniciar o Épico 03.  
> **Ressalva principal:** o ponto-a-ponto com um **LLM real** ainda **não** fechou (seção 4, item 1). Todo o resto foi provado.

---

## 1. Entregas

| Ticket | Entrega | Arquivos principais |
|:---|:---|:---|
| `TICK-02.1` | `legoloop fabricar "<nome>" --yes`: progresso `[1/5]` a `[5/5]`, `--json` (stdout puro), `--dry-run`, `--force`, `--category`, `--vault`, `--docs`, `--trust-removable`, `--link-harness`; exit 2 (nome ou categoria inválidos) e 1 (modelo fora do inventário, `ErrExists`, falha) | `internal/cli/fabricar.go` |
| `TICK-02.2` | Modo interativo (sem `--yes`): lista de drives para desmarcar, perguntas numeradas, até 3 entradas inválidas, resumo e `[s/N]`; EOF ou `N` aborta com exit 1 sem gravar | `internal/cli/fabricar_interactive.go` |
| `TICK-02.3` | `fabricar` registrado no dispatcher e no `help`; `--help` lista as flags; sem argumento imprime o uso e sai com 2; `canonize` intacto e sem despacho por flags | `internal/cli/cli.go`, `dispatch_fabricar_test.go` |
| `TICK-02.4` | Skill `legoloop-fabricar` (cria o diretório `skills/`): entrevista, prévia `--dry-run`, **stop-gate humano**, depois o comando; não cita `canonize` | `skills/legoloop-fabricar/SKILL.md`, `skill_fabricar_test.go` |

Aberturas mínimas no motor (sem mudar o comportamento do Épico 01): `fabricar.Deps.Progress` e a constante `Layers`; `fabricar.Deps.Confirm` e `ErrAborted` (a confirmação precisa ocorrer depois da Tríade pronta e antes de gravar). Na CLI: `Deps.Fabricar *FabricarDeps` (LLM, scanner, relógio, pastas de harness e drives injetáveis).

## 2. Checklist do `SG-E02`

| Item | Resultado |
|:---|:---|
| `legoloop fabricar --help` no terminal | **OK** no binário real; lista as 11 flags |
| `legoloop fabricar "x" --yes --dry-run --vault <TMP> --docs <TMP>` no terminal | **OK contra um 9Router falso local** (binário real, varredura real dos discos, 0 arquivos gravados). **Com o 9Router real da máquina: não fechou** (item 1 da seção 4) |
| `legoloop canonize` legado com saída idêntica | **OK**: `commands.go` e `internal/canon` sem nenhum diff no git; mensagem de `canonize` sem `--root` testada; `canonize` recusa as flags do `fabricar` |
| Skill listada pelo `assets.Scanner` | **OK**, por teste que usa o scanner real sobre a raiz do repositório |
| Testes | `internal/cli` + `interview`: **45 passando, 0 falhas, 0 pulados**; `internal/fabricar`: **68 passando, 0 falhas, 0 pulados** |
| `go test -count=1 ./...` | Verde em todos os 21 pacotes |
| `go vet ./...` | Limpo; cross-vet `GOOS=linux` e `darwin` OK |
| `funclen` | Limpo em `fabricar`, `plangen` e nos meus arquivos de `cli`. `interview/client.go:Chat` tem 53 linhas: já passava de 35 antes (52), a minha linha `stream:false` somou 1 |
| `gofmt` | Meus arquivos limpos; `gofmt -l .` global ainda lista os mesmos 22 arquivos preexistentes (não os reformatei) |
| `go.mod` e `go.sum` | Inalterados |
| `-race` | Indisponível aqui (exige CGO); fica para o CI Linux |
| Autorização explícita do usuário | **Pendente** |

## 3. Evidências além dos testes

* **Mutação (os testes detectam regressão):** (a) a confirmação passou a aceitar `N` como sim: 1 teste falhou; (b) limite de entradas inválidas 3 → 5: 1 teste falhou; (c) retirei o texto do stop-gate da skill: falhou; (d) fiz a skill citar `legoloop canonize`: falhou. Tudo restaurado e verde. Numa primeira leitura achei que (c) e (d) não falhavam, mas era o meu `grep` sobre saída sem `-v`; reexecutei com saída completa e ambos falham.
* **Binário real contra um 9Router falso** (arquivo temporário fora do repositório, já removido): `--dry-run` gravou 0 arquivos; a gravação real criou **8 arquivos**; `legoloop validate` aprovou o DAG (`5 nós válidos, sem ciclos`, exit 0); repetir sem `--force` saiu com **exit 1**; com `--force` saiu 0, criou 8 `.legoloop-bak` e **nenhum** `.legoloop-tmp`. Cofre e docs foram sempre diretórios temporários; o `~/.aidd` real não foi tocado.
* **Teste de ordem na skill:** a prévia com `--dry-run` pode vir antes do gate; a primeira invocação **que grava** tem de vir depois dele.

## 4. Achados e desvios (leia antes de autorizar)

1. **Não consegui fechar o ponto-a-ponto com um LLM real.** Fiz **5 chamadas reais** ao seu 9Router remoto (`https://9router.vpsconexao.org`), o que tem custo. Resultado: com `code-pro` a requisição passou de **4 minutos sem resposta** (estourou o prazo); com `code-fast` respondeu em segundos, mas o texto **violou o contrato** da SPEC ("não contém nenhum critério `AC-XX`"). O `specgen` é de **tentativa única**: qualquer desvio do modelo derruba a fabricação, sem artefato parcial (o contrato protege, mas o usuário fica sem resultado). Parei de insistir para não gastar mais. **Decisão sua:** (a) eu acrescento uma repetição limitada (até 2 retentativas devolvendo o erro de validação ao modelo, no passo da Tríade), que é o que eu recomendo; ou (b) você usa `--model` com um modelo que siga o formato; ou (c) ambos.
2. **Quatro defeitos reais, achados só ao rodar o binário, corrigidos com teste vermelho antes:**
   * `--help` mostrava `-model legoloop doctor`: o pacote `flag` toma o texto entre crases como **nome do valor**.
   * A chave do 9Router **não era enviada**, dando `401`. Agora segue a regra do `legoloop run`: só vai para a URL confiável; teste prova que **não vaza** para outra base nem aparece na saída.
   * `interview.RouterLLMClient` não enviava `"stream": false` e este router responde em SSE por padrão (`invalid character 'd'`). **É um defeito preexistente do cliente compartilhado**, que também afeta as chamadas de entrevista, SPEC e PLAN do **Studio** contra este router. Corrigi na origem (1 linha, com teste).
   * Prazo fixo de 45 s era curto para gerar a SPEC; o cliente do `fabricar` agora usa 5 min.
3. **Mudança de interpretação do `--model`:** a INTENT fala em "seletor de modelo cognitivo", então o modelo escolhido passou a ser o **cérebro** que conduz a fabricação (antes eu fixava `code-pro`, como o Studio). Efeito colateral: modelos que não são do 9Router (Ollama, `claude` CLI) são **recusados com explicação**, em vez de trocados em silêncio. Se você quiser Ollama como cérebro, é trabalho a mais.
4. **Opções numeradas:** a entrevista só expõe pergunta e palpite, então as opções são "aceitar o palpite", "responder com minhas palavras" e "cancelar".
5. **Modo interativo só foi provado com LLM falso.** O 9Router falso não simula a entrevista, então não rodei o interativo no binário real.
6. **Mensagem de `ErrExists` ambígua:** lista só `SKILL.md`, `...loop.json`, `...dag.json` (nomes, não caminhos). Funciona, mas é pouco útil; fica como melhoria pequena.
7. **`--docs` é a raiz do projeto**, e `docs/` é acrescentado; passar `--docs X/docs` gera `X/docs/docs/...`. O help diz isso, mas é fácil errar.
8. **Varredura real:** com a máquina real, `C:\` atinge o limite de 500 peças e ignora milhares de itens além da profundidade 4 (avisos no stderr); levou ~7 s na execução quente.
9. O primeiro comando em segundo plano desta sessão (o teste vermelho inicial) travou por tempo e foi encerrado pelo sistema; ele não é mais necessário e não afeta nenhum resultado.

## 5. Próximo passo

`SG-E02` pede sua autorização explícita para iniciar o **Épico 03** (Studio: `POST /api/fabricar/run`, modal e plotagem no Canvas). **Nada do Épico 03 foi iniciado.** Antes disso, preciso da sua decisão sobre o item 1 (repetição limitada no `specgen`/Tríade), porque o Studio herdaria a mesma fragilidade.
