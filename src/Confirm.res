type propShape = {"cognitoUser": Js.Nullable.t<Cognito.usr>, "search": Route.query}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

let code_max_length = 6
let password_max_length = 98

@react.component
let make = (~cognitoUser, ~search) => {
  let valErrInit = {
    open Route
    switch search {
    | VerificationCode => "CODE: 6-digit number only; "
    | ForgotPassword
    | ForgotUsername
    | Other => "CODE: 6-digit number only; PASSWORD: 8-98 length; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; "
    }
  }
  Js.log3("user", cognitoUser, search)

  let (code, setCode) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErrInit))
  let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)
  React.useEffect3(() => {
    open ErrorHook
    switch search {
    | VerificationCode => useError(code, Code, setValidationError)
    | ForgotPassword | ForgotUsername | Other =>
      useMultiError([(code, Code), (password, Password)], setValidationError)
    }
    None
  }, (code, password, search))

  let confirmregistrationCallback = (. err, res) =>
    switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
    | (_, Some(val)) => {
        setCognitoError(._ => None)
        Route.push(SignIn)
        Js.log2("conf res", val)
      }

    | (Some(ex), _) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoError(._ => Some(msg))
        | None => setCognitoError(._ => Some("unknown confirm error"))
        }
        Js.log2("conf problem", ex)
      }

    | _ => Js.Exn.raiseError("invalid cb argument")
    }

  open Cognito
  let confirmpasswordCallback = {
    onSuccess: str => {
      setCognitoError(._ => None)
      Route.push(SignIn)
      Js.log2("pw confirmed: ", str)
    },
    onFailure: err => {
      switch Js.Exn.message(err) {
      | Some(msg) => setCognitoError(._ => Some(msg))
      | None => setCognitoError(._ => Some("unknown confirm pw error"))
      }
      Js.log2("confirm pw problem: ", err)
    },
  }

  let on_Click = (. ()) => {
    switch validationError {
    | None =>
      switch Js.Nullable.isNullable(cognitoUser) {
      | false =>
        switch search {
        | VerificationCode =>
          cognitoUser->confirmRegistration(
            code
            ->Js.String2.slice(~from=0, ~to_=code_max_length)
            ->Js.String2.replaceByRe(%re("/\D/g"), ""),
            false,
            confirmregistrationCallback,
            Js.Nullable.null,
          )
        | ForgotPassword =>
          cognitoUser->confirmPassword(
            code
            ->Js.String2.slice(~from=0, ~to_=code_max_length)
            ->Js.String2.replaceByRe(%re("/\D/g"), ""),
            password
            ->Js.String2.slice(~from=0, ~to_=password_max_length)
            ->Js.String2.replaceByRe(%re("/\s/g"), ""),
            confirmpasswordCallback,
            Js.Nullable.null,
          )
        | ForgotUsername | Other => setCognitoError(._ => Some("unknown method - not submitting"))
        }
      | true => setCognitoError(._ => Some("null user - not submitting"))
      }
    | Some(_) => ()
    }
  }

  <Form
    ht="h-52"
    on_Click
    leg={switch search {
    | ForgotPassword => "Change password"
    | VerificationCode | ForgotUsername | Other => "Confirm code"
    }}
    validationError
    cognitoError>
    <Input
      value=code propName="code" autoComplete="one-time-code" inputMode="numeric" setFunc=setCode
    />
    {switch search {
    | ForgotPassword =>
      <Input value=password propName="password" autoComplete="new-password" setFunc=setPassword />
    | VerificationCode | ForgotUsername | Other => React.null
    }}
  </Form>
}
