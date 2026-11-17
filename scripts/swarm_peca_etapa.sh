#!/bin/bash
# Corpo do worker da peca legoloop-swarm: cada no do DAG roda a sua etapa com comandos swarm reais.
# Ambiente herdado: BIN, PROJ, OUT, DEMO; LEGOLOOP_SWARM_WORKER diz qual etapa.
set -euo pipefail
mkdir -p "$OUT"
py() { python -X utf8 "$@"; }
status_json() { "$BIN" swarm status -dir "$DEMO"; }

case "$LEGOLOOP_SWARM_WORKER" in
TICK-01-1)
  {
    echo "# Niveis topologicos"
    "$BIN" validate "$PROJ/processos/legoloop-swarm/pipeline.dag.json" | awk '/camada/ {n=$2; gsub(":","",n); printf "nivel %d (camada %d): %s\n", n-1, n, $3}'
    echo "Validacao: codigo 0"
  } > "$OUT/01_niveis.md"
  ;;
TICK-01-2)
  mkdir -p "$DEMO" && cd "$DEMO"
  [ -d .git ] || { git init -q && git config user.email e2e@legoloop.local && git config user.name e2e && echo "# demo" > README.md && git add . && git commit -qm init; }
  {
    echo "# Workers registrados"
    for pair in "w1 engenheiro-backend" "w2 engenheiro-tdd"; do
      set -- $pair
      "$BIN" swarm add -dir "$DEMO" -repo "$DEMO" -name "$1" -role "$2" -workspace share > "$OUT/add_$1.json"
      py -c "import json,sys;w=json.load(open(sys.argv[1]));print('- %s papel=%s worktree=%s'%(w['name'],w['role'],w['worktree_path']))" "$OUT/add_$1.json"
    done
  } > "$OUT/02_workers.md"
  ;;
TICK-01-3)
  for n in w1 w2; do "$BIN" swarm spawn -dir "$DEMO" -name $n -- "$BIN" swarm wait -timeout 120s > "$OUT/spawn_$n.json"; done
  {
    echo "# Processos"
    for n in w1 w2; do py -c "import json,sys;s=json.load(open(sys.argv[1]));print('- %s pid=%s status=%s'%(s['name'],s['pid'],s['status']))" "$OUT/spawn_$n.json"; done
    status_json | py -c "
import json,sys
st=json.load(sys.stdin)['workers']
pids={w['pid'] for w in st}
assert len(st)==2 and len(pids)==2 and all(w['status']=='working' for w in st),'workers nao paralelos'
print('paralelo: %d pids distintos vivos'%len(pids))"
  } > "$OUT/03_spawn.md"
  ;;
TICK-01-4)
  "$BIN" swarm post -dir "$DEMO" -to @w1 -body "go" > /dev/null
  for i in $(seq 1 80); do
    s=$(status_json | py -c "import json,sys;print({w['name']:w['status'] for w in json.load(sys.stdin)['workers']}['w1'])")
    [ "$s" = "stopped" ] && break; sleep 0.25
  done
  {
    echo "# Coordenacao"
    status_json | py -c "
import json,sys
for w in json.load(sys.stdin)['workers']:
    print('- %s status=%s%s'%(w['name'],w['status'],' concluido' if w['status']=='stopped' else ''))"
    "$BIN" swarm monitor -dir "$DEMO" -timeout 60s | py -c "import json,sys;print('parados:',len(json.load(sys.stdin)['stalled']))"
  } > "$OUT/04_coordenacao.md"
  ;;
TICK-01-5)
  for n in w1 w2; do "$BIN" swarm kill -dir "$DEMO" -repo "$DEMO" -name $n > /dev/null; done
  cd "$DEMO"
  {
    echo "# Encerramento"
    status_json | py -c "
import json,os,sys
for w in json.load(sys.stdin)['workers']:
    assert w['status']=='killed' and w.get('pid',0)==0,w
    assert not os.path.exists(w['worktree_path']),w['worktree_path']
    print('- %s status=%s pid=0 worktree removido'%(w['name'],w['status']))"
    echo "git worktree list: $(git worktree list | wc -l) entrada(s)"
  } > "$OUT/05_encerramento.md"
  ;;
*) echo "etapa desconhecida: $LEGOLOOP_SWARM_WORKER" >&2; exit 9 ;;
esac
