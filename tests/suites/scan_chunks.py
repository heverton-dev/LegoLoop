import os
import re

chunks_dir = r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks"

terms = ["function_declarations", "functionDeclarations", "CONTENT_LENGTH_EXCEEDS_THRESHOLD", "generativelanguage", "gemini", "failover", "combo", "kiro"]

results = {}
for fname in os.listdir(chunks_dir):
    if not fname.endswith(".js"):
        continue
    fpath = os.path.join(chunks_dir, fname)
    with open(fpath, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()
    found = [t for t in terms if t in content]
    if found:
        results[fname] = found

print("Matching chunks:")
for k, v in results.items():
    print(f"{k}: {v}")
