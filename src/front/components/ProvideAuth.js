import React, { createContext } from "react";
import useAuthState from "../hooks/useAuthState";

export const authContext = createContext();

export default function ProvideAuth({
  children
}) {
  const auth = useAuthState();
  return React.createElement(authContext.Provider, {
    value: auth
  }, children);
}
