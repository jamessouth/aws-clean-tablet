type apigwPayload = {
  action: string,
  gameno: string,
  data: string,
}

@react.component
let make: (
  ~playerGame: string,
  ~games: Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  ~send: (. option<string>) => unit,
  ~wsError: string,
  ~close: (. int, string) => unit,
  ~count: string,
  ~setLeaderData: (. array<Reducer.stat> => array<Reducer.stat>) => unit,
) => React.element