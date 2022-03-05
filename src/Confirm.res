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

  let (code, setCode) = React.useState(_ => "")
  let (password, setPassword) = React.useState(_ => "")

  let (validationError, setValidationError) = React.useState(_ => Some(valErr))

  let (submitClicked, setSubmitClicked) = React.useState(_ => false)

  React.useEffect3(() => {
    switch url.search {
    | "cd_un" => ErrorHook.useError(code, "CODE", setValidationError)
    | _ => ErrorHook.useMultiError([(code, "CODE"), (password, "PASSWORD")], setValidationError)
    }
    None
  }, (code, password, url.search))

  let confirmregistrationCallback = Signup.cbToOption(res =>
    switch res {
    | Ok(val) => {
        setCognitoError(_ => None)
        RescriptReactRouter.push("/signin")
        Js.log2("conf res", val)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoError(_ => Some(msg))
        | None => setCognitoError(_ => Some("unknown confirm error"))
        }
        Js.log2("conf problem", ex)
      }
    }
  )

  let confirmpasswordCallback: GetInfo.passwordPWCB = {
    onSuccess: str => {
      setCognitoError(_ => None)
      RescriptReactRouter.push("/signin")
      Js.log2("pw confirmed: ", str)
    },
    onFailure: err => {
      switch Js.Exn.message(err) {
      | Some(msg) => setCognitoError(_ => Some(msg))
      | None => setCognitoError(_ => Some("unknown confirm pw error"))
      }
      Js.log2("confirm pw problem: ", err)
    },
  }

  let onClick = _ => {
    setSubmitClicked(_ => true)
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
        | _ => (_ => Some("unknown method - not submitting"))->setCognitoError
        }
      | true => (_ => Some("null user - not submitting"))->setCognitoError
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
        {React.useMemo1(_ => {
          <Input
            value=code
            propName="code"
            autoComplete="one-time-code"
            inputMode="numeric"
            setFunc=setCode
          />
        }, [code])}
        {switch url.search {
        | "pw_un" => React.useMemo1(_ => {
            <Input
              value=password propName="password" autoComplete="new-password" setFunc=setPassword
            />
          }, [password])
        | _ => React.null
        }}
      </fieldset>
      <Button text="confirm" onClick />
    </form>
  </main>
}
