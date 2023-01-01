type propShape = {
  "leaderData": array<Reducer.stat>,
  "playerName": string,
  "send": (. option<string>) => unit,
  "setLeaderData": (. array<Reducer.stat> => array<Reducer.stat>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make: (
  ~leaderData: array<Reducer.stat>,
  ~playerName: string,
  ~send: (. option<string>) => unit,
  ~setLeaderData: (. array<Reducer.stat> => array<Reducer.stat>) => unit,
) => React.element
