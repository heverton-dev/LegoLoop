# PLAN-EP14: Refinamento Hiper-Profundo, Experiência N8N, Sidebar Ergonômica e Acervo

> **Status:** CONCLUÍDO COM 100% DE APROVAÇÃO  
> **Referência:** [SPEC_EP14_CANVAS_N8N_INVENTARIO.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/SPEC_EP14_CANVAS_N8N_INVENTARIO.md)

---

## Fases e Tickets

### Fase 1: Coleta Hiper-Profunda, 9Router VPS & Loops Kebab-case
- [x] **TICK-1401 (Micro):** Descoberta de .env no SO & 9Router VPS
  - *Arquivo:* `internal/provider/env.go`, `internal/provider/inventory.go`
  - *Teste:* `internal/provider/env_test.go`
  - *Resultado:* Lê `Desktop/9router/.env`, `~/.aidd/.env`, etc., e extrai `NINEROUTER_URL`/`KEY`; lista modelos da VPS com sucesso.
- [x] **TICK-1402 (Micro):** Varredura Hiper-Profunda Recursiva de Skills
  - *Arquivo:* `internal/assets/skills.go`
  - *Teste:* `internal/assets/skills_deep_test.go`
  - *Resultado:* Desce recursivamente em mono-repos de skills (`upstream/skills/**`) sem parar no primeiro `SKILL.md`; adiciona raízes do Gemini, Cline e Factory.
- [x] **TICK-1403 (Nano):** Sanitização Estrita de Loops em Kebab-case
  - *Arquivo:* `internal/loop/sanitize.go`, `internal/loop/store.go`, `internal/studio/loops_api.go`
  - *Teste:* `internal/loop/sanitize_test.go`
  - *Resultado:* `ToKebabCase(name)` normaliza espaços, acentos e caracteres especiais para lowercase-kebab; `store.Save` e API aplicam compulsoriamente.

### Fase 2: Backend de Tradução Não-Destrutiva e Detalhes
- [x] **TICK-1404 (Micro):** Cache e Endpoint de Tradução PT-BR sob Demanda
  - *Arquivo:* `internal/studio/translate.go`, `internal/studio/server.go`
  - *Teste:* `internal/studio/translate_test.go`
  - *Resultado:* Armazena tradução em `~/.aidd/cache/translations/{sha256}.json`; arquivo fonte nunca é alterado; validação anti-traversal.
- [x] **TICK-1405 (Micro):** Grafo Ontológico de Peças Compatíveis
  - *Arquivo:* `internal/assets/compat.go`, `internal/studio/assets_api.go`
  - *Teste:* `internal/assets/compat_test.go`
  - *Resultado:* Calcula peças compatíveis com base em `Inputs`/`Outputs` de skills, ferramentas de MCPs e gatilhos de hooks.

### Fase 3: Frontend do Inventário (Agrupamento, Paginação e Modal)
- [x] **TICK-1406 (Micro):** Agrupamento Semântico e Paginação no Inventário
  - *Arquivo:* `internal/studio/static/inventario.html`, `internal/studio/static/inventory.js`
  - *Resultado:* Deduplica visualmente por SHA-256 em 1 card com badges das origens; toggle "Agrupar cópias"; paginação (25/50/100).
- [x] **TICK-1407 (Micro):** Modal de Inspeção Rico com 4 Sub-abas
  - *Arquivo:* `internal/studio/static/inventario.html`, `internal/studio/static/inventory.js`
  - *Resultado:* 4 abas (Visão Geral, Peças Compatíveis, Ocorrências no SO, Código com botão Traduzir).

### Fase 4: Frontend do Canvas N8N e Sidebar Ergonômica
- [x] **TICK-1408 (Micro):** Refatoração Ergonômica da Sidebar Esquerda
  - *Arquivo:* `internal/studio/static/index.html`, `internal/studio/static/style.css`
  - *Resultado:* Largura 320px respirável, seções colapsáveis com setas, botão de ocultar sidebar, filtros organizados.
- [x] **TICK-1409 (Micro):** Canvas N8N (Pan Infinito 2D, Zoom Ctrl+Wheel, Ctrl +/- e Fit Ctrl+0)
  - *Arquivo:* `internal/studio/static/app.js`
  - *Resultado:* Pan infinito livre, zoom escalado centrado no cursor, Fit to screen.
- [x] **TICK-1410 (Micro):** Desconexão Rápida de Arestas e Encaixe Intermediário
  - *Arquivo:* `internal/studio/static/app.js`
  - *Resultado:* Botão 'x' em arestas no centro da curva, tecla `Del` desfaz conexão, suporte a soltar nó sobre aresta para inserir intermediariamente.

### Fase 5: Validação, Auditoria de Segurança e Build
- [x] **TICK-1411 (Micro):** Auditoria de Segurança (Secret Protection, Zero innerHTML & Anti-Traversal)
  - *Resultado:* `TestNoInnerHTMLInScripts` verde; `TestDetectAPIKeysNeverLeaksValues` verde; anti-traversal ativo.
- [x] **TICK-1412 (Nano):** Build Final e Verificação funclen <= 35 linhas
  - *Resultado:* `go run ./tools/funclen internal` 0 violações; `go build -o bin/legoloop.exe ./cmd/legoloop` compilado com sucesso.
