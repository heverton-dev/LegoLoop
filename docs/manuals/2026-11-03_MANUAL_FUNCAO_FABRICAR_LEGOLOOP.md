# 📘 MANUAL OPERACIONAL: FUNÇÃO FABRICAR (FÁBRICA DE PROCESSOS) NO LEGOLOOP

> **Data:** 2026-11-03  
> **Autoridade:** Ecossistema LegoLoop & AIDD (demanda `LEGOLOOP_FABRICAR`, `SG-FINAL` aprovado)  
> **Status:** Canônico  
> **Escopo:** Criar processos novos (skill + loop + pipeline + governança) a partir de uma ideia, pela CLI, pelo Chat ou pelo Studio  
> **Fontes de verdade:** `docs/specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md` · `internal/fabricar/` · relatório `docs/reports/2026-11-03_REPORT_LEGOLOOP_FABRICAR.md`

---

## 1. Visão Geral e Taxonomia

O **`legoloop fabricar`** transforma uma ideia ("quero um extrator de áudio") em um **processo** pronto para usar: descobre as melhores peças que você já tem nos seus discos, monta uma esteira com elas, escreve a governança (INTENT, SPEC, PLAN, ROUTE) e grava tudo de forma atômica.

### 1.1 O que ele é: uma **função**

Os termos abaixo já existem no LegoLoop; o `fabricar` não é nenhum deles, ele **produz** os três primeiros.

| Termo | O que é no LegoLoop | Onde vive | O `fabricar`... |
|:---|:---|:---|:---|
| **Peça** | Unidade encaixável no Canvas. Tem 8 naturezas: Trigger, Controle, **Skill**, MCP, **Loop**, Script, Hook, Gate | `SKILL.md`, scripts, MCPs, hooks, gates | **consome** peças existentes (skills dos seus discos) e **cria** uma peça Skill nova |
| **Loop** | Esteira/sub-grafo reutilizável, guardada no acervo e usada como sub-loop em outras esteiras | `~/.aidd/canonical/loops/<nome>.loop.json` | **cria** um |
| **Pipeline** | O DAG executável: nós, dependências, artefatos e gates. Roda com `legoloop run <caminho>` | `<nome>.dag.json` | **cria** um |
| **Função** | Capacidade do próprio binário LegoLoop, com uma ou mais "portas" (CLI, Chat, Studio). Não vira nó do Canvas e não é executada pelo Orchestrix | `legoloop fabricar`, `legoloop canonize`, `legoloop route`... | **é** uma |
| **Processo** | O produto do `fabricar`: o conjunto `legoloop-<slug>` = **1 Skill + 1 Loop + 1 Pipeline + 5 documentos de governança** | vários arquivos (seção 5) | é o que ele entrega |

> **Por que o manual se chama "FUNÇÃO FABRICAR" e não "…_LOOP" ou "…_PIPELINE":** os outros manuais da pasta (`MEGA_PIPELINE`, `ADDY_OSMANI_LOOP`...) documentam *uma esteira que se executa*. O `fabricar` não é executado como esteira: é Go dentro do binário. Só o **resultado** dele é loop, pipeline e peça.
>
> **Convenção de nomes de manuais** (a partir daqui): `MANUAL_FUNCAO_<NOME>` documenta uma função do binário; `MANUAL_<NOME>_LOOP` documenta um loop; `MANUAL_<NOME>_PIPELINE` documenta um pipeline; uma peça isolada se documenta no próprio `SKILL.md`.

### 1.2 Não confundir com `canonize`

| | `legoloop fabricar` | `legoloop canonize` |
|:---|:---|:---|
| Faz | **Cria** um processo novo | **Desduplica** arquivos iguais (SHA-256), trocando cópias por links para o cofre |
| Entrada | Uma ideia e um nome | Pastas (`--root`) |
| Saída | Skill, loop, DAG e documentos | Links e economia de disco |
| Pacote | `internal/fabricar` | `internal/canon` |

São independentes: usar um não altera o outro.

---

## 2. Quando Usar (e Quando Não)

**Use** quando você quer um processo novo e reaproveitável (skill + loop + DAG) montado com as suas próprias peças.

