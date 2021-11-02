@val @scope(("import", "meta", "env"))
external apiid: string = "VITE_APIID"
@val @scope(("import", "meta", "env"))
external region: string = "VITE_REGION"
@val @scope(("import", "meta", "env"))
external stage: string = "VITE_STAGE"

type t
type openEventHandler = unit => unit
type errorEventHandler = unit => unit
type messageEvent = {data: string}
type messageEventHandler = messageEvent => unit
type closeEvent = {
  code: int,
  reason: string,
  wasClean: bool,
}
type closeEventHandler = closeEvent => unit

@new external newWs: string => t = "WebSocket"
@set external onOpen: (Js.Nullable.t<t>, openEventHandler) => unit = "onopen"
@set external onError: (Js.Nullable.t<t>, errorEventHandler) => unit = "onerror"
@set external onMessage: (Js.Nullable.t<t>, messageEventHandler) => unit = "onmessage"
@set external onClose: (Js.Nullable.t<t>, closeEventHandler) => unit = "onclose"

@send external closeNiladic: (Js.Nullable.t<t>, unit) => unit = "close"
@send external closeCode: (Js.Nullable.t<t>, int) => unit = "close"
@send external closeReason: (Js.Nullable.t<t>, string) => unit = "close"
@send external closeCodeReason: (Js.Nullable.t<t>, int, string) => unit = "close"

@send external sendString: (Js.Nullable.t<t>, string) => unit = "send"

type return = {
  playerColor: string,
  wsConnected: bool,
  game: Reducer.game,
  games: Js.Nullable.t<array<Reducer.game>>,
  playerGame: string,
  currentWord: string,
  previousWord: string,
  connID: string,
  send: (. option<string>) => unit,
  close: (int, string) => unit,
  wsError: string,
  setWs: (. Js.Nullable.t<t> => Js.Nullable.t<t>) => unit,
}

type listGamesData = {
  list: array<Reducer.game>,
  connID: string,
}
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"

type modConnData = {modC: string}
@scope("JSON") @val
external parseModConn: string => modConnData = "parse"

type addGameData = {addG: Reducer.game}
@scope("JSON") @val
external parseAddGame: string => addGameData = "parse"

type modGameData = {modG: Reducer.game}
@scope("JSON") @val
external parseModGame: string => modGameData = "parse"

type remGameData = {remG: Reducer.game}
@scope("JSON") @val
external parseRemGame: string => remGameData = "parse"

type msgType =
  | InsertConn
  | ModifyConn
  | InsertGame
  | ModifyGame
  | RemoveGame
  | Other

let getMsgType = tag => {
  switch tag->Js.String2.slice(~from=2, ~to_=6) {
  | "list" => InsertConn
  | "modC" => ModifyConn
  | "addG" => InsertGame
  | "modG" => ModifyGame
  | "remG" => RemoveGame
  | _ => Other
  }
}

let useWs = (token, setToken) => {
  // Js.log2("wshook ", token)

  let emptyGame: Reducer.game = {
    ready: false,
    no: "",
    starting: false,
    players: [],
  }

  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)

  let (playerColor, _setPlayerColor) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  let (currentWord, _setCurrentWord) = React.Uncurried.useState(_ => "")
  let (previousWord, _setPreviousWord) = React.Uncurried.useState(_ => "")
  let (game, _setGame) = React.Uncurried.useState(_ => emptyGame)
  let (playerGame, setPlayerGame) = React.Uncurried.useState(_ => "")
  let (connID, setConnID) = React.Uncurried.useState(_ => "")

  let {initialState, reducer} = Reducer.appState()

  let (state, dispatch) = React.useReducer(reducer, initialState)

  React.useEffect1(() => {
    // Js.log2("effect ", token)
    switch token {
    | None => setWs(._ => Js.Nullable.null)
    | Some(token) =>
      setWs(._ =>
        Js.Nullable.return(
          newWs(`wss://${apiid}.execute-api.${region}.amazonaws.com/${stage}?auth=${token}`),
        )
      )
    }

    None
  }, [token])



  React.useEffect1(() => {
    switch Js.Nullable.isNullable(ws) {
    | true => ()
    | false =>
      ws->onOpen(e => {
        setWsConnected(._ => true)
        Js.log2("open", e)
      })
      ws->onError(e => {
        Js.log2("errrr", e)
        setWsError(._ => "temp error placehold")
      })

      ws->onMessage(({data}) => {
        Js.log2("msg", data)

        switch getMsgType(data) {
        | InsertConn => {
            let {list, connID} = parseListGames(data)
            Js.log3("parsedlistgames", list, connID)
            dispatch(ListGames(Js.Nullable.return(list)))
            setConnID(._ => connID)
          }
        | ModifyConn => {
            let {modC} = parseModConn(data)
            Js.log2("parsedmodconn", modC)
            setPlayerGame(._ => modC)
          }
        | InsertGame => {
            let {addG} = parseAddGame(data)
            Js.log2("parsedaddgame", addG)
            dispatch(AddGame(addG))
          }
        | ModifyGame => {
            let {modG} = parseModGame(data)
            Js.log2("parsedmodgame", modG)
            dispatch(UpdateGame(modG))
          }
        | RemoveGame => {
            let {remG} = parseRemGame(data)
            Js.log2("parsedremgame", remG)
            dispatch(RemoveGame(remG))
          }
        | Other => Js.log2("unknown json data", data)
        }
      })

      ws->onClose(({code, reason, wasClean}) => {
        Js.log4("close", code, reason, wasClean)
        setToken(_ => None)
      })
    }

    let cleanup = () => {
      setWsConnected(._ => false)
      setWsError(._ => "")
      switch Js.Nullable.isNullable(ws) {
      | true => ()
      | false => ws->closeCode(1000)
      }
    }
    Some(cleanup)
  }, [ws])

  let send = (. str) => {
    // let dict = Js.Dict.empty()
    // Js.Dict.set(dict, "name", Js.Json.string("John Doe"))
    // Js.Dict.set(dict, "age", Js.Json.number(30.0))
    // Js.Dict.set(dict, "likes", Js.Json.stringArray(["bucklescript", "ocaml", "js"]))

    // ws->sendString(Js.Json.stringify(Js.Json.object_(dict)))
    switch str {
    | None => ()
    | Some(s) => ws->sendString(s)
    }
  }

  let close = (code, reason) => ws->closeCodeReason(code, reason)

  {
    playerColor: playerColor,
    wsConnected: wsConnected,
    game: game,
    games: state.gamesList,
    playerGame: playerGame,
    currentWord: currentWord,
    previousWord: previousWord,
    connID: connID,
    send: send,
    close: close,
    wsError: wsError,
    setWs: setWs,
  }
}
