with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\6922.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

start = content.find("function o(a){")
end = content.find("function p(a){", start)
print("Function o exact string length:", end - start)
print(content[start:end])
