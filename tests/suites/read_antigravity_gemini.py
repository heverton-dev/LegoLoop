with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("b.length>0&&e.request.contents.push({role:f.role===k.z7.ASSISTANT")
print(content[pos:pos+2500])
