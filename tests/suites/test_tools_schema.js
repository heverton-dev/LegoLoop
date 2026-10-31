// Let's inspect tools that Claude Code typically sends
// and test normalizer behavior.

const toolsExample = [
  {
    name: "bash",
    description: "Run bash command",
    input_schema: {
      type: "object",
      properties: {
        command: { type: "string" },
        timeout: { type: "number" }
      },
      required: ["command"]
    }
  },
  {
    name: "mcp_tool_with_empty_object",
    description: "Tool with empty object parameter",
    input_schema: {
      type: "object",
      properties: {
        metadata: {
          type: "object" // NO PROPERTIES!
        },
        payload: {
          type: "object",
          properties: {} // EMPTY PROPERTIES!
        },
        items_list: {
          type: "array",
          items: {
            type: "object" // OBJECT IN ARRAY ITEMS WITHOUT PROPERTIES!
          }
        },
        union_param: {
          anyOf: [
            { type: "string" },
            { type: "object" }
          ]
        }
      }
    }
  }
];

console.log("Mock tools ready for test");
