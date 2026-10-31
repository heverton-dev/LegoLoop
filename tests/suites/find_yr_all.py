import os, re

chunks_dir = r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks"
for fname in os.listdir(chunks_dir):
    if not fname.endswith(".js"): continue
    fpath = os.path.join(chunks_dir, fname)
    with open(fpath, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()
    for m in re.finditer(r'Yr:(\(\)=>|\w+)', content):
        print(f"File {fname} pos {m.start()}: {content[m.start()-50:m.end()+150]}")
