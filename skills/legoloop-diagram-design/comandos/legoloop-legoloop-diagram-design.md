---
name: legoloop-legoloop-diagram-design
description: Executa o Processo legoloop-diagram-design pelo comando legoloop exec.
---

Execute o Processo legoloop-diagram-design com o comando:

    legoloop exec legoloop-diagram-design

Regras obrigatórias de execução:
1. Executar estritamente o DAG determinístico dos 6 nós canônicos:
   - e01 (inicializar-slug-e-anticorpos): cria a estrutura isolada em `docs/diagrams/<slug>/{src,assets,export,animation,anticorpos}/` e sincroniza o `anticorpos.json`.
   - e02 (ingestao-e-quadrantes): seleciona 1 dos 50 tipos visuais do catálogo e valida tokens editoriais (sem serifa, sem emojis).
   - e03 (geracao-svg-html): gera `src/<slug>.html` auto-contido com grid ortogonal, zero transbordamento e máscaras opacas.
   - e04 (micro-animacao): aplica animação CSS nativa contínua (@keyframes flow-dash) ou Remotion condicional.
   - e05 (gate-e-exportacao): executa validação mecânica e exportação com exit code 0.
   - e06 (fechamento-e-anticorpos): consolida achados em `anticorpos.json` e emite relatório em `docs/reports/`.
2. Proibição absoluta: NUNCA crie ou gere arquivos de diagramas soltos na raiz de `docs/diagrams/` ou fora da pasta do slug (`docs/diagrams/<slug>/`).
