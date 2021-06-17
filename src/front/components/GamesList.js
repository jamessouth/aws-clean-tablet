import React, { useEffect, useState } from "react";
import Game from "./Game";

const ce = React.createElement;
// const chk = String.fromCharCode(10003);

export default function GamesList({ action, games, ingame, leadertoken, send, user }) {
    // const [connectedWS, setConnectedWS] = useState(false);
    // const [ready, setReady] = useState(true);
    // const [startedNewGame, setStartedNewGame] = useState(false);
    // const [token, setToken] = useState("");
    // const [wsError, setWSError] = useState();

    // console.log('gamesss: ', Array.isArray(games));
    // useEffect(() => {

    // }, []);

    // useEffect(() => {

    // }, []);

    

    return ce(
        "ul",
        {
            className: "mx-auto mb-10 w-10/12",
        },
        games.map((game) =>
            ce(
                Game,
                {
                    action,
                    key: game.no,
                    game,
                    ingame,
                    leadertoken,
                    send,
                    user,
                }
            )
        )
    );
}
