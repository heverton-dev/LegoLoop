async function testCombo() {
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
              description: "Arbitrary environment variables dictionary"
            },
            options: {
              type: "object",
              properties: {}
            },
            timeout: {
              type: "number",
              minimum: 0
            }
          },
          required: ["command", "non_existent_key"]
        }
      }
    }
  ];

  const payload = {
    model: "code-fast",
    messages: [
      { role: "user", content: "Briefly confirm you can see tool execute_command." }
    ],
    tools: problematicTools,
    stream: false
  };

  console.log("Sending POST to", url, "with combo:", payload.model);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${apiKey}`,
      "X-9Router-Token-Saver": "off"
    },
    body: JSON.stringify(payload)
  });

  console.log(`HTTP STATUS: ${res.status} ${res.statusText}`);
  const json = await res.json();
  console.log("RESPONSE:", JSON.stringify(json.choices?.[0]?.message?.content || json, null, 2));
  if (res.status === 200) {
    console.log(">>> COMBO TEST PASSED: HTTP 200 OK! <<<");
  }
}

testCombo();
