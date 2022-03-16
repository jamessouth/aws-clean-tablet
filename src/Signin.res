@react.component
let make = (
  ~userpool,
  ~setCognitoUser,
  ~setToken,
  ~cognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~playerName,
) => {
  
  let (username, setUsername) = React.Uncurried.useState(_ => playerName)
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 characters; PASSWORD: 8-98 characters; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; ",
  ))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)

  React.useEffect2(() => {
    ErrorHook.useMultiError([(username, "USERNAME"), (password, "PASSWORD")], setValidationError)
    None
  }, (username, password))

  open Cognito
  let onClick = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None => {
        let cbs = {
          onSuccess: res => {
            setCognitoError(._ => None)
            Js.log2("signin result:", res)
            setToken(._ => Some(res.accessToken.jwtToken))
          },
          onFailure: ex => {
            switch Js.Exn.message(ex) {
            | Some(msg) => setCognitoError(._ => Some(msg))
            | None => setCognitoError(._ => Some("unknown signin error"))
            }

            setCognitoUser(._ => Js.Nullable.null)
            Js.log2("problem", ex)
          },
          newPasswordRequired: Js.Nullable.null,
          mfaRequired: Js.Nullable.null,
          customChallenge: Js.Nullable.null,
        }
        let authnData = {
          username: username,
          password: password,
          validationData: Js.Nullable.null,
          authParameters: Js.Nullable.null,
          clientMetadata: Js.Nullable.null,
        }
        let authnDetails = authenticationDetailsConstructor(authnData)

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => {
            let userdata = {
              username: username,
              pool: userpool,
            }
            let user = Js.Nullable.return(userConstructor(userdata))
            user->authenticateUser(authnDetails, cbs)
            setCognitoUser(._ => user)
          }
        | false => cognitoUser->authenticateUser(authnDetails, cbs)
        }
      }
    | Some(_) => ()
    }
  }

  <Form onClick leg="Sign in" submitClicked validationError cognitoError>
    <Input value=username propName="username" setFunc=setUsername />
    <Input value=password propName="password" autoComplete="current-password" setFunc=setPassword />
  </Form>
}
