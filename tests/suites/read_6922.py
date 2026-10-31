with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\6922.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("26922:")
print("Length of 6922.js:", len(content))
print("Position of 26922:", pos)
# Print around 4000 characters from 26922
print(content[pos:pos+5000])
