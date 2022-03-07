type t

type confirmRegistrationCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<t>) => unit

@send
external confirmRegistration: (
  Js.Nullable.t<Signup.usr>,
  string,
  bool,
  confirmRegistrationCB,
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "confirmRegistration"

@send
external confirmPassword: (
  Js.Nullable.t<Signup.usr>, //user object
  string, //conf code
  string, //new pw
  GetInfo.passwordPWCB, //cb obj
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "confirmPassword"

@react.component
let make = (~cognitoUser, ~cognitoError, ~setCognitoError) => {
  let url = RescriptReactRouter.useUrl()
  let valErr = switch url.search {
  | "cd_un" => "CODE: 6-digit number only; "
  | _ => "CODE: 6-digit number only; PASSWORD: 8-98 characters; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; "
  }
  Js.log3("user", cognitoUser, url)

  let (code, setCode) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")

  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErr))

  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)

  React.useEffect3(() => {
    switch url.search {
    | "cd_un" => ErrorHook.useError(code, "CODE", setValidationError)
    | _ => ErrorHook.useMultiError([(code, "CODE"), (password, "PASSWORD")], setValidationError)
    }
    None
  }, (code, password, url.search))

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
  

  let confirmpasswordCallback: GetInfo.passwordPWCB = {
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
        switch url.search {
        | "cd_un" =>
          cognitoUser->confirmRegistration(
            code,
            false,
            confirmregistrationCallback,
            Js.Nullable.null,
          )
        | "pw_un" =>
          cognitoUser->confirmPassword(code, password, confirmpasswordCallback, Js.Nullable.null)
        | _ => setCognitoError(._ => Some("unknown method - not submitting"))
        }
      | true => setCognitoError(._ => Some("null user - not submitting"))
      }
    | Some(_) => ()
    }
  }

  <main>
    <form className="w-4/5 m-auto relative">
      <fieldset className="flex flex-col justify-between h-52">
        <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
          {switch url.search {
          | "pw_un" => React.string("Change password")
          | _ => React.string("Confirm code")
          }}
        </legend>
        {switch submitClicked {
        | false => React.null
        | true => <Error validationError cognitoError />
        }}
        <Input
          value=code
          propName="code"
          autoComplete="one-time-code"
          inputMode="numeric"
          setFunc=setCode
        />
        {switch url.search {
        | "pw_un" =>
          <Input
            value=password propName="password" autoComplete="new-password" setFunc=setPassword
          />
        | _ => React.null
        }}
      </fieldset>
      <Button text="confirm" onClick />
    </form>
  </main>
}
