import React from "react";
import { Switch, Route, Link, useLocation } from "react-router-dom";
import LoginPage from "./LoginPage";
import Comp from "./Comp";
import PrivateRoute from "./PrivateRoute";

const ce = React.createElement;

export default function Routes() {
    return ce(
        Switch,
        null,
        ce(
            Route,
            {
                path: "/",
                exact: true,
            },
            ce(
                "div",
                {
                    className: "flex flex-col items-center",
                },
                    ce(
                        Link,
                        {
                            to: "/lobby",
                            className: "border border-blue-300 w-3/5 h-64 mb-28 font-luck decay-mask text-6xl",
                        },
                        "ENTER"
                    ),
                    ce(
                        Link,
                        {
                            to: "/leaderboards",
                        },
                        "Leaderboards"
                    )
            )
        ),
        ce(
            Route,
            {
                path: "/leaderboards",
            },
            ce(Comp)
        ),
        ce(
            Route,
            {
                path: "/login",
            },
            ce(LoginPage, null)
        ),
        ce(
            PrivateRoute,
            {
                path: "/lobby",
            },
            ce(
                "h3",
                null,
                "Games",
                ce(
                    Link,
                    {
                        to: "/leaderboards",
                    },
                    "Leaderboards"
                )
            )
        )
    );
}
