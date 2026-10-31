# 🚀 LegoLoop — Plano de Execução & Micro/Nano Tickets (Master Plan)

> **Metodologia:** Tracer Bullet Development + TDD Estrito + Build Auto (Addy Osmani)  
> **Estratégia de Execução:** Fatiamento vertical em 5 Épicos com Nano-Tickets atômicos.

---

## 🗺️ Mapa de Épicos

| Épico | Título | Objetivo Central |
|---|---|---|
| **EP-01** | **Harness Engineering & Canonizador de SO** | Script de varredura por hash SHA-256 e substituição dos 22k ativos por Symlinks transparentes sem quebra. |
| **EP-02** | **Core Runtime & DAG Runner (Go Engine)** | Engine compilada em Go que lê `pipeline.dag.json` e executa nós com gates determinísticos. |
| **EP-03** | **Orchestrix Cirúrgica & Provider-Awareness** | Despachante isolado por bloco (single-task context) com auto-discovery de Ollama, CLIs e 9Router. |
| **EP-04** | **Graph Engineering & Matriz de Compatibilidade** | Grafo semântico de nós e arestas (consome/produz) com cálculo de menor caminho. |
| **EP-05** | **Looping Studio UI (Canvas N8N no Symbiont)** | Interface visual drag-and-drop de nós, portas e exportação de esteiras prontas. |

---

## 📋 Micro e Nano Tickets

### EP-01: Harness Engineering & Canonização do SO
- [x] **TICK-101 (Nano):** Criar `internal/canon/hasher.go` para calcular SHA-256 e metadados de arquivos Markdown/JSON em worker pool.
- [x] **TICK-102 (Nano):** Implementar eleitor de versão canônica (critério: mais recente + maior cobertura de frontmatter).
- [x] **TICK-103 (Nano):** Criar executor de Symlinks seguro para Windows (`os.Symlink` com fallback para `mklink /D` e `mklink /H`).
- [x] **TICK-104 (Micro):** Implementar comando CLI `legoloop canonize --dry-run` e `--apply` com relatório de espaço economizado.

### EP-02: Core Runtime & DAG Runner
- [x] **TICK-201 (Nano):** Definir types em Go para `DAG`, `Node`, `RoutingConfig`, `HandoffContract` e `GateExpression`.
- [x] **TICK-202 (Nano):** Implementar parser e validador de ciclo (detecção de Deadlock em Grafo Acíclico Dirigido).
- [x] **TICK-203 (Nano):** Criar motor de validação de Gates (`file_exists`, `has_pattern`, `exit_code == 0`).
- [x] **TICK-204 (Micro):** Implementar runner sequencial e paralelo de nós com timeout individual e telemetria de tokens.

### EP-03: Orchestrix Cirúrgica & Provider-Awareness
- [x] **TICK-301 (Nano):** Criar scanner passivo de Ollama (`GET http://127.0.0.1:11434/api/tags`) com timeout de 200ms.
- [x] **TICK-302 (Nano):** Criar detector de CLIs ativas (`claude`, `agy`, `opencode`) no `PATH` do sistema.
- [x] **TICK-303 (Nano):** Criar conector de 9Router (`GET /v1/models` na porta `:20128` ou `:9000`) para listar combos ativos.
- [x] **TICK-304 (Micro):** Implementar a casca da **Orchestrix Isolada**: dispara processo efêmero com prompt único, sem herança de histórico.

### EP-04: Graph Engineering
- [x] **TICK-401 (Nano):** Definir ontologia de tipos de artefato (`Artifact:Briefing`, `Artifact:Spec`, `Artifact:Code`, `Artifact:Audit`).
- [x] **TICK-402 (Nano):** Construir índice de compatibilidade: mapeia quais nós podem se conectar a quais portas.
- [x] **TICK-403 (Micro):** Implementar algoritmo Dijkstra/A* para cálculo de rota ótima entre uma entrada e um objetivo.

### EP-05: Looping Studio UI (Canvas N8N)
- [x] **TICK-501 (Nano):** Integrar biblioteca de nós gráficos leves no frontend do Symbiont (`index.html`).
- [x] **TICK-502 (Nano):** Renderizar nós com crachás visuais: Handoff Fixo, Roteamento (Sugerido/Manual/Auto) e Status do Gate.
- [x] **TICK-503 (Micro):** Adicionar botão *"Materializar Esteira"*, que valida o grafo e salva `pipeline.dag.json` no projeto.

---

## 🔄 Protocolo de Implementação Autônoma (Build Auto do Addy Osmani)

Para executar os tickets acima, o loop segue o ciclo de 6 etapas inquebráveis:
1. **Build:** Compila o código em Go/TypeScript sem warnings.
2. **Test:** Roda suíte TDD unitária e de integração (cobertura mínima de 85%).
3. **Review:** Auditoria de padrões (Zero Stubs, tipos estritos, nomes expressivos).
4. **Web Performance:** Validação de tempo de inicialização (<100ms) e consumo de RAM (<45MB).
5. **Code Simplify:** Refatoração para eliminar duplicações e funções longas.
6. **Cloudflare Security Pass:** Auditoria de segurança contra injeção de comandos ou escalonamento de privilégios.
