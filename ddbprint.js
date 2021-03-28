const readable = process.stdin;

const chunks = [];

readable.on("readable", () => {
    let chunk;
    while (null !== (chunk = readable.read())) {
        chunks.push(chunk);
    }
});

readable.on("end", () => {
    const content = chunks.join("");
    const { Items } = JSON.parse(content);
    const data = Items.reduce((acc, itm) => {
        const obj = {};
        for (const p in itm) {
            if (p === "players") {
                if (Object.keys(itm[p].M).length === 0) {
                    acc.push({
                        pk: itm.pk.S,
                        sk: itm.sk.S,
                        players: 0,
                    });
                } else {
                    acc.push(
                        ...Object.keys(itm[p].M).map((x) => {
                            return {
                                pk: itm.pk.S,
                                sk: itm.sk.S,
                                players:
                                    itm[p].M[x].M.name.S +
                                    "_" +
                                    itm[p].M[x].M.ready.BOOL,
                            };
                        })
                    );
                }
                return acc;
            } else if (itm["pk"].S !== "GAME") {
                const val = Object.values(itm[p])[0];
                obj[p] =
                    val.length > 19
                        ? val.slice(0, 8) + "..." + val.slice(val.length - 8)
                        : val.slice(0);
            }
        }
        acc.push(obj);
        return acc;
    }, []).sort((a, b) => {
        const x = a["pk"].slice(0, 4);
        const y = b["pk"].slice(0, 4);
        if (x < y) {
            return -1;
        } else if (x > y) {
            return 1;
        } else {
            if (x === "CONN") {
                if (a["sk"] < b["sk"]) return -1;
                if (a["sk"] >= b["sk"]) return 1;
            } else if (x === "GAME") {
                if (a["sk"] < b["sk"]) {
                    return -1;
                } else if (a["sk"] > b["sk"]) {
                    return 1;
                } else {
                    if (a["players"] < b["players"]) return -1;
                    if (a["players"] >= b["players"]) return 1;
                }
            }
        }
    });
    console.table(data);
});
