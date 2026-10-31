import sqlite3, os, json, urllib.request, urllib.error

db_path = os.path.expandvars(r'%APPDATA%\9router\db\data.sqlite')
conn = sqlite3.connect(db_path)
c = conn.cursor()
c.execute("SELECT data FROM providerConnections WHERE provider='gemini' LIMIT 1")
apiKey = json.loads(c.fetchone()[0]).get("apiKey")

url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-3.8-flash:generateContent?key={apiKey}"

def test_schema(name, schema):
    body = {
        "contents": [
            { "role": "user", "parts": [{"text": "Hello"}] }
        ],
        "tools": [
            {
                "functionDeclarations": [
                    {
                        "name": "test_tool",
                        "description": "test tool",
                        "parameters": schema
                    }
                ]
            }
        ]
    }
    req = urllib.request.Request(
        url,
        data=json.dumps(body).encode("utf-8"),
        headers={"Content-Type": "application/json"}
    )
    print(f"\n--- Testing: {name} ---")
    try:
        with urllib.request.urlopen(req) as resp:
            print("STATUS:", resp.status)
            print("OK!")
    except urllib.error.HTTPError as e:
        print("HTTP ERROR:", e.code)
        err_body = e.read().decode("utf-8")
        print("Response:", err_body[:400])

# Case 1: Empty properties
test_schema("Case 1: properties is empty object {}", {
    "type": "OBJECT",
    "properties": {}
})

# Case 2: properties is array
test_schema("Case 2: properties is array []", {
    "type": "OBJECT",
    "properties": []
})

# Case 3: nested object without properties
test_schema("Case 3: nested object without properties", {
    "type": "OBJECT",
    "properties": {
        "meta": {
            "type": "OBJECT"
        }
    }
})

# Case 4: nested object with empty properties {}
test_schema("Case 4: nested object with empty properties {}", {
    "type": "OBJECT",
    "properties": {
        "meta": {
            "type": "OBJECT",
            "properties": {}
        }
    }
})

# Case 5: valid object with string property
test_schema("Case 5: valid object with string property", {
    "type": "OBJECT",
    "properties": {
        "meta": {
            "type": "STRING",
            "description": "metadata"
        }
    }
})

# Case 6: properties with array containing object
test_schema("Case 6: properties as array with object", {
    "type": "OBJECT",
    "properties": [{"name": "foo", "type": "STRING"}]
})
