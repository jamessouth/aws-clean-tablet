import React, { useEffect, useState } from "react";


const ce = React.createElement;


export default function GamesList({games}) {
    // const [connectedWS, setConnectedWS] = useState(false);
    // const [games, setGames] = useState(null);
    // const [startedNewGame, setStartedNewGame] = useState(false);
    // const [token, setToken] = useState("");
    // const [wsError, setWSError] = useState();
    
console.log('gamesss: ', Array.isArray(games));
    // useEffect(() => {

    // }, []);

    // useEffect(() => {

    // }, []);

    // function send() {

    //   }

    return ce(
        "ul",
        null,
        games.map(g => ce(
            "li",
            {key: g.conn},
            ce(
                "button",
                null,
                ce("p", null, `${g.type}#${g.no}`),
                ce("p", null, g.name)
            )
        ))
    );
}

