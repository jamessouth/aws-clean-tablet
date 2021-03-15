import React from "react";
import ReactDOM from "react-dom";
import App from "./App";

import "./styles/index.css";

if (CSS.paintWorklet) {
    CSS.paintWorklet.addModule('paint.js');
}

const ce = React.createElement;

ReactDOM.render(
    ce(
        React.Fragment,
        null,
        ce(
            "h1",
            {
                className: "text-6xl mt-11 text-center font-arch decay-mask",
            },
            "CLEAN TABLET"
        ),
        ce(App, null)
    ),
    document.querySelector("#root")
);
