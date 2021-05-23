import React, { useEffect, useState, useContext } from "react";
import { Prompt, useParams } from "react-router-dom";
import { authContext } from "./ProvideAuth";
import { wsContext } from "./ProvideWS";
import Scoreboard from "./Scoreboard";
import Word from "./Word";
import Form from "./Form";

const circ = String.fromCharCode(9862);
const ce = React.createElement;
export default function Play({history: {action}, user}) {
    const {
        connectedWS,
        game,
        leadertoken,
        playerColor,
        send,
        wsError,
        currentWord,
        previousWord,
    } = useContext(wsContext);

   
    
    const ANSWER_MAX_LENGTH = 12;// see also app.go

    const [answered, setAnswered] = useState(false);
    const [inputText, setInputText] = useState('');
    // const [hideWord, setHideWord] = useState(false);

    function sendAnswer() {
        send({
            action: "play",
            gameno: game.no,
            answer: inputText.slice(0, ANSWER_MAX_LENGTH),
            type: "game",
        });
        setAnswered(true);
        // setHideWord(true);
        setInputText("");
    }

    // const [count, setCount] = useState(5);
    // const [stopCount, setStopCount] = useState(false);

    useEffect(() => {
        setAnswered(false);
        // setHideWord(false);
    }, [currentWord]);

    // useEffect(() => {
    //     console.log("cnt play: ", count);
    //     if (count === 0) {
    //         setStopCount(true);
    //         if (game.leader === leadertoken) {
    //             send({
    //                 action: "play",
    //                 gameno: game.no,
    //                 type: "start",
    //             });
    //         }
    //     }
    // }, [count, game.leader, leadertoken]);

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



        game.playing && !(currentWord || previousWord)
        ? ce(
              "p",
              {
                  className: "text-yellow-200 text-2xl font-bold ",
              },
              "Get Ready",
              ce(
                  "span",
                  {
                        className: "animate-spin"
                  },
                  circ
              )
          )
        : null,

        ce(
            Word,
            {
                className: answered ? "animate-erase" : "",
                onAnimationEnd: () => sendAnswer(),
                playerColor,
                currentWord
            }
        ),


        ce(
            Form,
            {
                ANSWER_MAX_LENGTH,
                answered,
                inputText,
                onEnter: () => sendAnswer(),
                setInputText,
            }
        ),


        ce(Prompt, {
            when: true,
            message: "You will be ejected from the game!",
        })
    );
}
