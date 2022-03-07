@val @scope(("import", "meta", "env"))
external apiid: string = "VITE_APIID"
@val @scope(("import", "meta", "env"))
external region: string = "VITE_REGION"
@val @scope(("import", "meta", "env"))
external stage: string = "VITE_STAGE"

type t
type openEventHandler = unit => unit
type errorEventHandler = Dom.errorEvent => unit
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

@val external document: Dom.document = "document"
@get external body: Dom.document => Dom.htmlBodyElement = "body"
@set external setClassName: (Dom.htmlBodyElement, string) => unit = "className"

type listGamesData = {listGms: array<Reducer.listGame>}
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"

type modConnData = {modConn: string, color: string, leader: bool}
@scope("JSON") @val
external parseModConn: string => modConnData = "parse"

type addGameData = {addGame: Reducer.listGame}
@scope("JSON") @val
external parseAddGame: string => addGameData = "parse"

type modListGameData = {mdLstGm: Reducer.listGame}
@scope("JSON") @val
external parseModListGame: string => modListGameData = "parse"

type modLiveGameData = {mdLveGm: Reducer.liveGame}
@scope("JSON") @val
external parseModLiveGame: string => modLiveGameData = "parse"

type rmvGameData = {rmvGame: Reducer.listGame}
@scope("JSON") @val
external parseRmvGame: string => rmvGameData = "parse"

type msgType =
  | InsertConn
  | ModifyConn
  | InsertGame
  | ModifyListGame
  | ModifyLiveGame
  | RemoveGame
  | Other

let getMsgType = tag => {
  switch tag->Js.String2.slice(~from=2, ~to_=9) {
  | "listGms" => InsertConn
  | "modConn" => ModifyConn
  | "addGame" => InsertGame
  | "mdLstGm" => ModifyListGame
  | "mdLveGm" => ModifyLiveGame
  | "rmvGame" => RemoveGame
  | _ => Other
  }
}

type revokeTokenCallback = Js.Exn.t => unit

@send
external signOut: (Js.Nullable.t<Signup.usr>, Js.Nullable.t<revokeTokenCallback>) => unit =
  "signOut"

let useWs = (token, setToken, cognitoUser, setCognitoUser, setPlayerName) => {
  // Js.log2("wshook ", token)

  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)

  let (playerGame, setPlayerGame) = React.Uncurried.useState(_ => "")
  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "transparent")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")

  let (leader, setLeader) = React.Uncurried.useState(_ => false)

  let (reducer, initialState) = Reducer.appState()

  let (state, dispatch) = React.Uncurried.useReducer(reducer, initialState)

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
        body(document)->setClassName("bodchmob bodchtab bodchbig")
      })
      ws->onError(e => {
        Js.log2("errrr", e)
        setWsError(._ => "temp error placehold")
      })

      ws->onMessage(({data}) => {
        Js.log2("msg", data)

        switch getMsgType(data) {
        | InsertConn => {
            let {listGms} = parseListGames(data)
            Js.log2("parsedlistgames", listGms)
            dispatch(.ListGames(Js.Nullable.return(listGms)))
          }
        | ModifyConn => {
            let {modConn, color, leader} = parseModConn(data)
            Js.log4("parsedmodconn", modConn, color, leader)
            setPlayerGame(._ => modConn)
            setPlayerColor(._ => color)
            setLeader(._ => leader)
          }
        | InsertGame => {
            let {addGame} = parseAddGame(data)
            Js.log2("parsedaddgame", addGame)
            dispatch(.AddGame(addGame))
          }
        | ModifyListGame => {
            let {mdLstGm} = parseModListGame(data)
            Js.log2("parsedmodlistgame", mdLstGm)
            dispatch(.UpdateListGame(mdLstGm))
          }
        | ModifyLiveGame => {
            let {mdLveGm} = parseModLiveGame(data)
            Js.log2("parsedmodlivegame", mdLveGm)
            dispatch(.UpdateLiveGame(mdLveGm))
          }
        | RemoveGame => {
            let {rmvGame} = parseRmvGame(data)
            Js.log2("parsedremgame", rmvGame)
            dispatch(.RemoveGame(rmvGame))
          }
        | Other => Js.log2("unknown json data", data)
        }
      })

      ws->onClose(({code, reason, wasClean}) => {
        Js.log4("close", code, reason, wasClean)
        setToken(._ => None)
        setWsConnected(._ => false)
        setWsError(._ => "")

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => ()
        | false => cognitoUser->signOut(Js.Nullable.null)
        }
        setCognitoUser(._ => Js.Nullable.null)
        setPlayerName(._ => "")
        setPlayerColor(._ => "transparent")

        setPlayerGame(._ => "")
        setLeader(._ => false)
        setWs(._ => Js.Nullable.null)
        dispatch(.(ResetPlayerState: Reducer.action))
        body(document)->setClassName("bg-no-repeat bg-center bg-cover bodmob bodtab bodbig")
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

  (playerGame, playerColor, wsConnected, state.game, state.gamesList, leader, send, close, wsError)
}
