# 📊 REPORT FINAL: Hierarquia de Iniciativas no LegoLoop Mega Pipeline

> **Data:** 2026-11-07  
> **Demanda:** `docs/intents|specs|plans|routes/2026-11-07_*_HIERARQUIA_INICIATIVA`  
> **Estado:** 6 tickets `PASSED`, Stop-Gates `SG-E01`, `SG-E02` autorizados. Entrega e homologação final prontas.  
> **Protocolo:** STE-PT (Situação, Teoria/Tese, Evidência/Execução)  

---

## 1. Resumo Executivo & Diagnóstico

* **Situação Encontrada:** A matriz de decomposição do LegoLoop operava apenas com 3 níveis (Épico ➔ Micro-ticket ➔ Nano-ticket). Para demandas de grande porte corporativo ou multi-sistemas (ex.: mudanças conjuntas em Go, Studio Web, CLI e Infraestrutura), a ausência de um nível superior agregador gerava acúmulo artificial de épicos ou forçava a fragmentação de demandas sem uma visão de portfólio unificada.
* **Tese de Resolução:** Formalização do nível canônico **Iniciativa (`I01`, `I02`...)** acima de **Épico (`E01`, `E02`...)**, mantido estritamente **opcional e cumulativo**:
  1. **SSOT Íntegro:** Atualização coordenada de `TEMPLATE_INTENT.md`, `TEMPLATE_SPEC.md`, `TEMPLATE_PLAN.md` e `TEMPLATE_ROUTE.json` em `docs/templates/`.
  2. **Motores Go sem Quebra:** Evolução de `internal/plangen` (`meta.go`, `route.go`, `plangen.go`) com regexes resilientes e emissão de `initiative_id` e `initiative_ref`.
  3. **Retrocompatibilidade de 100%:** Planos e rotas sem iniciativa continuam sendo processados normalmente sem erros ou valores espúrios.
  4. **Normativa da Skill:** Atualização da tabela de 4 níveis e guia de classificação na skill `aidd-legoloop-mega-pipeline/SKILL.md`.

---

## 2. Entregas por Épico

### Épico 01 — Padronização de Templates & Schemas Canônicos (SSOT)
* `TICK-01.1`: Inclusão do metadado `Iniciativa` nos templates canônicos markdown e da chave `initiative_id` (isomórfica em todos os nós) e `initiative_ref` no `TEMPLATE_ROUTE.json`.
* `TICK-01.2`: Garantia de conformidade em `docs/templates/templates_test.go` (10/10 testes passando) e prova de mutação mecânica via `TestRouteKeysAreExposedForEngines`.

### Épico 02 — Suporte a Iniciativas no Montador de Rotas (`internal/plangen`)
* `TICK-02.1`: Suporte a `InitiativeRef`, `Initiatives` e mapeamento pai-filho em `internal/plangen/meta.go` e `assemble.go`.
* `TICK-02.2`: Emissão precisa de nós de rota JSON com `initiative_id` e `initiative_ref` em `internal/plangen/route.go`.
* Suíte de testes `internal/plangen` expandida para 20 testes com `initiative_test.go` cobrindo cenários com e sem iniciativa e entradas malformadas.

### Épico 03 — Governança da Skill & Recompilação Go
* `TICK-03.1`: Atualização da skill `aidd-legoloop-mega-pipeline/SKILL.md` com a matriz de 4 níveis (Iniciativa > 50k tokens, esforço high) e sincronização nos harnesses.
* `TICK-03.2`: Recompilação de `bin/legoloop.exe`, instalação global via `go install`, e validação de invariantes.

---

## 3. Homologação Mecânica Final

| Verificação | Comando | Resultado | Evidência |
| :--- | :--- | :--- | :--- |
| **Suíte Global de Testes** | `go test -count=1 ./...` | `PASS` | 100% verde em todos os 32 pacotes do projeto |
| **Suíte de Templates** | `go test -v ./docs/templates/...` | `PASS` | 10/10 testes aprovados |
| **Suíte Plangen** | `go test -v ./internal/plangen/...` | `PASS` | 20/20 testes aprovados |
| **Prova de Mutação** | Invalidação intencional da chave no template | `PASS` | Captura imediata com falha controlada |
| **Análise Estática (Vet)** | `go vet ./internal/plangen/... ./docs/templates/...` | `PASS` | 0 avisos |
| **Formatação (gofmt)** | `gofmt -l internal/plangen docs/templates` | `PASS` | Saída vazia (100% formatado) |
| **Limite de Linhas (funclen)** | `go run ./tools/funclen internal/plangen` | `PASS` | Todas as funções $\le 35$ linhas |
| **Integridade de Dependências** | `git diff --exit-code -- go.mod go.sum` | `PASS` | Inalterado (apenas Go stdlib) |
| **Compilação e Instalação** | `go build -o bin/legoloop.exe ./cmd/legoloop` | `PASS` | Binário atualizado no repositório e no PATH |
