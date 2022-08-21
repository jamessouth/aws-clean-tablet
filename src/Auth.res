// type propShape = {
//   "close": (. int, string) => unit,
//   "count": string,
//   "games": Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
//   "playerGame": string,
//   "send": (. option<string>) => unit,
//   "wsError": string,
//   "setLeaderData": (. array<Reducer.stat> => array<Reducer.stat>) => unit,
// }

// @val
// external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

// @module("react")
// external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
//   propShape,
// > = "lazy"

@react.component
let make = (~token, ~setToken, ~cognitoUser, ~setCognitoUser, ~setAppName, ~setAppColor, ~ppp) => {
  let {path} = RescriptReactRouter.useUrl()
Js.log3("u345876l", path, ppp)
  let initialState: Reducer.state = {
    gamesList: Js.Nullable.null,
    players: [],
    sk: "",
    oldWord: "",
    word: "",
    showAnswers: false,
    winner: "",
  }

  let (
    playerGame,
    playerName,
    playerColor,
    endtoken,
    count,
    wsConnected,
    players,
    sk,
    showAnswers,
    winner,
    oldWord,
    word,
    games,
    leaderData,
    setLeaderData,
    send,
    resetConnState,
    close,
    wsError,
  ) = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, initialState)


React.useEffect1(() => {
  setAppName(._ => playerName)
  None
}, [playerName])

React.useEffect1(() => {
  setAppColor(._ => playerColor)
  None
}, [playerColor])




  let load = Loading.lazy_(() =>
    Loading.import_("./Loading.bs")->Promise.then(comp => {
      Promise.resolve({"default": comp["make"]})
    })
  )

  let loading1 = React.createElement(load, Loading.makeProps(~label="games...", ()))

  let loading2 = React.createElement(load, Loading.makeProps(~label="game...", ()))

  let lobby = React.createElement(
    Lobby.lazy_(() =>
      Lobby.import_("./Lobby.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Lobby.makeProps(~playerGame, ~games, ~send, ~wsError, ~close, ~count, ~setLeaderData, ()),
  )

  let play = React.createElement(
    Play.lazy_(() =>
      Play.import_("./Play.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Play.makeProps(
      ~players,
      ~sk,
      ~showAnswers,
      ~winner,
      ~isWinner={winner != ""},
      ~oldWord,
      ~word,
      ~playerColor,
      ~send,
      ~playerName,
      ~endtoken,
      ~resetConnState,
      (),
    ),
  )

  let leaders = React.createElement(
    Leaders.lazy_(() =>
      Leaders.import_("./Leaders.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Leaders.makeProps(~leaderData, ~playerName, ()),
  )

open Web
  {
    switch path {
    | list{"lobby"} =>
      switch wsConnected {
      | false => {
          body(document)->setClassName("bodchmob bodchtab bodchbig")
          <React.Suspense fallback=React.null> loading1 </React.Suspense>
        }

      | true => {
          body(document)->classList->removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")
        <Protected token>
            <React.Suspense fallback=React.null> lobby </React.Suspense>
          </Protected>
        }
      }
    | list{"play", gameno} =>
      switch wsConnected {
      | true =>
        switch Js.Array2.length(players) > 0 && gameno == sk {
        | true =>
        <Protected token>
         <React.Suspense fallback=React.null> play </React.Suspense>
          </Protected>

        | false => <React.Suspense fallback=React.null> loading2 </React.Suspense>
        }

      | false =>
        <p className="text-center text-stone-100 font-anon text-lg">
          {React.string("not connected...")}
        </p>
      }

    | list{"leaderboard"} => {
        body(document)->classList->addClassList3("bodleadmob", "bodleadtab", "bodleadbig")

        <Protected token>

        <React.Suspense fallback=React.null> leaders </React.Suspense>
          </Protected>
      }


| _ => <div> {React.string("other")} </div> // <PageNotFound/>

    }
  }
}
