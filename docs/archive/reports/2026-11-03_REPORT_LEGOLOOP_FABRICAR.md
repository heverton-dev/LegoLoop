# 📊 REPORT FINAL: `legoloop fabricar` — fábrica de processos (CLI, Chat e Studio)

> **Data:** 2026-11-03  
> **Demanda:** `docs/intents|specs|plans|routes/2026-11-03_*_LEGOLOOP_FABRICAR`  
> **Estado:** 18 tickets `PASSED`, `SG-E01`, `SG-E02` e `SG-E03` autorizados. `SG-FINAL` aprovado pelo usuário.  
> **Relatórios por épico:** `..._REPORT_BUILD_EPICO01_FABRICAR.md`, `..._EPICO02_...`, `..._EPICO03_...`

---

## 1. O que foi entregue

Um processo novo passa a nascer por **três portas**, todas sobre o mesmo núcleo Go (`internal/fabricar`) e gerando os mesmos arquivos: `SKILL.md` (Chat), `.loop.json` (Studio), `.dag.json` (CLI, `legoloop run <caminho>`) e a governança `INTENT`, `SPEC`, `PLAN`, `ROUTE` e `TEMPLATE`.

| Porta | Como se usa |
|:---|:---|
| **CLI** | `legoloop fabricar "<nome>" [--yes] [--goal] [--model] [--category] [--vault] [--docs] [--force] [--dry-run] [--trust-removable] [--link-harness] [--json]`. Sem `--yes`, entrevista numerada, lista de drives e `[s/N]` |
| **Chat** | Skill `legoloop-fabricar`: conduz a entrevista, faz uma prévia com `--dry-run`, **para no stop-gate humano** e só então chama o comando |
| **Studio** | Botão **+ Novo Processo**: nome, categoria, modelo, entrevista, fabricação, plotagem no Canvas e Inspetor com a Tríade; `POST /api/fabricar/run` |

`canonize` (desduplicação) continua separado e **sem nenhum diff** em `commands.go` e `internal/canon`.

## 2. `TICK-04.1` — homologação mecânica final

| Verificação | Resultado |
|:---|:---|
| `go test -count=1 ./...` | Verde em **21 pacotes**. 478 testes de topo + 67 subtestes passando, 0 falhas. **1 pulado, preexistente** (`assets`: symlink exige administrador no Windows; não é meu) |
| `go vet ./...` | Limpo; cross-vet `GOOS=linux` e `darwin` limpos nos pacotes tocados |
| `go run ./tools/funclen internal/fabricar` | Limpo (também `plangen`, `specgen`, `provider` e os meus arquivos de `cli` e `studio`) |
| `git diff --exit-code -- go.mod go.sum` | Inalterados (stdlib pura) |
| `go build` e `go install` | OK; o binário em `~/go/bin` é **idêntico byte a byte** ao de `bin/legoloop.exe` |
| `gofmt -l .` | **Vazio.** Eu formatei 22 arquivos preexistentes: a mudança é só de espaçamento e linhas em branco (`git diff -w --ignore-blank-lines` mostra apenas a minha linha em `static_test.go`) |
| **Smoke hermético** (9Router falso, tudo em `scratch/smoke-fabricar`) | `--dry-run`: exit 0 e **0 arquivos**. Gravação real: exit 0 e **8 arquivos**. `legoloop validate`: `5 nós válidos, sem ciclos`, exit 0. Repetir o nome: exit 1. **Nenhum** `.legoloop-tmp` nem `.legoloop-bak`. Diretório removido ao final |
| Cofre real `~/.aidd` | **Intocado**: nenhum arquivo de teste (`teste-smoke`, `extrator-de-audio`, `primeiro-fluxo`, `segundo-fluxo`) e nenhum `.legoloop-*` |
| `-race` | **Indisponível aqui** (exige CGO); obrigatório no CI Linux |

Testes por pacote tocado: `fabricar` 71, `cli` 66, `studio` 130, `interview` 7, `provider` 17, `specgen` 12, `plangen` 18; 0 falhas e 0 pulados em todos.

## 3. Homologação visual do modal (navegador real)

