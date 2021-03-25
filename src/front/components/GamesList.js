import React, { useEffect, useState } from "react";


const ce = React.createElement;


export default function GamesList({games, ingame, send}) {
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
                className: "mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8",
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
            g.players.map(s => ce(
            // ["aaa", "bbb", "ccc", "ddd", "fff", "zzz", "ooo", "ttt"]
                "p",
                {
                    key: s,
                    // className: "text-right",
                },
                s.split("#", 1)[0]
            )),
            ce(
                "button",
                {// border border-white border-solid
                    className: "w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                    disabled: (!!ingame && ingame !== g.no) || (!ingame && g.players.length > 7),
                    onClick: () => {
                        send({
                            action: "lobby",
                            game: `${g.no}`,
                            type: !!ingame && ingame === g.no ? "leave" : "join",
                        });
                    },
                },
                !!ingame && ingame === g.no ? "leave" : "join"
            ),
            ce(
                "button",
                {// border border-white border-solid
                    className: "w-1/2 bottom-0 h-8 right-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                    disabled: (!!ingame && ingame !== g.no) || !ingame || g.players.length < 3,
                    onClick: () => {
                        send({
                            action: "lobby",
                            game: `${g.no}`,
                            type: !!ingame ? "leave" : "join",
                        });
                        
                    },
                },
                // !!ingame ? "leave" : "join"
                "ready"
            )
        ))
    );
}