**Não use** para: rodar um processo que já existe (`legoloop run`), editar uma skill existente, desduplicar arquivos (`canonize`), ou executar a ideia em si (o `fabricar` cria o processo; ele não o executa).

---

## 3. As Três Portas

### 3.1 CLI

```bash
# Autônomo: sem perguntas. O objetivo é --goal ou, sem ele, o próprio nome.
legoloop fabricar "extrator de audio" --yes --goal "Extrair o áudio de vídeos locais em mp3, em lote"

# Prévia: calcula tudo e NÃO grava nada. Faça sempre antes da primeira vez.
legoloop fabricar "extrator de audio" --yes --dry-run

# Interativo (sem --yes): lista os drives, faz a entrevista e pede confirmação.
legoloop fabricar "extrator de audio"
```

No modo interativo, o terminal conduz três momentos:
1. **Drives:** lista os discos detectados; digite os números a desmarcar (separados por vírgula) ou Enter para manter todos. Desmarcar todos é inválido.
2. **Entrevista:** cada pergunta traz `[1] Aceitar o palpite`, `[2] Responder com minhas palavras`, `[0] Cancelar`. Vai até 95% de confiança ou 12 turnos.
3. **Resumo e `[s/N]`:** mostra os caminhos que serão gravados. Enter vale **N**. Tudo antes disso é descartável: nada foi gravado.

Entrada inválida reaparece até 3 vezes e depois aborta; EOF no stdin também aborta (exit 1, nada gravado).

### 3.2 Chat (skill `legoloop-fabricar`)

Digite `/legoloop-fabricar` ou descreva o processo que quer. A skill (em `skills/legoloop-fabricar/SKILL.md`):
1. conduz a entrevista, uma pergunta por vez (objetivo, nome, categoria e modelo);
2. roda a **prévia** `legoloop fabricar "<nome>" --yes --goal "..." --dry-run`;
3. **para no stop-gate humano**: mostra o resumo e só continua com um "sim" explícito;
4. roda o mesmo comando sem `--dry-run` e informa os caminhos.

A skill não duplica a lógica: ela só conduz a conversa. Quem varre, funde, gera e grava é o comando.

### 3.3 Studio

Em `legoloop studio` (padrão `http://127.0.0.1:7878`), botão **+ Novo Processo** na barra superior:
1. **Nome**, **Categoria** e **Modelo que conduz a fabricação** (só modelos do 9Router ficam habilitados; os outros aparecem desabilitados com o motivo).
2. **Iniciar Entrevista** (usa as rotas `/api/interview/*`) ou **Fabricar sem entrevista** (usa o nome como objetivo).
3. Ao concluir a entrevista, **Fabricar processo**. A fabricação leva de alguns segundos a alguns minutos (varredura dos discos + LLM).
4. Ao terminar: o modal fecha, o processo aparece em **Loops**, o Canvas planta os nós e o **Inspetor** abre com a Tríade e os caminhos gerados.

Duas confirmações protegem você:
* **Canvas ocupado:** se já houver nós na tela, o Studio pergunta "Substituir a esteira do Canvas?" (o foco inicial fica em **Cancelar**). Cancelando, o processo continua criado e fica em **Loops**; `Ctrl+Z` desfaz uma substituição.
* **Nome repetido:** o erro aparece no próprio modal ("já existe um processo com este nome") e só então surge **Sobrescrever o processo existente**, que pede confirmação e guarda o original como `.legoloop-bak`.

> Se você já tinha o Studio aberto antes de atualizar o binário, **reinicie-o** para ver o modal.

---

## 4. As Cinco Camadas

