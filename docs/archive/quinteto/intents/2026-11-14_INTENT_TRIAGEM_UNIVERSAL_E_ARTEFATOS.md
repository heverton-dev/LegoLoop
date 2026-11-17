# INTENT — Peça de Triagem Universal e Artefatos do Fabricar

- **Data:** 2026-11-14
- **Origem:** Pausa Estratégica da Auditoria do Mega-Pipeline (`docs/reports/2026-11-14_AUDITORIA_MEGA_PIPELINE_VIVO.md`)

## 1. Problema
Pedidos em linguagem natural ou fluxos existentes não possuem um roteador único de entrada antes de acionar LLMs ou nós de execução. Além disso, o `legoloop fabricar` gera `template` (poluindo a aplicação) e não gera automaticamente a Peça de Triagem, o Manual, o Livro e as Fichas Técnicas dos nano-processos.

## 2. Objetivo
1. Projetar e implementar a **Peça de Triagem / Roteador Universal de Entrada** com decisão híbrida (determinística por regras/artefatos + probabilística apenas na ambiguidade).
2. Operar a triagem nos 3 pontos de entrada: **CLI** (`legoloop triagem`), **Chat/Skill** e **Studio**.
3. Atualizar `legoloop fabricar` para remover `template` e gerar **13 artefatos** (`skill`, `loop`, `dag`, `processo`, `intent`, `spec`, `plan`, `route`, `dag_iso`, `triagem`, `manual`, `livro`, `fichas_tecnicas`).
4. Adequar e sanear todos os processos existentes em `loops/` (atualmente 10 diretórios) removendo legados inúteis ou quebrados e gerando suas Peças de Triagem e artefatos faltantes.

## 3. Critério de Sucesso
- Todo processo criado via `legoloop fabricar` sai pronto com os 13 artefatos validados.
- `legoloop triagem "<entrada>"` resolve deterministicamente comandos/slugs/estado em disco e só chama LLM quando há ambiguidade.
- Studio e Chat/Skill consomem o mesmo motor de triagem.
- Todos os processos em `loops/` estão saneados e possuem Peça de Triagem válida.
