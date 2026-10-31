with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    code = f.read()

import re

# find where Yr is defined / exported
for m in re.finditer(r'Yr:(\(\)=>|\w+)', code):
    print(f"Export Yr: {code[m.start()-50:m.end()+50]}")

# search for function definition of Yr or the identifier it maps to
pos = code.find('Yr:')
if pos != -1:
    snippet = code[pos:pos+200]
    print("Snippet around Yr:", snippet)
