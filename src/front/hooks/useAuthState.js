import { useState } from "react";
import { AuthState } from "@aws-amplify/ui-components";

export default function useAuthState() {
    const [user, setUser] = useState(localStorage.getItem("user"));
    const [authState, setAuthState] = useState(user ? AuthState.SignedIn : AuthState.SignIn);

    return {
        authState,
        setAuthState,
        user,
        setUser,
    };
}
