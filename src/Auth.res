@val @scope(("import", "meta", "env"))
external apiid: string = "VITE_APIID"
@val @scope(("import", "meta", "env"))
external region: string = "VITE_REGION"
@val @scope(("import", "meta", "env"))
external stage: string = "VITE_STAGE"
%%raw(`import './css/lobby.css'`)

type propShape = {
  "token": option<string>,
  "setToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setWsError": (. string => string) => unit,
  "route": Route.t,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

type listGamesData = {listGms: array<Reducer.listGame>, name: string}
type modConnData = {modConn: string, color: string, endtoken: string}
type countdownData = {cntdown: string}
type addGameData = {addGame: Reducer.listGame}
type modListGameData = {mdLstGm: Reducer.listGame}
type modPlayersData = {
  players: array<Reducer.livePlayer>,
  sk: string,
  showAnswers: bool,
  winner: string,
}
type wordData = {newword: string}
type rmvGameData = {rmvGame: Reducer.listGame}
type leadersData = {leaders: array<Reducer.stat>}
type msgType =
  | InsertConn
  | ModifyConn
  | InsertGame
  | ModifyListGame
  | Countdown
  | ModifyPlayers
  | Word
  | RemoveGame
  | Leaders
  | Other
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"
@scope("JSON") @val
external parseModConn: string => modConnData = "parse"
@scope("JSON") @val
external parseCountdown: string => countdownData = "parse"
@scope("JSON") @val
external parseAddGame: string => addGameData = "parse"
@scope("JSON") @val
external parseModListGame: string => modListGameData = "parse"
@scope("JSON") @val
external parseModPlayers: string => modPlayersData = "parse"
@scope("JSON") @val
external parseWord: string => wordData = "parse"
@scope("JSON") @val
external parseRmvGame: string => rmvGameData = "parse"
@scope("JSON") @val
external parseLeaders: string => leadersData = "parse"
let getMsgType = tag =>
  switch tag->Js.String2.slice(~from=2, ~to_=9) {
  | "listGms" => InsertConn
  | "modConn" => ModifyConn
  | "addGame" => InsertGame
  | "mdLstGm" => ModifyListGame
  | "cntdown" => Countdown
  | "players" => ModifyPlayers
  | "newword" => Word
  | "rmvGame" => RemoveGame
  | "leaders" => Leaders
  | _ => Other
  }

let initialState: Reducer.state = {
  gamesList: Js.Nullable.null,
  players: [],
  sk: "",
  oldWord: "",
  word: "",
  showAnswers: false,
  winner: "",
}
@react.component
let make = (~token, ~setToken, ~cognitoUser, ~setCognitoUser, ~setWsError, ~route) => {
  Js.log2("u345876l", route)

  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)
  let (playerGame, setPlayerGame) = React.Uncurried.useState(_ => "")
  let (playerName, setPlayerName) = React.Uncurried.useState(_ => "")
  let (endtoken, setEndtoken) = React.Uncurried.useState(_ => "")
  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "transparent")
  let (count, setCount) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)

  let (leaderData, setLeaderData) = React.Uncurried.useState(_ => [])
  let (state, dispatch) = React.Uncurried.useReducerWithMapState(
    Reducer.reducer,
    initialState,
    Reducer.init,
  )

  let {players, sk, showAnswers, winner, oldWord, word, gamesList: games} = state

  let resetConnState = (. ()) => {
    dispatch(. ResetPlayerState(initialState))
    setLeaderData(._ => [])
    setEndtoken(._ => "")
    setCount(._ => "")
  }

  let wsorigin = `wss://${apiid}.execute-api.${region}.amazonaws.com`

  open Web
  let logAndDisconnect = (~msg: string, ~data: string, ~code: int) => {
    switch Js.Json.stringifyAny({
      Lobby.action: "logging",
      gameno: msg,
      aW5mb3Jt: data,
    }) {
    | None => ()
    | Some(s) => ws->sendString(s)
    }
    switch Js.Json.stringifyAny({
      Lobby.action: "lobby",
      gameno: switch playerGame == "" {
      | true => "discon"
      | false => playerGame
      },
      command: fromLobbyCommandToString(Disconnect),
    }) {
    | None => ()
    | Some(s) => ws->sendString(s)
    }
    ws->closeCodeReason(code, msg)
  }

  React.useEffect1(() => {
    switch token {
    | None => setWs(._ => Js.Nullable.null)
    | Some(token) => setWs(._ => Js.Nullable.return(newWs(`${wsorigin}/${stage}?auth=${token}`)))
    }
    None
  }, [token])

  React.useEffect1(() => {
    switch Js.Nullable.isNullable(ws) {
    | true => ()
    | false =>
      ws->onOpen(e => {
        setWsConnected(. _ => true)
        Js.log2("open", e)
      })
      ws->onError(e => {
        Js.log2("errrr", e)
        setWsError(. _ => "websocket error: connection closed")
      })

      ws->onMessage(({data, origin}) => {
        switch origin == wsorigin {
        | true => ()
        | false => logAndDisconnect(~msg="wrong origin", ~data=origin, ~code=4002)
        }

        switch Js.String2.length(data) > 3000 {
        | true =>
          logAndDisconnect(
            ~msg="excessive json data",
            ~data=Js.String2.slice(data, ~from=0, ~to_=3000),
            ~code=4004,
          )

        | false => ()
        }

        switch Js.String2.match_(data, %re(`/(\"\w+\"\:).+\1/g`)) {
        | None => ()
        | Some(_) => logAndDisconnect(~msg="duplicate keys", ~data, ~code=4003)
        }

        Js.log3("msg", data, origin)

        switch getMsgType(data) {
        | InsertConn => {
            let {listGms, name} = parseListGames(data)
            Js.log3("parsedlistgames", listGms, name)
            setPlayerName(. _ => name)
            dispatch(. ListGames(Js.Nullable.return(listGms)))
          }

        | ModifyConn => {
            let {modConn, color, endtoken} = parseModConn(data)
            Js.log4("parsedmodconn", modConn, color, endtoken)
            setPlayerGame(. _ => modConn)
            setPlayerColor(. _ => color)
            setEndtoken(. _ => endtoken)
          }

        | Countdown => {
            let {cntdown} = parseCountdown(data)
            Js.log2("parsedCountdown", cntdown)
            setCount(. _ => cntdown)
          }

        | InsertGame => {
            let {addGame} = parseAddGame(data)
            Js.log2("parsedaddgame", addGame)
            dispatch(. AddGame(addGame))
          }

        | ModifyListGame => {
            let {mdLstGm} = parseModListGame(data)
            Js.log2("parsedmodlistgame", mdLstGm)
            dispatch(. UpdateListGame(mdLstGm))
          }

        | ModifyPlayers => {
            let {players, sk, showAnswers, winner} = parseModPlayers(data)
            Js.log3("parsedmodplayers", players, sk)
            Js.log3("parsedmodplayers 2", showAnswers, winner)
            dispatch(. UpdatePlayers(players, sk, showAnswers, winner))
          }

        | Word => {
            let {newword} = parseWord(data)
            Js.log2("parsedword", newword)
            dispatch(. UpdateWord(newword))
          }

        | RemoveGame => {
            let {rmvGame} = parseRmvGame(data)
            Js.log2("parsedremgame", rmvGame)
            dispatch(. RemoveGame(rmvGame))
          }

        | Leaders => {
            let {leaders} = parseLeaders(data)
            Js.log2("parsedleaders", leaders)
            setLeaderData(. _ => leaders)
          }

        | Other => {
            Js.log2("unknown json data", data)

            logAndDisconnect(~msg="unknown json data", ~data, ~code=4005)
          }
        }
      })

      ws->onClose(({code, reason, wasClean}) => {
        open Cognito
        Js.log4("close", code, reason, wasClean)
        setWsConnected(. _ => false)
        // setWsError(. _ => "")

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => ()
        | false => cognitoUser->signOut(Js.Nullable.null)
        }
        setCognitoUser(. _ => Js.Nullable.null)
        setWs(. _ => Js.Nullable.null)
        setPlayerName(. _ => "")

        resetConnState(.)
        setToken(. _ => None)
      })
    }

    let cleanup = () => {
      Js.log("cleanup")
      switch Js.Nullable.isNullable(ws) {
      | true => ()
      | false => ws->closeCode(1000)
      }
    }
    Some(cleanup)
  }, [ws])

  let send = (. str) => {
    switch str {
    | None => ()
    | Some(s) => ws->sendString(s)
    }
  }

  let close = (. code, reason) => ws->closeCodeReason(code, reason)

  let leaders = React.createElement(
    Leaders.lazy_(() =>
      Leaders.import_("./Leaders.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Leaders.makeProps(~leaderData, ~playerName, ()),
  )

  <>
    {switch route {
    | Route.Leaderboard => React.null
    | Home | SignIn | SignUp | GetInfo(_) | Confirm(_) | Lobby | Play(_) | Other =>
      <header className="mb-10 newgmimg:mb-12">
        <p
          className="font-flow text-stone-100 text-2xl newgmimg:text-4xl h-10 font-bold text-center">
          {React.string(playerName)}
        </p>
        <h1
          style={ReactDOM.Style.make(~backgroundColor={playerColor}, ())}
          className="text-6xl mt-11 mx-auto px-6 text-center font-arch decay-mask text-stone-100">
          {React.string("CLEAN TABLET")}
        </h1>
      </header>
    }}
    {switch route {
    | Lobby =>
      switch wsConnected {
      | false => {
          body(document)->setClassName("bodchmob bodchtab bodchbig")

          <Loading label="games..." />
        }

      | true => {
          body(document)->classList->removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")
          <Lobby playerGame games send close count setLeaderData />
        }
      }
    | Play({play}) =>
      switch wsConnected {
      | true =>
        switch Js.Array2.length(players) > 0 && play == sk {
        | true =>
          <Play
            players
            sk
            showAnswers
            winner
            isWinner={winner != ""}
            oldWord
            word
            playerColor
            send
            playerName
            endtoken
            resetConnState
          />

        | false => <Loading label="game..." />
        }

      | false =>
        <p className="text-center text-stone-100 font-anon text-lg">
          {React.string("not connected...")}
        </p>
      }

    | Leaderboard => {
        body(document)->classList->addClassList3("bodleadmob", "bodleadtab", "bodleadbig")

        <React.Suspense fallback=React.null> leaders </React.Suspense>
      }

    | Home | SignIn | SignUp | GetInfo(_) | Confirm(_) | Other =>
      <div className="text-center text-stone-100 text-4xl"> {React.string("page not found")} </div>
    }}
  </>
}
