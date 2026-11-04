# 🧩 LegoLoop

> **The Composable Looping & Harness Fabric**  
> *Transforme suas skills, MCPs e modelos em esteiras visuais no estilo N8N com Orchestrix cirúrgicas de altíssima economia de tokens.*

---

## ⚡ Por que LegoLoop?

Os agentes de IA tradicionais sofrem de dois problemas fatais:
1. **Inchaço de Contexto:** Conversas monolíticas chegam a 150k tokens, acumulam ruído e estouram as cotas de modelos.
2. **Engessamento:** Os fluxos de trabalho são rígidos e não se adaptam às ferramentas que o desenvolvedor realmente possui na máquina.

O **LegoLoop** resolve isso tratando tudo como **peças Lego**:
* **Canvas Visual (N8N Style):** Ligue nós de Concepção, Engenharia Reversa, Design e TDD com arrastar e soltar.
* **Orchestrix Cirúrgica por Bloco:** Cada nó roda isolado com contexto zero (2k a 5k tokens), lê apenas o artefato do disco e morre ao concluir.
* **Fábrica de Processos (`legoloop fabricar`):** descreve-se o que se quer, e o LegoLoop varre os discos, escolhe as melhores peças e gera o processo completo (skill, loop e DAG) mais a governança INTENT, SPEC, PLAN e ROUTE. Pela CLI, pelo Chat ou pelo Studio.
* **Canonização de SO sem Quebras:** Centraliza seus ativos em `~/.aidd/canonical/` e substitui duplicatas exatas por links (symlink; hardlink quando o Windows não permite symlink sem admin), com ~80% menos espaço em disco e rollback.
* **Convergência Total:** Integra as regras do **Ecossistema AIDD**, o motor de alta performance do **SYMBIONT** e a blindagem de modelos do **9Router**.

---

## 📂 Estrutura do Projeto

```
LegoLoop/
├── cmd/legoloop/            # Binário (main curto; lógica em internal/cli)
├── internal/
│   ├── canon/               # Canonização: hash SHA-256, eleitor, links com rollback (desduplicação)
│   ├── dag/                 # Contrato pipeline.dag.json v2: tipos de nó, triggers, slots, ciclos
│   ├── gate/                # Motor de gates (file_exists, has_pattern, exit_code...)
│   ├── runner/              # Runner por camadas, paralelo, timeout, telemetria
│   ├── provider/            # Inventário passivo: Ollama, 9Router, CLIs, chaves (TrustedRouterKey)
│   ├── orchestrix/          # Orchestrix isolada por bloco (claude -p, Ollama, 9Router, harnesses, script, sub-loop)
│   ├── graph/               # Ontologia, compatibilidade de portas, rota ótima (Dijkstra)
│   ├── assets/              # Inventário do SO: skills, MCPs, hooks, scripts, gates, loops, harnesses
│   ├── suggest/             # Sugestões determinísticas (gates, scripts, hooks, vizinhos), sem LLM
│   ├── loop/                # Acervo de loops (.loop.json), portas derivadas, hierarquia de sub-loops
│   ├── interview/           # Entrevista socrática até 95% de confiança; renderiza a INTENT; cliente do 9Router
│   ├── specgen/             # Gera a SPEC (LLM) conforme TEMPLATE_SPEC; valida o contrato; ErrContract
│   ├── plangen/             # Gera PLAN e ROUTE conforme os templates (Generate com LLM, Render sem LLM)
│   ├── fabricar/            # FÁBRICA DE PROCESSOS: nome, drives, scanner, fusão, Tríade, gravação atômica
│   ├── reputation/          # Matriz de reputação das peças (tiers ouro, prata...)
│   ├── vault/               # Cofre de credenciais (AES-256-GCM, mascaramento)
│   ├── board/ · wave/       # Integração com GitHub Projects e planejador de ondas sem sobreposição
│   ├── studio/              # Looping Studio (canvas + API local + modal Novo Processo) embutido no binário
│   └── cli/                 # Comandos
├── skills/legoloop-fabricar/  # Skill do Chat que conduz a entrevista e chama `legoloop fabricar`
├── docs/
│   ├── templates/           # Templates canônicos v2 (INTENT, SPEC, PLAN, ROUTE e peças) + conformidade em Go
│   ├── intents/ · specs/ · plans/ · routes/   # Cascata de governança (prefixo ISO de data)
│   ├── reports/ · manuals/ · references/ · sessions/ · archive/
├── loops/                   # Loops de exemplo (addy-osmani, impeccable-design, matt-pocock, mega-pipeline)
├── tools/funclen/           # Gate de qualidade: nenhuma função de produção com mais de 35 linhas
└── AGENTS.md · MEMORY.md (diário de implementação) · README.md
```

