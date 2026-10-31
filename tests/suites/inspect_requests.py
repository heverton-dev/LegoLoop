import sqlite3
import os
import json

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()

c.execute("SELECT name FROM sqlite_master WHERE type='table'")
tables = [row[0] for row in c.fetchall()]

if "requestDetails" in tables:
    c.execute("PRAGMA table_info(requestDetails)")
    cols = [col[1] for col in c.fetchall()]
    print("requestDetails cols:", cols)
    c.execute("SELECT * FROM requestDetails ORDER BY rowid DESC LIMIT 5")
    rows = c.fetchall()
    for r in rows:
        d = dict(zip(cols, r))
        print("--- Row ---")
        for k, v in d.items():
            if isinstance(v, str) and len(v) > 200:
                print(f"{k}: {v[:200]}... (len {len(v)})")
            else:
                print(f"{k}: {v}")
