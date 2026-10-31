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
* **Canonização de SO sem Quebras:** Centraliza seus ativos em `~/.aidd/canonical/` e usa Symlinks (`mklink`) para que seus projetos antigos continuem funcionando com 80% menos espaço em disco.
* **Convergência Total:** Integra as regras do **Ecossistema AIDD**, o motor de alta performance do **SYMBIONT** e a blindagem de modelos do **9Router**.

---

## 📂 Estrutura do Projeto

```
LegoLoop/
├── cmd/legoloop/            # Binário (main de 7 linhas; lógica em internal/cli)
├── internal/
│   ├── canon/               # EP-01 Canonização: hash SHA-256, eleitor, links com rollback
│   ├── dag/                 # EP-02 Contrato pipeline.dag.json, validação e ciclos
│   ├── gate/                # EP-02 Motor de gates (file_exists, has_pattern, exit_code...)
│   ├── runner/              # EP-02 Runner por camadas, paralelo, timeout, telemetria
│   ├── provider/            # EP-03 Inventário passivo: Ollama, 9Router, CLIs, chaves
│   ├── orchestrix/          # EP-03 Orchestrix isolada por bloco (claude -p, Ollama, 9Router, script)
│   ├── graph/               # EP-04 Ontologia, compatibilidade de portas, rota ótima (Dijkstra)
│   ├── studio/              # EP-05 Looping Studio (canvas + API local) embutido no binário
│   └── cli/                 # Comandos
├── SPEC.md · PLAN.md · AGENTS.md · MEMORY.md (diário de implementação)
├── docs/specs/ · docs/references/
└── secoes/                  # Histórico factual gravado via faz-session
```

---

## 🛠️ Início Rápido

```bash
go build -trimpath -ldflags="-s -w" -o bin/legoloop.exe ./cmd/legoloop

legoloop doctor                                   # modelos executáveis, harnesses, 9Router
legoloop canonize --root DIR --dry-run            # plano + economia estimada (não toca o disco)
legoloop canonize --root DIR --apply              # substitui duplicatas exatas por links
legoloop route --from Artifact:URL --to Artifact:Spec
legoloop validate pipeline.dag.json
legoloop run pipeline.dag.json --var idea="..."   # executa os nós com gates
legoloop studio                                   # http://127.0.0.1:7878
```

## 🧠 Executores (roteamento por nó)

| `routing.model` | Executor | Observação medida (31/10/2026) |
|---|---|---|
| combo do 9Router (`code-pro`, `code-fast`…) | 9Router `/v1/chat/completions` | requer `NINEROUTER_URL` (https) e `NINEROUTER_KEY` (ou `NINEROUTER_API_KEY`) |
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
| `NINEROUTER_URL` só https ou http-loopback; a chave vai apenas para a base confiável | credencial não trafega em claro nem para portas genéricas |
| Studio só escuta em loopback, valida Host/Origin, exige JSON, limita corpo e taxa | anti DNS-rebinding e anti-CSRF |

**Riscos aceitos:** com fallback por hardlink (Windows sem Developer Mode), um editor que grave *no próprio arquivo* altera todas as cópias linkadas — é a natureza do hardlink; o runner do LegoLoop sempre grava por temporário + rename. No Windows, processos netos de um nó podem ficar órfãos após timeout (sem Job Object); o runner não trava.
