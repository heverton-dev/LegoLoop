import sqlite3, sys

db = r"C:\Users\trcnologia\AppData\Roaming\orca\agent-session-domain.db"
con = sqlite3.connect(f"file:{db}?mode=ro", uri=True)
cur = con.cursor()

cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
tables = [r[0] for r in cur.fetchall()]
print("TABLES:", tables)

for t in tables:
    cur.execute(f"PRAGMA table_info({t})")
    cols = [c[1] for c in cur.fetchall()]
    for col in cols:
        try:
            cur.execute(f"SELECT rowid FROM {t} WHERE {col} LIKE '%maxMessageTokens%' LIMIT 5")
            rows = cur.fetchall()
            if rows:
                print(f"HIT table={t} col={col} rows={rows}")
                for (rid,) in rows:
                    cur.execute(f"SELECT * FROM {t} WHERE rowid=?", (rid,))
                    data = cur.fetchone()
                    s = str(data)
                    idx = s.find("maxMessageTokens")
                    print("  SNIP:", s[max(0, idx-200):idx+200])
        except Exception:
            pass
con.close()
