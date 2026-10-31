const path = require('path');
const Database = require(path.join(process.env.APPDATA, '9router/runtime/node_modules/better-sqlite3'));
const db = new Database(path.join(process.env.APPDATA, '9router/db/data.sqlite'));

const tables = db.prepare("SELECT name FROM sqlite_master WHERE type='table'").all();
console.log('TABLES:', tables.map(t => t.name));

for (const t of tables) {
  try {
    const count = db.prepare(`SELECT count(*) as c FROM "${t.name}"`).get();
    console.log(`Table ${t.name}: ${count.c} rows`);
  } catch (e) {
    console.log(`Table ${t.name}: error ${e.message}`);
  }
}
