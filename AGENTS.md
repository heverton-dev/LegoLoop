# 🤖 AGENTS.md — Regras de Operação para Agentes no LegoLoop

> Este arquivo é a Lei de Conduta para qualquer LLM ou agente operando neste repositório.

---

## 1. Leis Invioláveis (Zero Tolerância) & As 5 Prisões Mecânicas

1. **Zero Stubs:** Nunca crie funções vazias com `TODO`, `pass`, `// implementar depois` ou mocks simulados. O código entregue deve ser real, testado e funcional.
2. **Graph-First Obrigatório (Lei #14):**
   * SEMPRE consulte o grafo semântico via `codebase-memory-mcp` (`search_code`, `get_code_snippet`, `trace_path`) ANTES de ler arquivos extensos (> 100 linhas) ou rodar buscas cegas no disco.
   * Só use `read_file` pontual com `StartLine`/`EndLine` sobre o trecho localizado pelo grafo.
   * O grafo mora em `.codebase-memory/` (artefato gerado: **nunca edite à mão**). Para atualizá-lo: `codebase-memory-mcp cli index_repository '{"repo_path":"<raiz>"}'`. Reindexe ao fim de cada demanda.
3. **Economia Severa de Tokens & Rule 10 (Silent Executor):**
   * **Thinking:** Estilo telegráfico (Caveman), direto ao ponto, sem repetições.
   * **Teto de Saída:** Respostas normais $\le 300$ tokens; explicações técnicas $\le 600$ tokens.
   * **Sem Polling em Background Tasks:** Proibido rodar loops de `manage_task(status)`. Aguarde a notificação de *Reactive Wakeup* do harness.
   * Não repita artefatos ou diffs no corpo da mensagem se já estiverem salvos no disco.
4. **Orchestrix Cirúrgica por Bloco & Banimento de Subagentes:**
   * Nunca carregue ferramentas ou MCPs desnecessários para a tarefa atual.
   * Proibido usar agente de harness (`invoke_subagent`, fork, workflow livre) que clona o system prompt inteiro (queima de tokens). Toda execução pós-plano deve seguir o mapa formal em `docs/routes/` com execução estéril por bloco.
   * Exceção: o subagente efêmero (`isolation.subagents_allowed: "ephemeral"`) é permitido quando cumpre as seis condições:
     1. O contexto contém só os artefatos do ticket.
     2. O ticket lista as ferramentas permitidas.
     3. Um subagente executa um ticket.
     4. A saída fica em disco e passa pelo gate mecânico.
     5. O subagente não cria sub-subagentes.
     6. O `write_boundary` é `strict`.
   * Execute cada nó em contexto estéril; leia apenas o artefato de entrada e produza o artefato de saída contratado.
5. **TDD Estrito & Recompilação Fresca:**
   * O teste unitário deve falhar antes de qualquer implementação de código de produção (Red-Green-Refactor).
   * Garanta que todos os testes passem com código limpo e sem warnings.
   * **Recompilação Obrigatória Pós-Ticket:** Ao finalizar qualquer ticket com alteração em Go, execute `go build -o bin/legoloop.exe ./cmd/legoloop` e, ao fechar o épico, `go install ./cmd/legoloop`, para que o binário no PATH permaneça sempre atualizado. Um Studio já aberto continua no binário antigo: reinicie-o.
6. **Quinteto de Governança ISO de Artefatos:**
   * TODOS os artefatos do Quinteto gerados em `docs/` (`intents/`, `specs/`, `plans/`, `routes/`, `dags/`, `reports/`, `sessions/`) DEVEM obrigatoriamente utilizar o prefixo de data no formato ISO: `YYYY-MM-DD_<NOME>.<ext>`.
   * Isso garante ordenação cronológica natural e busca imediata no SO.
   * **INTENT, SPEC, PLAN, ROUTE e DAG** compõem o **Quinteto de Governança ISO** e seguem os templates canônicos v2 em `docs/templates/` (conformidade testada em `docs/templates/templates_test.go`). O `DAG.json` em `docs/dags/` é o espelho executável 1:1 da ROUTE, projetado atomicamente para `processos/<slug>/pipeline.dag.json` e com paridade SHA-256 auditada mecanicamente no fechamento (`legoloop fechar`). A Skill (`SKILL.md`) é a manifestação em prosa fiel dos nós do DAG.
7. **Idioma:** SEMPRE utilize o idioma Português do Brasil (PT-BR) em documentações, planos, especificações e respostas.
8. **Verificação de Verdade (lições da demanda `fabricar`):**
   * **Rode o sistema real antes de declarar pronto** (binário e, se houver UI, o navegador). LLM falso e fixtures só provam o que o autor imaginou: o validador de SPEC recusava a SPEC correta de um LLM real e nenhum teste unitário via isso.
   * **Teste novo só vale se uma mutação do código o faz falhar.** Mutou, não falhou ⇒ o teste está errado. Restaure o código depois de cada mutação.
   * **Portão de qualidade de cada ticket:** `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal` (nenhuma função > 35 linhas), `gofmt -l .` com saída **vazia** (não falha por exit code: confira a saída), `git diff --exit-code -- go.mod go.sum` (só stdlib). `-race` só roda com CGO: obrigatório no CI Linux.
   * **Hermético por padrão:** testes e smoke usam LLM falso, raízes em `t.TempDir()`/`scratch/` e **nunca** tocam o cofre real `~/.aidd`. Limpe o que criou.
   * **Segredos:** a chave do 9Router só vai para a URL confiável (`provider.TrustedRouterKey`), nunca em log, saída ou resposta de API. Respostas do Studio nunca trazem caminho absoluto (`~` e `<projeto>`).
   * **Stop-gates (`SG-*`) são humanos:** nunca marque um gate como aprovado sem autorização explícita do usuário. Commit e push também só quando pedidos.
   * **Processos:** nunca `taskkill /IM legoloop.exe` (mata o Studio do usuário). Encerre só o PID ou a porta que você mesmo abriu.
   * **Declare o que não foi provado.** Relatórios dizem o que ficou aberto (ex.: `-race`, fluxo com LLM real) em vez de arredondar para "pronto".
9. **Diagramação Visual Impecável (Zero Quebras em Mermaid e Tabelas):**
   * **Proibido HTML em Mermaid:** NUNCA use `<br>`, `<br/>`, `<b>` ou qualquer tag HTML em nós ou arestas. Isso quebra o renderizador do chat e expõe o código cru.
   * **Aspas Obrigatórias:** Todo rótulo de nó deve estar entre aspas duplas: `["Label"]`, `{"Decisão"}`.
   * **Arestas Seguras:** Nunca use caminhos de arquivos ou barras (`/`) em rótulos de arestas (`-->|"Aresta Segura"|`).
   * **Tabelas Enxutas:** Alinhamento estrito (`:---`, `:---:`) e colunas concisas para evitar rolagem horizontal na janela do chat.
10. **Comunicação Canônica STE-PT (modo `prova`):**
   * Fonte única: `internal/ste` (`legoloop stept manual`). Este item resume o contrato e não o substitui.
   * Estrutura fixa da resposta técnica: **Resultado** (1 frase), **Fatos comprovados** (`fato | prova: arquivo:linha ou comando e resultado`; sem prova, escrever "Não verificado"), **Solução passo a passo**, **Sugestão técnica**, **Stop-Gate Humano** (`[1] seguir | [2] personalizar | [0] cancelar`).
   * Stop-Gate conforme a pausa ativa: `destrutivo` (padrão), `sempre` ou `nunca`.
   * A Tríade `[S][T][E]` vale só se o modo ativo for `hibrido` ou `stept`.
   * Proibido preâmbulos, saudações ou re-resumos de arquivos no chat.
11. **Protocolo Inviolável de Fechamento & Canonização Universal (Zero Entregas pela Metade):**
   * Nenhuma implementação de processo, peça ou pipeline é declarada concluída sem disparar a cadeia obrigatória de fechamento:
     1. **Manual Operacional:** em `docs/manuals/` (`YYYY-MM-DD_MANUAL_FUNCAO_<NOME>_LEGOLOOP.md`).
     2. **Livro Didático LegoPress:** em `docs/books/<slug>/` (`livro.json`, `partes/`, `.md`, `.pdf` compilado e auditado com `legoloop legopress auditar`).
     3. **Skill Canônica Universal:** em `skills/<nome>/SKILL.md` (distribuída universalmente para todos os harnesses via `legoloop instalar-skills`).
     4. **Loop / Grafo DAG:** em `processos/<slug>/pipeline.dag.json` (e `.loop.json`), validado por `legoloop validate`.
     5. **Recompilação e Instalação Fresca:** `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop` para manter o PATH e o runtime do usuário atualizados.
     6. **Atualização da Governança Viva:** `AGENTS.md`, `MEMORY.md`, `README.md`.
     7. **Reindexação do Grafo Semântico:** via `codebase-memory-mcp cli index_repository`.
     8. **Portões Mecânicos de Qualidade & Commit/Push:** testes 100% verdes, gofmt limpo, go vet sem warnings, funclen <= 35 linhas, commit e push rastreáveis.
     9. **Relatório de Conclusão & Telemetria Estrita:** em `docs/reports/` (`YYYY-MM-DD_REPORT_<NOME>.md`) compilado conforme `TEMPLATE_REPORT.md` com balanço auditável de tempo, recursos/tokens e custo financeiro USD/BRL orçado vs realizado.
12. **Inviolabilidade Taxonômica ASD-STE100 (Lei #15):**
   * Regra 1: um termo, um significado; um significado, um termo. Zero sinônimo em raciocínio, planos, comandos de terminal, código, JSON e Markdown.
   * SSOT: `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json` (embutido em `internal/taxonomy`). O `tools/taxlinter` aplica a lei.

     | Conceito | Termo canônico | Sinônimos proibidos |
     | :--- | :--- | :--- |
     | Fluxo de alto nível | Processo | esteira, fluxo, pipeline, loop |
     | Grafo de execução | DAG | grafo de execução, árvore de execução, cadeia de dependência |
     | Unidade reutilizável | Peça | componente, bloco, módulo, brick |
     | Instrução de agente | Skill | habilidade, capacidade, macro, prompt de agente |
     | Vértice de execução unitária | Nó | etapa de grafo, passo de grafo, task |
     | Unidade de planejamento | Ticket | tarefa, card, issue |
     | Passo interativo CLI | Etapa | fase, estágio |
     | Armazém canônico (`~/.aidd/canonical/`) | Cofre | repositório central, banco canônico, storage |
     | Armazém de credenciais | Cofre de Segredos | keystore, cofre de senhas, secrets |
     | Tabela de aptidão e reputação | Matriz | tabela de aptidão, grade |
     | Catálogo físico em disco | Inventário | catálogo físico, almoxarifado |
     | Barreira mecânica binária | Gate | portão automático, trava automática |
     | Ponto de aprovação humana | Stop-Gate | portão humano, aprovação manual, checkpoint humano |
     | Criação de processo | Fabricar | criar processo, gerar pipeline, construir processo |
     | Deduplicação por hash SHA-256 | Canonizar | desduplicar, normalizar hash, deduplicar |
     | Homologação e conclusão | Fechar | concluir processo, encerrar ciclo, finalizar entrega |

   * Exceção de marca: o nome do projeto "LegoLoop" fica mantido. O termo "loop" dentro da marca não conta como violação.
   * A tabela completa (16 domínios) vive no dicionário. Em divergência, o dicionário vence.
   * Portão de qualidade de cada ticket: `go run ./tools/taxlinter internal cmd loops skills docs/templates` (exit 1 com violação). O portão cobre os alvos já normalizados. Alvo com violação herdada exige ticket de normalização antes de entrar no portão.
13. **Migração de Demandas Multi-Estágio (Lei 10):** Toda demanda multi-estágio é um **Processo** e tem sua **Matriz de Compatibilidade** de Motores em `loop.json`. A migração **nunca altera código Go** (extensão `pessoal-*` no componente só pode ser renomeada; o núcleo é intocável). Cada troca de motor num passo de `pessoal-*` tem 4 etapas: (1) **Diagnóstico** `leitura-matriz` via `legoloop migracao status <slug> --json`, (2) **Oráculo** com **Prova por Contradição Controlada** (semântica divergente ⇒ `esperado ≠ obtido` ⇒ evidência falsificado: `<evidencia>`, com grep/verbo do caminho real), (3) **Migração** `legoloop migracao trocar <slug> --etapa <ID> --motor <novo> --yes`, (4) **Homologação** E2E `rodar-processo real --exec` com `INTEGRIDADE=OK` (não é portão cego: resposta real e motivo do payload).
14. **Contradição Controlada de Legado (`steps-execucao.json`):** O legado que a Matriz bloqueia não é descartado: **figura como contraprova** no histórico (`leitura-matriz --gravar` em `~/.aidd/migracao/<slug>-historico.json`) e ganha **Vínculo Vivo** (remendo de `<PROJETO>/steps-execucao.json` por `<slug>`, transformando-se em pistas da próxima demanda). Drop de legado do disco exige **Stop-Gate humano**, e validação on-demand confirma `obsoleto=true` + `dono=<slug>` antes de qualquer `migracao drop-legado`.
15. **Lexicais de Diagnóstico de Nós (`legoloop lexica`):** Toda avaliação de nó de **linguagem natural** usa **critério de separação de nós** com score `{cor, forma, opacidade, tamanho, elemento, palavra, consistencia}` e uma **palavra-lei** por nó (objetivo claro). Nó avaliado sem verbo aplica **gravação com restrição** e a **Palavra-Lei**; nó que falha vira **objeto de lacuna** com palavra-chave e recomendação de peça de sistema em `pipeline.dag.json`. O diagnóstico real roda com `legoloop lexica avaliar <alvo> --json` sobre alvos de `web-qa`/`design-linter`/`cartoon-shorts`/`fabricar`/`executar`/`validar`.
16. **Achado Corrigido no Ato, Universalidade e Governança Viva:**
   * Todo achado de análise ou varredura é corrigido na execução, com teste Red antes. Decisão técnica sem consenso para antes da correção e sobe ao humano com recomendação.
   * `routing` declara preferência, nunca exigência: nós `nano` e `micro` degradam para o próximo candidato homologado. `tooling.runtime` vale `legoloop`. Nenhum `command` fixa shell de um só SO.
   * Toda implementação nova ou concluída atualiza `AGENTS.md`, `README.md`, `MEMORY.md` e o grafo (`index_repository`) no mesmo ciclo.

---

## 2. Os 7 Protocolos Industriais e as 5 Prisões Mecânicas

### 🔒 As 5 Prisões Mecânicas (Invioláveis em Runtime)
1. **Path Jail:** Proibição absoluta de escape de diretório (`..`, symlinks externos, caminhos absolutos desautorizados). Operações em disco limitadas à raiz sanitizada.
2. **Schema Validation:** Todos os manifestos JSON (`pipeline.dag.json`, `route.v2.json`, `spec.v1.json`, `livro.json`) validados estritamente contra JSON Schema formal (`schemas/dag.v1.json`, etc.).
3. **Exit-Code Gate:** Passos automatizados só avançam com confirmação mecânica de `exit_code == 0` do sistema operacional.
4. **Sterile Context:** Cada bloco executa isolado sem herança poluída de prompt, consumindo estritamente `context_in` e emitindo `context_out`.
5. **Test Mutation:** Teste unitário só é certificado se a mutação forçada do código sob teste o fizer falhar mecanicamente (Green-Red-Green comprovado).

### 📐 Os 7 Protocolos Industriais
1. **STANAG 2014:** Ordens estruturadas de missão (Situação, Missão, Execução, Administração, Comando/Sinal).
2. **ASD-STE100 (Simplified Technical English) a 80%:** Frases curtas (<20 palavras), voz ativa, sem ambiguidades, 1 instrução por sentença.
3. **EARS (Easy Approach to Requirements Syntax):** Critérios de aceitação binários em formato Ubíquo, Orientado a Evento, Orientado a Estado ou Opcional.
4. **FMEA (Failure Mode and Effects Analysis):** Matriz de risco RPN prévia (Severidade × Ocorrência × Detecção) antes de tocar em código.
5. **ICAO Readback/Hearback:** Confirmação explícita no terminal com exibição de blast radius antes de comandos destrutivos.
6. **Poka-Yoke & Disjuntor Andon Cord (Jidoka):** Interrupção automática imediata ao menor erro mecânico para impedir efeito cascata.
7. **SBAR (Situation, Background, Assessment, Recommendation):** Formato universal de escalada de incidentes e rotação de contexto.

---

## 3. Comandos & Rotas de Governança

* **Validar conformidade e testes:** `go test -count=1 ./...` e `go test -count=1 -tags e2e ./internal/e2e/...` (ou `npm test` quando houver front-end com testes JS; a lógica pura do Studio roda em `node`).
* **Fábrica & Ciclo de Vida dos 7 Verbos de Processo:**
  * **Fabricar (criação completa):** `legoloop fabricar "<nome>" [--yes] [--goal T] [--model M] [--category build|reverse|web|mobile] [--dry-run] [--json]` (Chat: `/legoloop-fabricar`; Studio: **+ Novo Processo**). **Não confundir com `legoloop canonize`** (desduplicação por SHA-256).
  * **Matriz de Peças (reputação/aptidão):** `legoloop matriz ver [--tier ouro|prata|bronze] [--json]` (Chat: `/legoloop-matriz`).
  * **Listar (inventário e vínculos):** `legoloop listar [--json]` (Chat: `/legoloop-listar`; Studio: **Processos**).
  * **Refinar (substituição atômica de passos):** `legoloop refinar "<nome>" --etapa ID [--peca P] [--json]` (Chat: `/legoloop-refinar`).
  * **Remover (desvinculação com rollback):** `legoloop remover "<nome>" [--dry-run] [--force] [--json]` (Chat: `/legoloop-remover`).
  * **Compor (sub-loops hierárquicos):** `legoloop compor "<alvo>" "<sub-loop>" [--pos N] [--json]` (Chat: `/legoloop-compor`).
  * **Fechar (canonização determinística em 9 etapas com telemetria):** `legoloop fechar "<slug>" [--dry-run] [--skip-build] [--json]` (Chat: `/legoloop-fechar`).
  * **Video-REA (engenharia reversa e dissecação com Matriz Original e Adoção Soberana):** `legoloop video-rea [run|ingest|forensics|curate|gate|govern|book] <args>` (Chat: `/legoloop-video-rea`).
  * **Cartoon-Shorts (pipeline determinístico 1080x1920 30fps):** `legoloop cartoon-shorts --theme "<tema>"` (Chat: `/legoloop-cartoon-shorts`).
* **Instalação Canônica das Skills:** `legoloop instalar-skills [--vault DIR] [--dry-run] [--json]` (distribui as 39 skills nativas embutidas para o cofre e cria vínculos nos harnesses).
* **Demanda N10 Foundry Harness & Motores Globais (09/11/2026):** [Manual Operacional](docs/manuals/2026-11-09_MANUAL_FUNCAO_N10_FOUNDRY_HARNESS_LEGOLOOP.md) · [Livro Didático](docs/books/n10-foundry/livro.json) · [Skill](skills/legoloop-n10-foundry/SKILL.md) · [DAG](processos/legoloop-n10-foundry/pipeline.dag.json) · [Relatório final](docs/reports/2026-11-09_REPORT_N10_FOUNDRY_GLOBAL.md)
* **Gateway de LLM no AGY por variáveis de ambiente (13/11/2026):** [Manual Operacional](docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md) · [Livro Didático](docs/books/gateway-llm-agy/livro.json) · [Skill](skills/legoloop-gateway-agy/SKILL.md). Nunca faça rodízio de contas Google no `agy`: o `gw` bloqueia `agy next|failover|devNN|sync|probe --all|refresh all`.
* **Demanda Mega-Pipeline REDC (07/11/2026):** [INTENT](docs/intents/2026-11-07_INTENT_CANAL_REDC.md) · [SPEC](docs/specs/2026-11-07_SPEC_CANAL_REDC.md) · [PLAN](docs/plans/2026-11-07_PLAN_CANAL_REDC.md) · [ROUTE](docs/routes/2026-11-07_ROUTE_CANAL_REDC.json) · [Relatório final](docs/reports/2026-11-07_REPORT_CANAL_REDC.md)
* **Demanda Mega-Pipeline Audiovisual Soberano (07/11/2026):** [INTENT](docs/intents/2026-11-07_INTENT_VIDEO_MASTER.md) · [SPEC](docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md) · [PLAN](docs/plans/2026-11-07_PLAN_VIDEO_MASTER.md) · [ROUTE](docs/routes/2026-11-07_ROUTE_VIDEO_MASTER.json) · [Relatório final](docs/reports/2026-11-07_REPORT_VIDEO_MASTER.md)
* **Governança Lazy Universal de MCPs:** `legoloop sync-harnesses [--dry-run|--apply]` (conecta Claude, MIMO, Cursor e Windsurf ao Meta-MCP Router) e `legoloop mcp-router` (servidor STDIO leve com bootstrap de apenas 2 tools).
* **Alternador de Governança STE-PT:** `legoloop stept [status|set|on|off|scan|manual] [--harness H] [--modo M] [--pausa P] [--json]` (auditoria e chaveamento em todos os harnesses).
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all` (saída em `~/Desktop/sessoes-agentes/`; o script pode ignorar `--id` e exportar o transcript ativo).
* **Consultar diretórios:** [docs/templates/](docs/templates) (templates canônicos v2, SSOT) · [docs/manuals/](docs/manuals) (manuais operacionais) · [docs/intents/](docs/intents) · [docs/routes/](docs/routes) (mapas de rota) · [docs/reports/](docs/reports) (relatórios; capturas em `docs/reports/img/`)
* **Consultar a Spec Master:** [docs/specs/2026-10-30_SPEC_MASTER.md](docs/specs/2026-10-30_SPEC_MASTER.md)
* **Consultar o Plano de Tickets:** [docs/plans/2026-10-30_PLAN_MASTER.md](docs/plans/2026-10-30_PLAN_MASTER.md)
* **Demanda Memória Concorrente e Anticorpos (16/11/2026):** [Manual Operacional](docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md) · [Livro Didático](docs/books/memoria-anticorpos/livro.json) · [Skill](skills/legoloop-memoria-anticorpos/SKILL.md) · [Relatório de integração](docs/reports/2026-11-16_RELATORIO_INTEGRACAO_MEMORIA_ANTICORPOS.md)
* **Demanda Nós e Peças: unit_no, universalidade e templates (20/11/2026) — mais recente:** [Relatório](docs/reports/2026-11-20_NOS_E_PECAS.md) · [Template](docs/templates/TEMPLATE_PIECE_UNIT.json) · [Apoio](docs/templates/TEMPLATE_PIECE_UNIT.md). Nó unit_no: `legoloop run <dag> --allow-unit-no` (alias legado `--allow-unit-tasks`).
* **Demanda Mega Plano Produção (16-17/11/2026):** [INTENT](docs/intents/2026-11-16_INTENT_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md) · [SPEC](docs/specs/2026-11-16_SPEC_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md) · [PLAN](docs/plans/2026-11-16_PLAN_MEGA_PLANO_PRODUCAO.md) · [ROUTE](docs/routes/2026-11-16_ROUTE_MEGA_PLANO_PRODUCAO.json) · [Relatório final](docs/reports/2026-11-17_REPORT_MEGA_PLANO_PRODUCAO.md)
* **Demandas anteriores:** [docs/DEMANDAS_ANTERIORES.md](docs/DEMANDAS_ANTERIORES.md)
* **Diário de evolução e retomada de sessão:** [MEMORY.md](MEMORY.md)

* **Processos Canonizados:** [docs/PROCESSOS_CANONIZADOS.md](docs/PROCESSOS_CANONIZADOS.md) (registro gravado pelo `legoloop fechar`; não acrescente entradas aqui, o agy trunca este arquivo acima de 24000 bytes)
