with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\6922.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

start = content.find("26922:")
end = content.find("33077:", start)
print("Start:", start, "End:", end)
print(content[start:start+300])
print("...")
print(content[end-200:end])
