import React, { createContext } from "react";
import useWSState from "../hooks/useWSState";

export const wsContext = createContext();

export default function ProvideWS({
  children
}) {
  const ws = useWSState();
  return React.createElement(wsContext.Provider, {
    value: ws
  }, children);
}
