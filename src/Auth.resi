type propShape = {
  "authToken": option<string>,
  "setAuthToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "route": Route.t,
  "idToken": option<string>,
  "setIdToken": (. option<string> => option<string>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make: (
  ~authToken: option<string>,
  ~setAuthToken: (. option<string> => option<string>) => unit,
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~route: Route.t,
  ~idToken: option<string>,
  ~setIdToken: (. option<string> => option<string>) => unit,
) => React.element
