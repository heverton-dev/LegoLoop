# Manual Didático: legoloop-kraften-deploy

- **Data:** 2026-11-26
- **Versão:** 1.0.0
- **Processo:** `legoloop-kraften-deploy`
- **Natureza:** determinístico (zero LLM)
- **Motor:** `internal/kraften` (Go nativo) + CLI `legoloop kraften-deploy`

## Sumário

1. **Macro — Fundação:** por que deploy determinístico e qual a fronteira universal do `.zip`.
2. **Meso — Arquitetura:** os 5 nós, os artefatos intermediários e os contratos `DeployConfig`/`DeployReceipt`.
3. **Micro — Passo a passo:** cada subcomando, seus gates e as provas binárias.
