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

type userSession = {
  accessToken: accessToken,
}

type usr = {
  username: string,
}

type signupOk = {
  user: usr,
}

type signUpCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<signupOk>) => unit











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








@new @module("amazon-cognito-identity-js")
external userAttributeConstructor: attributeDataInput => attributeData = "CognitoUserAttribute"
