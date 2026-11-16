#!/bin/bash
# Corpo de worker de LLM para `legoloop swarm run`: cada no roda `claude -p` via 9Router, somente leitura.
# Ambiente: OUT (pasta dos artefatos), NINEROUTER_URL e NINEROUTER_KEY (so por nome), MODEL (padrao t-agente-128k).
# LEGOLOOP_SWARM_WORKER (a, b ou c) escolhe o prompt. O cwd e o worktree do worker.
set -uo pipefail
mkdir -p "$OUT"
NAME="$LEGOLOOP_SWARM_WORKER"
MODEL="${MODEL:-t-agente-128k}"
export ANTHROPIC_BASE_URL="$NINEROUTER_URL" ANTHROPIC_AUTH_TOKEN="$NINEROUTER_KEY"

FORMATO='Responda somente em Markdown, em portugues. Cada risco comeca com "Risco N:" e traz: arquivo:linha de prova (confirme com Read, nao invente linhas), cenario de falha e correcao sugerida.'
case "$NAME" in
a)
  PROMPT="Voce e auditor de qualidade. Leia internal/board/swarm_run.go e internal/board/swarm_runtime.go neste repositorio (somente leitura). Liste de 3 a 5 riscos reais de RunLevels: corrida, vazamento de worker, tempo maximo, codigo de saida. $FORMATO"
  ;;
b)
  PROMPT="Voce e engenheiro de testes. Leia internal/cli/swarm_run.go, internal/board/swarm_dag.go e a funcao swarmRunDAG em internal/mcp/swarm_tools.go neste repositorio (somente leitura). Liste de 3 a 5 riscos reais de entrada e de borda do swarm run: DAG invalido, no sem papel, nome de worker, comando ausente, timeout. $FORMATO"
  ;;
c)
  PROMPT="Voce e engenheiro backend. Consolide em uma unica lista ordenada por gravidade os riscos abaixo, descartando duplicatas e mantendo arquivo:linha. Confirme com Read que cada linha citada existe; remova o que nao se confirmar. Termine com uma secao 'Recomendacao' com um unico passo no imperativo. Responda somente em Markdown, em portugues.

=== Riscos do auditor a ===
$(cat "$OUT/a_riscos.md")

=== Riscos do auditor b ===
$(cat "$OUT/b_riscos.md")"
  ;;
*) echo "worker desconhecido: $NAME" >&2; exit 9 ;;
esac

case "$NAME" in a) ART=a_riscos.md ;; b) ART=b_riscos.md ;; c) ART=c_consolidado.md ;; esac
date +%s%N > "$OUT/$NAME.start"
for tentativa in 1 2; do
  timeout 220 claude -p "$PROMPT" --model "$MODEL" --no-session-persistence \
    --allowedTools "Read Grep Glob" --disallowedTools "Bash Edit Write NotebookEdit WebFetch WebSearch" \
    --strict-mcp-config --mcp-config '{"mcpServers":{}}' </dev/null > "$OUT/$ART.tmp" 2> "$OUT/$NAME.stderr"
  if [ $? -eq 0 ] && [ -s "$OUT/$ART.tmp" ]; then
    mv "$OUT/$ART.tmp" "$OUT/$ART"
    date +%s%N > "$OUT/$NAME.end"
    echo "$NAME tentativa=$tentativa modelo=$MODEL" > "$OUT/$NAME.meta"
    exit 0
  fi
  sleep 3
done
date +%s%N > "$OUT/$NAME.end"
exit 1