| # | Camada | O que faz | Detalhes que importam |
|:--|:---|:---|:---|
| 1 | **Clarificação** | Modela a intenção | Interativo: entrevista socrática até ≥ 95% ou 12 turnos. Autônomo: não pergunta nada, usa `--goal` (ou o nome) e a INTENT nasce `DRAFT` com confiança **0%**, sem inventar número |
| 2 | **Varredura** | Acha peças nos discos | Agnóstica a SO (letras no Windows, `/proc/mounts` no Linux, `/Volumes` no macOS). Por raiz: profundidade 4, até 20.000 itens, até 500 peças e prazo de 60 s. Estourar um limite vira **aviso**, não erro |
| 3 | **Fusão** | Escolhe e ordena as peças | Só skills entram como passos. Até **5 passos**. Nomes repetidos: fica a melhor. Peças `Untrusted` (mídia removível/rede) e de tier tóxico ficam de fora. Sem nenhuma utilizável: passos padrão `aidd-spec` e `aidd-tickets` |
| 4 | **Tríade** | Escreve a governança | INTENT (renderizada), SPEC (**LLM**, validada contra o contrato) e PLAN + ROUTE (**sem LLM**, pelo `plangen.Render`). Se a SPEC violar o contrato, o motivo é devolvido ao modelo: até **3 tentativas**, só para violação de contrato (nunca para falha de rede) |
| 5 | **Template e gravação** | Gera o gabarito e grava | `TEMPLATE_<NOME>.md` e gravação **atômica** de tudo (ou nada) |

O progresso aparece como `[1/5] Clarificação...` até `[5/5] Template e gravação...`.

---

## 5. O Que É Gerado (Anatomia de um Processo)

Para `legoloop fabricar "extrator de audio"` o nome canônico é `legoloop-extrator-de-audio`:

| Meio / documento | Caminho (padrão) | Para quê |
|:---|:---|:---|
| **Skill** (peça, Chat) | `~/.aidd/canonical/skills/legoloop-extrator-de-audio/SKILL.md` | aparece no inventário do Studio; vira `/legoloop-extrator-de-audio` no Chat depois de ligada a um harness com `--link-harness` |
| **Loop** (acervo, Studio) | `~/.aidd/canonical/loops/legoloop-extrator-de-audio.loop.json` | carrega no Canvas; usável como sub-loop |
| **Pipeline** (CLI) | `~/.aidd/canonical/loops/legoloop-extrator-de-audio.dag.json` | `legoloop run <caminho>` |
| INTENT | `<docs>/docs/intents/AAAA-MM-DD_INTENT_EXTRATOR_DE_AUDIO.md` | o que se quer e por quê |
| SPEC | `<docs>/docs/specs/AAAA-MM-DD_SPEC_EXTRATOR_DE_AUDIO.md` | contratos e critérios binários de aceite |
| PLAN | `<docs>/docs/plans/AAAA-MM-DD_PLAN_EXTRATOR_DE_AUDIO.md` | tickets e stop-gates |
| ROUTE | `<docs>/docs/routes/AAAA-MM-DD_ROUTE_EXTRATOR_DE_AUDIO.json` | espelho 1:1 do PLAN, executável |
| TEMPLATE | `<docs>/docs/templates/TEMPLATE_EXTRATOR_DE_AUDIO.md` | gabarito das saídas de cada passo |

A data vem do relógio injetado (determinismo); `<docs>` é a raiz do projeto (`--docs`, padrão: pasta atual) e o `docs/` é acrescentado por baixo.

### 5.1 O pipeline (`.dag.json`)

É uma **cadeia sequencial**: cada passo depende do anterior e consome o artefato dele.
* Nós `n01_<peça>`, `n02_<peça>`... com `type: skill` e `ref: skills/<peça>`.
* Artefato por passo `out/NN_<peça>.md` e gate `file_exists('out/NN_<peça>.md')`.
* Roteamento **Sugerido**: o modelo é o que você escolheu em `--model` e o fallback é `code-fast`.
* É um DAG v2 válido: `legoloop validate <caminho>` deve sair com 0 (`N nós válidos, sem ciclos`).

### 5.2 A skill (`SKILL.md`)

Frontmatter lido pelo Studio e pelo grafo: `name` (igual ao diretório), `description` (até 300 caracteres, começando pelo gatilho de uso), `consumes`, `produces`, `artifact`, `gate` e `suggested_model`. O corpo traz o contrato de I/O, o procedimento (validar o grafo e aplicar cada skill em ordem), as regras negativas e o critério de parada.

### 5.3 Regra de nomes

O nome é normalizado de forma determinística e idempotente: `legoloop-` + kebab-case, sem acento, até 64 caracteres.

| Você digita | Vira |
|:---|:---|
| `Processo do Joãozinho` | `legoloop-processo-do-joaozinho` |
| `legoloop_teste` | `legoloop-teste` |
| `  Ação   Rápida!! ` | `legoloop-acao-rapida` |
| `../../etc/passwd` | `legoloop-etc-passwd` (nunca sai das raízes) |
| `""`, `---`, `legoloop`, `日本語` | **erro** de nome inválido (exit 2) |

