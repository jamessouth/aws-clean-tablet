import React, { useContext } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { authContext } from "./ProvideAuth";
//import { AuthState } from "@aws-amplify/ui-components";

//import {
//    withAuthenticator,
//    AmplifyAuthenticator,
//    AmplifyFormSection,
//    AmplifyForgotPassword,
//    AmplifyConfirmSignUp,
//    AmplifySignIn,
//    AmplifySignUp,
//} from "@aws-amplify/ui-react";

const ce = React.createElement;

export default function LoginPage() {
  let history = useHistory();
  let location = useLocation();
  let auth = useContext(authContext);
  let { from } = location.state || {
    from: {
      pathname: "/",
    },
  };

  const handleAuthChange = (authState, userData) => {
    console.log("login: ", authState, userData);
    auth.setAuthState(authState);

    if (authState === "signedin") {
      auth.setUser(userData.username);
      localStorage.setItem("user", userData.username);
    }

    return history.replace(from);
  };

  switch (auth.authState) {
    case "signup":
      return ce(AmplifySignUp, {
        handleAuthStateChange: handleAuthChange,
      });

    case "confirmSignUp":
      return ce(AmplifyConfirmSignUp, {
        handleAuthStateChange: handleAuthChange,
      });

    case "forgotpassword":
      return ce(AmplifyForgotPassword, {
        handleAuthStateChange: handleAuthChange,
      });

    default:
      return ce(AmplifySignIn, {
        handleAuthStateChange: handleAuthChange,
      });
  }
}
