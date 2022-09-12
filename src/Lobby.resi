type leaderPayload = {
  action: string,
  info: string,
}

@react.component
let make: (
  ~playerGame: string,
  ~games: Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  ~send: (. option<string>) => unit,
  ~wsError: string,
  ~close: (. int, string) => unit,
  ~count: string,
  ~setLeaderData: (. 'a => array<'b>) => unit,
) => React.element
