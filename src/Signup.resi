@react.component
let make: (
  ~userpool: Cognito.poolData,
  ~setCognitoUser: (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
) => React.element
