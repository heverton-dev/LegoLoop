const fs = require('fs');
const content = fs.readFileSync('C:/Users/trcnologia/Desktop/aidd-symbiont/index.html', 'utf8');
const regex = /code:\s*"([^"]+)",\s*title:\s*"([^"]+)"/g;
let m;
while ((m = regex.exec(content)) !== null) {
  console.log(`${m[1]}: ${m[2]}`);
}
