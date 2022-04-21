@val @scope(("import", "meta", "env"))
external apiid: string = "VITE_APIID"
@val @scope(("import", "meta", "env"))
external region: string = "VITE_REGION"
@val @scope(("import", "meta", "env"))
external stage: string = "VITE_STAGE"

let useWs = (token, setToken, cognitoUser, setCognitoUser, initialState) => {
  let (ws, setWs) = React.Uncurried.useState(_ => Js.Nullable.null)
  let (playerGame, setPlayerGame) = React.Uncurried.useState(_ => "")
  let (playerName, setPlayerName) = React.Uncurried.useState(_ => "")
  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "transparent")
  let (playerIndex, setPlayerIndex) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  let (leader, setLeader) = React.Uncurried.useState(_ => false)
  let (rawLeaders, setRawLeaders) = React.Uncurried.useState(_ => [])
  let (state, dispatch) = React.Uncurried.useReducerWithMapState(
    Reducer.reducer,
    initialState,
    Reducer.init,
  )

  let resetConnState = _ => {
    dispatch(. ResetPlayerState(initialState))
    setPlayerColor(._ => "transparent")
    setPlayerIndex(._ => "")
    setPlayerGame(._ => "")
    setLeader(._ => false)
    setRawLeaders(._ => [])
  }

  open Web
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
        open Json
        Js.log2("msg", data)

        switch getMsgType(data) {
        | InsertConn => {
            let {listGms, name, returning} = parseListGames(data)
            Js.log4("parsedlistgames", listGms, name, returning)
            setPlayerName(._ => name)
            switch returning {
            | true => resetConnState()
            | false => ()
            }
            dispatch(. ListGames(Js.Nullable.return(listGms)))
          }
        | ModifyConn => {
            let {modConn, color, leader, index} = parseModConn(data)
            Js.log2("parsedmodconn", modConn)
            Js.log3(color, leader, index)
            setPlayerGame(._ => modConn)
            setPlayerColor(._ => color)
            setPlayerIndex(._ => index)
            setLeader(._ => leader)
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
        | ModifyLiveGame => {
            let {mdLveGm} = parseModLiveGame(data)
            Js.log2("parsedmodlivegame", mdLveGm)
            dispatch(. UpdateLiveGame(mdLveGm))
          }
        | RemoveGame => {
            let {rmvGame} = parseRmvGame(data)
            Js.log2("parsedremgame", rmvGame)
            dispatch(. RemoveGame(rmvGame))
          }
        | Leaders => {
            let {leaders} = parseLeaders(data)
            Js.log2("parsedleaders", leaders)
            setRawLeaders(._ => leaders)
          }
        | Other => Js.log2("unknown json data", data)
        }
      })

      ws->onClose(({code, reason, wasClean}) => {
        open Cognito
        Js.log4("close", code, reason, wasClean)
        setToken(._ => None)
        setWsConnected(._ => false)
        setWsError(._ => "")

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => ()
        | false => cognitoUser->signOut(Js.Nullable.null)
        }
        setCognitoUser(._ => Js.Nullable.null)
        setWs(._ => Js.Nullable.null)
        setPlayerName(._ => "")

        resetConnState()
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

  (playerGame, playerName, playerColor, playerIndex, wsConnected, state.game, state.gamesList, leader, rawLeaders, send, close, wsError)
}
