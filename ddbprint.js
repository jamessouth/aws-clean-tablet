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
          if (Array.isArray(val)) {
              obj[p] = val.length;
          } else {
              obj[p] = val.length > 19 ? val.slice(0, 8) + "..." + val.slice(val.length - 8) : val.slice(0);
          }
      }
      return obj;
  });
  console.table(data);
});