Feita por mim, no Chrome, em **ambiente isolado**: Studio em porta alternativa, `USERPROFILE` e projeto temporários, 9Router falso e varredura real dos discos. **A homologação humana do `SG-E03` continua sendo sua.**

Fluxo: abrir o modal → nomear → escolher o modelo → entrevistar (98%) → fabricar → plotar.

* O modal abre com o painel de entrevista **oculto** (defeito de CSS achado aqui e corrigido), e o dropdown lista só modelos do 9Router habilitados; os demais aparecem desabilitados com o motivo.
* Ao concluir: modal fechado, **5 nós e 4 arestas** no Canvas, Inspetor com INTENT, SPEC, PLAN, ROUTE, TEMPLATE, Skill, Loop e DAG (caminhos mascarados como `<projeto>` e `~`) e a peça listada em **Loops** sem recarregar a página.
* **Esteira ocupada:** com o Canvas preenchido, plotar abre o diálogo do próprio Studio ("Substituir a esteira do Canvas?") com o foco em **Cancelar**. Cancelar preserva a esteira e o processo continua criado, em Loops.
* **Duplicata:** o erro aparece dentro do modal ("já existe um processo com este nome") e só então surge **Sobrescrever o processo existente**, com confirmação. Confirmado, foram criados 8 `.legoloop-bak` e nenhum `.legoloop-tmp`.
* Console sem erros e sem violação de CSP.
* Capturas: `docs/reports/img/2026-11-03_E03_modal_plotagem.jpg` e `docs/reports/img/2026-11-03_E03_confirmar_substituicao.jpg`.

## 4. Critérios de aceite da SPEC (12 de 12 marcados)

| AC | Evidência principal |
|:---|:---|
| `E01.1` nome | `name_test` com tabela, fuzz de 4×20 s, idempotência |
| `E01.2` drives | `drives_test`: sondagem injetada, unidade que trava, `ParseMounts` Linux e macOS |
| `E01.3` pacote | `engine_test`: 3 meios, `assets.Scanner` real, `loop.Store`, `dag.Parse`, `legoloop validate`, `RouteKeys`, determinismo e raízes |
| `E01.4` atomicidade | `materialize_test` com `FS` que falha na 3ª escrita, `ErrExists`, `Force` com `.legoloop-bak`, `DryRun`; mutação do rollback detectada |
| `E01.5` scanner | `scanner_test`: laço por symlink/junction (rodou de verdade no Windows), limites viram avisos, peça removível `Untrusted`, texto hostil nunca vira `command`; mutação detectada |
| `E02.1` a `E02.3` | `fabricar_test`, `fabricar_interactive_test` (14 casos), `dispatch_fabricar_test`; binário real |
| `E02.4` skill | `skill_fabricar_test`: `assets.Scanner`, frontmatter, ordem gate antes de gravar, sem citar `canonize` |
| `E03.1` modal | `fabricar_modal_test` (HTML, CSP, sem diálogo nativo/emoji/inline) e navegador real |
| `E03.2` plotagem | `fabricar_catalog_test`, `fabricar_plot_test` e navegador real |
| `E03.3` API | `fabricar_api_test`: 405, 403, 421, 415, 400, 409, `../`, `Result` em JSON, resposta sem caminhos absolutos |

## 5. Defeitos reais achados ao rodar o sistema de verdade (todos corrigidos com teste)

1. **Validador da SPEC recusava a SPEC correta.** A regex exigia `AC-` colado ao colchete, mas o template v2 coloca o ID entre crases. Com um LLM real a fabricação **sempre** falhava; meus testes usavam o formato sem crases e nunca viram. Corrigido em `specgen` **e** `plangen` (o aviso de critério sem ticket nunca disparava para SPECs reais).
2. **Chave do 9Router nunca era enviada** (`401`), na CLI e nas rotas de entrevista, SPEC e PLAN do Studio. Agora só vai à URL confiável (`provider.TrustedRouterKey`), com teste de que **não vaza**.
3. **`interview.RouterLLMClient` não enviava `"stream": false`**; este router responde em SSE por padrão. Defeito preexistente que afetava também o Studio.
4. Prazo fixo de 45 s curto demais para gerar a SPEC; `--help` mostrava `-model legoloop doctor` (crases no `flag`); painel `hidden` aparecendo (CSS `display` vence `hidden`).
5. **Dois testes meus estavam errados** e só as mutações revelaram: a checagem de caminhos no JSON (variante de barra escapada era um no-op) e a minha leitura de uma mutação por um `grep` sem `-v`. Corrigidos.

