# 3. Micro — Passo a Passo Operacional

1. Simule o plano: `legoloop forjar --dry-run --perfil padrao --json`.
2. Revise as alterações pendentes no Stop-Gate.
3. Aplique com backup automático: `legoloop forjar --apply --perfil padrao --json`.
4. Verifique o orçamento de tokens: `legoloop forjar --verificar --orcamento docs/governanca/orcamento_tokens.json --json`.
5. Caso deseje reverter: `legoloop forjar --desfazer <AAAA-MM-DD>`.
