import React, { useEffect, useState } from "react";


const ce = React.createElement;


export default function GamesList({games}) {
    // const [connectedWS, setConnectedWS] = useState(false);
    // const [games, setGames] = useState(null);
    // const [startedNewGame, setStartedNewGame] = useState(false);
    // const [token, setToken] = useState("");
    // const [wsError, setWSError] = useState();
    
// console.log('gamesss: ', Array.isArray(games));
    // useEffect(() => {

    // }, []);

    // useEffect(() => {

    // }, []);

    function send(text) {
        ws.send(JSON.stringify({
            action: text,
        }));
    }

    return ce(
        "ul",
        {
            className: "mx-auto mb-10 w-10/12"
        },
        games.map(g => ce(
            "li",
            {
                key: g.no,
                className: "mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-6 relative pb-6",
            },
            ce(
                "p",
                {
                    className: "text-xs col-span-2"
                },
                `${g.no}`
            ),
            ce(
                "p",
                {
                    className: "text-xs col-span-2"
                },
                "players"
            ),
            // g.players
            ["aaa", "bbb", "ccc", "ddd", "fff", "zzz", "ooo", "ttt"].map(s => ce(
                "p",
                {
                    key: s,
                    // className: "text-right",
                },
                s
                // .split("#", 1)[0]
            )),
            ce(
                "button",
                {
                    className: "w-full absolute pt-36",
                    onClick: () => {
                        send({
                            action: "lobby",
                            game: `${g.no}`,
                          });
                    },
                },
                "join"
            )
        ))
    );
}

