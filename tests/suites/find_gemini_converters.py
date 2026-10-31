with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

import re
# find calls where converters are registered
matches = re.findall(r'ao\([^)]+\)', content)
print("Registration calls (ao):", len(matches))
for m in matches:
    if 'GEMINI' in m or 'ANTIGRAVITY' in m:
        print(m)
