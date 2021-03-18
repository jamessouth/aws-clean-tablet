const readable = process.stdin;

const chunks = [];

readable.on('readable', () => {
  let chunk;
  while (null !== (chunk = readable.read())) {
    chunks.push(chunk);
  }
});

readable.on('end', () => {
  const content = chunks.join('');
  const { Items } = JSON.parse(content);
  const data = Items.map(i => {
      const obj = {};
      for (const p in i) {
        const val = Object.values(i[p])[0];
        const half = val.length > 19 ? val.length / 2 : 0;
        obj[p] = val.slice(half);
      }
      return obj;
  });
  console.table(data);
});