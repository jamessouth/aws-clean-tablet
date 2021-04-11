import React, { useEffect, useState } from "react";
import { Switch, Route, Redirect, Link, useLocation } from "react-router-dom";

const ce = React.createElement;
const chk = String.fromCharCode(10003);

export default function Game({ game, ingame, send }) {
    const gameReady = game.ready;
    // const [connectedWS, setConnectedWS] = useState(false);
    const [ready, setReady] = useState(true);
    const [count, setCount] = useState(5);
    // const [startedNewGame, setStartedNewGame] = useState(false);
    // const [token, setToken] = useState("");
    const [startGame, setStartGame] = useState(false);

    // console.log('gamesss: ', Array.isArray(games));

    useEffect(() => {
        let id;
        if (gameReady) {
            id = setInterval(() => {
                setCount((c) => c - 1);
            }, 1000);
        }

        return () => {
            setCount(5);
            clearInterval(id);
        };
    }, [gameReady]);

    useEffect(() => {
        console.log("cntfffff: ", count);
        if (ingame === game.no && count === 0) {
            
            send({
                action: "play",
                game: `${game.no}`,
                type: "start",
            });
            
            setStartGame(true);
        }
    }, [count]);

    return ce(
        "li",
        {
            className:
                "mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8",
        },

        startGame ? ce(Redirect, { to: `game/${game.no}`, push: true }) : null,

        ce(
            "p",
            {
                className: "text-xs col-span-2",
            },
            `${game.ready}`
        ),
        ce(
            "p",
            {
                className: "text-xs col-span-2",
            },
            "players"
        ), // ["aaa", "bbb", "ccc", "ddd", "fff", "zzz", "ooo", "ttt"]
        Object.keys(game.players).map((p) => {
            const plr = game.players[p].M;
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
                              className:
                                  "text-green-200 text-2xl font-bold leading-3",
                          },
                          chk
                      )
                    : null
            );
        }),
        ce(
            "p",
            {
                className:
                    "absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4",
            },
            `${count}`
        ),
        ce(
            "button",
            {
                className:
                    "w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                disabled:
                    (!!ingame && ingame !== game.no) ||
                    (!ingame && Object.keys(game.players).length > 7),
                onClick: () => {
                    send({
                        action: "lobby",
                        game: `${game.no}`,
                        type: !!ingame && ingame === game.no ? "leave" : "join",
                    });
                    if (!!ingame && ingame === game.no) {
                        setReady(true);
                    }
                },
            },
            !!ingame && ingame === game.no ? "leave" : "join"
        ),
        ce(
            "button",
            {
                className:
                    "w-1/2 bottom-0 h-8 right-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                disabled:
                    (!!ingame && ingame !== game.no) ||
                    !ingame ||
                    Object.keys(game.players).length < 3,
                onClick: () => {
                    send({
                        action: "lobby",
                        game: `${game.no}`,
                        type: "ready",
                        value: ready,
                    });
                    setReady(!ready);
                },
            },
            ready ? "ready" : "not ready"
        )
    );
}
