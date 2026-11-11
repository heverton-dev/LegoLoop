# TICKET: passo lauren-unslop-ship antes do gate dos nós de normalização (TICK-03.3)

Data: 2026-11-10. Origem: sessão orquestradora da trilha Lei #15. Destino: sessão de engenharia do motor/gerador.

## Fato
- O motor aceita `context_out` de diretório (commit 3fb92ee).
- O TICK-03.3 rodou 2 vezes (3 tentativas cada). Custo acumulado: US$ 2,14.
- Rodada 2: restaram 2 violações e `gofmt`. Rodada 3: restou 1 violação.
- A violação restante é um comentário narrativo: `internal/studiodoc/types.go:41`
  (`// Summary é o resumo do loop retornado na listagem do acervo.`).
- Ao falhar o gate, o nó reverte as 149 edições. O trabalho se perde.
- O sub-loop `lauren-unslop-ship` (loops/lauren-tan.loop.json) remove comentários
  narrativos (`unslop`, `no-comments`). A regra remove o comentário; não renomeia o termo.

## Tarefa
1. Anexe um passo de higiene (unslop + no-comments + gofmt) antes do gate dos nós
   de normalização taxonômica. Decida o dono: `tools/route2dag` (passo no DAG) ou
   o prompt do nó em `internal/orchestrix`.
2. Mantenha o gate como juiz final. O passo de higiene não substitui o gate.
3. Avalie manter as edições parciais após falha do gate (hoje o nó reverte tudo).
   Registre a decisão e o motivo.
4. Teste com o DAG de um nó só (`TICK-03.3`) e confirme a violação em types.go:41 resolvida.
5. Commite só os seus arquivos. Sem push.

## Restrições
- O claude do nó não tem shell. Use Edit/Write, não renome de pacote.
- Limites de custo: micro com effort low e teto US$ 0,80; macro com medium e US$ 2,50.

## Retomada da trilha (orquestrador)
Rode `medicao.json` do TICK-03.3. Se passar, rode `derivado.json` até o SG-E03.
