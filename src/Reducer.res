
type player = {
    name: string,
    connid: string,
    ready: bool,
    color: string,
    score: int
}


type game = {
    leader: string,
    no: string,
    starting: bool,
    players: array<player>
}


type state = {
    games: array<game>
}


type returnVal = {
  setToken: string => unit,
  token: option<string>,
}


let mergeGame = (arr, ni) => {
    let list = [...arr]
    for i in 0 to Js.Array2.length(list) - 1 {
        switch (list[i].no == ni.no, ni.starting) {
        | (true, true) => Js.Array2.spliceInPlace(list, ~pos=i, ~remove=1, ~add=[])
        | (true, false) => {
            list[i] = ni
            list
        }
        | (false, _) => [ni, ...list]
        }
    
    }

}

let appState= () => {
  Js.log("appState")
  
  

  let (playerColor, setPlayerColor) = React.Uncurried.useState(_ => "")



  let return = {
    setToken: saveToken,
    token: token,
  }

  return
}
