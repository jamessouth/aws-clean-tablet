import React, { useContext, useEffect, useState } from "react";
import GamesList from "./GamesList";
import { wsContext } from "./ProvideWS";

const ce = React.createElement;


export default function Lobby({history: {action}, user}) {
    const {
        connectedWS,
        games,
        ingame,
        leadertoken,
        playing,
        send,
        wsError
    } = useContext(wsContext);

    // const [startedNewGame, setStartedNewGame] = useState(false);
    
    
    console.log("lobbbbbbbb ", games, ingame, action, user);


    const startBtnStyles = " mx-auto mb-8 w-1/2 bg-smoke-100 text-gray-700";

    return wsError
        ? ce("p", null, "not connected: connection error")
        : !connectedWS || !games
        ? ce(
            "p",
            null,
            "loading games..."
        )
        : ce(
            "div",
            {
                className: "flex flex-col mt-8",
            },
            ce(
                "button",
                {
                    className: !ingame ? `visible${startBtnStyles}` : `invisible${startBtnStyles}`,
                    type: "button",
                    onClick: () => {
                        // setStartedNewGame(true);
                        send({
                            action: "lobby",
                            game: "new",
                            type: "join",
                        });
                    },
                },
                "start a new game"
            ),
            games.length < 1
            ? ce(
                "p",
                null,
                "no games found. start a new one!"
            )
            : ce(
                GamesList,
                {
                    action,
                    games,
                    ingame,
                    leadertoken,
                    send: val => send(val),
                    user,
                }
            )
        );
}
