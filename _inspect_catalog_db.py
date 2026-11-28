import sqlite3

db = r"C:\Users\trcnologia\AppData\Roaming\orca\agent-model-catalog.db"
con = sqlite3.connect(f"file:{db}?mode=ro", uri=True)
cur = con.cursor()
cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
tables = [r[0] for r in cur.fetchall()]
print("TABLES:", tables)
for t in tables:
    cur.execute(f"PRAGMA table_info({t})")
    cols = [c[1] for c in cur.fetchall()]
    print(f"-- {t}: {cols}")
    cur.execute(f"SELECT * FROM {t} LIMIT 20")
    for row in cur.fetchall():
        s = str(row)
        if "code-fast" in s or "maxMessage" in s or "token" in s.lower():
            print("  ROW:", s[:500])
con.close()
