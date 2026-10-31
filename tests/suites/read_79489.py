with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8635.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

pos = content.find("79489:")
pos_next = content.find("87155:", pos)
print(content[pos:pos_next if pos_next != -1 else pos+15000])
