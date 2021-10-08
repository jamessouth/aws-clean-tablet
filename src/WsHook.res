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
  wasClean: bool
}
type closeEventHandler = closeEvent => unit

@new external newWs: string => t = "WebSocket"
@set external onOpen: (Js.Nullable.t<t>, openEventHandler) => unit = "onopen"
@set external onError: (Js.Nullable.t<t>, errorEventHandler) => unit = "onerror"
@set external onMessage: (Js.Nullable.t<t>, messageEventHandler) => unit = "onmessage"
@set external onClose: (Js.Nullable.t<t>, closeEventHandler) => unit = "onclose"

type returnVal = {
  setToken: string => unit,
  token: option<string>,
}

let useWs = token => {
  Js.log("wshook")

  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)

  // let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "")
  // let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  // let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  // let (currentWord, setCurrentWord) = React.Uncurried.useState(_ => "")
  // let (previousWord, setPreviousWord) = React.Uncurried.useState(_ => "")
  // let (game, setGame) = React.Uncurried.useState(_ => Js.Nullable.null)
  // let (ingame, setIngame) = React.Uncurried.useState(_ => "")
  // let (leadertoken, setLeadertoken) = React.Uncurried.useState(_ => "")

  // let {initialState, reducer} = Reducer.appState()

  // let (state, dispatch) = React.useReducer(reducer, initialState)

  React.useEffect1(() => {
    switch token {
    | None => ()
    | Some(token) => {
        setWs(._ => Js.Nullable.return(newWs(`wss://${apiid}.execute-api.${region}.amazonaws.com/${stage}?auth=${token}`)))

        switch Js.Nullable.isNullable(ws) {
        | true => ()
        | false =>
          ws->onOpen((e) => Js.log2("open", e))
          ws->onError((e) => Js.log2("error", e))
          ws->onMessage(({data}) => Js.log2("msg", data))
          ws->onClose(({code, reason, wasClean}) => Js.log4("close", code, reason, wasClean))
        }
      }
    }

    None
  }, [token])
}
