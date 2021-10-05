
@new @module("amazon-cognito-identity-js")
external userConstructor: Types.userDataInput => Signup.usr = "CognitoUser"


type callback = {
  onFailure: Js.Exn.t => unit,
  newPasswordRequired: Js.Nullable.t<
    (array<Types.attributeData>, array<Types.attributeData>) => unit,
  >,
  mfaRequired: Js.Nullable.t<(string, string) => unit>,
  customChallenge: Js.Nullable.t<string => unit>,
  onSuccess: Signup.userSession => unit,
}

type authDetails = {
  @as("ValidationData") validationData: Js.Nullable.t<array<Types.attributeData>>,
  @as("Username") username: string,
  @as("Password") password: string,
  @as("AuthParameters") authParameters: Js.Nullable.t<array<Types.attributeData>>,
  @as("ClientMetadata") clientMetadata: Js.Nullable.t<Signup.clientMetadata>,
}

@new @module("amazon-cognito-identity-js")
external authenticationDetailsConstructor: authDetails => authDetails = "AuthenticationDetails"


@send
external authenticateUser: (Js.Nullable.t<Signup.usr>, authDetails, callback) => unit =
  "authenticateUser"


type returnVal = {
  handleSubmit: ReactEvent.Form.t => unit
}

let useSignin = (username, password, userpool, setCognitoErr, setToken, setCognitoUser) => {
  Js.log("signinhook")


    let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    let authnData = {
      username: username,
      password: password,
      validationData: Js.Nullable.null,
      authParameters: Js.Nullable.null,
      clientMetadata: Js.Nullable.null,
    }
    let authnDetails = authenticationDetailsConstructor(authnData)
    let userdata: Types.userDataInput = {
      username,
      pool: userpool,
    }
    let cbs = {
      onSuccess: res => {
        Js.log2("signin result:", res)
        setToken(_ => Some(res.accessToken.jwtToken))
      },
      onFailure: ex => {
        switch Js.Exn.message(ex) {
        | Some(msg) => (_ => Some(msg))->setCognitoErr
        | None => (_ => None)->setCognitoErr
        }

        Js.log2("problem", ex)
      },
      newPasswordRequired: Js.Nullable.null,
      mfaRequired: Js.Nullable.null,
      customChallenge: Js.Nullable.null,
    }
    let user = Js.Nullable.return(userConstructor(userdata))
    user->authenticateUser(authnDetails, cbs)
    setCognitoUser(_ => user)
  }

  let return: returnVal = {
    handleSubmit: handleSubmit
  }

  return

}
