const path = require('path');
const { sanitizeGeminiSchema } = require('./gemini_sanitizer');
const { execSync } = require('child_process');

const apiKey = execSync(`python "${path.join(__dirname, 'get_key.py')}"`, { encoding: 'utf-8' }).trim();

async function runTest() {
  const model = "gemini-3.5-flash-lite";
  const url = `https://generativelanguage.googleapis.com/v1beta/models/${model}:generateContent?key=${apiKey}`;

  // Test suite with problematic schemas sent by Claude Code & MCPs
  const testSchemas = [
    {
      name: "Empty root schema {}",
      raw: {}
    },
    {
      name: "Root schema with empty properties map {}",
      raw: { type: "object", properties: {} }
    },
    {
      name: "Properties as array [ { name: 'q', type: 'string' } ]",
      raw: {
        type: "object",
        properties: [
          { name: "q", type: "string" }
        ]
      }
    },
    {
      name: "Nested empty object parameter (e.g. metadata/filters/args)",
      raw: {
        type: "object",
        properties: {
          metadata: {
            type: "object",
            description: "Arbitrary dictionary"
          },
          tags: {
            type: "array",
            items: {
              type: "object" // empty object items
            }
          }
        }
      }
    },
    {
      name: "Complex union and invalid fields",
      raw: {
        type: "object",
        $schema: "http://json-schema.org/draft-07/schema#",
        title: "TestTool",
        additionalProperties: false,
        properties: {
          param1: {
            anyOf: [
              { type: "string" },
              { type: "null" }
            ],
            description: "Nullable string"
          },
          param2: {
            oneOf: [
              { type: "object", properties: { key: { type: "string" } } },
              { type: "string" }
            ]
          },
          param3: {
            type: ["number", "null"],
            default: 10,
            minimum: 1
          }
        },
        required: ["param1", "non_existent_prop"]
      }
    }
  ];

  for (const tc of testSchemas) {
    console.log(`\n========================================`);
    console.log(`TEST: ${tc.name}`);
    const sanitized = sanitizeGeminiSchema(tc.raw, true);
    console.log("Sanitized Schema:", JSON.stringify(sanitized, null, 2));

    const body = {
      contents: [{ role: "user", parts: [{ text: "Hello, what tools do you have?" }] }],
      tools: [
        {
          functionDeclarations: [
            {
              name: "sample_tool",
              description: "A tool to test schema sanitization",
              parameters: sanitized
            }
          ]
        }
      ]
    };

    try {
      const resp = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
      });
      console.log(`HTTP Status: ${resp.status} ${resp.statusText}`);
      if (!resp.ok) {
        const errText = await resp.text();
        console.error("API Error Response:", errText);
      } else {
        const data = await resp.json();
        console.log("SUCCESS! Model response received:", data.candidates?.[0]?.content?.parts?.[0]?.text?.slice(0, 80));
      }
    } catch (e) {
      console.error("Network/Fetch error:", e.message);
    }
  }
}

runTest();
