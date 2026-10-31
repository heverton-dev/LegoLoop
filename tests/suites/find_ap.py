with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8499.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("function ap(")
if pos == -1:
    pos = content.find("ap=")
print("pos of ap:", pos)
print(content[pos:pos+3000])
