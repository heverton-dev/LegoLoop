import sqlite3
import os
import json

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()

c.execute("SELECT id, provider, data FROM providerConnections WHERE provider='gemini' LIMIT 1")
row = c.fetchone()
data = json.loads(row[2])
apiKey = data.get("apiKey")
masked = apiKey[:6] + "..." + apiKey[-4:] if apiKey else "None"
print(f"Provider: {row[1]}, Key: {masked}")
