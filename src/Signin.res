let className = "text-gray-700 mt-14 bg-warm-gray-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

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
  Js.log("sinin")
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

  <main>
    <form className="w-4/5 m-auto relative">
      <fieldset className="flex flex-col items-center justify-around h-72">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {React.string("Sign in")}
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
      </fieldset>
      <Button textTrue="submit" textFalse="submit" textProp=true onClick disabled=false className />
    </form>
  </main>
}
