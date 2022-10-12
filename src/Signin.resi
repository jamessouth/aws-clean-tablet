type propShape = {
  "userpool": Cognito.poolData,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setAuthToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
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
  ~userpool: Cognito.poolData,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~setAuthToken: (. option<string> => option<string>) => unit,
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
  ~setIdToken: (. option<string> => option<string>) => unit,
) => React.element
