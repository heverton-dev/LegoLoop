const path = require('path');
const chunk6922 = require('C:/Users/trcnologia/AppData/Roaming/npm/node_modules/9router/app/.next-cli-build/server/chunks/6922.js');

const yrMod = {};
const mockRequire = function(id) {
  if (chunk6922.modules[id]) {
    const exports = {};
    chunk6922.modules[id](null, exports, mockRequire);
    return exports;
  }
  return {};
};
mockRequire.d = (obj, getters) => {
  for (const [k, v] of Object.entries(getters)) {
    Object.defineProperty(obj, k, { get: v });
  }
};

chunk6922.modules[26922](null, yrMod, mockRequire);

const mockTools = [
  {
    name: "bash",
    input_schema: {
      type: "object",
      properties: {
        command: { type: "string", description: "Command" },
        timeout: { type: "number", minimum: 0 }
      },
      required: ["command"]
    }
  },
  {
    name: "mcp_complex_tool",
    input_schema: {
      type: "object",
      properties: {
        metadata: {
          type: "object",
          description: "Arbitrary dictionary metadata"
        },
        options: {
          type: "object",
          properties: {}
        },
        items: {
          type: "array",
          items: {
            type: "object"
          }
        },
        union_field: {
          anyOf: [
            { type: "string" },
            { type: "null" }
          ]
        }
      },
      required: ["metadata", "non_existent"]
    }
  },
  {
    name: "tool_with_array_properties",
    input_schema: {
      type: "object",
      properties: [
        { name: "query", type: "string", description: "Search query" },
        { name: "limit", type: "integer", default: 10 }
      ]
    }
  }
];

for (const t of mockTools) {
  console.log(`\n--- Sanitizing tool: ${t.name} ---`);
  const sanitized = yrMod.Yr(t.input_schema);
  console.log(JSON.stringify(sanitized, null, 2));
}
