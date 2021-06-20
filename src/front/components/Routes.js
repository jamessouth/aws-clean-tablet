import React, { useContext } from "react";
import { Switch, Route, Redirect, Link, useLocation } from "react-router-dom";
import LoginPage from "./LoginPage";
import Comp from "./Comp";
import Play from "./Play";
import PrivateRoute from "./PrivateRoute";
import Lobby from "./Lobby";
import { authContext } from "./ProvideAuth";
import ProvideWS from "./ProvideWS";
import AuthButton from "./AuthButton";

const ce = React.createElement;

export default function Routes() {
    let auth = useContext(authContext);

    return ce(
        Switch,
        null,
        ce(
            Route,
            {
                path: "/",
                exact: true,
            },
            auth.user
                ? ce(Redirect, { to: "/lobby" })
                : ce(
                      "div",
                      {
                          className: "flex flex-col items-center",
                      },
                      ce(
                          Link,
                          {
                              to: "/lobby",
                              className:
                                  "w-3/5 border border-smoke-100 block font-fred decay-mask text-5xl leading-12rem sm:mt-16 sm:text-8xl sm:leading-12rem",
                          },
                          "ENTER"
                      ),
                      ce(
                          Link,
                          {
                              to: "/leaderboards",
                              className:
                                  "w-3/5 border border-smoke-100 mb-28 mt-10 block text-xl sm:mt-16 sm:text-2xl",
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
            ce(LoginPage)
        ),
        ce(
            ProvideWS,
            null,
            ce(AuthButton),
            ce(
                PrivateRoute,
                {
                    path: "/lobby",
                },
                Lobby//"div", null, ce()
            ),
            ce(
                PrivateRoute,
                {
                    path: "/game/:gameno",
                },
                Play
            )
        )
    );
}

// ce(
//     Link,
//     {
//         to: "/leaderboards",
//     },
//     "Leaderboards"
// )
