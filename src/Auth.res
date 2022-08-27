type propShape = {
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setToken": (. option<string> => option<string>) => unit,
  // "subpath": list<string>,
  "token": option<string>,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make = (~token, ~setToken, ~cognitoUser, ~setCognitoUser) => {
  // ~subpath,

  let {path} = RescriptReactRouter.useUrl()
  Js.log2("u345876l", path)

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
  <>
    {switch path {
    | list{"auth", "leaderboard"} => React.null
    | _ =>
      <header className="mb-10 newgmimg:mb-12">
        <p className="font-flow text-stone-100 text-4xl h-10 font-bold text-center">
          {React.string(playerName)}
        </p>
        <h1
          style={ReactDOM.Style.make(~backgroundColor={playerColor}, ())}
          className="text-6xl mt-11 mx-auto px-6 text-center font-arch decay-mask text-stone-100">
          {React.string("CLEAN TABLET")}
        </h1>
      </header>
    }}
    {switch path {
    | list{"auth", "lobby"} =>
      switch wsConnected {
      | false => {
          body(document)->setClassName("bodchmob bodchtab bodchbig")
          <React.Suspense fallback=React.null> loading1 </React.Suspense>
        }

      | true => {
          body(document)->classList->removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")

          <React.Suspense fallback=React.null> lobby </React.Suspense>
        }
      }
    | list{"auth", "play", gameno} =>
      switch wsConnected {
      | true =>
        switch Js.Array2.length(players) > 0 && gameno == sk {
        | true => <React.Suspense fallback=React.null> play </React.Suspense>

        | false => <React.Suspense fallback=React.null> loading2 </React.Suspense>
        }

      | false =>
        <p className="text-center text-stone-100 font-anon text-lg">
          {React.string("not connected...")}
        </p>
      }

    | list{"auth", "leaderboard"} => {
        body(document)->classList->addClassList3("bodleadmob", "bodleadtab", "bodleadbig")

        <React.Suspense fallback=React.null> leaders </React.Suspense>
      }

    | _ => <div> {React.string("other")} </div> // <PageNotFound/>
    }}
  </>
}
