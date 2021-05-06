import React, { useEffect, useState } from "react";
import { Switch, Route, Redirect, Link, useLocation } from "react-router-dom";

const ce = React.createElement;
const chk = String.fromCharCode(10003);
// const dot = String.fromCharCode(8228);

export default function Game({ game, ingame, leadertoken, send, user }) {
    const gameReady = !!game.leader;
    const leaderName = game.leader.split("_")[0];
    // const [connectedWS, setConnectedWS] = useState(false);
    const [ready, setReady] = useState(true);
    const [count, setCount] = useState(500);
    const [leader, setLeader] = useState(false);
    const [startGame, setStartGame] = useState(false);

    const chkstyl = " text-2xl font-bold leading-3";

    console.log("game.l: ", game.leader, "ldrtkn", leadertoken);

    useEffect(() => {
        if (game.leader !== "" && game.leader === leadertoken) {
            setLeader(true);
        }
    }, [game.leader, leadertoken]);

    useEffect(() => {
        let id;
        if (gameReady && game.no === ingame) {
            id = setInterval(() => {
                setCount((c) => c - 1);
            }, 1000);
        }

        return () => {
            setCount(500);
            clearInterval(id);
        };
    }, [gameReady, game.no, ingame]);

    useEffect(() => {
        console.log("cntfffff: ", count);
        if (ingame === game.no && count === 0) {
            setStartGame(true);
            if (leader) {
                send({
                    action: "play",
                    gameno: game.no,
                    type: "start",
                });
            }
        }
    }, [count, game.no, ingame, leader]);

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
            `${game.no}`
        ),
        ce(
            "p",
            {
                className: "text-xs col-span-2",
            },
            "players"
        ), // ["aaa", "bbb", "ccc", "ddd", "fff", "zzz", "ooo", "ttt"]
        game.players.map((p) => {
            return ce(
                "p",
                {
                    key: p.connid.S,
                },
                p.name.S,
                p.ready.BOOL
                    ? ce(
                          "span",
                          {
                              className: leaderName === p.name.S ? `text-red-200${chkstyl}` : `text-green-200${chkstyl}`
                          },
                          chk
                      )
                    : null
            );
        }),

        gameReady && ingame !== game.no
            ? ce(
                  "p",
                  {
                      className:
                          "absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4",
                  },
                  "Starting soon..."
              )
            : gameReady && ingame === game.no
            ? ce(
                  "p",
                  {
                      className:
                          "absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4",
                  },
                  `${count}`
              )
            : null,

        ce(
            "button",
            {
                className:
                    "w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70",
                disabled:
                    (!!ingame && ingame !== game.no) ||
                    (!ingame && game.players.length > 7),
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
                    game.players.length < 3,
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
