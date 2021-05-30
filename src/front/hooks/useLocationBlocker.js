import { useHistory } from "react-router-dom";
import { useEffect } from "react";

export default function useLocationBlocker() {
  const history = useHistory();
  useEffect(() =>
      history.block(
        (location, action) => {
            console.log('fret: ', location, action, action !== "PUSH" ||
            getLocationId(location) !== getLocationId(history.location), getLocationId(history.location), getLocationId(location));
          return action !== "PUSH" ||
          getLocationId(location) !== getLocationId(history.location)
        }), []);
}

function getLocationId({ pathname, search, hash }) {
  return pathname + (search ? "?" + search : "") + (hash ? "#" + hash : "");
}