type propShape = {"cognitoUser": Js.Nullable.t<Cognito.usr>, "search": Route.query}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make: (~cognitoUser: Js.Nullable.t<Cognito.usr>, ~search: Route.query) => React.element
