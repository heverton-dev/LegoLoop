# 🔬 RELATÓRIO DE PROVAS POR MUTAÇÃO DELIBERADA — ÉPICO 04 (`TICK-04.3`)

> **Data:** 2026-11-06  
> **Status:** PASSED (100% das mutações capturadas pelos testes)  
> **Escopo:** Execução de 5 mutações deliberadas nos pontos críticos de governança do LegoLoop (Fusão, gates de qualidade, defasagem, rollback e cérebro plugável) para validar a eficácia da suíte de testes (Lei #8 do `AGENTS.md`).  
> **Spec Ref:** `docs/specs/2026-11-04_SPEC_FABRICAR_V2.md#AC-E04.3`  
> **Plan Ref:** `docs/plans/2026-11-04_PLAN_FABRICAR_V2.md#TICK-04.3`

---

## 1. Resumo Executivo das Cinco Mutações

| ID | Ponto Crítico | Arquivo Alvo | Mutação Deliberada | Teste que Capturou | Saída da Falha (Erro Real) | Resultado |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **M1** | Fusão e Aptidão | `internal/matriz/rank.go` | Usar `TotalScore` antes de `Aptidao` no desempate de peças | `TestRankAptidaoEhFracaoDosCasosPedidos` | `rank_test.go:128: ordem por aptidão; aptidão 0 nunca é escolhida: [um dois]` | `MUTADO & MORTO` (Verde restaurado) |
| **M2** | Gate Trivial | `internal/gate/calls.go` | Adicionar `"file_exists": true` em `contentCalls` (aceitar existence isolado) | `TestGateTrivialChecksContent` | `calls_test.go:38: ChecksContent("file_exists('a.md')") = true, queria false` | `MUTADO & MORTO` (Verde restaurado) |
| **M3** | Detecção de Defasagem | `internal/matriz/stale.go` | Comentar `s.Alteradas = append(s.Alteradas, r)` (ignorar peça alterada) | `TestStaleness_NovaAlteradaERemovida` | `stale_test.go:42: stale = {Alteradas:[]}` (esperava peça alterada) | `MUTADO & MORTO` (Verde restaurado) |
| **M4** | Rollback Atômico | `internal/processo/vinculo.go` | Comentar `desfazer(criados)` e `desfazer(...)` ao falhar link | `TestVinculo/falha_no_2o_link_desfaz_o_1o` | `vinculo_test.go:103: o vínculo ... devia ter sido desfeito` | `MUTADO & MORTO` (Verde restaurado) |
| **M5** | Cérebro Plugável | `internal/brain/pick.go` | Forçar escolha automática sempre do 9Router (`SourceRouter/code-fast`) | `TestPick_EscolhaAutomaticaPorInventario/só_Ollama` | `brain_test.go:43: Name() = "9router/code-fast", quer "ollama/qwen2.5"` | `MUTADO & MORTO` (Verde restaurado) |

---

## 2. Detalhamento e Evidências das Falhas

### Mutação M1 — Score no Lugar da Aptidão (Fusão)
* **Hipótese:** Se a ordenação priorizar o Score reputacional antes da aptidão de casos de uso, peças de alto score com aptidão zero ou inferior ultrapassariam peças 100% adequadas à esteira.
* **Mutação:** Alteração na função `comparar(a, b Ranked)` em `rank.go` colocando a comparação de `Rating.TotalScore` antes de `Aptidao`.
* **Evidência:** `TestRankAptidaoEhFracaoDosCasosPedidos` falhou com exit 1.

### Mutação M2 — Aceitar `file_exists` Isolado (Gate Trivial)
* **Hipótese:** Gates devem verificar conteúdo ou aprovação humana (`exit_code`, `word_count`, `has_binary_criteria`, etc.). `file_exists` isolado cria falsa sensação de sucesso (stub de arquivo vazio).
* **Mutação:** Registro de `"file_exists": true` no mapa `contentCalls` de `calls.go`.
* **Evidência:** `TestGateTrivialChecksContent` falhou em 7 casos de teste com exit 1.

### Mutação M3 — Ignorar Peça Alterada (Defasagem)
* **Hipótese:** Se um script ou skill tiver seu código modificado no disco mas seu caminho mantido, a matriz precisa detectar como alterada para invalidação de cache incremental.
* **Mutação:** Remoção da adição de peças à lista `s.Alteradas` em `stale.go`.
* **Evidência:** 5 testes em `internal/matriz` falharam (`TestStaleness_NovaAlteradaERemovida`, `TestIncremental_PecasDoMesmoArquivoDeConfiguracaoTemRefProprio`, etc.).

### Mutação M4 — Não Desfazer Links no Rollback
* **Hipótese:** Se o processo falha ao vincular no segundo harness, o link criado no primeiro harness deve ser desfeito atomicamente para não poluir o ambiente do desenvolvedor.
* **Mutação:** Comentadas as chamadas a `desfazer(...)` em `vinculo.go`.
* **Evidência:** `TestVinculo/falha_no_2o_link_desfaz_o_1o` falhou apontando o link remanescente não desfeito.

### Mutação M5 — Escolher Sempre o 9Router
* **Hipótese:** O LegoLoop deve ser desacoplado de qualquer provedor proprietário: máquinas apenas com Ollama, Claude CLI ou harness agy devem funcionar perfeitamente de ponta a ponta.
* **Mutação:** Hardcode forçando `provider.SourceRouter` e `code-fast` na seleção automática de `pick.go`.
* **Evidência:** `TestPick_EscolhaAutomaticaPorInventario` e `TestPick_OllamaConversaPeloOrchestrix` falharam imediatamente.

---

## 3. Conclusão da Análise de Mutação

Todos os 5 pontos críticos do LegoLoop possuem cobertura de testes com alta densidade e capacidade de detecção de mutações. Zero testes silenciosos ou tolerantes a regressões de governança. O critério `AC-E04.3` está **100% atendido**.
