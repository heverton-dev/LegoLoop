# 🧩 LegoLoop

> **The Composable Looping & Industrial Harness Fabric**  
> *Fundição Autônoma de Software e esteiras determinísticas com Orchestrix cirúrgicas de altíssima economia de tokens, Protocolos Industriais e 6 Motores de Fronteira N10.*

---

## ⚡ Por que LegoLoop?

Os agentes de IA tradicionais sofrem de dois problemas fatais:
1. **Inchaço de Contexto:** Conversas monolíticas chegam a 150k tokens, acumulam ruído e estouram as cotas de modelos.
2. **Engessamento:** Os fluxos de trabalho são rígidos e não se adaptam às ferramentas que o desenvolvedor realmente possui na máquina.

O **LegoLoop** resolve isso tratando tudo como **peças Lego** sob governança industrial:
* **Canvas Visual (Looping Studio):** Ligue nós de Concepção, Engenharia Reversa, Design e TDD com arrastar e soltar em `http://127.0.0.1:7878`.
* **Orchestrix Cirúrgica por Bloco:** Cada nó roda isolado com contexto estéril (2k a 5k tokens), lê apenas o artefato contratado do disco e morre ao concluir.
* **As 3 Superfícies de Criação de Pipelines:**
  1. **Terminal CLI:** Comandos industriais diretos (`legoloop fabricar`, `legoloop matriz`, `legoloop listar`, etc.).
  2. **Agente em Linguagem Natural:** 14 skills canônicas universais nos 5 harnesses (`/legoloop-fabricar`, `/legoloop-fechar`, etc.).
  3. **Studio LegoLoop (Dupla Superfície):** Modo Cockpit com telemetria profunda para engenharia e Modo Autopilot com cartões executivos simples para operadores leigos.
* **Ciclo de Vida Canônico:** Esteiras completas governadas pelos verbos **Fabricar** (`fabricar`), **Matriz de Peças** (`matriz`), **Listar** (`listar`), **Refinar** (`refinar`), **Remover** (`remover`), **Compor** (`compor`) e **Fechar** (`fechar`).
* **Protocolos Industriais & 5 Prisões Mecânicas:** STANAG 2014, ASD-STE100, EARS, FMEA, Readback ICAO, Andon Cord/Jidoka e SBAR com isolamento de contexto e validação rígida de schemas.
* **Cérebros Plugáveis & Resilientes:** Suporte multi-provedor (9Router VPS, Gemini via harness `agy`, Ollama local) para conduzir entrevistas socráticas e geração da Tríade.
* **Canonização de SO sem Quebras:** Centraliza seus ativos em `~/.aidd/canonical/` e substitui duplicatas exatas por links com ~80% menos espaço em disco e rollback.

---

## 📂 Estrutura do Projeto

