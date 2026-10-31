import os
import re

chunks_dir = r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks"

def extract_contexts(fname, terms):
    fpath = os.path.join(chunks_dir, fname)
    with open(fpath, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()
    print(f"\n================ {fname} ================")
    for t in terms:
        for m in re.finditer(re.escape(t), content):
            start = max(0, m.start() - 200)
            end = min(len(content), m.end() + 200)
            print(f"--- Term '{t}' at pos {m.start()} ---")
            print(content[start:end])
            print("---------------------------------------")

extract_contexts("8499.js", ["functionDeclarations"])
extract_contexts("458.js", ["generativelanguage"])
