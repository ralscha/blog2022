const fs = require('fs');
const archiver = require('archiver');
const {join} = require("node:path");

const packageJson = JSON.parse(fs.readFileSync(join(__dirname, 'package.json'), 'utf8'));
const name = packageJson.name;
const version = packageJson.version;

const output = fs.createWriteStream(join(__dirname, `${name}-${version}.zip`));
const archive = archiver('zip', {
  zlib: { level: 9 }
});

output.on('close', () => {
  console.log(`Zip file created successfully. Total bytes: ${archive.pointer()}`);
});

archive.on('error', err => {
  throw err;
});

archive.pipe(output);
archive.directory('www/browser/', '/', {});
archive.finalize();
