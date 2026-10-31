import sqlite3, os, json
conn = sqlite3.connect(os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite'))
c = conn.cursor()
c.execute("SELECT data FROM providerConnections WHERE provider='gemini' LIMIT 1")
row = c.fetchone()
print(json.loads(row[0])["apiKey"])
