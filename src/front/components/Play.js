import React, { useEffect, useState, useContext } from "react";
import { Prompt, useParams } from "react-router-dom";
import { authContext } from "./ProvideAuth";
import { wsContext } from "./ProvideWS";
import Scoreboard from "./Scoreboard";

const ce = React.createElement;
export default function Play({history: {action}, user}) {
    const {
        connectedWS,
        game,
        send,
        wsError
    } = useContext(wsContext);

    const [count, setCount] = useState(5);

    console.log("playyyyyy ", game, action, user);

    // console.log('props: ', history, location);
    const { gameno } = useParams();

    return ce(
        "div",
        null,
        ce(
            Scoreboard,
            {
                playerName: user,
                players: game.players,
            }
        ),
        ce(Prompt, {
            when: true,
            message: "You will be ejected from the game!",
        })
    );
}
