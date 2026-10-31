import sqlite3, os
conn = sqlite3.connect(os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite'))
c = conn.cursor()
c.execute("SELECT id, name, key FROM apiKeys LIMIT 2")
for r in c.fetchall():
    print(r)
