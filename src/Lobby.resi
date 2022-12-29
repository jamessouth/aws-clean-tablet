type apigwAction =
  | Answer
  | Lobby
  | Logging
  | Query

type lobbyGameno =
  | Gameno({no: string})
  | Newgame

type lobbyCommand =
  | Custom({cv: string})
  | Join
  | Leaders
  | Leave
  | ListGames

type lobbyPayload = {
  act: apigwAction,
  gn: lobbyGameno,
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
  ~playerListGame: string,
  ~games: Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  ~send: (. option<string>) => unit,
  ~close: (. int, string) => unit,
  ~count: string,
  ~setLeaderData: (. array<Reducer.stat> => array<Reducer.stat>) => unit,
) => React.element
