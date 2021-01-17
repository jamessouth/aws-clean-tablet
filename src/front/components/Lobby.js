import React, { useEffect, useState } from "react";
// import {
//     withAuthenticator,
//     AmplifyAuthenticator,
//     AmplifySignOut,
//     AmplifyAuthFields,
//     AmplifySignUp,
// } from "@aws-amplify/ui-react";
// import { useHistory } from "react-router-dom";

const ce = React.createElement;

export default function Lobby() {
    const [connected, setConnected] = useState(false);
    // const [error, setError] = useState("");
    console.log("lobbbbbbbb ", connected);

    useEffect(() => {
        const ws = new WebSocket(process.env.CT_WS);
        console.log("pojoihuh", Date.now());

        ws.addEventListener(
            "open",
            function (e) {
                setConnected(true);
                console.log(e, Date.now());
            },
            false
        );

        ws.addEventListener(
            "error",
            function (e) {
                // setConnected(false);
                // console.log('eeee: ', e);
                // setError(e.message);
                console.log(e, Date.now());
            },
            false
        );

        ws.addEventListener(
            "close",
            function (e) {
                setConnected(false);
                // setError(false);
                console.log(e, Date.now());
            },
            false
        );

        return function cleanup() {
            console.log("cleanup");
            setConnected(false);
            // setError("");
            ws.close(1000);
        };
    }, []);

    // error
        // ? ce("p", null, "connection error, please try again")
        // : 
    return connected
        ? ce(
              React.Fragment,
              null,
              ce(
                  "div",
                  {
                      className: "flex flex-col mt-8",
                  },
                  ce(
                      "button",
                      {
                          className:
                              "mx-auto mb-8 h-40 w-1/2 bg-smoke-100 text-gray-700",
                      },
                      "start a new game"
                  ),
                  ce(
                      "button",
                      {
                          className:
                              "mx-auto mb-8 h-40 w-1/2 bg-gray-100 text-gray-700",
                      },
                      "join"
                  ),
                  ce(
                      "button",
                      {
                          className:
                              "mx-auto mb-8 h-40 w-1/2 bg-gray-100 text-gray-700",
                      },
                      "join"
                  ),
                  ce(
                      "button",
                      {
                          className:
                              "mx-auto mb-8 h-40 w-1/2 bg-gray-100 text-gray-700",
                      },
                      "join"
                  )
              )
          )
        : ce("p", null, "not connected: connection error");
}
