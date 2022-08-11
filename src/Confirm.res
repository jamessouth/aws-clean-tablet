type propShape = {
  "cognitoError": option<string>,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "search": string,
  "setCognitoError": (. option<string> => option<string>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make = (~cognitoUser, ~cognitoError, ~setCognitoError, ~search) => {
  let valErrInit = switch search {
  | "cd_un" => "CODE: 6-digit number only; "
  | _ => "CODE: 6-digit number only; PASSWORD: 8-98 length; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; "
  }
  Js.log3("user", cognitoUser, search)

  let (code, setCode) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErrInit))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  let code_max_length = 6
  let password_max_length = 98

  React.useEffect3(() => {
    switch search {
    | "cd_un" => ErrorHook.useError(code, "CODE", setValidationError)
    | _ => ErrorHook.useMultiError([(code, "CODE"), (password, "PASSWORD")], setValidationError)
    }
    None
  }, (code, password, search))

  let confirmregistrationCallback = (. err, res) =>
    switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
    | (_, Some(val)) => {
        setCognitoError(._ => None)
        RescriptReactRouter.push("/signin")
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
      RescriptReactRouter.push("/signin")
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

  let onClick = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None =>
      switch Js.Nullable.isNullable(cognitoUser) {
      | false =>
        switch search {
        | "cd_un" =>
          cognitoUser->confirmRegistration(
            code
            ->Js.String2.slice(~from=0, ~to_=code_max_length)
            ->Js.String2.replaceByRe(%re("/\D/g"), ""),
            false,
            confirmregistrationCallback,
            Js.Nullable.null,
          )
        | "pw_un" =>
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
        | _ => setCognitoError(._ => Some("unknown method - not submitting"))
        }
      | true => setCognitoError(._ => Some("null user - not submitting"))
      }
    | Some(_) => ()
    }
  }

  <Form
    ht="h-52"
    onClick
    leg={switch search {
    | "pw_un" => "Change password"
    | _ => "Confirm code"
    }}
    submitClicked
    validationError
    cognitoError>
    <Input
      value=code propName="code" autoComplete="one-time-code" inputMode="numeric" setFunc=setCode
    />
    {switch search {
    | "pw_un" =>
      <Input value=password propName="password" autoComplete="new-password" setFunc=setPassword />
    | _ => React.null
    }}
  </Form>
}
