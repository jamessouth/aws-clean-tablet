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
  let (endtoken, setEndtoken) = React.Uncurried.useState(_ => Js.Nullable.undefined)
  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "transparent")
  let (count, setCount) = React.Uncurried.useState(_ => "")
  let (wsConnected, setWsConnected) = React.Uncurried.useState(_ => false)
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")
  // let (leader, setLeader) = React.Uncurried.useState(_ => false)
  let (leaderData, setLeaderData) = React.Uncurried.useState(_ => [])
  let (state, dispatch) = React.Uncurried.useReducerWithMapState(
    Reducer.reducer,
    initialState,
    Reducer.init,
  )

  let resetConnState = _ => {
    dispatch(. ResetPlayerState(initialState))
    setLeaderData(._ => [])
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
            let {modConn, color, endtoken} = parseModConn(data)
            Js.log4("parsedmodconn", modConn, color, endtoken)
            setPlayerGame(._ => modConn)
            setPlayerColor(._ => color)
            setEndtoken(._ => endtoken)
          }


        | Countdown => {
            let {cntdown} = parseCountdown(data)
            Js.log2("parsedCountdown", cntdown)
            setCount(._ => cntdown)
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
            setLeaderData(._ => leaders)
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
        setCount(._ => "")

        resetConnState()
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

  (playerGame, playerName, playerColor, endtoken, count, wsConnected, state.players, state.sk, state.showAnswers, state.winner, state.oldWord, state.word, state.gamesList, leaderData, send, close, wsError)
}
