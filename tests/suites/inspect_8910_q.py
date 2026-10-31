with open(r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8910.js", "r", encoding="utf-8", errors="ignore") as f:
    content = f.read()

start = content.find("async function q(")
mid = content.find("let p=null,r=null,s=null;", start)
print(content[start:mid])
