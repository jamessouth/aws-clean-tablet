type listGamesData = {listGms: array<Reducer.listGame>, returning: bool}
type modConnData = {modConn: string, color: string, index: string,leader: bool}
type addGameData = {addGame: Reducer.listGame}
type modListGameData = {mdLstGm: Reducer.listGame}
type modLiveGameData = {mdLveGm: Reducer.liveGame}
type rmvGameData = {rmvGame: Reducer.listGame}
type msgType =
  InsertConn | ModifyConn | InsertGame | ModifyListGame | ModifyLiveGame | RemoveGame | Other
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"
@scope("JSON") @val
external parseModConn: string => modConnData = "parse"
@scope("JSON") @val
external parseAddGame: string => addGameData = "parse"
@scope("JSON") @val
external parseModListGame: string => modListGameData = "parse"
@scope("JSON") @val
external parseModLiveGame: string => modLiveGameData = "parse"
@scope("JSON") @val
external parseRmvGame: string => rmvGameData = "parse"
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
