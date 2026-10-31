with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\4572.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("12557:")
print(content[pos:pos+4000])
