import sqlite3
import os

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()

c.execute("PRAGMA table_info(providerConnections)")
cols = [col[1] for col in c.fetchall()]
print("providerConnections cols:", cols)

c.execute(f"SELECT {', '.join(cols[:5])} FROM providerConnections")
for r in c.fetchall():
    print(r)
