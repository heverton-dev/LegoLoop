import os

fpath = r"C:\Users\trcnologia\AppData\Roaming\npm\node_modules\9router\app\.next-cli-build\server\chunks\8910.js"
with open(fpath, "r", encoding="utf-8") as f:
    content = f.read()

target = 'a[0]!==o[0]&&g.info("COMBO",`auto-switch for [${[...b].join(",")}] → ${a[0]}`),o=a}}let p=null,r=null,s=null;'

if target not in content:
    print("Target not found in 8910.js!")
    exit(1)

replacement = '''a[0]!==o[0]&&g.info("COMBO",`auto-switch for [${[...b].join(",")}] → ${a[0]}`),o=a}}let estTokens=(function(e){if(!e||"object"!=typeof e)return 0;let t=0;if(Array.isArray(e.messages))for(let a of e.messages)"string"==typeof a.content?t+=a.content.length:Array.isArray(a.content)&&a.content.forEach(a=>{"string"==typeof a.text&&(t+=a.text.length)});if(Array.isArray(e.contents))for(let a of e.contents)Array.isArray(a.parts)&&a.parts.forEach(a=>{"string"==typeof a.text&&(t+=a.text.length)});if(Array.isArray(e.input))for(let a of e.input)"string"==typeof a.content?t+=a.content.length:Array.isArray(a.content)&&a.content.forEach(a=>{"string"==typeof a.text&&(t+=a.text.length)});"string"==typeof e.system&&(t+=e.system.length);Array.isArray(e.tools)&&(t+=JSON.stringify(e.tools).length);return Math.ceil(t/3.5)})(a);if(estTokens>100000){let is1M=m=>{let l=m.toLowerCase();return(l.includes("gemini")||l.includes("flash"))&&!l.startsWith("kr/")&&!l.includes("kiro")};let isKr=m=>{let l=m.toLowerCase();return l.startsWith("kr/")||l.includes("kiro")};let ones=o.filter(is1M),others=o.filter(m=>!is1M(m)&&!isKr(m)),krs=o.filter(isKr);ones.sort((a,b)=>{let fa=a.includes("flash")?1:0,fb=b.includes("flash")?1:0;return fb-fa});let reord=[...ones,...others,...krs];if(reord.length>0&&reord[0]!==o[0]){g.info("COMBO",`large-context detected (~${estTokens} tokens) → prioritizing 1M models (${reord[0]})`);o=reord;}}let p=null,r=null,s=null;'''

new_content = content.replace(target, replacement, 1)

with open(fpath, "w", encoding="utf-8") as f:
    f.write(new_content)

print(f"Updated 8910.js successfully. Original len: {len(content)}, New len: {len(new_content)}")
