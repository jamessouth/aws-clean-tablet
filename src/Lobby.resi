type apigwAction = 
| Answer
| End
 | Leaders
  | Lobby

type lobbyGameno =
  | Gameno({no: string})
  | Discon
  | Newgame
  | None

type lobbyCommand =
  | Disconnect
  | Endtoken({et: string})
  | Join
  | Leave
  | None
  | Ready
  | Unready

type lobbyPayload = {
  act: apigwAction,
  gn: lobbyGameno,
  cmd: lobbyCommand,
}


let apigwActionToString: apigwAction => string
let lobbyGamenoToString: lobbyGameno => string
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
