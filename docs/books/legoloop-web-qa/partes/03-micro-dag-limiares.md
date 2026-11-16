# Parte III — Anatomia do DAG, Viewports e Limiares

## Topologia em 8 Camadas

O grafo `loops/legoloop-web-qa/pipeline.dag.json` (`schema_version: 2`) organiza 12 nós em 8 camadas topológicas, garantindo isolamento estrito (`orchestrix_hermetic`) e `write_boundary: strict`:

1. **Camada 1 (`n01_preflight`):** Verifica presença real das 8 peças no PATH, `node_modules` local e cofre de skills.
2. **Camada 2 (`n02_boot_e_descoberta`):** Sobe o servidor local se `--cmd` e `--port` forem fornecidos e descobre rotas via crawl HTML + `/sitemap.xml` + lista explícita `--routes`.
3. **Camada 3 (Execução Paralela):** Dispara simultaneamente cinco auditores independentes: `n03_links_lychee`, `n04_obscura_e_playwright`, `n05_renderlint_e_design`, `n06_axe_acessibilidade` e `n07_unlighthouse_rotas`.
4. **Camada 4 (`n08_lost_pixel_visual`):** Compara capturas nos viewports `desktop (1440x900)` e `mobile (390x844)` após `n04_obscura_e_playwright`.
5. **Camada 5 (`n09_sg_baseline_visual`):** Stop-Gate humano obrigatório sempre que `.lostpixel/baseline` ainda não existir.
6. **Camada 6 (`n10_relatorio_consolidado`):** Consolida métricas e achados em `docs/qa-web/YYYY-MM-DD_<alvo>/`.
7. **Camada 7 (`n11_entrevista_correcao`):** Conduz entrevista socrática de afunilamento lendo o relatório gerado.
8. **Camada 8 (`n12_governanca_build`):** Materializa o pacote `INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM` para o ciclo `BUILD`.

## Matriz de Limiares e Códigos de Saída

A função `evaluateThresholds` em `tools/webqa/webqa.mjs` aplica critérios binários sem margem subjetiva:

- **`exit_code == 0` (`PASSED`):** Zero links quebrados, zero erros de console/rota, zero violações `critical`/`serious` no `axe-core`, zero violações no `renderlint`/`aidd-linter-design`/`obscura`, nota mínima >= 90 no `unlighthouse`, baseline existente e `diffCount == 0` no `lost-pixel`.
- **`exit_code == 2` (`NEEDS_BASELINE_STOP_GATE`):** Todas as métricas funcionais passaram, mas não há baseline prévia do `lost-pixel`. O processo para e solicita confirmação humana no nó `n09_sg_baseline_visual`.
- **`exit_code == 1` (`FAILED`):** Qualquer violação de limiar acima ou ferramenta ausente no preflight.
