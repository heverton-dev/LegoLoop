import sqlite3
import os
import json

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()

tables = [row[0] for row in c.execute("SELECT name FROM sqlite_master WHERE type='table'")]
print("TABLES:", tables)

for t in tables:
    c.execute(f'SELECT count(*) FROM "{t}"')
    count = c.fetchone()[0]
    print(f"{t}: {count} rows")

print("\n--- COMBOS ---")
try:
    c.execute("SELECT id, name, models FROM combos")
    for row in c.fetchall():
        print(f"ID: {row[0]}, Name: {row[1]}, Models: {row[2]}")
except Exception as e:
    print("Error querying combos:", e)

print("\n--- PROVIDERS ---")
try:
    c.execute("SELECT id, name, type FROM providers")
    for row in c.fetchall():
        print(row)
except Exception as e:
    print("Error querying providers:", e)
