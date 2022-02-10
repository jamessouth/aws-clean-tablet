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
let make = (
  ~userpool,
  ~cognitoUser,
  ~setCognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~usernameFuncList,
) => {
  // let pwInput = React.useRef(Js.Nullable.null)

  let url = RescriptReactRouter.useUrl()
  Js.log2("url", url)

  let (username, setUsername) = React.useState(_ => "")
  let (usernameError, setUsernameError) = React.useState(_ => Some("username: 3-10 characters; "))

  let (submitClicked, setSubmitClicked) = React.useState(_ => false)

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

  let onClick = _ => {
    setSubmitClicked(_ => true)
    switch usernameError {
    | None => {
        let userdata: Types.userDataInput = {
          username: username,
          pool: userpool,
        }
        setCognitoUser(._ => Js.Nullable.return(userConstructor(userdata)))
      }
    | Some(_) => ()
    }
  }

  switch url.search {
  | "code" | "pw" =>
    <main>
      <form className="w-4/5 m-auto relative">
        <fieldset className="flex flex-col items-center justify-around h-52">
          <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
            {React.string("Enter username")}
          </legend>
          {switch submitClicked {
          | false => React.null
          | true => <Error validationError=usernameError cognitoError />
          }}
          <Input
            submitClicked
            value=username
            setFunc=setUsername
            setErrorFunc=setUsernameError
            funcList=usernameFuncList
            propName="username"
            validationError=usernameError
          />
        </fieldset>
        <Button text="submit" onClick />
      </form>
    </main>
  | _ => <div> {React.string("other")} </div>
  }
}
