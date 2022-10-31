type propShape = {
  "token": option<string>,
  "setToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setWsError": (. string => string) => unit,
  "route": Route.t,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make: (
  ~token: option<string>,
  ~setToken: (. option<string> => option<string>) => unit,
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~setWsError: (. string => string) => unit,
  ~route: Route.t,
) => React.element
