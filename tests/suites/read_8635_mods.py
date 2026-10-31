with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8635.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

import re

# find module ids in 8635.js
modules = re.findall(r'(\d+):\s*\(', content)
print("Modules in 8635.js:", modules[:20])

for m in re.finditer(r'(\d+):\s*\([a-z,\s]*\)\s*=>', content):
    mod_id = m.group(1)
    snippet = content[m.start():m.start()+250]
    print(f"\n--- Module {mod_id} ---")
    print(snippet)