```
LegoLoop/
├── cmd/legoloop/            # Binário (main curto; lógica em internal/cli)
├── internal/
│   ├── canon/               # Canonização: hash SHA-256, eleitor, links com rollback (desduplicação)
│   ├── dag/                 # Contrato pipeline.dag.json v2: tipos de nó, triggers, slots, ciclos
│   ├── gate/                # Motor de gates (file_exists, has_pattern, exit_code, gates da Tríade)
│   ├── runner/              # Runner por camadas, paralelo, timeout, telemetria
│   ├── provider/            # Inventário passivo: Ollama, 9Router, CLIs, chaves (TrustedRouterKey)
│   ├── brain/               # Cérebros plugáveis (9Router, Antigravity/Gemini, Ollama) com auto-seleção
│   ├── ste/                 # Catálogo formal de mensagens e avisos de risco no padrão STE-PT
│   ├── orchestrix/          # Orchestrix isolada por bloco (claude -p, Ollama, 9Router, harnesses, script, sub-loop)
│   ├── graph/               # Ontologia, compatibilidade de portas, rota ótima (Dijkstra)
│   ├── assets/              # Inventário do SO: skills, MCPs, hooks, scripts, gates, loops, harnesses
│   ├── suggest/             # Sugestões determinísticas (gates, scripts, hooks, vizinhos), sem LLM
│   ├── loop/                # Acervo de loops (.loop.json), portas derivadas, hierarquia de sub-loops
│   ├── processo/            # Store canônica de processos (~/.aidd) e gestão de vínculos nos 5 harnesses
│   ├── matriz/              # Matriz de peças: vocabulário fechado, reputação/ranking e defasagem
│   ├── listar/              # Inventário consolidado de processos, lacunas e sugestões de evolução
│   ├── refinar/             # Refatoração e substituição atômica de passos com preservação de versão
│   ├── remover/             # Remoção atômica de processos com desvinculação limpa e rollback
│   ├── compor/              # Composição hierárquica e encaixe de sub-loops reutilizáveis
│   ├── interview/           # Entrevista socrática até 95% de confiança; renderiza a INTENT
│   ├── specgen/             # Gera a SPEC (LLM) conforme TEMPLATE_SPEC; valida o contrato; ErrContract
│   ├── plangen/             # Gera PLAN e ROUTE conforme os templates (Generate com LLM, Render sem LLM)
│   ├── fabricar/            # FÁBRICA DE PROCESSOS: nome, drives, scanner, fusão, Tríade, gravação atômica
│   ├── reputation/          # Matriz de reputação das peças (tiers ouro, prata...)
│   ├── vault/               # Cofre de credenciais (AES-256-GCM, mascaramento)
│   ├── board/ · wave/       # Integração com GitHub Projects e planejador de ondas sem sobreposição
│   ├── studio/              # Looping Studio (canvas + API REST + modais Novo Processo e Processos)
│   ├── e2e/                 # Testes E2E herméticos dos cenários de referência (tag e2e)
│   ├── redc/                # MEGA-PIPELINE REDC: engenharia reversa de canais, DNA, ganchos e Brand Swap
│   ├── mcp/                 # Meta-MCP Router STDIO: bootstrap mínimo (2 tools), busca Lazy e catálogo de skills
│   └── cli/                 # Comandos da CLI, alternador STE-PT e sincronização de harnesses
├── skills/                  # Skills canônicas embutidas no binário (//go:embed via skills/embed.go)
│   ├── legoloop-fabricar/   # Criação de processos (/legoloop-fabricar)
│   ├── legoloop-matriz/     # Consulta à matriz de peças (/legoloop-matriz)
│   ├── legoloop-listar/     # Inspeção de processos e vínculos (/legoloop-listar)
│   ├── legoloop-refinar/    # Refinamento de passos (/legoloop-refinar)
│   ├── legoloop-remover/    # Remoção e desvinculação (/legoloop-remover)
│   ├── legoloop-compor/     # Composição como sub-loop (/legoloop-compor)
│   └── legoloop-canal-redc/ # Engenharia reversa profunda de canais e clonagem (/legoloop-canal-redc)
├── docs/
│   ├── templates/           # Templates canônicos v2 (INTENT, SPEC, PLAN, ROUTE e peças) + conformidade em Go
│   ├── manuals/             # Manuais operacionais dos 6 comandos principais
│   ├── intents/ · specs/ · plans/ · routes/   # Cascata de governança (prefixo ISO de data)
│   ├── reports/ · references/ · sessions/ · archive/
│   └── reports/img/         # Evidências visuais de homologação
├── loops/                   # Loops de exemplo (addy-osmani, impeccable-design, matt-pocock, mega-pipeline)
├── tools/funclen/           # Gate de qualidade: nenhuma função de produção com mais de 35 linhas
└── AGENTS.md · MEMORY.md (diário de implementação) · README.md
```

---

## 🛠️ Início Rápido

