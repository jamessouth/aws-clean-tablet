@react.component
let make: (
  ~userpool: Cognito.poolData,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  ~setToken: (. option<string> => option<string>) => unit,
  ~cognitoUser: Js.Nullable.t<Cognito.usr>,
) => React.element
