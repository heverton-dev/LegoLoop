# Parte V — Guardrails Negativos, FMEA e Apêndices

## Regras Negativas Invioláveis

1. **Proibição de Fallback Falso:** Se `renderlint`, `lychee`, `axe`, `unlighthouse` ou `lost-pixel` não existirem no ambiente, o runner interrompe no `n01_preflight`. É proibido retornar `exit_code == 0` simulando ferramentas ausentes.
2. **Proibição de Edição Cega:** Nenhum agente pode alterar componentes React/HTML diretamente após ler `RELATORIO_QA_WEB.md` sem concluir `n11_entrevista_correcao` e `n12_governanca_build`.
3. **Path Jail em `docs/qa-web/`:** Slugs contendo `..`, `/` ou `\\` são rejeitados antes de criar diretórios em disco.
4. **Proteção da Baseline Visual:** O snapshot inicial do `lost-pixel` nunca se autopromove a baseline oficial sem confirmação humana em `n09_sg_baseline_visual`.

## Matriz de Riscos (FMEA Simplificada)

| Nó                          | Modo de Falha                         | Efeito                      | Contenção Mecânica                      |
| :-------------------------- | :------------------------------------ | :-------------------------- | :-------------------------------------- |
| `n01_preflight`                | Ferramenta ausente no PATH            | Auditoria incompleta        | Abort imediato com comando exato de instalação |
| `n02_boot_e_descoberta`          | Porta ocupada ou servidor não responde | Travamento de pipeline      | Timeout de 30s no healthcheck e `kill()` no `finally` |
| `n08_lost_pixel_visual`          | Mudança legítima de layout confundida com bug | Falso positivo contínuo     | Stop-Gate humano `n09_sg_baseline_visual`    |
| `n10_relatorio_consolidado`     | Slug malicioso no `--alvo`              | Escrita fora do repositório | `sanitizeTargetSlug` + verificação `startsWith(qaRoot)` |

## Apêndice: Rastreabilidade dos Artefatos

- **Skill Canônica:** `skills/legoloop-web-qa/SKILL.md`
- **DAG Executivo (v2):** `processos/legoloop-web-qa/pipeline.dag.json` e `docs/dags/2026-11-15_DAG_WEB_QA.json`
- **Rota v2:** `docs/routes/2026-11-15_ROUTE_WEB_QA.json`
- **SPEC e PLAN:** `docs/specs/2026-11-15_SPEC_WEB_QA.md` e `docs/plans/2026-11-15_PLAN_WEB_QA.md`
- **Fichas Técnicas:** `docs/fichas/legoloop-web-qa/n01.md` a `n12.md`
- **Runner e Suíte TDD:** `tools/webqa/webqa.mjs` e `tools/webqa/webqa.test.mjs`