```bash
go build -trimpath -ldflags="-s -w" -o bin/legoloop.exe ./cmd/legoloop
go install ./cmd/legoloop                         # deixa o binário atualizado no PATH

# Diagnóstico e Desduplicação
legoloop doctor                                   # modelos executáveis, harnesses, 9Router
legoloop canonize --root DIR --dry-run            # plano + economia estimada (não toca o disco)
legoloop canonize --root DIR --apply              # substitui duplicatas exatas por links

# Fábrica & Ciclo de Vida de Processos
legoloop fabricar "extrator de audio" --yes       # cria processo (skill + loop + dag + docs)
legoloop matriz ver --tier ouro                   # consulta matriz de peças e vocabulário
legoloop listar                                   # lista processos do cofre e vínculos
legoloop refinar "extrator-audio" --etapa e1 --peca aidd-rea
legoloop remover "extrator-audio" --dry-run
legoloop compor "esteira-master" "extrator-audio"
legoloop fechar "extrator-audio" --dry-run        # canonização determinística em 9 etapas
legoloop instalar-skills                          # distribui as 14 skills canônicas nos harnesses

# Mega-Pipelines Audiovisuais & Engenharia Reversa de Canais
legoloop canal redc <dir> --canal meu-canal --json   # analisa 100 vídeos e extrai Dossiê DNA
legoloop canal gerar meu-canal --tema "Novo Tema"     # gera vídeo estrutural com Brand Swap
legoloop video-motion                                # pipeline determinístico de edição e motion
legoloop video-rea                                   # pipeline universal de engenharia reversa audiovisual
legoloop cartoon-shorts --theme "Origem da Pizza"    # pipeline determinístico de cartoon shorts

# Execução e Studio Visual
legoloop route --from Artifact:URL --to Artifact:Spec
legoloop validate pipeline.dag.json
legoloop run pipeline.dag.json --var idea="..."   # executa os nós com gates
legoloop studio                                   # http://127.0.0.1:7878
legoloop forge                                    # abre o Studio e dispara a primeira varredura
legoloop import <github-url>                      # importa e canoniza uma skill do GitHub
legoloop upgrade <arquivo> --diagnose             # diagnóstico para o Padrão Ouro

# Governança Lazy Universal de MCPs e Alternador STE-PT
legoloop sync-harnesses --dry-run                 # simula alinhamento dos harnesses ao Meta-MCP Router
legoloop sync-harnesses --apply                   # conecta Claude, MIMO, Cursor e Windsurf com backup
legoloop mcp-router                               # servidor STDIO com Lazy discovery (mcp_search, call_mcp_tool)
legoloop stept [status|set|on|off|scan|manual]   # auditoria e alternância de governança STE-PT nos harnesses
```

## 🏭 Fábrica & Ciclo de Vida de Processos

O LegoLoop cobre todo o ciclo de vida de esteiras autônomas através de **7 verbos canônicos**:

1. **`fabricar`:** Concepção e criação completa de novos processos em **cinco camadas**: Clarificação (entrevista socrática até 95% de confiança) → Varredura dos discos → Fusão (curadoria no vocabulário fechado) → Tríade (INTENT, SPEC, PLAN e ROUTE) → Template e gravação atômica.
2. **`matriz`:** Classificação das peças em vocabulário rígido, pontuação de reputação/aptidão e detecção de defasagem técnica (`gerar`, `status`, `vocab`, `perfil`, `ver`).
3. **`listar`:** Auditoria de integridade do cofre `~/.aidd/canonical/`, inventário de vínculos ativos nos 5 harnesses e diagnóstico de lacunas.
4. **`refinar`:** Substituição segura e atômica de peças e passos obsoletos, com elevação semântica de versão e preservação de backup.
5. **`remover`:** Remoção atômica de processos com desvinculação limpa de atalhos e hardlinks nos harnesses, com suporte a `--dry-run` e rollback.
6. **`compor`:** Encaixe de processos existentes como sub-loops modulares dentro de esteiras maiores com prevenção de ciclos.
7. **`fechar`:** Cadeia mecânica obrigatória de canonização em 9 etapas: Manual, Livro LegoPress, Skill nos harnesses, Validação DAG, Rebuild do binário, Governança viva, Reindexação semântica, Portões de qualidade e Relatório final com telemetria estrita.

