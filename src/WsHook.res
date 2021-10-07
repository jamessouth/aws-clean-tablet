

type t
type messageEvent = {data: string}
type messageEventHandler = messageEvent => unit

@new external newWs: string => t = "WebSocket"
@set external onMessage: (Js.Nullable.t<t>, messageEventHandler) => unit = "onmessage"


type returnVal = {
  setToken: string => unit,
  token: option<string>,
}

let useWs= (token) => {
  Js.log("wshook")
  
  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)

  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  let (currentWord, setCurrentWord) = React.Uncurried.useState(_ => "")
  let (previousWord, setPreviousWord) = React.Uncurried.useState(_ => "")
  let (game, setGame) = React.Uncurried.useState(_ => Js.Nullable.null)
  let (ingame, setIngame) = React.Uncurried.useState(_ => "")
  let (leadertoken, setLeadertoken) = React.Uncurried.useState(_ => "")

  let {initialState, reducer} = Reducer.appState()

  let (state, dispatch) = React.useReducer(reducer, initialState)

  React.useEffect1(() => {
    switch token {
    | None => ()
    | Some(token) => {
      // let mystr = `wss://${process.env.CT_APIID}.execute-api.${process.env.CT_REGION}.amazonaws.com/${process.env.CT_STAGE}?auth=${token}`
      let sock = Js.Nullable.return(newWs("mystr"))
      setWs(. _ => sock)

switch Js.Nullable.isNullable(ws) {
| true => ()
| false => ws->onMessage(({data}) => Js.log(data))
}
      
    }
    }












    None
  }, [token])


}