---

## 🛠️ Início Rápido

```bash
go build -trimpath -ldflags="-s -w" -o bin/legoloop.exe ./cmd/legoloop
go install ./cmd/legoloop                         # deixa o binário atualizado no PATH

legoloop doctor                                   # modelos executáveis, harnesses, 9Router
legoloop canonize --root DIR --dry-run            # plano + economia estimada (não toca o disco)
legoloop canonize --root DIR --apply              # substitui duplicatas exatas por links
legoloop route --from Artifact:URL --to Artifact:Spec
legoloop validate pipeline.dag.json
legoloop run pipeline.dag.json --var idea="..."   # executa os nós com gates
legoloop fabricar "extrator de audio" --yes       # cria um processo novo (skill + loop + dag + docs)
legoloop studio                                   # http://127.0.0.1:7878
legoloop forge                                    # abre o Studio e dispara a primeira varredura
legoloop import <github-url>                      # importa e canoniza uma skill do GitHub
legoloop upgrade <arquivo> --diagnose             # diagnóstico para o Padrão Ouro
legoloop matrix --tier ouro                       # matriz de reputação das peças
```

## 🏭 Fábrica de Processos (`legoloop fabricar`)

Gera o processo `legoloop-<nome>` em **cinco camadas**: Clarificação (entrevista até 95% de confiança ou 12 turnos; no modo autônomo usa `--goal` ou o próprio nome) → Varredura dos discos (agnóstica a SO, com limites e prazos) → Fusão (curadoria das melhores peças, no máximo 5 passos) → Tríade (INTENT, SPEC pelo LLM, PLAN e ROUTE sem LLM) → Template e gravação atômica.

| Porta | Como |
|---|---|
| **CLI** | `legoloop fabricar "<nome>" [--yes] [--goal T] [--model M] [--category build\|reverse\|web\|mobile] [--vault DIR] [--docs DIR] [--force] [--dry-run] [--trust-removable] [--link-harness H] [--json]`. Sem `--yes`: perguntas numeradas, lista de drives e `[s/N]`. Saída 0 ok, 1 falha, 2 uso inválido |
| **Chat** | `/legoloop-fabricar`: a skill conduz a entrevista, faz a prévia com `--dry-run`, **para num stop-gate humano** e só então chama o comando |
| **Studio** | Botão **+ Novo Processo**: nome, categoria, modelo, entrevista, fabricação, plotagem no Canvas e Inspetor com a Tríade. `POST /api/fabricar/run` |

**O que é gerado:** `SKILL.md` (Chat) e `.loop.json` + `.dag.json` (Studio e CLI, `legoloop run <caminho>`) no cofre `~/.aidd/canonical/`, mais `INTENT`, `SPEC`, `PLAN`, `ROUTE` e `TEMPLATE` em `docs/`. Tudo conforme os templates v2 e validado por `legoloop validate`.

* `--model` é o **cérebro** que conduz a fabricação e precisa ser servido pelo 9Router (modelos de outro executor são recusados com explicação).
* A SPEC gerada por LLM é validada contra o contrato; se violar, a Tríade reenvia ao modelo o motivo (até 3 tentativas, só em violação de contrato).
* Processo já existente ⇒ erro; `--force` guarda o original como `.legoloop-bak`. No Studio, "Sobrescrever" só aparece depois do erro de duplicata, com confirmação.
* `fabricar` **cria** processos; `canonize` **desduplica** arquivos. São funções diferentes e independentes.

Manual completo: [`docs/manuals/2026-11-03_MANUAL_FUNCAO_FABRICAR_LEGOLOOP.md`](docs/manuals/2026-11-03_MANUAL_FUNCAO_FABRICAR_LEGOLOOP.md).

## 🧠 Executores (roteamento por nó)

| `routing.model` | Executor | Observação medida (31/10/2026) |
|---|---|---|
| combo do 9Router (`code-pro`, `code-fast`…) | 9Router `/v1/chat/completions` | requer `NINEROUTER_URL` (https) e `NINEROUTER_KEY` (ou `NINEROUTER_API_KEY`); o cliente envia `"stream": false` porque o 9Router transmite em SSE por padrão |
| tag do Ollama (`llama3.1:8b`) | Ollama `/api/generate` | ~260–750 tokens de entrada por nó |
| `haiku` / `sonnet` / `opus` | `claude -p` isolado | usa a autenticação do **CLI** (`claude auth status`) |
| `agy:<modelo>` / `agy:default` | Antigravity `agy -p` | ~18,7k tokens fixos de sistema por chamada |
| `mimo:<provedor/modelo>` / `mimo:default` | MiMo Code `mimo run` | ~98k tokens fixos de sistema por chamada |
| `opencode:<provedor/modelo>` | `opencode run` | mesmo formato do mimo |
| `omp:<modelo>` / `omp:default` | `omp -p` | `omp:default` usa o 9Router ⇒ precisa de `NINEROUTER_KEY` |

