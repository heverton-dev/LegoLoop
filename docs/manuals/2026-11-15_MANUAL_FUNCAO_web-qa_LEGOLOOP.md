# Manual Operacional: `legoloop-web-qa` (v1.1.0)

## 1. Propósito e Visão Arquitetural

O processo `legoloop-web-qa` executa auditoria determinística ponta a ponta de aplicações web em diretório local (`--project <dir>` com `--port` opcional), por comando explícito (`--cmd` + `--port`) ou publicadas (`--url`), combinando oito ferramentas reais especializadas e convertendo os achados em correções governadas por contrato (`INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM -> BUILD`).

Com `--project <dir>` (por exemplo, `--project .` para auditar o próprio LegoLoop Studio), `webqa.mjs` valida Path Jail/diretório existente e detecta automaticamente:
- **LegoLoop (`go.mod` + `cmd/legoloop` ou `bin/legoloop.exe`):** sobe o Studio em `127.0.0.1:<port>` (porta padrão `7878` ou livre).
- **Node (`package.json` com `scripts.dev` ou `scripts.start`):** inicia `npm run dev` ou `npm start` na porta indicada (padrão `4173` ou livre).
- **Estático (`index.html`):** sobe servidor HTTP Node em `127.0.0.1:<port>` com Path Jail (padrão `4173` ou livre).

## 2. O Arsenal de 8 Ferramentas Obrigatórias

| Ferramenta | Camada | Papel Técnico | Limiar de Reprovação |
| :--- | :--- | :--- | :--- |
| `playwright` | Browser / Runtime | Navegação real, `storageState`, erros de console e HTTP >= 400 | `consoleErrors > 0` ou `routeErrors > 0` |
| `obscura` | Browser Headless Rust | Fetch e renderização rápida de `localhost` com `--allow-private-network` | `errors > 0` |
| `aidd-linter-design` | Design Tokens | Consistência visual e aderência a tokens do design system | `violations > 0` |
| `renderlint` | Comportamento UI | Cliques mortos (`dead-control`), sobreposição (`obscured-content`), overflow | `violations > 0` |
| `axe-core` | Acessibilidade | Auditoria WCAG automatizada (contraste, ARIA, teclado) | `critical > 0` ou `serious > 0` |
| `unlighthouse` | Qualidade Multi-Rota | Varredura Lighthouse paralela em todas as rotas descobertas | `minScore < 90` |
| `lychee` | Integridade | Verificação de links, âncoras e assets quebrados | `brokenLinks > 0` |
| `lost-pixel` | Regressão Visual | Diff de screenshots em `1440x900` e `390x844` | `diffCount > 0` ou Stop-Gate sem baseline |

## 3. Fluxo de Dados e Camadas do DAG (`pipeline.dag.json`)

1. **Camada 1 (`n01_preflight`):** Bloqueia imediatamente com `exit_code == 1` caso qualquer das 8 ferramentas falte no ambiente.
2. **Camada 2 (`n02_boot_e_descoberta`):** Sobe o servidor local (opcional) e consolida rotas via crawl HTML + `sitemap.xml` + `--routes`.
3. **Camada 3 (Paralela):** Executa `n03_links_lychee`, `n04_obscura_e_playwright`, `n05_renderlint_e_design`, `n06_axe_acessibilidade` e `n07_unlighthouse_rotas`.
4. **Camada 4 (`n08_lost_pixel_visual`):** Compara screenshots com a baseline ou cria o snapshot inicial.
5. **Camada 5 (`n09_sg_baseline_visual`):** Stop-Gate humano obrigatório quando não há baseline prévia.
6. **Camada 6 (`n10_relatorio_consolidado`):** Grava `RELATORIO_QA_WEB.json` e `RELATORIO_QA_WEB.md` em `docs/qa-web/YYYY-MM-DD_<alvo>/`.
7. **Camada 7 (`n11_entrevista_correcao`):** Inicia entrevista socrática de afunilamento lendo o relatório gerado.
8. **Camada 8 (`n12_governanca_build`):** Emite os artefatos `INTERVIEW`, `INTENT`, `SPEC`, `PLAN`, `ROUTE`, `DAG` e `SWARM` para execução em `BUILD`.

## 4. Comandos de Operação e Verificação

```bash
legoloop validate processos/legoloop-web-qa/pipeline.dag.json
node --test tools/webqa/webqa.test.mjs
node tools/webqa/webqa.mjs preflight --json
node tools/webqa/webqa.mjs run --url http://127.0.0.1:3000 --alvo meu-app
node tools/webqa/webqa.mjs run --project . --alvo legoloop-studio
node tools/webqa/webqa.mjs run --project ./meu-app --port 4173 --alvo meu-app
```
