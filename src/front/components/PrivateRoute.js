import React, { useContext } from "react";
import { Route, Redirect } from "react-router-dom";
import { authContext } from "./ProvideAuth";

const ce = React.createElement;

export default function PrivateRoute({ children, ...rest }) {
    console.log('ccxxxx: ', children);
    let auth = useContext(authContext);
    return ce(Route, {
        ...rest,
        render: ({ location, history }) =>
            auth.user
                ? ce(children, {history, user: auth.user})
                : ce(Redirect, {
                      to: {
                          pathname: "/login",
                          state: {
                              from: location,
                          },
                      },
                  }),
    });
}
