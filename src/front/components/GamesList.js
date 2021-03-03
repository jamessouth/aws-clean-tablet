import React, { useEffect, useState } from "react";


const ce = React.createElement;


export default function GamesList({games}) {
    // const [connectedWS, setConnectedWS] = useState(false);
    // const [games, setGames] = useState(null);
    // const [startedNewGame, setStartedNewGame] = useState(false);
    // const [token, setToken] = useState("");
    // const [wsError, setWSError] = useState();
    
console.log('gamesss: ', Array.isArray(games));
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
            className: "m-auto w-10/12"
        },
        games.map(g => ce(
            "li",
            {
                key: g.conn,
                className: "mb-8",
            },
            ce(
                "button",
                {
                    className: "w-full h-full"
                },
                ce(
                    "p",
                    {
                        className: "text-xs"
                    },
                    `${g.type}#${g.no}`
                ),
                ce("p", null, g.name)
            )
        ))
    );
}