`:default` usa o modelo padrão configurado no próprio harness. No modo **Sugerido**, se o modelo principal falhar em execução, o `fallback` assume e o relatório mostra a rota efetiva (coluna MODELO) e o motivo do descarte.

```bash
# 9Router na VPS (a chave fica só no ambiente; nunca no pipeline):
setx NINEROUTER_URL https://9router.vpsconexao.org
setx NINEROUTER_KEY <sua-chave>
```

## 🔒 Segurança (padrões seguros)

| Comportamento | Por quê |
|---|---|
| Nós `script` só rodam com `--allow-scripts` | um `pipeline.dag.json` de repositório não confiável não executa comandos sem consentimento |
| `claude -p` roda num diretório temporário vazio, com `--safe-mode --setting-sources user --tools ""` | o repositório não injeta `.claude/settings.json`, hooks ou `CLAUDE.md` |
| Artefatos nunca em `.git`, `.claude`, `.vscode`, `.github`… nem através de diretório linkado ou nome 8.3 | saída do LLM não vira hook/config executável |
| `canonize` ignora repositórios git (use `--include-git` para incluir) | link commitado quebraria o repositório para colaboradores |
| `NINEROUTER_URL` só https ou http-loopback; a chave vai apenas para a base confiável (`provider.TrustedRouterKey`) | credencial não trafega em claro nem para portas genéricas, e nenhuma resposta do Studio a expõe |
| Studio só escuta em loopback, valida Host/Origin, exige JSON, limita corpo e taxa | anti DNS-rebinding e anti-CSRF |
| `fabricar` grava de forma **atômica** (temporário + rename, rollback em falha) e nunca fora das raízes injetadas; nomes são sanitizados (`../` vira slug) | falha no meio não deixa processo pela metade |
| Peças de mídia removível/rede entram como `Untrusted` e ficam **fora da fusão** (`--trust-removable` para aceitar); MCP, hook, script e loop varridos não viram passo por padrão | texto hostil numa `SKILL.md` alheia nunca vira `command`, `slot` nem `hook` do DAG gerado |
| `POST /api/fabricar/run`: o cliente envia só nome, objetivo, categoria, modelo e `force`; caminhos, drives, link em harness e confiança em mídia são do servidor; respostas com `~` e `<projeto>`, nunca caminho absoluto | o navegador nunca decide onde gravar |
| `--link-harness` é opt-in (o cofre não é lido pelos harnesses) | nada é ligado a `~/.claude` sem pedido explícito |

**Riscos aceitos:** com fallback por hardlink (Windows sem Developer Mode), um editor que grave *no próprio arquivo* altera todas as cópias linkadas — é a natureza do hardlink; o runner do LegoLoop sempre grava por temporário + rename. No Windows, processos netos de um nó podem ficar órfãos após timeout (sem Job Object); o runner não trava. Um HD externo USB que se anuncia como disco fixo é tratado como confiável pelo scanner do `fabricar`.

## ✅ Qualidade e verificação

```bash
go test -count=1 ./...            # 21 pacotes; lógica JS pura do Studio roda em node
go vet ./...
go run ./tools/funclen internal   # nenhuma função > 35 linhas
gofmt -l .                        # deve sair vazio (o exit code é 0 mesmo listando: confira a saída)
git diff --exit-code -- go.mod go.sum   # só biblioteca padrão
```

`go test -race` exige CGO (indisponível no Windows sem `gcc`): roda no CI Linux. Estado em 03/11/2026: 478 testes de topo + 67 subtestes, 0 falhas, 1 pulado preexistente (symlink sem admin no Windows).

## 📚 Governança

Cada demanda segue a cascata **INTENT → SPEC → PLAN → ROUTE → BUILD (TDD) → relatório**, com stop-gates humanos e artefatos em `docs/` com prefixo `YYYY-MM-DD_`. Os formatos estão em `docs/templates/` (templates v2, testados por `go test ./docs/templates`). A demanda mais recente é a `fabricar` (`docs/*/2026-11-03_*_LEGOLOOP_FABRICAR.*`, relatório final em `docs/reports/2026-11-03_REPORT_LEGOLOOP_FABRICAR.md`). Manuais operacionais em `docs/manuals/`. Regras para agentes em [`AGENTS.md`](AGENTS.md); diário de evolução em [`MEMORY.md`](MEMORY.md).
