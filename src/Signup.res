let className = "text-gray-700 mt-14 bg-warm-gray-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

@react.component
let make = (~userpool, ~setCognitoUser, ~cognitoError, ~setCognitoError) => {
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  Js.log("url22")

  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 characters; PASSWORD: 8-98 characters; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; EMAIL: 5-99 characters; enter a valid email address.",
  ))

  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)

  React.useEffect3(() => {
    ErrorHook.useMultiError(
      [(username, "USERNAME"), (password, "PASSWORD"), (email, "EMAIL")],
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
        RescriptReactRouter.push("/confirm?cd_un")

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

  let onClick = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None => {
        let emailData = {
          name: "email",
          value: email,
        }
        let emailAttr = userAttributeConstructor(emailData)
        userpool->signUp(
          username,
          password,
          Js.Nullable.return([emailAttr]),
          Js.Nullable.null,
          signupCallback,
          Js.Nullable.null,
        )
      }
    | Some(_) => ()
    }
  }

  <main>
    <form className="w-4/5 m-auto relative">
      <fieldset className="flex flex-col items-center justify-around h-72">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {React.string("Sign up")}
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
        <Input value=username propName="username" setFunc=setUsername />
        <Input
          value=password propName="password" autoComplete="current-password" setFunc=setPassword
        />
        <Input value=email propName="email" inputMode="email" setFunc=setEmail />
      </fieldset>
      <Button textTrue="create" textFalse="create" textProp=true onClick disabled=false className />
    </form>
  </main>
}