---

## 6. Referência da CLI

```
legoloop fabricar "<nome>" [opções]
```

| Flag | Padrão | Efeito |
|:---|:---|:---|
| `--yes` | desligado | modo autônomo (não pergunta nada) |
| `--goal TEXTO` | o próprio nome | objetivo em linguagem natural |
| `--model M` | primeiro do inventário | modelo que conduz a fabricação; precisa ser do **9Router** e estar no inventário |
| `--category` | `build` | `build`, `reverse`, `web` ou `mobile` |
| `--vault DIR` | `~/.aidd/canonical` | cofre que recebe skill, loop e DAG |
| `--docs DIR` | pasta atual | **raiz do projeto**; os documentos vão em `DIR/docs/...` |
| `--force` | desligado | sobrescreve um processo existente, guardando o original como `.legoloop-bak` |
| `--dry-run` | desligado | calcula tudo, imprime os caminhos e **não grava nada** |
| `--trust-removable` | desligado | aceita peças de mídia removível/rede na fusão |
| `--link-harness H` | nenhum | liga a skill à pasta de skills do harness (ex.: `claude`); repetível; **opt-in** |
| `--json` | desligado | imprime o resultado como JSON (o progresso vai para o stderr) |

**Códigos de saída:** `0` sucesso · `1` falha de execução (modelo fora do inventário, processo já existente, LLM indisponível, SPEC fora do contrato, cancelado) · `2` uso inválido (nome ou categoria inválidos, flag desconhecida, sem nome).

---

## 7. API do Studio

`POST /api/fabricar/run` (passa pelo middleware de segurança do Studio: só loopback, Host permitido, mesma origem, `application/json`, corpo ≤ 1 MiB, limite de taxa).

Corpo aceito (qualquer outro campo é **400**): `{"name", "goal", "category", "model", "force"}`. O cliente **nunca** envia caminhos, drives, confiança em mídia removível nem link em harness: isso é decisão do servidor.

| Status | Significado |
|:---|:---|
| `200` | sucesso; corpo é o resultado em JSON, com caminhos mascarados (`~`, `<projeto>`) |
| `400` | nome/categoria/modelo inválidos, campo desconhecido, ou modelo fora do inventário |
| `403` / `421` / `415` | origem cruzada / Host inválido / content-type diferente de JSON |
| `405` | método diferente de POST |
| `409` | `code: "exists"` (processo já existe) ou `code: "busy"` (já há uma fabricação em andamento) |
| `422` | o modelo escolhido não pode conduzir a fabricação (não é do 9Router) |
| `502` | o modelo respondeu fora do contrato mesmo após as tentativas |
| `500` | falha genérica; nenhum arquivo parcial; o erro bruto não sobe ao navegador |

Só uma fabricação roda por vez. Ela termina mesmo que a aba seja fechada (teto de 10 minutos).

---

## 8. Segurança e Garantias

| Garantia | Como |
|:---|:---|
| **Tudo ou nada** | gravação por arquivos temporários + `rename`; falha no meio desfaz tudo (sem `.legoloop-tmp` restante) |
| **Sem sobrescrita silenciosa** | processo existente ⇒ erro; `--force` guarda `.legoloop-bak` |
| **Nada fora das raízes** | nome sanitizado; todos os arquivos ficam sob o cofre e a pasta de docs |
| **Peça alheia não vira comando** | nomes das peças varridas só viram `ref`/ID; texto hostil numa `SKILL.md` nunca vira `command`, `slot` ou `hook`. MCP, hook, script e loop varridos não viram passo por padrão |
| **Mídia removível é suspeita** | peças de mídia removível/rede entram como `Untrusted` e ficam fora da fusão até `--trust-removable` |
| **Chave do 9Router protegida** | só vai para a URL confiável configurada (`NINEROUTER_URL`); nunca aparece em log, saída ou resposta de API |
| **Cofre real protegido em teste** | os testes e o smoke usam LLM falso e diretórios temporários |
| **Harness só com pedido** | o cofre não é lido pelos harnesses; `--link-harness` liga a skill, e só quando você pede |

