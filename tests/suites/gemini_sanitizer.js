// Comprehensive Gemini Schema Sanitizer
function sanitizeGeminiSchema(schema, isRoot = true) {
  if (!schema || typeof schema !== "object") {
    return isRoot
      ? {
          type: "object",
          properties: {
            reason: {
              type: "string",
              description: "Brief explanation of why you are calling this tool"
            }
          }
        }
      : { type: "string" };
  }

  // Clone to avoid mutating external references
  const node = Array.isArray(schema) ? {} : { ...schema };

  // 1. Const to Enum
  if (node.const !== undefined && !node.enum) {
    node.enum = [node.const];
    delete node.const;
  }

  // 2. allOf flattening
  if (Array.isArray(node.allOf) && node.allOf.length > 0) {
    const mergedProps = {};
    const mergedReq = [];
    for (const sub of node.allOf) {
      if (sub && typeof sub === "object") {
        const sanitizedSub = sanitizeGeminiSchema(sub, false);
        if (sanitizedSub.properties) {
          Object.assign(mergedProps, sanitizedSub.properties);
        }
        if (Array.isArray(sanitizedSub.required)) {
          for (const r of sanitizedSub.required) {
            if (!mergedReq.includes(r)) mergedReq.push(r);
          }
        }
      }
    }
    delete node.allOf;
    if (Object.keys(mergedProps).length > 0) {
      node.properties = { ...(node.properties || {}), ...mergedProps };
    }
    if (mergedReq.length > 0) {
      node.required = [...(node.required || []), ...mergedReq];
    }
  }

  // 3. anyOf / oneOf resolution
  const unions = node.anyOf || node.oneOf;
  if (Array.isArray(unions) && unions.length > 0) {
    // Check if nullable
    const hasNull = unions.some(u => u && (u.type === "null" || (Array.isArray(u.type) && u.type.includes("null"))));
    if (hasNull) node.nullable = true;

    const nonNull = unions.filter(u => u && u.type !== "null");
    delete node.anyOf;
    delete node.oneOf;

    if (nonNull.length > 0) {
      // Pick best candidate: prefer object with properties, then array, then string/primitive
      let chosen = nonNull[0];
      let maxScore = -1;
      for (const cand of nonNull) {
        let score = 0;
        if (cand && typeof cand === "object") {
          const t = cand.type || (cand.properties ? "object" : cand.items ? "array" : "");
          if (t === "object" || cand.properties) score = 3;
          else if (t === "array" || cand.items) score = 2;
          else if (t) score = 1;
        }
        if (score > maxScore) {
          maxScore = score;
          chosen = cand;
        }
      }
      Object.assign(node, chosen);
    }
  }

  // 4. prefixItems -> items
  if (Array.isArray(node.prefixItems) && node.prefixItems.length > 0) {
    const valid = node.prefixItems.filter(p => p && p.type !== "null");
    if (!node.items) {
      node.items = valid.length > 0 ? valid[0] : { type: "string" };
    }
    delete node.prefixItems;
  }

  // 5. Type normalization (array types like ["string", "null"])
  if (Array.isArray(node.type)) {
    if (node.type.includes("null")) node.nullable = true;
    const nonNullType = node.type.filter(t => t !== "null");
    node.type = nonNullType.length > 0 ? nonNullType[0] : "string";
  }

  // Normalize type string
  if (typeof node.type === "string") {
    node.type = node.type.toLowerCase();
  }

  // Infer missing type
  if (!node.type) {
    if (node.properties || isRoot) {
      node.type = "object";
    } else if (node.items) {
      node.type = "array";
    } else if (node.enum) {
      node.type = "string";
    } else {
      node.type = "string";
    }
  }

  // 6. Handle array properties (e.g. properties: [{name: "x", type: "string"}])
  if (Array.isArray(node.properties)) {
    const propsObj = {};
    for (let i = 0; i < node.properties.length; i++) {
      const p = node.properties[i];
      if (p && typeof p === "object") {
        const key = p.name || p.key || p.id || String(i);
        propsObj[key] = p;
      }
    }
    node.properties = propsObj;
  }

  // 7. Clean forbidden fields
  const forbidden = [
    "name", "key", "id", "minimum", "maximum",
    "minLength", "maxLength", "exclusiveMinimum", "exclusiveMaximum",
    "minItems", "maxItems", "format", "multipleOf", "uniqueItems", "contains",
    "unevaluatedProperties", "unevaluatedItems", "contentSchema", "prefixItems",
    "additionalItems", "default", "examples", "$schema", "$defs", "definitions",
    "const", "$ref", "$comment", "deprecated", "readOnly", "writeOnly",
    "additionalProperties", "propertyNames", "patternProperties", "enumDescriptions",
    "anyOf", "oneOf", "allOf", "not", "dependencies", "dependentSchemas",
    "dependentRequired", "title", "optional", "if", "then", "else",
    "contentMediaType", "contentEncoding", "cornerRadius", "fillColor",
    "fontFamily", "fontSize", "fontWeight", "gap", "padding", "strokeColor",
    "strokeThickness", "textColor", "errorMessage", "errorMessages",
    "x-errorMessage", "x-errorMessages", "markdownDescription",
    "x-intellij-html-description", "x-taplo-info", "x-taplo", "doNotSuggest",
    "suggestSortText", "minProperties", "maxProperties"
  ];
  for (const f of forbidden) {
    delete node[f];
  }
  for (const k of Object.keys(node)) {
    if (k.startsWith("x-") || k.startsWith("$")) {
      delete node[k];
    }
  }

  // 8. Enum normalization
  if (Array.isArray(node.enum)) {
    node.enum = node.enum.map(x => String(x));
    node.type = "string";
  }

  // 9. Object handling
  if (node.type === "object") {
    if (isRoot) {
      // Root object: MUST have non-empty properties map
      if (!node.properties || typeof node.properties !== "object" || Object.keys(node.properties).length === 0) {
        node.properties = {
          reason: {
            type: "string",
            description: "Brief explanation of why you are calling this tool"
          }
        };
      } else {
        // Sanitize children
        const sanitized = {};
        for (const [k, v] of Object.entries(node.properties)) {
          sanitized[k] = sanitizeGeminiSchema(v, false);
        }
        node.properties = sanitized;
      }
    } else {
      // Nested object inside a property:
      // If it has NO properties or empty properties {}, Gemini v1beta rejects type: "object".
      // Convert it to type: "string" (JSON string representation).
      if (!node.properties || typeof node.properties !== "object" || Object.keys(node.properties).length === 0) {
        node.type = "string";
        node.description = node.description
          ? `${node.description} (JSON-encoded object string)`
          : "JSON-encoded object string";
        delete node.properties;
        delete node.required;
      } else {
        // Nested object with actual properties: sanitize children
        const sanitized = {};
        for (const [k, v] of Object.entries(node.properties)) {
          sanitized[k] = sanitizeGeminiSchema(v, false);
        }
        node.properties = sanitized;
      }
    }
  }

  // 10. Array handling
  if (node.type === "array") {
    if (!node.items || typeof node.items !== "object") {
      node.items = { type: "string" };
    } else {
      node.items = sanitizeGeminiSchema(node.items, false);
    }
  }

  // 11. Required fields filtering
  if (Array.isArray(node.required)) {
    if (node.properties && typeof node.properties === "object") {
      node.required = node.required.filter(
        r => typeof r === "string" && Object.prototype.hasOwnProperty.call(node.properties, r)
      );
    } else {
      delete node.required;
    }
    if (node.required && node.required.length === 0) {
      delete node.required;
    }
  } else {
    delete node.required;
  }

  return node;
}

// Export for testing
module.exports = { sanitizeGeminiSchema };
