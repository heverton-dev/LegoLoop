with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

import re
for m in re.finditer(r'GH:\s*\(', content):
    pos = m.start()
    print("Found GH at pos:", pos)
    print(content[pos:pos+1500])
