import { useEffect, useState, useReducer } from "react";
import { initialState, reducer } from '../reducers/appState';
import { Auth } from "@aws-amplify/auth";

let ws;

export default function useWSState() {
    const [connectedWS, setConnectedWS] = useState(false);
    const [wsError, setWSError] = useState();
    const [token, setToken] = useState("");
    // const [count, setCount] = useState(null);
    const [game, setGame] = useState(null);
    const [ingame, setInGame] = useState(false);
    const [leadertoken, setLeadertoken] = useState("");
    const [playing, setPlaying] = useState(false);
    const [
        {
          games,
        },
        dispatch
      ] = useReducer(reducer, initialState);

    console.log('wsstate: ', connectedWS, wsError, !!token && token[0], !!games && games, ingame);
    
    useEffect(() => {
        console.log('usewsst111: ', );
        async function getToken() {
            const user = await Auth.currentUserPoolUser();
            const {
                signInUserSession: {
                    accessToken: { jwtToken },
                },
            } = user;
            // console.log("nnnn: ", jwtToken);
            setToken(jwtToken);
        }
        getToken();
    }, []);
    
    useEffect(() => {
        console.log('usewsst222: ', );
        if (token) {
            console.log('usewsst333: ', );
            ws = new WebSocket(`wss://${process.env.CT_APIID}.execute-api.${process.env.CT_REGION}.amazonaws.com/${process.env.CT_STAGE}?auth=${token}`);
            console.log("pojoihuh", token[0]);

            ws.addEventListener(
                "open",
                function (e) {
                    setConnectedWS(true);
                    console.log(e, Date.now());
                },
                false
            );//note: remove listeners????

            ws.addEventListener(
                "error",
                function (e) {
                    // setConnectedWS(false);
                    // console.log('eeee: ', e);
                    setWSError(e);
                    console.log(e, Date.now());
                },
                false
            );
            ws.addEventListener(
                "message",
                function (e) {
                    const {
                        type,
                        game,
                        games,
                        ingame,
                        leadertoken,
                        playing,
                  
                        // winners,
                        // word
                    } = JSON.parse(e.data);

                    console.log("wsstate json parse", type, game, games, ingame, leadertoken, playing);

                    switch (type) {
                        case "games":
                            dispatch({ type: "games", games });
                            break;
                        case "playing":
                            setGame(game);
                            break;
                        case "user":
                            setInGame(ingame);
                            setLeadertoken(leadertoken);
                            setPlaying(playing);
                            break;
                        default:
                            console.log('no case found: ', e.data);
                    }
                },
                false
            );

            ws.addEventListener(
                "close",
                function (e) {
                    // setConnectedWS(false);
                    // setError(false);
                    console.log(e, Date.now());
                },
                false
            );
            return function cleanup() {
                console.log("cleanup");
                setConnectedWS(false);
                setWSError(null);
                setToken("");
                ws.close(1000);
            };
        }
    }, [token]);


    function send(obj) {
        // if (!hasJoined) {
        //   ws.send(JSON.stringify({
        //     name: text,
        //   }));
        // } else {
        //   setAnswered(true);
        //   setSubmitSignal(false);
        //   setShowSVGTimer(false);
          ws.send(JSON.stringify(obj));
        // }
    }

    return {
        connectedWS,
        game,
        games,
        ingame,
        leadertoken,
        playing,
        send,
        wsError
    };
}
