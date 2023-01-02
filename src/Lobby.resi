type apigwAction =
  | Answer
  | Lobby
  | Logging
  | Query

type lobbyCommand =
  | Custom({cv: string})
  | Join
  | Leaders
  | Leave
  | ListGames

type lobbyPayload = {
  act: apigwAction,
  gn?: string,
  cmd: lobbyCommand,
}

type payloadOutput = {
  action: string,
  gameno?: string,
  aW5mb3Jt?: string,
  command?: string,
  log?: string,
}

let payloadToObj: lobbyPayload => option<string>

@react.component
let make: (
  ~wsConnected: bool,
  ~playerListGame: string,
  ~games: Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  ~send: (. option<string>) => unit,
  ~close: (. int, string) => unit,
  ~count: string,
) => React.element
