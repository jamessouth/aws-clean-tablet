import React, { useContext } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { authContext } from "./ProvideAuth";
import { wsContext } from "./ProvideWS";




const ce = React.createElement;
export default function Comp() {
  const {
    connectedWS,
    games,
    ingame,
    send,
    wsError
} = useContext(wsContext);

  // console.log('props: ', history, location);




  return ce(
    "p",
    null,
    "gggghhhhuuuuhh" + ingame
  );
}
