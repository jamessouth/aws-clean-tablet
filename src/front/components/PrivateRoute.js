import React, { useContext } from "react";
import { Route, Redirect } from "react-router-dom";
import { authContext } from "./ProvideAuth";

const ce = React.createElement;

export default function PrivateRoute({ children, ...rest }) {
    let auth = useContext(authContext);
    return ce(Route, {
        ...rest,
        render: ({ location }) =>
            auth.user
                ? children
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
