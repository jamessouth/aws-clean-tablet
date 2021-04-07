import React, { useContext } from "react";
import { Prompt, useParams } from "react-router-dom";
import { authContext } from "./ProvideAuth";
import { wsContext } from "./ProvideWS";

const ce = React.createElement;
export default function Play() {
    //   const {
    //     connectedWS,
    //     games,
    //     ingame,
    //     send,
    //     wsError
    // } = useContext(wsContext);

    // console.log('props: ', history, location);
    const { gameno } = useParams();

    return ce(
        "div",
        null,
        ce("p", null, "gggghhhhuuuuhh " + gameno),
        ce(Prompt, {
            when: true,
            message: "You will be ejected from the game!",
        })
    );
}
