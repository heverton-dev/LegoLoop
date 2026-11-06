#!/bin/bash
# Teste interativo do agy no Orca: um modelo por vez, lendo a tela. Resultado em agy_inter.jsonl
SP="C:/Users/trcnologia/AppData/Local/Temp/claude/C--Users-trcnologia-Desktop-LegoLoop/bd04864d-dcbf-46a1-97c0-5666d6259c92/scratchpad/aud"
WT="id:29725318-bbe5-48f8-8556-d17d803da6f1::C:/Users/trcnologia/orca/workspaces/LegoLoop/e03-aud-agy"
OUT="$SP/agy_inter.jsonl"; : > "$OUT"
PROMPT="Quanto e 17 vezes 23? Responda apenas o numero."

caso() {
  modelo="$1"; effort="$2"
  inicio=$(date +%s)
  args=(--agent antigravity --worktree "$WT" --model "$modelo" --prompt "$PROMPT" --json)
  [ -n "$effort" ] && args+=(--effort "$effort")
  saida=$(orca agent launch "${args[@]}" 2>&1)
  h=$(echo "$saida" | node -e 'let s="";process.stdin.on("data",d=>s+=d).on("end",()=>{try{const j=JSON.parse(s);const o=(j.result.launch&&j.result.launch.outcome)||j.result.outcome||{};console.log(o.handle||"")}catch(e){console.log("")}})')
  tela=""
  if [ -n "$h" ]; then
    for i in 1 2 3 4 5 6; do
      sleep 10
      tela=$(orca terminal read --terminal "$h" --screen --json 2>&1 | node -e 'let s="";process.stdin.on("data",d=>s+=d).on("end",()=>{try{const j=JSON.parse(s);console.log((j.result.terminal.tail||[]).join(" // ").replace(/https:\/\/[^" ]+/g,"<link>").replace(/[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+/g,"<email>"))}catch(e){console.log("")}})')
      echo "$tela" | grep -qE "391|Eligibility|quota|ERROR|invalid model" && break
    done
    orca terminal close --terminal "$h" --json > /dev/null 2>&1
  fi
  node -e '
    const t=process.argv[4]||"";const ok=/391/.test(t.replace(/Quanto e 17 vezes 23\? Responda apenas o numero\./g,""));
    const el=/Eligibility/.test(t), qt=/quota/i.test(t), inv=/invalid model/.test(t);
    console.log(JSON.stringify({modelo:process.argv[1],effort:process.argv[2]||"",segundos:+process.argv[3],correto:ok,elegibilidade:el,cota:qt,invalido:inv,tela:t.replace(/[─▄▀█▌▐]+/g,"").replace(/\s+/g," ").slice(-420)}))' "$modelo" "$effort" "$(( $(date +%s)-inicio ))" "$tela" >> "$OUT"
  # encerra só os agy.exe criados depois do início deste caso (são meus)
  powershell -NoProfile -Command "\$t=[DateTimeOffset]::FromUnixTimeSeconds($inicio).LocalDateTime; Get-CimInstance Win32_Process | Where-Object { \$_.Name -eq 'agy.exe' -and \$_.CreationDate -ge \$t } | ForEach-Object { Stop-Process -Id \$_.ProcessId -Force -ErrorAction SilentlyContinue }" > /dev/null 2>&1
}

caso gemini-3.8-flash-low ""
caso gemini-3.8-flash-medium ""
caso gemini-3.8-flash-high ""
caso gemini-3.8-flash high
echo FIM > "$SP/agy_inter_fim.txt"
