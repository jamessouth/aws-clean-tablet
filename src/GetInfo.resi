@react.component
let make: (
  ~userpool: Cognito.poolData,
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~setShowName: (. Js.String2.t => Js.String2.t) => unit,
  ~search: Route.query,
) => React.element