| Superfície | Como Operar |
|---|---|
| **CLI Industrial** | Comandos dedicados (`fabricar`, `matriz`, `listar`, `refinar`, `remover`, `compor`, `fechar`, `instalar-skills`). Saída estruturada via `--json` para scripts |
| **Agente em Linguagem Natural** | 14 Skills canônicas dedicadas (`/legoloop-fabricar`, `/legoloop-matriz`, `/legoloop-listar`, `/legoloop-refinar`, `/legoloop-remover`, `/legoloop-compor`, `/legoloop-fechar`, etc.) |
| **Studio (Dupla Superfície)** | Interface visual completa com botões **+ Novo Processo** e **Processos** na topbar: Modo Cockpit para PhDs e Modo Autopilot com Cartões Executivos limpos para operadores |

**O que é gerado na fabricação:** `SKILL.md` (Chat) e `.loop.json` + `.dag.json` (Studio e CLI, `legoloop run <caminho>`) no cofre `~/.aidd/canonical/`, mais `INTENT`, `SPEC`, `PLAN`, `ROUTE` e `TEMPLATE` em `docs/`. Tudo conforme os templates v2 e validado por `legoloop validate`.

* `--model` é o **cérebro** que conduz a fabricação e suporta múltiplos provedores (9Router VPS, Gemini via harness `agy`, Ollama local com auto-detecção em `internal/brain`).
* A SPEC gerada por LLM é validada contra o contrato; se violar, a Tríade reenvia ao modelo o motivo (até 3 tentativas, só em violação de contrato).
* Processo já existente ⇒ erro; `--force` guarda o original como `.legoloop-bak`. No Studio, "Sobrescrever" só aparece depois do erro de duplicata, com confirmação.
* `fabricar` **cria** processos; `canonize` **desduplica** arquivos. São funções diferentes e independentes.

