type listGamesData = {listGms: array<Reducer.listGame>, name: string, returning: bool}
type modConnData = {modConn: string, color: string, index: string}
type addGameData = {addGame: Reducer.listGame}
type modListGameData = {mdLstGm: Reducer.listGame}
type modLiveGameData = {mdLveGm: Reducer.liveGame}
type rmvGameData = {rmvGame: Reducer.listGame}
type leadersData = {leaders: array<Reducer.stat>}
type msgType =
  InsertConn | ModifyConn | InsertGame | ModifyListGame | ModifyLiveGame | RemoveGame | Leaders | Other
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
@scope("JSON") @val
external parseLeaders: string => leadersData = "parse"
let getMsgType = tag => {
  switch tag->Js.String2.slice(~from=2, ~to_=9) {
  | "listGms" => InsertConn
  | "modConn" => ModifyConn
  | "addGame" => InsertGame
  | "mdLstGm" => ModifyListGame
  | "mdLveGm" => ModifyLiveGame
  | "rmvGame" => RemoveGame
  | "leaders" => Leaders
  | _ => Other
  }
}
