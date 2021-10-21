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

@send external close: (Js.Nullable.t<t>, unit) => unit = "close"
@send external closeCode: (Js.Nullable.t<t>, int) => unit = "close"
@send external closeReason: (Js.Nullable.t<t>, string) => unit = "close"
@send external closeCodeReason: (Js.Nullable.t<t>, int, string) => unit = "close"

@send external sendString: (Js.Nullable.t<t>, string) => unit = "send"

type return = {
  playerColor: string,
  wsConnected: bool,
  game: Reducer.game,
  games: Js.Nullable.t<array<Reducer.game>>,
  ingame: string,
  leadertoken: string,
  currentWord: string,
  previousWord: string,
  send: option<string> => unit,
  wsError: string,
}

type listGamesData = {
  listGames: array<Reducer.game>,
  connID: string
}
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"


type modConnData = {
  modConnGm: string
}
// @scope("JSON") @val
external parseModConn: string => modConnData = "parse"




let useWs = (token) => {
  // Js.log2("wshook ", token)

  let emptyGame: Reducer.game = {
    leader: None,
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
  let (modConnGm, setModConnGm) = React.Uncurried.useState(_ => "")
  // let (leadertoken, _setLeadertoken) = React.Uncurried.useState(_ => "")

  let {initialState, reducer} = Reducer.appState()

  let (state, dispatch) = React.useReducer(reducer, initialState)

  React.useEffect1(() => {
    // Js.log2("effect ", token)
    switch token {
    | None => ()
    | Some(token) =>
      setWs(._ =>
        Js.Nullable.return(
          newWs(`wss://${apiid}.execute-api.${region}.amazonaws.com/${stage}?auth=${token}`),
        )
      )
    }

    None
  }, [token])

//  msg {"ingame":"163470931","leadertoken":"test_HfOJPg=","type":"user"}

//  msg {"games":{"no":"163470931","players":[{"name":"test","connid":"HfOJPg=","ready":false,"score":0}]},"type":"games"}


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

        switch data->Js.String2.slice(from=2, to_=11) {
        | "listGames" => {
            let {listGames, connID} = parseListGames(data)
            Js.log2("parsedlistgames", listGames, connID)
            dispatch(ListGames(Js.Nullable.return(listGames)))
            setConnID(_ => connID)
        }
        | "modConnGm" => {
            let {modConnGm} = parseModConn(data)
            Js.log2("parsedmodconn", modConnGm)
            setModConnGm(._ => modConnGm)

        }
        }
        


      })




      ws->onClose(({code, reason, wasClean}) => {
        Js.log4("close", code, reason, wasClean)
      })
    }

    let cleanup = () => {
      setWsConnected(._ => false)
      setWsError(._ => "")
      // setToken(_ => None)
      switch Js.Nullable.isNullable(ws) {
      | true => ()
      | false => ws->closeCode(1000)
      }
    }
    Some(cleanup)
  }, [ws])

  let send = str => {
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

  {
    playerColor: playerColor,
    wsConnected: wsConnected,
    game: game,
    games: state.gamesList,
    ingame: ingame,
    leadertoken: leadertoken,
    currentWord: currentWord,
    previousWord: previousWord,
    send: send,
    wsError: wsError,
  }
}