**Manuais Operacionais Completos em [`docs/manuals/`](docs/manuals/):**
* [`MANUAL_FUNCAO_FABRICAR_LEGOLOOP.md`](docs/manuals/2026-11-03_MANUAL_FUNCAO_FABRICAR_LEGOLOOP.md)
* [`MANUAL_FUNCAO_MATRIZ_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_MATRIZ_LEGOLOOP.md)
* [`MANUAL_FUNCAO_LISTAR_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_LISTAR_LEGOLOOP.md)
* [`MANUAL_FUNCAO_REFINAR_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_REFINAR_LEGOLOOP.md)
* [`MANUAL_FUNCAO_REMOVER_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_REMOVER_LEGOLOOP.md)
* [`MANUAL_FUNCAO_COMPOR_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_COMPOR_LEGOLOOP.md)
* [`MANUAL_FUNCAO_LEGOPRESS_LEGOLOOP.md`](docs/manuals/2026-11-06_MANUAL_FUNCAO_LEGOPRESS_LEGOLOOP.md)
* [`MANUAL_FUNCAO_PIECE_UPGRADER_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_PIECE_UPGRADER_LEGOLOOP.md)
* [`MANUAL_FUNCAO_LOOP_MAP_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_LOOP_MAP_LEGOLOOP.md)
* [`MANUAL_FUNCAO_MEGA_PIPELINE_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_MEGA_PIPELINE_LEGOLOOP.md)
* [`MANUAL_FUNCAO_VIDEO_REA_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_VIDEO_REA_LEGOLOOP.md)
* [`MANUAL_FUNCAO_VIDEO_MOTION_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_VIDEO_MOTION_LEGOLOOP.md)
* [`MANUAL_FUNCAO_VIDEO_MASTER_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_VIDEO_MASTER_LEGOLOOP.md)
* [`MANUAL_FUNCAO_FECHAR_LEGOLOOP.md`](docs/manuals/2026-11-07_MANUAL_FUNCAO_FECHAR_LEGOLOOP.md)

### 🎬 Mega-Pipelines Audiovisuais & Engenharia Reversa (`video-rea`)
* **Engenharia Reversa Audiovisual (`legoloop video-rea`):** Disseca URLs ou mídias locais criando Target Workspaces isolados em `docs/videorea/<slug>/`.
  * **Fase 1 (Forense & Matriz Original):** Extração de transcrição e OCR rápido em `dossier.json`, alimentando a nova camada `original_matrix.json` (stack real do autor, ferramentas pagas vs gratuitas explícitas e comandos implícitos de terminal/código).
  * **Fase 2 (Curadoria & Soberania):** Substituição de ferramentas comerciais por peças Tier Ouro do cofre e **Adoção Soberana Imediata** de utilitários gratuitos vazados no vídeo (custo R$ 0,00).
  * **Fase 3 & 4 (Governança, Livro e Loop):** Fabricação da Tríade ISO local (`governance/`), Livro Didático LegoPress (`book/`), DAG com Golden Fixture de 3s (`loop/`), com projeção canônica na raiz (`loops/` e `skills/`) mediante aprovação do Stop-Gate humano.

## 🏭 Fundição Autônoma N10 & Protocolos Industriais

O LegoLoop opera como uma fundição de engenharia com blindagem de grau aeroespacial e automotivo:

### Os 7 Protocolos Industriais
1. **STANAG 2014:** Ordens estruturadas de missão (Situação, Missão, Execução, Administração, Comando/Sinal).
2. **ASD-STE100 (Simplified Technical English) a 80%:** Frases curtas (<20 palavras), voz ativa, sem ambiguidades, 1 instrução por sentença.
3. **EARS (Easy Approach to Requirements Syntax):** Critérios de aceitação binários em formato Ubíquo, Orientado a Evento, Orientado a Estado ou Opcional.
4. **FMEA (Failure Mode and Effects Analysis):** Matriz de risco RPN prévia (Severidade × Ocorrência × Detecção) antes de tocar em código.
5. **ICAO Readback/Hearback:** Confirmação explícita no terminal com exibição de blast radius antes de comandos destrutivos.
6. **Poka-Yoke & Disjuntor Andon Cord (Jidoka):** Interrupção automática imediata ao menor erro mecânico para impedir efeito cascata.
7. **SBAR (Situation, Background, Assessment, Recommendation):** Formato universal de escalada de incidentes e rotação de contexto.

### As 5 Prisões Mecânicas
1. **Path Jail:** Proibição de escape de diretório (`..`, symlinks externos, caminhos absolutos desautorizados).
2. **Schema Validation:** Validação mecânica estrita de todos os JSONs contra esquemas formais em `schemas/`.
3. **Exit-Code Gate:** Passos só avançam com confirmação mecânica de `exit_code == 0` do sistema operacional.
4. **Sterile Context:** Cada bloco executa isolado sem herança poluída de prompt, consumindo estritamente `context_in` e emitindo `context_out`.
5. **Test Mutation:** Teste unitário só é certificado se a mutação forçada do código o fizer falhar.

### Os 6 Motores de Fronteira N10
1. **Mutation Testing:** Testagem de mutação sintática de código.
2. **Git Worktrees com AST:** Execução isolada em branches efêmeras com merge por AST limpo.
3. **Memória Episódica de Anticorpos SBAR:** Cofre persistente de vacinas contra erros já diagnosticados.
4. **Cascade Speculative Execution:** Execução antecipada e paralela de nós prováveis.
5. **Flight Recorder:** Gravador de caixa preta com telemetria determinística de tempo, memória e tokens.
6. **Fuzzing Nativo:** Fuzzing embutido em Go para estresse contínuo de parsers e validadores.

## 🛡️ Governança & Alternador STE-PT (`legoloop stept`)

O LegoLoop inclui governança determinística e alternador universal de regras para os 5 harnesses do sistema (**Google Gemini/Antigravity**, **Claude Code**, **OMP CLI**, **OpenCode** e **MimoCode**), além do executável standalone `stept.exe` no PATH:

* **Gerenciamento Unificado:** `legoloop stept [status|set|on|off|scan|manual]`
* **Modos Suportados:** `hibrido` (padrão recomendado), `stept`, `caveman`, `off`.
* **Políticas de Pausa:** `destrutivo` (pausa em deleções/resets), `sempre` (pausa em toda alteração com efeito), `nunca` (informativo).
* **Contrato Canônico da Tríade STE-PT (Hierarquia Visual de Bold Limpa):**
  * **Pré-Ação (antes de alterar arquivos):**
    ```markdown
    **[S] Situação:** <Cenário factual do problema ou contexto encontrado>
    **[T] Teoria:** <Hipótese técnica e plano de ação antes de alterar código>
    **[E] Execução:**
    · **Ação:** <O que vai ser feito com verbo imperativo e arquivos tocados>
    · **Rollback:** <Comando exato de recuperação caso falhe>
    · **Decisão:** **[1] seguir** | **[2] personalizar** | **[0] cancelar**
    ```
  * **Pós-Ação (após executar e testar):**
    ```markdown
    **[S] Situação:** <O que foi executado no sistema>
    **[T] Tese:** <A tese técnica comprovada pela alteração>
    **[E] Evidência:**
    · **Testes:** <Comando de teste executado e resultado real>
    · **Build:** <Status de compilação ou verificação real do sistema>
    · **Verificação:** <Resultado mensurável final>
    ```

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
go test -count=1 ./...                         # 29 pacotes internos; lógica JS pura do Studio roda em node
go test -count=1 -tags e2e ./internal/e2e/...  # suíte hermética E2E (3 cenários de referência e 2 cérebros)
go vet ./...                                   # verificação estática (Windows, Linux e Darwin)
go run ./tools/funclen internal                # nenhuma função > 35 linhas
gofmt -l .                                     # deve sair vazio (o exit code é 0 mesmo listando: confira a saída)
git diff --exit-code -- go.mod                 # só biblioteca padrão (sem go.sum)
```

`go test -race` exige CGO (indisponível no Windows sem `gcc`): roda no CI Linux. Estado em 06/11/2026: 29 pacotes internos + pacote E2E 100% verdes, 0 falhas, 1 pulado preexistente (symlink sem admin no Windows).

## 📚 Governança

Cada demanda segue a cascata **INTENT → SPEC → PLAN → ROUTE → BUILD (TDD) → relatório**, com stop-gates humanos e artefatos em `docs/` com prefixo `YYYY-MM-DD_`. Os formatos estão em `docs/templates/` (templates v2.2.0, testados por `go test ./docs/templates`). Demandas recentes: N10 Foundry Harness & Protocolos Industriais (`docs/reports/2026-11-09_REPORT_N10_FOUNDRY_GLOBAL.md`), Governança de Rota e Nós v2.2.0 (`docs/*/2026-11-07_*_ROUTE_V2_GOVERNANCE.*`), Hierarquia Canônica de 4 Níveis (`docs/*/2026-11-07_*_HIERARQUIA_INICIATIVA.*`), Governança Lazy de MCPs (`docs/*/2026-11-06_*_GOVERNANCA_LAZY_UNIVERSAL.*`) e Alternador STE-PT (`docs/*/2026-11-05_*_STEPT_SWITCH.*`). Manuais operacionais em `docs/manuals/`. Regras para agentes em [`AGENTS.md`](AGENTS.md); diário de evolução em [`MEMORY.md`](MEMORY.md).
