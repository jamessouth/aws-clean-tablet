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
let make = (~userpool, ~setCognitoUser, ~cognitoErr, ~setCognitoErr) => {
  // let (unVisited, setUnVisited) = React.useState(_ => false)

  let (validationError, setValidationError) = React.useState(_ => None)
  let (username, setUsername) = React.useState(_ => "")

  let (password, setPassword) = React.useState(_ => "")
  let (email, setEmail) = React.useState(_ => "")

  // let (cognitoResult, setCognitoResult) = React.useState(_ => false)

  let signupCallback = cbToOption(res =>
    switch res {
    | Ok(val) => {
        setCognitoErr(_ => None)
        setCognitoUser(._ => Js.Nullable.return(val.user))
        RescriptReactRouter.push("/confirm")

        Js.log2("res", val.user.username)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoErr(_ => Some(msg))
        | None => setCognitoErr(_ => Some("unknown signup error"))
        }

        Js.log2("problem", ex)
      }
    }
  )

  let usernameError = UsernameValidation.useUsernameValidation(username)
  let passwordError = PasswordValidation.usePasswordValidation(password)

  let onClick = _ => {
    switch (usernameError, passwordError) {
    | (None, None) => {
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
          Js.Nullable.null,
        )
      }
    | (Some(err), _) | (_, Some(err)) => setValidationError(_ => Some(err))
    }
  }

  <main>
    <form className="w-4/5 m-auto relative">
      <fieldset className="flex flex-col items-center justify-around h-72">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {React.string("Sign up")}
        </legend>
        {switch validationError {
        | Some(err) =>
          <span
            className="absolute right-0 top-0 text-sm text-warm-gray-100 bg-red-600 font-anon w-3/4 leading-4 p-1">
            {React.string(err)}
          </span>
        | None => React.null
        }}
        <Username username setUsername />
        <Password password setPassword />
        <Email email setEmail />
      </fieldset>
      {switch cognitoErr {
      | Some(msg) =>
        <span
          className="text-sm text-warm-gray-100 absolute bg-red-500 text-center w-full left-1/2 transform max-w-lg -translate-x-1/2">
          {React.string(msg)}
        </span>
      | None => React.null
      }}
      <button
        type_="button"
        className="text-gray-700 mt-14 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"
        onClick>
        {React.string("create")}
      </button>
    </form>
  </main>
}
