import sqlite3
import os

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()

c.execute("SELECT model, provider, status, promptTokens, completionTokens, timestamp FROM usageHistory ORDER BY timestamp DESC LIMIT 20")
for r in c.fetchall():
    print(r)
