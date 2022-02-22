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
  Js.log3("user", cognitoUser, url)
  let (password, setPassword) = React.useState(_ => "")
  let (showPassword, setShowPassword) = React.useState(_ => false)

  let (pwErr, setPwErr) = React.useState(_ => None)
  let (showVerifCode, setShowVerifCode) = React.useState(_ => false)
  let (verifCode, setVerifCode) = React.useState(_ => "")

  let confirmregistrationCallback = Signup.cbToOption(res =>
    switch res {
    | Ok(val) => {
        setCognitoError(_ => None)
        RescriptReactRouter.push("/signin")
        Js.log2("conf rego res", val)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoError(_ => Some(msg))
        | None => setCognitoError(_ => Some("unknown confirm rego error"))
        }

        Js.log2("conf rego problem", ex)
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
    switch Js.Nullable.isNullable(cognitoUser) {
    | false =>
      switch url.search {
      | "code" =>
        cognitoUser->confirmRegistration(
          verifCode,
          false,
          confirmregistrationCallback,
          Js.Nullable.null,
        )
      | "pw" =>
        cognitoUser->confirmPassword(verifCode, password, confirmpasswordCallback, Js.Nullable.null)
      | _ => (_ => Some("unknown method - not submitting"))->setCognitoError
      }

    | true => (_ => Some("null user - not submitting"))->setCognitoError
    }
  }

  switch url.search {
  | "cd_un" | "pw_un" =>
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
          <Input //code
            submitClicked
            value=password
            setFunc=setPassword
            setErrorFunc=setPasswordError
            funcList=passwordFuncList
            propName="code"
            autoComplete="one-time-code"
            inputMode="numeric"
            validationError
          />
          {switch url.search {
          | "pw_un" =>
            <Input
              submitClicked
              value=password
              setFunc=setPassword
              setErrorFunc=setPasswordError
              funcList=passwordFuncList
              propName="password"
              autoComplete="new-password"
              toggleProp=showPassword
              toggleButton
              validationError
            />

          | _ => React.null
          }}
        </fieldset>
        <Button text="confirm" onClick />
      </form>
    </main>
  | _ => <div> {React.string("other")} </div>
  }
}
