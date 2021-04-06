import React, { createContext, useContext } from "react";
import useWSState from "../hooks/useWSState";
import { authContext } from "./ProvideAuth";

export const wsContext = createContext();

export default function ProvideWS({
  children
}) {
  let auth = useContext(authContext);
  console.log('provWS: ', auth.user);
  const ws = auth.user && useWSState();
  return React.createElement(wsContext.Provider, {
    value: ws
  }, children);
}
