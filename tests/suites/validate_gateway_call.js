async function testGateway() {
  const url = "http://127.0.0.1:20128/v1/chat/completions";
  const apiKey = "sk-4bc1e878b525f20a-qo8rt7-7ca55e91";

  const problematicTools = [
    {
      type: "function",
      function: {
        name: "execute_command",
        description: "Executes a system shell command with arbitrary environment variables",
        parameters: {
          type: "object",
          title: "ExecuteCommandParams",
          properties: {
            command: { type: "string", description: "The shell command to run" },
            env: {
              type: "object",
              description: "Arbitrary environment variables dictionary" // EMPTY OBJECT NODE!
            },
            options: {
              type: "object",
              properties: {} // EMPTY PROPERTIES OBJECT!
            },
            timeout: {
              type: "number",
              minimum: 0 // FORBIDDEN KEY IN GEMINI
            }
          },
          required: ["command", "non_existent_key"] // NON-EXISTENT REQUIRED KEY
        }
      }
    },
    {
      type: "function",
      function: {
        name: "search_data",
        description: "Searches data with array properties schema",
        parameters: {
          type: "object",
          // PROPERTIES AS ARRAY (Anthropic/MCP format bug)
          properties: [
            { name: "query", type: "string", description: "Search keyword" },
            { name: "filters", type: "object", description: "Filter dict" }
          ]
        }
      }
    },
    {
      type: "function",
      function: {
        name: "list_items",
        description: "Lists items with array of objects",
        parameters: {
          type: "object",
          properties: {
            tags: {
              type: "array",
              items: {
                type: "object" // EMPTY OBJECT IN ARRAY ITEMS!
              }
            },
            union_field: {
              anyOf: [
                { type: "string" },
                { type: "null" }
              ]
            }
          }
        }
      }
    }
  ];

  const payload = {
    model: "gemini/gemini-3.5-flash-lite",
    messages: [
      { role: "user", content: "List the tools you have available and explain what execute_command does." }
    ],
    tools: problematicTools,
    stream: false
  };

  console.log("Sending POST to", url);
  console.log("Model:", payload.model);
  console.log("Tools count:", payload.tools.length);

  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${apiKey}`,
        "X-9Router-Token-Saver": "off"
      },
      body: JSON.stringify(payload)
    });

    console.log(`\nHTTP STATUS: ${res.status} ${res.statusText}`);
    const json = await res.json();
    console.log("RESPONSE JSON:", JSON.stringify(json, null, 2));

    if (res.status === 200 && json.choices && json.choices.length > 0) {
      console.log("\n>>> VALIDATION PASSED: HTTP 200 OK! Tools sanitized successfully and accepted by Gemini! <<<");
    } else {
      console.error("\n>>> VALIDATION FAILED! <<<");
    }
  } catch (err) {
    console.error("Fetch error:", err);
  }
}

testGateway();
