@send external focus: Dom.element => unit = "focus"

@new @module("amazon-cognito-identity-js")
external userAttributeConstructor: Types.attributeDataInput => Types.attributeData =
  "CognitoUserAttribute"

type clientMetadata = {key: string}
// type cdd = {
//     @as("AttributeName") attributeName: string,
//     @as("DeliveryMedium") deliveryMedium: string,
//     @as("Destination") destination: string
// }

// type clnt = {
//     endpoint: string,
//     fetchOptions: {}
// }

// type pl = {
//     advancedSecurityDataCollectionFlag: bool,
//     client: clnt,
//     clientId: string,
//     storage: {"length": float},
//     userPoolId: string
// }

type accessToken = {jwtToken: string}

type userSession = {
  // @as("IdToken") idToken: idToken,
  // @as("RefreshToken") refreshToken: string,
  accessToken: accessToken,
  // @as("ClockDrift") clockDrift: int
}

type usr = {
  // @as("Session") session: Js.Nullable.t<userSession>,
  // authenticationFlowType: string,
  // client: clnt,
  // keyPrefix: string,
  // pool: pl,
  // signInUserSession: Js.Nullable.t<string>,
  // storage: {"length": float},
  // userDataKey: string,
  username: string,
}

type signupOk = {
  // codeDeliveryDetails: cdd,
  user: usr,
  // userConfirmed: bool,
  // userSub: string
}
// type signupResult = result<signupOk, Js.Exn.t>

type signUpCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<signupOk>) => unit

@send
external signUp: (
  Types.poolData,
  string,
  string,
  Js.Nullable.t<array<Types.attributeData>>,
  Js.Nullable.t<array<Types.attributeData>>,
  signUpCB,
  Js.Nullable.t<clientMetadata>,
) => unit = "signUp"

let cbToOption = (f, . err, res) =>
  switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
  | (Some(err), _) => f(Error(err))
  | (_, Some(res)) => f(Ok(res))
  | _ => invalid_arg("invalid argument for cbToOption")
  }

@react.component
let make = (
  ~userpool,
  ~setCognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~usernameFuncList,
  ~passwordFuncList,
  ~emailFuncList,
) => {
  let (username, setUsername) = React.useState(_ => "")
  let (password, setPassword) = React.useState(_ => "")
  let (email, setEmail) = React.useState(_ => "")

  let (usernameError, setUsernameError) = React.useState(_ => Some("username: 3-10 characters; "))
  let (passwordError, setPasswordError) = React.useState(_ => Some(
    "password: 8-98 characters; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; ",
  ))
  let (emailError, setEmailError) = React.useState(_ => Some(
    "email: 5-99 characters; enter a valid email address.",
  ))

  let (validationError, setValidationError) = React.useState(_ => Some(
    "username: 3-10 characters; ",
  ))

  let (submitClicked, setSubmitClicked) = React.useState(_ => false)
  let (showPassword, setShowPassword) = React.useState(_ => false)

  React.useEffect3(() => {
    switch (usernameError, passwordError, emailError) {
    | (None, None, None) => setValidationError(_ => None)
    | (Some(err), _, _) | (_, Some(err), _) | (_, _, Some(err)) =>
      setValidationError(_ => Some(err))
    }
    None
  }, (usernameError, passwordError, emailError))

  let toggleButton = React.useMemo1(
    _ => <Toggle toggleProp=showPassword toggleSetFunc=setShowPassword />,
    [showPassword],
  )

  let signupCallback = cbToOption(res =>
    switch res {
    | Ok(val) => {
        setCognitoError(_ => None)
        setCognitoUser(._ => Js.Nullable.return(val.user))
        RescriptReactRouter.push("/confirm?code")

        Js.log2("res", val.user.username)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoError(_ => Some(msg))
        | None => setCognitoError(_ => Some("unknown signup error"))
        }

        Js.log2("problem", ex)
      }
    }
  )

  let onClick = _ => {
    setSubmitClicked(_ => true)
    switch validationError {
    | None => {
        let emailData: Types.attributeDataInput = {
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
          Js.Nullable.return({key: "sprbl"}),
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
        | true => <Error validationError cognitoError />
        }}
        <Input
          submitClicked
          value=username
          setFunc=setUsername
          setErrorFunc=setUsernameError
          funcList=usernameFuncList
          propName="username"
          validationError
        />
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
        <Input
          submitClicked
          value=email
          setFunc=setEmail
          setErrorFunc=setEmailError
          funcList=emailFuncList
          propName="email"
          validationError
        />
      </fieldset>
      <Button text="create" onClick/>
    </form>
  </main>
}