Limite conhecido: um HD externo USB que se anuncia como disco fixo é tratado como confiável pelo scanner.

---

## 9. Usando o Processo Fabricado

```bash
legoloop validate  ~/.aidd/canonical/loops/legoloop-extrator-de-audio.dag.json
legoloop run       ~/.aidd/canonical/loops/legoloop-extrator-de-audio.dag.json --var idea="..."
```

* **No Chat:** a skill `legoloop-extrator-de-audio` só chega a um harness (ex.: `~/.claude/skills`) com `--link-harness`, porque o cofre não é lido por eles.
* **No Studio:** o loop aparece em **Loops**; arraste-o como **sub-loop** para outras esteiras ou abra o processo no Canvas.
* **Nós `script`** só rodam com `--allow-scripts` (regra geral do `run`).
* **Governança:** os documentos gerados seguem os templates canônicos v2 (`docs/templates/`) e podem ser revisados e aprovados como qualquer demanda (INTENT → SPEC → PLAN → ROUTE).

---

## 10. Solução de Problemas

| Mensagem | Causa e saída |
|:---|:---|
| `nome inválido` (exit 2) | o nome não gera nenhuma letra ou número; escolha outro |
| `modelo "X" fora do inventário` | rode `legoloop doctor` e use um dos modelos listados |
| `o modelo "X" não é servido pelo 9Router` | o cérebro precisa ser do 9Router; escolha um com executor `9router` |
| `a Tríade precisa de um LLM: suba o 9Router` | o 9Router está offline ou sem `NINEROUTER_URL`; veja `legoloop doctor` |
| `processo já existe` | use outro nome, ou `--force` (CLI) / **Sobrescrever** (Studio) |
| `especificação gerada violou o contrato` | o modelo errou o formato mesmo após 3 tentativas; tente outro modelo com `--model` |
| `já existe uma fabricação em andamento` | espere a anterior terminar |
| Demora de 30 s a vários minutos | normal: varredura real dos discos + geração da SPEC; `--dry-run` mostra o mesmo custo sem gravar |
| Avisos `limite de 500 peças atingido`, `além da profundidade 4` | o scanner parou num limite e seguiu; o resultado vale, com menos peças |
| O modal **+ Novo Processo** não aparece | o Studio aberto usa o binário antigo; reinicie-o |

---

## 11. Limites Conhecidos

* A **entrevista** do Studio usa o modelo `code-pro` fixo; o modelo escolhido no dropdown conduz a Tríade.
* Modelos que **não são do 9Router** (Ollama, `claude` CLI) não podem conduzir a fabricação.
* O modo **autônomo** não entrevista, então a INTENT gerada nasce `DRAFT` (confiança 0%) até você revisá-la.
* O Studio decide "esteira não salva" por **Canvas não vazio**: pode perguntar a mais, nunca a menos.
* O fluxo completo com **LLM real** foi validado em `--dry-run`; a repetição da SPEC está provada com LLM falso.
* `go test -race` só roda onde há CGO (CI Linux).

---

## 12. Para Quem Mantém o Código

| Peça | Onde |
|:---|:---|
| Orquestrador e camadas | `internal/fabricar/engine.go` (`Fabricar`, `Deps`, `Result`) |
| Nome, drives, scanner, fusão, gravação | `name.go`, `drives.go`, `scanner.go`, `merger.go`, `materialize.go` |
| Tríade e repetição da SPEC | `triad.go`, `repair.go` (`specgen`, `plangen.Render`, `interview.RenderIntent`) |
| Emissores | `emit_skill.go`, `emit_loop.go`, `emit_dag.go`, `template.go` |
| Cérebro e chave | `brain.go` (`BrainFor`), `internal/provider/trusted.go` (`TrustedRouterKey`) |
| CLI | `internal/cli/fabricar.go`, `fabricar_interactive.go` |
| Skill do Chat | `skills/legoloop-fabricar/SKILL.md` |
| Studio | `internal/studio/fabricar_api.go`, `static/fabricar.js`, `static/app.js`, `static/inspector.js` |

Verificação: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal/fabricar`, `gofmt -l .` (vazio). A SPEC tem 12 critérios binários (`AC-E01.1` a `AC-E03.3`), todos atendidos e rastreados no relatório final.
