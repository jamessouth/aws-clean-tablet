

type t
type messageEvent = {data: string}
type messageEventHandler = messageEvent => unit

@new external new_: string => t = "WebSocket"
@set external set_onmessage: (t, messageEventHandler) => unit = "onmessage"


type returnVal = {
  setToken: string => unit,
  token: option<string>,
}

let useWs= () => {
  Js.log("wshook")
  
  

  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  let (currentWord, setCurrentWord) = React.Uncurried.useState(_ => "")
  let (previousWord, setPreviousWord) = React.Uncurried.useState(_ => "")
  let (game, setGame) = React.Uncurried.useState(_ => Js.Nullable.null)
  let (ingame, setIngame) = React.Uncurried.useState(_ => "")
  let (leadertoken, setLeadertoken) = React.Uncurried.useState(_ => "")

  let saveToken = token => {
    localStorage->Dom.Storage2.setItem("token", token)
    setToken(._ => Some(token))
  }

  let return = {
    setToken: saveToken,
    token: token,
  }

  return
}
