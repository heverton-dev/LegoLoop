# 2. Meso — Arquitetura e Perfis

A execução organiza-se em cinco Etapas encadeadas no DAG `processos/legoloop-forjar/pipeline.dag.json`:
1. `e01` (`detectar-harnesses`): descobre os harnesses instalados no `HOME`.
2. `e02` (`simular-forja`): calcula todas as alterações sem escrever.
3. `e03` (`aplicar-forja`): grava backup, aplica blocos gerenciados ligados ao `AGENTS.md`, faz merge de `settings.json` e gera lançadores por sistema operacional.
4. `e04` (`verificar-orcamento`): audita os tetos definidos em `docs/governanca/orcamento_tokens.json`.
5. `e05` (`emitir-relatorio`): consolida o veredito determinístico.
