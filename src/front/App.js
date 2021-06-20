import Amplify, { Auth } from "@aws-amplify/auth";
import awsExports from "../aws-exports";
import React, { createContext, useEffect, useState } from "react";

// import { arch, div, h1, winner } from './styles/index.css';
import {
    BrowserRouter as Router,
    Link,
    Switch,
    Route,
    Redirect,
} from "react-router-dom";

import ProvideAuth from "./components/ProvideAuth";
// import AuthButton from "./components/AuthButton";

import Routes from "./components/Routes";

Amplify.configure(awsExports);

const ce = React.createElement;



export default function App() {
    return ce(
        ProvideAuth,
        null,
        ce(
            Router,
            null,
            ce(
                "div",
                {
                    className: "mt-8",
                },
                // ce(AuthButton),
                ce(Routes)
            )
        )
    );
}
