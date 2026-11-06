"""Exporta só nome, tipo e membros dos combos (sem chaves) para combos.json. Somente leitura."""
import json, sqlite3, pathlib

DB = r"C:/Users/trcnologia/AppData/Roaming/9router/db/data.sqlite"
OUT = pathlib.Path(__file__).with_name("combos.json")
con = sqlite3.connect(f"file:{DB}?mode=ro&immutable=1", uri=True)
res = []
for nome, tipo, models, criado in con.execute("select name, kind, models, createdAt from combos order by createdAt, name"):
    res.append({"nome": nome, "tipo": tipo, "membros": json.loads(models), "criado": criado})
OUT.write_text(json.dumps(res, ensure_ascii=False, indent=1), encoding="utf-8")
print(len(res), "combos em", OUT.name)
