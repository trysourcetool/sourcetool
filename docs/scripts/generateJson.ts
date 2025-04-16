import fs from 'fs';
import path from 'path';

const getNestedFiles = (dir: string) => {
  const files = fs.readdirSync(dir);
  return files
    .map((file) => {
      const filePath = path.join(dir, file);
      if (fs.statSync(filePath).isDirectory()) {
        return getNestedFiles(filePath);
      }
      return {
        path: filePath,
        filename: file,
      };
    })
    .flat()
    .filter((file) => file.filename.endsWith('.md'));
};

const files = getNestedFiles('../docs/docs');

const json = files.map((file) => {
  const content = fs.readFileSync(file.path, 'utf8');
  return {
    title: file.filename.replace('.md', ''),
    path: file.path.replace('../docs/docs', 'docs'),
    content: content.replace(/---\n.+\n---\n+/g, ''),
  };
});

fs.mkdirSync(path.join('./build/assets/json'), { recursive: true });
fs.writeFileSync(
  path.join('./build/assets/json', 'docs.json'),
  JSON.stringify(json),
);
