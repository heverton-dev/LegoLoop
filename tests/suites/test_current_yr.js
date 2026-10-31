// Test current Yr (function o) behavior on complex schemas
const path = require('path');

// Let's load chunk 6922
// Since it's a webpack chunk, let's extract function o and its dependencies
const fs = require('fs');
const content = fs.readFileSync('C:/Users/trcnologia/AppData/Roaming/npm/node_modules/9router/app/.next-cli-build/server/chunks/6922.js', 'utf-8');

// Let's run a test with the exact code
const f = ["minLength","maxLength","exclusiveMinimum","exclusiveMaximum","minItems","maxItems","format","multipleOf","uniqueItems","contains","unevaluatedProperties","unevaluatedItems","contentSchema","prefixItems","additionalItems","default","examples","$schema","$defs","definitions","const","$ref","$comment","deprecated","readOnly","writeOnly","additionalProperties","propertyNames","patternProperties","enumDescriptions","anyOf","oneOf","allOf","not","dependencies","dependentSchemas","dependentRequired","title","optional","deprecated","if","then","else","contentMediaType","contentEncoding","cornerRadius","fillColor","fontFamily","fontSize","fontWeight","gap","padding","strokeColor","strokeThickness","textColor","errorMessage","errorMessages","x-errorMessage","x-errorMessages","markdownDescription","x-intellij-html-description","x-taplo-info","x-taplo","doNotSuggest","suggestSortText","minProperties","maxProperties"];

function n(a){let b=0,c=-1;for(let d=0;d<a.length;d++){let e=a[d],f=0,g=e.type;"object"===g||e.properties?f=3:"array"===g||e.items?f=2:g&&"null"!==g&&(f=1),f>c&&(c=f,b=d)}return b}

function currentYr(a){return a&&"object"==typeof a&&(!function a(b){if(b&&"object"==typeof b)for(let c of(void 0===b.const||b.enum||(b.enum=[b.const],delete b.const),Object.values(b)))c&&"object"==typeof c&&a(c)}(a),!function a(b){if(b&&"object"==typeof b)for(let c of(b.enum&&Array.isArray(b.enum)&&(b.enum=b.enum.map(a=>String(a)),b.type||(b.type="string")),Object.values(b)))c&&"object"==typeof c&&a(c)}(a),!function a(b){if(b&&"object"==typeof b){if(b.allOf&&Array.isArray(b.allOf)){let a={};for(let c of b.allOf)if(c.properties&&(a.properties||(a.properties={}),Object.assign(a.properties,c.properties)),c.required&&Array.isArray(c.required))for(let b of(a.required||(a.required=[]),c.required))a.required.includes(b)||a.required.push(b);delete b.allOf,a.properties&&(b.properties={...b.properties,...a.properties}),a.required&&(b.required=[...b.required||[],...a.required])}for(let c of Object.values(b))c&&"object"==typeof c&&a(c)}}(a),!function a(b){if(b&&"object"==typeof b){if(Array.isArray(b.prefixItems)&&b.prefixItems.length>0){let a=b.prefixItems.filter(a=>a&&"null"!==a.type);b.items||1!==a.length?!b.items&&a.length>1&&(b.items={anyOf:a}):b.items=a[0],delete b.prefixItems}for(let c of Object.values(b))c&&"object"==typeof c&&a(c)}}(a),!function a(b){if(b&&"object"==typeof b){if(b.anyOf&&Array.isArray(b.anyOf)&&b.anyOf.length>0){let a=b.anyOf.filter(a=>a&&"null"!==a.type);if(a.length>0){let c=n(a),d=a[c];delete b.anyOf,Object.assign(b,d)}}if(b.oneOf&&Array.isArray(b.oneOf)&&b.oneOf.length>0){let a=b.oneOf.filter(a=>a&&"null"!==a.type);if(a.length>0){let c=n(a),d=a[c];delete b.oneOf,Object.assign(b,d)}}for(let c of Object.values(b))c&&"object"==typeof c&&a(c)}}(a),!function a(b){if(b&&"object"==typeof b){if(b.type&&Array.isArray(b.type)){let a=b.type.filter(a=>"null"!==a);b.type=a.length>0?a[0]:"string"}for(let c of Object.values(b))c&&"object"==typeof c&&a(c)}}(a),!function a(b){if(b&&"object"==typeof b)for(let c of(b.properties&&!b.type&&(b.type="object"),Object.values(b)))c&&"object"==typeof c&&a(c)}(a),!function a(b){if(b&&"object"==typeof b)for(let c of("array"!==b.type||b.items||(b.items={type:"string"}),Object.values(b)))c&&"object"==typeof c&&a(c)}(a),!function a(b,c){if(b&&"object"==typeof b){if(Array.isArray(b)){for(let d of b)a(d,c);return}for(let d of Object.keys(b)){if(c.includes(d)||d.startsWith("x-")){delete b[d];continue}let e=b[d];e&&"object"==typeof e&&a(e,c)}}}(a,f),!function a(b){if(b&&"object"==typeof b){if(b.required&&Array.isArray(b.required)&&b.properties){let a=b.required.filter(a=>Object.prototype.hasOwnProperty.call(b.properties,a));0===a.length?delete b.required:b.required=a}for(let c of Object.values(b))c&&"object"==typeof c&&a(c)}}(a),!function a(b){if(b&&"object"==typeof b){if(0===Object.keys(b).length){b.type="object",b.properties={reason:{type:"string",description:"Brief explanation of why you are calling this tool"}},b.required=["reason"];return}for(let c of("object"!==b.type||b.properties&&0!==Object.keys(b.properties).length||(b.properties={reason:{type:"string",description:"Brief explanation of why you are calling this tool"}},b.required=["reason"]),Object.values(b)))c&&"object"==typeof c&&a(c)}}(a)),a}

// Test cases:
const testCases = [
  // 1. Tool parameter that is an object without properties (like MCP metadata, args, headers)
  {
    name: "Object parameter without properties",
    input: {
      type: "object",
      properties: {
        meta: {
          type: "object"
        }
      }
    }
  },
  // 2. Tool parameter where properties is an array
  {
    name: "Properties as array",
    input: {
      type: "object",
      properties: [
        { name: "query", type: "string" }
      ]
    }
  },
  // 3. Properties where a property has type: ["object", "null"]
  {
    name: "Union type with object",
    input: {
      type: "object",
      properties: {
        config: {
          type: ["object", "null"]
        }
      }
    }
  },
  // 4. AnyOf with object without properties
  {
    name: "AnyOf with empty object",
    input: {
      type: "object",
      properties: {
        filter: {
          anyOf: [
            { type: "string" },
            { type: "object" }
          ]
        }
      }
    }
  },
  // 5. Array of objects where items has no properties
  {
    name: "Array with empty object items",
    input: {
      type: "object",
      properties: {
        list: {
          type: "array",
          items: {
            type: "object"
          }
        }
      }
    }
  }
];

for (const tc of testCases) {
  console.log("=== TEST:", tc.name);
  const cloned = JSON.parse(JSON.stringify(tc.input));
  try {
    const res = currentYr(cloned);
    console.log(JSON.stringify(res, null, 2));
  } catch (e) {
    console.error("Error:", e);
  }
}
