import sqlite3, json, re

db = r"C:\Users\trcnologia\AppData\Roaming\orca\agent-model-catalog.db"
con = sqlite3.connect(f"file:{db}?mode=ro", uri=True)
cur = con.cursor()
cur.execute("SELECT key, seq, value FROM rows ORDER BY seq DESC LIMIT 3")
for key, seq, value in cur.fetchall():
    print("KEY:", key, "SEQ:", seq)
    s = value
    for m in re.finditer(r"maxMessageTokens|messageTokens|maxOutputTokens|contextWindow|context_window", s):
        print("  FIELD HIT:", s[max(0, m.start()-150):m.start()+150])
    if "code-fast" in s:
        i = s.find("code-fast")
        print("  code-fast SNIP:", s[max(0,i-300):i+500])
    try:
        d = json.loads(s)
        models = d.get("discovered", {}).get("models", [])
        print("  MODEL KEYS:", [sorted(m.keys()) for m in models][:3])
    except Exception as e:
        print("  parse err", e)
con.close()
