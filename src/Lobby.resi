type apigwAction =
  | Answer
  | End
  | Leaders
  | Lobby
  | Logging

type lobbyGameno =
  | Gameno({no: string})
  | Discon
  | Newgame

type lobbyCommand =
  | Custom({cv: string})
  | Disconnect
  | Join
  | Leave
  | Ready
  | Unready

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
  ~playerGame: string,
  ~games: Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  ~send: (. option<string>) => unit,
  ~close: (. int, string) => unit,
  ~count: string,
  ~setLeaderData: (. array<Reducer.stat> => array<Reducer.stat>) => unit,
) => React.element
