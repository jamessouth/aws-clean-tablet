type poolData
type poolDataInput = {
  @as("UserPoolId") userPoolId: string,
  @as("ClientId") clientId: string,
  @as("AdvancedSecurityDataCollectionFlag") advancedSecurityDataCollectionFlag: bool,
}
type attributeData
type attributeDataInput = {
  @as("Name") name: string,
  @as("Value") value: string,
}
type userData
type userDataInput = {
  @as("Username") username: string,
  @as("Pool") pool: poolData,
}
type clientMetadata = {key: string}
type accessToken = {jwtToken: string}
type userSession = {accessToken: accessToken}
type usr = {username: string}
type signupOk = {user: usr}
type t
type signUpCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<signupOk>) => unit
type signInCB = {
  onFailure: Js.Exn.t => unit,
  newPasswordRequired: Js.Nullable.t<(array<attributeData>, array<attributeData>) => unit>,
  mfaRequired: Js.Nullable.t<(string, string) => unit>,
  customChallenge: Js.Nullable.t<string => unit>,
  onSuccess: userSession => unit,
}
type authDetails = {
  @as("ValidationData") validationData: Js.Nullable.t<array<attributeData>>,
  @as("Username") username: string,
  @as("Password") password: string,
  @as("AuthParameters") authParameters: Js.Nullable.t<array<attributeData>>,
  @as("ClientMetadata") clientMetadata: Js.Nullable.t<clientMetadata>,
}
type passwordPWCB = {
  onFailure: Js.Exn.t => unit,
  onSuccess: string => unit,
}
type confirmRegistrationCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<t>) => unit
type revokeTokenCB = Js.Exn.t => unit
@send
external signOut: (Js.Nullable.t<usr>, Js.Nullable.t<revokeTokenCB>) => unit = "signOut"
@send
external signUp: (
  poolData,
  string,
  string,
  Js.Nullable.t<array<attributeData>>,
  Js.Nullable.t<array<attributeData>>,
  signUpCB,
  Js.Nullable.t<clientMetadata>,
) => unit = "signUp"
@send
external confirmRegistration: (
  Js.Nullable.t<usr>,
  string,
  bool,
  confirmRegistrationCB,
  Js.Nullable.t<clientMetadata>,
) => unit = "confirmRegistration"
@send
external confirmPassword: (
  Js.Nullable.t<usr>,
  string, //conf code
  string, //new pw
  passwordPWCB, //cb obj
  Js.Nullable.t<clientMetadata>,
) => unit = "confirmPassword"
@send
external authenticateUser: (Js.Nullable.t<usr>, authDetails, signInCB) => unit = "authenticateUser"
@new @module("amazon-cognito-identity-js")
external userAttributeConstructor: attributeDataInput => attributeData = "CognitoUserAttribute"
@new @module("amazon-cognito-identity-js")
external authenticationDetailsConstructor: authDetails => authDetails = "AuthenticationDetails"
@new @module("amazon-cognito-identity-js")
external userConstructor: userDataInput => usr = "CognitoUser"
@new @module("amazon-cognito-identity-js")
external userPoolConstructor: poolDataInput => poolData = "CognitoUserPool"
