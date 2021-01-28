import React, { useEffect, useState } from "react";
import { Auth } from "@aws-amplify/auth";

const ce = React.createElement;
export default function Lobby() {
    const [connectedWS, setConnectedWS] = useState(false);
    const [token, setToken] = useState("");
    console.log("lobbbbbbbb ", connectedWS);

    useEffect(() => {
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
        if (token) {
            const ws = new WebSocket(`${process.env.CT_WS}?auth=${token}`);
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
                    // setError(e.message);
                    console.log(e, Date.now());
                },
                false
            );
            ws.addEventListener(
                "message",
                function (e) {
                    // setConnectedWS(false);
                    // console.log('eeee: ', e);
                    // setError(e.message);
                    console.log("mmmm", e);
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
                // setError("");
                ws.close(1000);
            };
        }
    }, [token]);

    // error
    // ? ce("p", null, "connection error, please try again")
    // :
    return connectedWS
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
