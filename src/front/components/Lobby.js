import React, { useEffect } from "react";
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
    // const history = useHistory();
    console.log("lobbbbbbbb");

    useEffect(() => {
        const ws = new WebSocket(process.env.CT_WS);
        console.log("pojoihuh", Date.now());

        ws.addEventListener(
            "open",
            function (e) {
                console.log(e, Date.now());
            },
            false
        );

        ws.addEventListener(
            "error",
            function (e) {
                console.log(e, Date.now());
            },
            false
        );

        ws.addEventListener(
            "close",
            function (e) {
                console.log(e, Date.now());
            },
            false
        );

        return function cleanup() {
            console.log("cleanup");
            ws.close(1000);
        };
    }, []);

    return ce(
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
    );
}
