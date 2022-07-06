type listGamesData = {listGms: array<Reducer.listGame>, name: string, returning: bool}
type modConnData = {modConn: string, color: string}
type countdownData = {cntdown: string}
type addGameData = {addGame: Reducer.listGame}
type modListGameData = {mdLstGm: Reducer.listGame}


type modPlayersData = {players: array<Reducer.livePlayer>, sk: string, showAnswers: bool, winner: string}
type wordData = {newword: string}


type rmvGameData = {rmvGame: Reducer.listGame}
type leadersData = {leaders: array<Reducer.stat>}
type msgType =
  | InsertConn
  | ModifyConn
  | InsertGame
  | ModifyListGame
  | Countdown
  | ModifyPlayers
  | Word
  | RemoveGame
  | Leaders
  | Other
@scope("JSON") @val
external parseListGames: string => listGamesData = "parse"
@scope("JSON") @val
external parseModConn: string => modConnData = "parse"
@scope("JSON") @val
external parseCountdown: string => countdownData = "parse"
@scope("JSON") @val
external parseAddGame: string => addGameData = "parse"
@scope("JSON") @val
external parseModListGame: string => modListGameData = "parse"
@scope("JSON") @val

external parseModPlayers: string => modPlayersData = "parse"
@scope("JSON") @val
external parseWord: string => wordData = "parse"
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
  | "cntdown" => Countdown
  | "players" => ModifyPlayers
  | "newword" => Word
  | "rmvGame" => RemoveGame
  | "leaders" => Leaders
  | _ => Other
  }
}
