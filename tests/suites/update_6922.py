import os

fpath = r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\6922.js"
with open(fpath, "r", encoding="utf-8") as f:
    content = f.read()

start_marker = "function o(a){"
end_marker = "function p(a){"

start_pos = content.find(start_marker)
end_pos = content.find(end_marker, start_pos)

if start_pos == -1 or end_pos == -1:
    print("Markers not found!")
    exit(1)

new_fn_o = """function o(a,isRoot=!0){if(!a||"object"!=typeof a)return isRoot?{type:"object",properties:{reason:{type:"string",description:"Brief explanation of why you are calling this tool"}}}:{type:"string"};let b=Array.isArray(a)?{}:structuredClone(a);void 0===b.const||b.enum||(b.enum=[b.const],delete b.const);if(Array.isArray(b.allOf)&&b.allOf.length>0){let a={},c=[];for(let d of b.allOf)if(d&&"object"==typeof d){let e=o(d,!1);e.properties&&Object.assign(a,e.properties),Array.isArray(e.required)&&e.required.forEach(a=>{c.includes(a)||c.push(a)})}delete b.allOf,Object.keys(a).length>0&&(b.properties={...b.properties||{},...a}),c.length>0&&(b.required=[...b.required||[],...c])}let c=b.anyOf||b.oneOf;if(Array.isArray(c)&&c.length>0){c.some(a=>a&&("null"===a.type||Array.isArray(a.type)&&a.type.includes("null")))&&(b.nullable=!0);let a=c.filter(a=>a&&"null"!==a.type);delete b.anyOf,delete b.oneOf;if(a.length>0){let c=n(a),d=a[c];Object.assign(b,d)}}if(Array.isArray(b.prefixItems)&&b.prefixItems.length>0){let a=b.prefixItems.filter(a=>a&&"null"!==a.type);b.items||(b.items=a.length>0?a[0]:{type:"string"}),delete b.prefixItems}if(Array.isArray(b.type)){b.type.includes("null")&&(b.nullable=!0);let a=b.type.filter(a=>"null"!==a);b.type=a.length>0?a[0]:"string"}"string"==typeof b.type&&(b.type=b.type.toLowerCase()),b.type||(b.properties||isRoot?b.type="object":b.items?b.type="array":b.enum?b.type="string":b.type="string");if(Array.isArray(b.properties)){let a={};for(let c=0;c<b.properties.length;c++){let d=b.properties[c];if(d&&"object"==typeof d){let e=d.name||d.key||d.id||String(c);a[e]=d}}b.properties=a}let d=["name","key","id","minimum","maximum",...f];for(let a of d)delete b[a];for(let a of Object.keys(b))(a.startsWith("x-")||a.startsWith("$"))&&delete b[a];Array.isArray(b.enum)&&(b.enum=b.enum.map(a=>String(a)),b.type="string");if("object"===b.type)if(isRoot){if(!b.properties||"object"!=typeof b.properties||0===Object.keys(b.properties).length)b.properties={reason:{type:"string",description:"Brief explanation of why you are calling this tool"}};else{let a={};for(let[c,d]of Object.entries(b.properties))a[c]=o(d,!1);b.properties=a}}else if(!b.properties||"object"!=typeof b.properties||0===Object.keys(b.properties).length)b.type="string",b.description=b.description?`${b.description} (JSON-encoded object string)`:"JSON-encoded object string",delete b.properties,delete b.required;else{let a={};for(let[c,d]of Object.entries(b.properties))a[c]=o(d,!1);b.properties=a}"array"===b.type&&(b.items&&"object"==typeof b.items?b.items=o(b.items,!1):b.items={type:"string"});Array.isArray(b.required)?(b.properties&&"object"==typeof b.properties&&(b.required=b.required.filter(a=>"string"==typeof a&&Object.prototype.hasOwnProperty.call(b.properties,a))),(!b.required||0===b.required.length)&&delete b.required):delete b.required;return b}"""

new_content = content[:start_pos] + new_fn_o + content[end_pos:]

with open(fpath, "w", encoding="utf-8") as f:
    f.write(new_content)

print(f"Updated 6922.js successfully. Original len: {len(content)}, New len: {len(new_content)}")
