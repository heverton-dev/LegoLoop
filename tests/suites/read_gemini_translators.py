with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("(0,d.kz)(e.h.OPENAI,e.h.GEMINI")
print(content[pos-2500:pos+1500])
