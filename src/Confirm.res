let className = "text-gray-700 mt-14 bg-warm-gray-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

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
        | true =>
          switch (validationError, cognitoError) {
          | (Some(error), _) | (_, Some(error)) =>
            <span
              className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
              {React.string(error)}
            </span>
          | (None, None) => React.null
          }
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
      <Button
        textTrue="confirm" textFalse="confirm" textProp=true onClick disabled=false className
      />
    </form>
  </main>
}
