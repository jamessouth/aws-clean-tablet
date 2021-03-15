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
                className: "mb-8 w-10/12 mx-auto",
            },
            ce(
                "p",
                {
                    className: "text-xs"
                },
                `${g.no}`
            ),
            g.players.map(s => ce(
                "p",
                {
                    key: s,
                    className: "text-right",
                },
                s.split("#", 1)[0]
            )),
            ce(
                "button",
                {
                    className: "w-full h-full",
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

