import React, { useEffect, useState } from "react";

const ce = React.createElement;
const chk = String.fromCharCode(10003);

export default function GamesList({ games, ingame, send }) {
    // const [connectedWS, setConnectedWS] = useState(false);
    const [ready, setReady] = useState(true);
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
            className: "mx-auto mb-10 w-10/12",
        },
        games.map((g) =>
            ce(
                "li",
                {
                    key: g.no,
                    className:
                        "mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8",
                },
                ce(
                    "p",
                    {
                        className: "text-xs col-span-2",
                    },
                    `${g.ready}`
                ),
                ce(
                    "p",
                    {
                        className: "text-xs col-span-2",
                    },
                    "players"
                ), // ["aaa", "bbb", "ccc", "ddd", "fff", "zzz", "ooo", "ttt"]
                Object.keys(g.players).map((p) => {
                    const plr = g.players[p].M;
                    return ce(
                        "p",
                        {
                            key: plr.connid.S,
                        },
                        plr.name.S,
                        plr.ready.BOOL
                        ? ce(
                                "span",
                                {
                                    className: "text-green-200 text-2xl font-bold leading-3",
                                },
                                chk
                            )
                        : null
                    );
                }),
                ce(
                    "button",
                    {
                        className:
                            "w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                        disabled:
                            (!!ingame && ingame !== g.no) ||
                            (!ingame && Object.keys(g.players).length > 7),
                        onClick: () => {
                            send({
                                action: "lobby",
                                game: `${g.no}`,
                                type:
                                    !!ingame && ingame === g.no
                                        ? "leave"
                                        : "join",
                            });
                            if (!!ingame && ingame === g.no) {
                                setReady(true);
                            }
                        },
                    },
                    !!ingame && ingame === g.no ? "leave" : "join"
                ),
                ce(
                    "button",
                    {
                        className:
                            "w-1/2 bottom-0 h-8 right-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                        disabled:
                            (!!ingame && ingame !== g.no) ||
                            !ingame ||
                            Object.keys(g.players).length < 3,
                        onClick: () => {
                            send({
                                action: "lobby",
                                game: `${g.no}`,
                                type: "ready",
                                value: ready,
                            });
                            setReady(!ready);
                        },
                    },
                    ready ? "ready" : "not ready"
                )
            )
        )
    );
}
