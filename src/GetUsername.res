@new @module("amazon-cognito-identity-js")
external userConstructor: Types.userDataInput => Signup.usr = "CognitoUser"

type cdd = {
  @as("AttributeName") attributeName: string,
  @as("DeliveryMedium") deliveryMedium: string,
  @as("Destination") destination: string,
}

type clnt = {
  endpoint: string,
  fetchOptions: {.},
}

type pl = {
  advancedSecurityDataCollectionFlag: bool,
  client: clnt,
  clientId: string,
  storage: {"length": float},
  userPoolId: string,
}

type usr = {
  @as("Session") session: Js.Nullable.t<string>,
  authenticationFlowType: string,
  client: clnt,
  keyPrefix: string,
  pool: pl,
  signInUserSession: Js.Nullable.t<string>,
  storage: {"length": float},
  userDataKey: string,
  username: string,
}

type signupOk = {
  codeDeliveryDetails: cdd,
  user: usr,
  userConfirmed: bool,
  userSub: string,
}

type passwordPWCB = {
  onFailure: Js.Exn.t => unit,
  onSuccess: string => unit,
}

@send
external forgotPassword: (
  Js.Nullable.t<Signup.usr>, //user object
  passwordPWCB, //cb obj
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "forgotPassword"

@react.component
let make = (~userpool, ~cognitoUser, ~setCognitoUser, ~cognitoError, ~setCognitoError) => {
  // let pwInput = React.useRef(Js.Nullable.null)

  let url = RescriptReactRouter.useUrl()
  Js.log2("url", url)

  let (unVisited, setUnVisited) = React.useState(_ => false)

  let (unErr, setUnErr) = React.useState(_ => None)

  let (disabled, _setDisabled) = React.useState(_ => false)
  let (username, setUsername) = React.useState(_ => "")

  
  // let (cognitoResult, setCognitoResult) = React.useState(_ => false)

  let checkUnForbiddenChars = un => {
    let r = %re("/\W/")

    switch Js.String2.match_(un, r) {
    | Some(_) => (_ => Some("alphanumeric..."))->setUnErr
    | None => (_ => None)->setUnErr
    }
  }

  let checkUnMaxLength = un => {
    switch un->Js.String2.length > 10 {
    | true => (_ => Some("too long..."))->setUnErr
    | false => un->checkUnForbiddenChars
    }
  }

  let checkNoUnWhitespace = un => {
    let r = %re("/\s/")

    switch Js.String2.match_(un, r) {
    | Some(_) => (_ => Some("no whitespace..."))->setUnErr
    | None => un->checkUnMaxLength
    }
  }

  let checkUnLength = un => {
    switch un->Js.String2.length < 4 {
    | true => (_ => Some("too short..."))->setUnErr
    | false => un->checkNoUnWhitespace
    }
  }

  let onChange = e => {
    let value = ReactEvent.Form.target(e)["value"]
    setUsername(_ => value)
  }

  let onBlur = _e => {
    setUnVisited(_ => true)
  }

  let forgotPWcb = {
    onSuccess: str => {
      setCognitoError(_ => None)
      Js.log2("forgot pw initiated: ", str)
      // RescriptReactRouter.push("/confirm")
    },
    onFailure: err => {
      switch Js.Exn.message(err) {
      | Some(msg) => setCognitoError(_ => Some(msg))
      | None => setCognitoError(_ => Some("unknown forgot pw error"))
      }
      Js.log2("forgot pw problem: ", err)
    },
  }

  React.useEffect1(() => {
    switch Js.Nullable.isNullable(cognitoUser) {
    | true => ()
    | false =>
      switch url.search {
      | "pw" => cognitoUser->forgotPassword(forgotPWcb, Js.Nullable.null)
      | _ => ()
      }

      RescriptReactRouter.push(`/confirm?${url.search}`)
    }
    None
  }, [cognitoUser])

  let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    let userdata: Types.userDataInput = {
      username: username,
      pool: userpool,
    }
    setCognitoUser(._ => Js.Nullable.return(userConstructor(userdata)))
  }

  React.useEffect2(() => {
    switch unVisited {
    | true => username->checkUnLength
    | false => (_ => None)->setUnErr
    }
    None
  }, (username, unVisited))

  // React.useEffect5(() => {
  //   switch (unErr, pwErr, username->Js.String2.length < 4, password->Js.String2.length < 8, email->Js.String2.length < 3) {
  //   | (None, None, false, false, false) => (_ => false)->setDisabled
  //   | (Some(_), _, _, _, _)
  //   | (_, Some(_), _, _, _)
  //   | (_, _, true, _, _)
  //   | (_, _, _, true, _)
  //   | (_, _, _, _, true) => (_ => true)->setDisabled
  //   }

  //   None
  // }, (unErr, pwErr, username, password, email))

  switch url.search {
  | "code" | "pw" =>
    <main>
      <form className="w-4/5 m-auto" onSubmit={handleSubmit}>
        <fieldset className="flex flex-col justify-between h-52">
          <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
            {React.string("Enter username")}
          </legend>
          <div className="relative">
            <label
              className={switch (unVisited, unErr) {
              | (true, Some(_)) => "block text-2xl text-red-500 font-bold font-flow"
              | (false, _) | (true, None) => "block text-2xl text-warm-gray-100 font-flow"
              }}
              htmlFor="username">
              {"username:"->React.string}
            </label>
       

            {switch (unVisited, unErr) {
          | (true, Some(err)) =>
            <span
              className="absolute right-0 text-lg text-warm-gray-100 bg-red-500 font-anon font-flow h-30 w-2/3 z-10">
              {err->React.string}
            </span>
          | (false, _) | (true, None) => React.null
          }}
            <input
              autoComplete="username"
              autoFocus=true
              className={switch (unVisited, unErr) {
              | (
                  true,
                  Some(_),
                ) => "h-8 w-full text-xl outline-none text-red-500 bg-transparent border-b-1 border-red-500"
              | (false, _)
              | (
                true,
                None,
              ) => "h-8 w-full text-xl outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
              }}
              id="username"
              minLength=4
              name="username"
              onBlur
              onChange
              // placeholder="Enter username"
              required=true
              spellCheck=false
              type_="text"
              value={username}
            />
          </div>
        </fieldset>
        {switch cognitoError {
        | Some(msg) =>
          <span
            className="text-sm text-warm-gray-100 absolute bg-red-500 text-center w-full left-1/2 transform max-w-lg -translate-x-1/2">
            {msg->React.string}
          </span>
        | None => React.null
        }}
        <button
          disabled
          className="text-gray-700 mt-10 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7">
          {"submit"->React.string}
        </button>
      </form>
    </main>
  | _ => <div> {"other"->React.string} </div>
  }
}
