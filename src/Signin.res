type propShape = {
  "userpool": Cognito.poolData,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

let username_max_length = 10
let password_max_length = 98

@react.component
let make = (~userpool, ~setCognitoUser, ~setToken, ~cognitoUser) => {
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 length; PASSWORD: 8-98 length; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; ",
  ))
  let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)
  React.useEffect2(() => {
    ErrorHook.useMultiError([(username, Username), (password, Password)], setValidationError)
    None
  }, (username, password))

  open Cognito
  let on_Click = (. ()) => {
    switch validationError {
    | None => {
        let cbs = {
          onSuccess: res => {
            setCognitoError(._ => None)
            Js.log2("signin result:", res)
          setToken(._ => Some(res.idToken.jwtToken))
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
          username: username
          ->Js.String2.slice(~from=0, ~to_=username_max_length)
          ->Js.String2.replaceByRe(%re("/\W/g"), ""),
          password: password
          ->Js.String2.slice(~from=0, ~to_=password_max_length)
          ->Js.String2.replaceByRe(%re("/\s/g"), ""),
          validationData: Js.Nullable.null,
          authParameters: Js.Nullable.null,
          clientMetadata: Js.Nullable.null,
        }
        let authnDetails = authenticationDetailsConstructor(authnData)

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => {
            let userdata = {
              username: username
              ->Js.String2.slice(~from=0, ~to_=username_max_length)
              ->Js.String2.replaceByRe(%re("/\W/g"), ""),
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

  <Form on_Click leg="Sign in" validationError cognitoError>
    <Input value=username propName="username" setFunc=setUsername />
    <Input value=password propName="password" autoComplete="current-password" setFunc=setPassword />
  </Form>
}
