@new @module("amazon-cognito-identity-js")
external userConstructor: Types.userDataInput => Signup.usr = "CognitoUser"

type callback = {
  onFailure: Js.Exn.t => unit,
  newPasswordRequired: Js.Nullable.t<
    (array<Types.attributeData>, array<Types.attributeData>) => unit,
  >,
  mfaRequired: Js.Nullable.t<(string, string) => unit>,
  customChallenge: Js.Nullable.t<string => unit>,
  onSuccess: Signup.userSession => unit,
}

type authDetails = {
  @as("ValidationData") validationData: Js.Nullable.t<array<Types.attributeData>>,
  @as("Username") username: string,
  @as("Password") password: string,
  @as("AuthParameters") authParameters: Js.Nullable.t<array<Types.attributeData>>,
  @as("ClientMetadata") clientMetadata: Js.Nullable.t<Signup.clientMetadata>,
}

@new @module("amazon-cognito-identity-js")
external authenticationDetailsConstructor: authDetails => authDetails = "AuthenticationDetails"

@send
external authenticateUser: (Js.Nullable.t<Signup.usr>, authDetails, callback) => unit =
  "authenticateUser"

@react.component
let make = (~userpool, ~setCognitoUser, ~setToken, ~cognitoUser) => {
  let (_cognitoErr, setCognitoErr) = React.useState(_ => None)
  let (showPassword, setShowPassword) = React.useState(_ => false)
  let (disabled, setDisabled) = React.useState(_ => true)
  let (username, setUsername) = React.useState(_ => "")
  let (password, setPassword) = React.useState(_ => "")
  let onChange = (func, e) => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->func
  }

  let onSubmit = e => {
    e->ReactEvent.Form.preventDefault
    let cbs = {
      onSuccess: res => {
        Js.log2("signin result:", res)
        setToken(._ => Some(res.accessToken.jwtToken))
      },
      onFailure: ex => {
        switch Js.Exn.message(ex) {
        | Some(msg) => (_ => Some(msg))->setCognitoErr
        | None => (_ => None)->setCognitoErr
        }

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
        let userdata: Types.userDataInput = {
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

  let onClick = _e => {
    (prev => !prev)->setShowPassword
  }

  React.useEffect2(() => {
    switch (username->Js.String2.length > 3, password->Js.String2.length > 7) {
    | (true, true) => (_ => false)->setDisabled
    | (false, true) | (true, false) | (false, false) => (_ => true)->setDisabled
    }

    None
  }, (username, password))

  <main>
    <form className="w-4/5 m-auto" onSubmit>
      <fieldset className="flex flex-col items-center justify-around h-80">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {"Sign in"->React.string}
        </legend>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
            {"username:"->React.string}
          </label>
          <input
            autoComplete="username"
            autoFocus=true
            className="h-6 w-full text-xl pl-1 text-left font-anon outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="username"
            minLength=4
            name="username"
            onChange={onChange(setUsername)}
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={username}
          />
        </div>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="password">
            {"password:"->React.string}
          </label>
          <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            onChange={onChange(setPassword)}
            // placeholder="Enter password"
            required=true
            spellCheck=false
            type_={switch showPassword {
            | true => "text"
            | false => "password"
            }}
            value={password}
          />
          <button
            type_="button"
            className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 cursor-pointer"
            onClick>
            {switch showPassword {
            | true => "hide"->React.string
            | false => "show"->React.string
            }}
          </button>
        </div>
      </fieldset>
      <button
        disabled
        className="text-gray-700 mt-16 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7">
        {"submit"->React.string}
      </button>
    </form>
  </main>
}