## 6. Validação com o seu 9Router real

* Feita **uma execução `--dry-run`** com objetivo real e o modelo padrão (`code-fast`): **exit 0 em ~35 s**, 5 camadas e 8 caminhos, nada gravado. Antes da correção da regex (item 1 da seção 5), a mesma execução falhava.
* **Ressalva:** não sei se a repetição limitada chegou a ser acionada nessa execução. Ela está provada com LLM falso (corrige na 2ª tentativa, desiste na 3ª sem gravar, não repete falha de transporte). Ao longo do BUILD foram feitas **pelo menos 10 chamadas reais** ao seu router (o número exato é incerto, pois cada falha de contrato dispara até 3 tentativas), com custo.

## 7. Decisões tomadas e limites conhecidos

* **Decidido por você:** a plotagem pergunta antes de substituir uma esteira (implementado). **Decidido por mim, a seu pedido:** `force` no modal só aparece **depois** de um 409 de duplicata, com confirmação explícita; não há botão destrutivo à vista.
* **`--model` é o cérebro:** conduz a entrevista (CLI) e a Tríade; modelos que não são do 9Router (Ollama, `claude` CLI) são recusados com explicação.
* **Entrevista no Studio ainda usa `code-pro` fixo**, e não o modelo do dropdown (`/api/interview/start` só recebe `idea`); o modelo escolhido conduz a Tríade.
* **A plotagem avalia "esteira não salva" por Canvas não vazio**, pois o Studio não rastreia alterações. Pode perguntar a mais, nunca a menos.
* **Entrevista interativa na CLI** provada só com LLM falso (o router falso não simula a entrevista).
* **Varredura real** leva ~7 s a quente e até ~60 s a frio; `C:\` atinge o limite de 500 peças (avisos no stderr).
* **Mensagem de `ErrExists`** lista só nomes de arquivo, sem caminhos.
* **O Studio em execução (PID 14484)** usa o binário **antigo** (o `go install` trocou o arquivo, não o processo): reinicie para ver o modal.
* **`faz-session`** rodou duas vezes (ID do PLAN `1da38997-…` e ID desta sessão `872376ce-…`), ambas com exit 0, e exportou para `Desktopsessoes-agentesLegoLoop`. Os dois arquivos têm o **mesmo conteúdo** (diferem só nas 14 linhas que citam o ID): o script parece ignorar `--id` e exportar o transcript ativo. Não validei isso no código do script; o histórico exportado cobre esta sessão, mas o rótulo do ID não é confiável.
* Em um momento usei `taskkill /IM legoloop.exe`, que encerra **todo** processo com esse nome. Depois passei a encerrar só por porta. Se você tinha um Studio aberto naquela hora, ele pode ter sido encerrado.

## 8. Checklist do `SG-FINAL`

| Item | Estado |
|:---|:---|
| `TICK-04.1` verde (regressão, vet, funclen, `go.mod`, build, instalação, smoke) | **Sim** |
| `gofmt -l .` vazio | **Sim** (22 arquivos preexistentes formatados) |
| `scratch/smoke-fabricar` removido; sem `.legoloop-tmp`/`.legoloop-bak` fora de um `Force` deliberado; cofre real intocado | **Sim** |
| Este relatório, com a homologação visual | **Sim** |
| `AC-E01.1` a `AC-E03.3` marcados | **Sim** (12 de 12) |
| Aprovação humana da demanda | **Aprovada pelo usuário** em 2026-11-03 (resposta "sim"), junto com o commit local e a exportação desta sessão |

Nada foi commitado nem enviado; todas as mudanças estão no diretório de trabalho.
