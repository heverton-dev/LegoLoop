import json, io, re

d = json.load(io.open(r"C:\Users\trcnologia\.gw\registry.json", encoding="utf-8"))

def walk(obj, path=""):
    if isinstance(obj, dict):
        for k, v in obj.items():
            walk(v, f"{path}.{k}")
    elif isinstance(obj, list):
        for i, v in enumerate(obj):
            walk(v, f"{path}[{i}]")
    else:
        s = str(obj)
        if "code-fast" in path or (isinstance(obj, str) and "code-fast" in s) or "token" in path.lower():
            print(path, "=", s[:200])

walk(d)

print("\n=== TOP KEYS ===")
print(json.dumps({k: (list(v.keys()) if isinstance(v, dict) else type(v).__name__) for k, v in d.items()}, indent=1))
h = d.get("harnesses", {})
print("\n=== HARNESS KEYS ===", list(h.keys()) if isinstance(h, dict) else h)
