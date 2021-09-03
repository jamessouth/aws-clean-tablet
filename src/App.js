import React, { createContext, useEffect, useState } from "react";

// import { arch, div, h1, winner } from './styles/index.css';
import {
  BrowserRouter as Router,
  Link,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";

import ProvideAuth from "./components/ProvideAuth";
// import AuthButton from "./components/AuthButton";

import Routes from "./components/Routes";


const ce = React.createElement;



export default function App() {
  return ce(
    ProvideAuth,
    null,
    ce(
      Router,
      null,
      ce(
        "div",
        {
          className: "mt-8",
        },
        // ce(AuthButton),
        ce(Routes)
      )
    )
  );
}

// export default function App() {
//   return <div className="mt-8"></div>
// }
