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
        leadertoken,
        send,
        wsError
    } = useContext(wsContext);

    const [count, setCount] = useState(5);
    const [stopCount, setStopCount] = useState(false);

    useEffect(() => {
        let id;
        if (game.playing) {
            id = setInterval(() => {
                setCount((c) => c - 1);
            }, 1000);
        }

        if (stopCount) {
            clearInterval(id);
        };
        return () => {
            clearInterval(id);
        };
    }, [game.playing, stopCount]);

    useEffect(() => {
        console.log("cnt play: ", count);
        if (count === 0) {
            setStopCount(true);
            if (game.leader === leadertoken) {
                send({
                    action: "play",
                    gameno: game.no,
                    type: "start",
                });
            }
        }
    }, [count, game.leader, leadertoken]);

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



        game.playing
        ? ce(
              "p",
              {
                  className:
                      "absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4",
              },
              `${count}`
          )
        : null,


        ce(Prompt, {
            when: true,
            message: "You will be ejected from the game!",
        })
    );
}
