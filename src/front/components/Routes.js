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
                    "div",
                    {
                        className: "w-3/5 border border-smoke-100 mb-28 mt-8 sm:mt-16",
                    },
                    ce(
                        Link,
                        {
                            to: "/lobby",
                            className:
                                "block font-luck decay-mask text-6xl leading-12rem sm:text-8xl sm:leading-12rem",
                        },
                        "ENTER"
                    )
                ),
                ce(
                    "div",
                    null,
                    ce(
                        Link,
                        {
                            to: "/leaderboards",
                        },
                        "Leaderboards"
                    )
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
