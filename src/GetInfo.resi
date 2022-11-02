type propShape = {
  "userpool": Cognito.poolData,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setRetrievedUsername": (. Js.String2.t => Js.String2.t) => unit,
  "search": Route.query,
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
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~setRetrievedUsername: (. Js.String2.t => Js.String2.t) => unit,
  ~search: Route.query,
) => React.element
