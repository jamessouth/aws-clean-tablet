type propShape = {
  "close": (. int, string) => unit,
  "count": string,
  "games": Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  "playerGame": string,
  "send": (. option<string>) => unit,
  "wsError": string,
  "setLeaderData": (. array<Reducer.stat> => array<Reducer.stat>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

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
