# 📋 SPEC: legoloop-web-qa

> **Status:** Aprovada após Entrevista Socrática (97% de confiança)  
> **Versão:** 1.1.0  
> **Template:** TEMPLATE_SPEC.md v2.2.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Natureza do processo: software (`web-qa`).
2. Objetivo: executar auditoria determinística, sem mocks ou stubs, em aplicações web locais (comando + porta) ou URLs publicadas utilizando o arsenal completo de 8 ferramentas obrigatórias: `playwright`, `obscura`, `aidd-linter-design`, `renderlint`, `axe-core`, `unlighthouse`, `lychee` e `lost-pixel`.
3. Viewports obrigatórios: `desktop (1440x900)` e `mobile (390x844)`, com suporte opcional a `storageState` do Playwright para rotas autenticadas.
4. O relatório consolidado (`RELATORIO_QA_WEB.json` + `RELATORIO_QA_WEB.md`) é gravado sob `docs/qa-web/YYYY-MM-DD_<alvo>/` e funciona como insumo obrigatório para a entrevista socrática de correção e geração da cadeia `INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM` para ciclo `BUILD`.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* **Sempre:** validar o grafo com `legoloop validate loops/legoloop-web-qa/pipeline.dag.json` antes de executar.
* **Sempre:** rodar o preflight real (`node tools/webqa/webqa.mjs preflight --json`) e bloquear imediatamente se faltar qualquer uma das 8 peças.
* **Nunca:** simular resultados de ferramenta ausente (Zero Stubs/Mocks) ou editar código-fonte automaticamente sem passar pela entrevista socrática e cadeia `INTERVIEW..BUILD`.
* **Nunca:** promover baseline inicial do `lost-pixel` sem aprovação explícita no Stop-Gate humano `n09_sg_baseline_visual`.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
* Runner determinístico em `tools/webqa/webqa.mjs`, coberto por suíte unitária TDD em `tools/webqa/webqa.test.mjs` (`node --test`).
* Path Jail estrito impede gravação fora de `docs/qa-web/YYYY-MM-DD_<alvo>/`.
* Execução em 8 camadas topológicas com paralelismo na camada 3 (`n03_links_lychee`, `n04_obscura_e_playwright`, `n05_renderlint_e_design`, `n06_axe_acessibilidade`, `n07_unlighthouse_rotas`).

## 4. Contratos de Dados & Schemas Formais
| Nó | Tipo | Ferramenta / Peça | Produz | Critério de Reprovação |
| :--- | :--- | :--- | :--- | :--- |
| `n01_preflight` | `script` | `webqa.mjs preflight` | `scratch/webqa/preflight.json` | Qualquer uma das 8 peças ausente |
| `n02_boot_e_descoberta` | `script` | `webqa.mjs discover` | `scratch/webqa/routes_manifest.json` | App inacessível ou falha de crawl |
| `n03_links_lychee` | `script` | `lychee` | `scratch/webqa/lychee.json` | `brokenLinks > 0` |
| `n04_obscura_e_playwright` | `script` | `obscura` + `playwright` | `scratch/webqa/playwright_obscura.json` | `consoleErrors > 0` ou `routeErrors > 0` |
| `n05_renderlint_e_design` | `skill` | `renderlint` + `aidd-linter-design` | `scratch/webqa/renderlint_design.json` | Botão morto, sobreposição ou violação de token |
| `n06_axe_acessibilidade` | `script` | `axe-core` | `scratch/webqa/axe.json` | Violação `critical > 0` ou `serious > 0` |
| `n07_unlighthouse_rotas` | `script` | `unlighthouse` | `scratch/webqa/unlighthouse.json` | Nota mínima `< 90` |
| `n08_lost_pixel_visual` | `script` | `lost-pixel` | `scratch/webqa/lost_pixel.json` | `diffCount > 0` ou ausência de baseline |
| `n09_sg_baseline_visual` | `stop_gate` | Operador Humano | `scratch/webqa/baseline_aprovada.json` | Rejeição humana do snapshot inicial |
| `n10_relatorio_consolidado` | `script` | `webqa.mjs run` | `docs/qa-web/RELATORIO_QA_WEB.{json,md}` | Falha ao gravar ou violação de Path Jail |
| `n11_entrevista_correcao` | `skill` | `aidd-grill` | `scratch/webqa/INTERVIEW_CORRECAO.md` | Entrevista sem decisões por achado |
| `n12_governanca_build` | `skill` | `legoloop-fabricar` | `docs/qa-web/PACOTE_BUILD_SWARM.json` | Cadeia `INTERVIEW..SWARM` incompleta |

## 5. Critérios Binários de Aceite por Nó / Ticket (EARS / BDD)

### Épico 01 — legoloop-web-qa
* [ ] `AC-n01_preflight`: **Quando** `n01_preflight` executar, **então** verifica as 8 ferramentas obrigatórias e retorna `exit_code == 0` somente se todas estiverem presentes.
* [ ] `AC-n02_boot_e_descoberta`: **Quando** URL ou comando+porta forem informados, **então** descobre rotas via crawl/sitemap + lista opcional em viewports `1440x900` e `390x844`.
* [ ] `AC-n03_links_lychee`: **Quando** `lychee` varrer as rotas, **então** falha com `exit_code != 0` se houver qualquer link ou asset quebrado.
* [ ] `AC-n04_obscura_e_playwright`: **Quando** `obscura` e `playwright` visitarem as rotas, **então** reprovam em erro de console, exceção não tratada ou HTTP >= 400.
* [ ] `AC-n05_renderlint_e_design`: **Quando** `renderlint` e `aidd-linter-design` rodarem, **então** reprovam em `dead-control`, `obscured-content`, `text-overflow` ou quebra de tokens.
* [ ] `AC-n06_axe_acessibilidade`: **Quando** `axe-core` auditar WCAG, **então** reprova com `exit_code != 0` se `critical > 0` ou `serious > 0`.
* [ ] `AC-n07_unlighthouse_rotas`: **Quando** `unlighthouse` auditar todas as rotas, **então** reprova se qualquer categoria ficar abaixo de `90`.
* [ ] `AC-n08_lost_pixel_visual`: **Quando** `lost-pixel` rodar sem baseline prévia, **então** gera snapshot inicial e aciona `n09_sg_baseline_visual`; com baseline, reprova em `diffCount > 0`.
* [ ] `AC-n09_sg_baseline_visual`: **Quando** houver snapshot inicial, **então** bloqueia promoção até confirmação humana explícita.
* [ ] `AC-n10_relatorio_consolidado`: **Quando** a auditoria terminar, **então** grava `RELATORIO_QA_WEB.json` e `RELATORIO_QA_WEB.md` em `docs/qa-web/YYYY-MM-DD_<alvo>/`.
* [ ] `AC-n11_entrevista_correcao`: **Quando** o relatório existir, **então** conduz entrevista socrática de afunilamento sobre como corrigir cada achado.
* [ ] `AC-n12_governanca_build`: **Quando** a entrevista concluir, **então** gera `INTERVIEW`, `INTENT`, `SPEC`, `PLAN`, `ROUTE`, `DAG` e `SWARM` para ciclo `BUILD`.
