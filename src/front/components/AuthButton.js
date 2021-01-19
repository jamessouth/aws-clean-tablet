import React, { useContext } from "react";
import { useHistory } from "react-router-dom";
import { authContext } from "../App";

import { AmplifySignOut } from "@aws-amplify/ui-react";

export default function AuthButton() {
    let history = useHistory();
    let auth = useContext(authContext);

    const handleAuthChange = (authState, userData) => {
        console.log("chgout b4", authState, userData);
        auth.setAuthState(authState);
        auth.setUser(userData);
        localStorage.removeItem("user");
        return history.push("/");
    };

    const ce = React.createElement;

    return auth.user
        ? ce(
              React.Fragment,
              null,
              ce("p", null, "Welcome! " + auth.user),
              ce(AmplifySignOut, {
                  handleAuthStateChange: handleAuthChange,
              })
          )
        : null;
}
