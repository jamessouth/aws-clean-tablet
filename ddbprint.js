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
  // const data = Items.map(i => {
  //     const obj = {};
  //     for (const p in i) {
  //         const val = Object.values(i[p])[0];
  //         console.log(val);
  //         if (p === "players") {
  //             obj[p] = Object.keys(i[p].M).length;
  //         } else {
  //             obj[p] = val.length > 19 ? val.slice(0, 8) + "..." + val.slice(val.length - 8) : val.slice(0);
  //         }
  //     }
  //     return obj;
  // });
  const data = Items.reduce((acc, itm) => {
    for (const p in itm) {
      if (p === "players") {
        acc.push(...Object.keys(itm[p].M).map(x => {
          // console.log('xxx: ', itm[p].M[x]);
          return {
          
          pk: itm.pk.S,
          sk: itm.sk.S,
          players: itm[p].M[x].M.name.S + "_" + itm[p].M[x].M.ready.BOOL
        }}));
      } else if (p === "pk" && itm[p].S !== "GAME") {
        acc.push(itm);
      }
    }
    return acc;
  }, []);
  console.log('jj: ', data);
  console.table(data);
});