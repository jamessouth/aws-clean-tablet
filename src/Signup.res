type propShape = {
  "userpool": Cognito.poolData,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

let email_max_length = 99
let username_max_length = 10
let password_max_length = 98

@react.component
let make = (~userpool, ~setCognitoUser) => {
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 length; PASSWORD: 8-98 length; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; EMAIL: 5-99 length; enter a valid email address.",
  ))

  React.useEffect3(() => {
    ErrorHook.useMultiError(
      [(username, Username), (password, Password), (email, Email)],
      setValidationError,
    )
    None
  }, (username, password, email))

  open Cognito
  let signupCallback = (. err, res) =>
    switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
    | (_, Some(val)) => {
        setCognitoError(._ => None)
        setCognitoUser(._ => Js.Nullable.return(val.user))
        Route.push(Confirm({search: VerificationCode}))

        Js.log2("res", val.user.username)
      }

    | (Some(ex), _) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoError(._ => Some(msg))
        | None => setCognitoError(._ => Some("unknown signup error"))
        }

        Js.log2("problem", ex)
      }

    | _ => Js.Exn.raiseError("invalid cb argument")
    }

  let on_Click = (. ()) => {
    switch validationError {
    | None => {
        let emailData = {
          name: "email",
          value: email->Js.String2.slice(~from=0, ~to_=email_max_length),
        }
        let emailAttr = userAttributeConstructor(emailData)
        userpool->signUp(
          username
          ->Js.String2.slice(~from=0, ~to_=username_max_length)
          ->Js.String2.replaceByRe(%re("/\W/g"), ""),
          password
          ->Js.String2.slice(~from=0, ~to_=password_max_length)
          ->Js.String2.replaceByRe(%re("/\s/g"), ""),
          Js.Nullable.return([emailAttr]),
          Js.Nullable.null,
          signupCallback,
          Js.Nullable.null,
        )
      }

    | Some(_) => ()
    }
  }

  <Form on_Click leg="Sign up" validationError cognitoError>
    <Input value=username propName="username" setFunc=setUsername />
    <Input value=password propName="password" autoComplete="new-password" setFunc=setPassword />
    <Input value=email propName="email" inputMode="email" setFunc=setEmail />
  </Form>
}
