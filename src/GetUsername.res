@new @module("amazon-cognito-identity-js")
external userAttributeConstructor: Types.attributeDataInput => Types.attributeData =
  "CognitoUserAttribute"

type clientMetadata = {key: string}

// type usr = {
//   // @as("Session") session: Js.Nullable.t<userSession>,
//   // authenticationFlowType: string,
//   // client: clnt,
//   // keyPrefix: string,
//   // pool: pl,
//   // signInUserSession: Js.Nullable.t<string>,
//   // storage: {"length": float},
//   // userDataKey: string,
//   username: string,
// }

type signupOk = {
  // codeDeliveryDetails: cdd,
  user: Signup.usr,
  // userConfirmed: bool,
  // userSub: string
}

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

// type usr = {
//   @as("Session") session: Js.Nullable.t<string>,
//   authenticationFlowType: string,
//   client: clnt,
//   keyPrefix: string,
//   pool: pl,
//   signInUserSession: Js.Nullable.t<string>,
//   storage: {"length": float},
//   userDataKey: string,
//   username: string,
// }

// type signupOk = {
//   codeDeliveryDetails: cdd,
//   user: usr,
//   userConfirmed: bool,
//   userSub: string,
// }

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

let cbToOption = (f, . err, res) =>
  switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
  | (Some(err), _) => f(Error(err))
  | (_, Some(res)) => f(Ok(res))
  | _ => invalid_arg("invalid argument for cbToOption")
  }

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

  // let forgotPWcb = {
  //   onSuccess: str => {
  //     setCognitoError(_ => None)
  //     Js.log2("forgot pw initiated: ", str)
  //     // RescriptReactRouter.push("/confirm")
  //   },
  //   onFailure: err => {
  //     switch Js.Exn.message(err) {
  //     | Some(msg) => setCognitoError(_ => Some(msg))
  //     | None => setCognitoError(_ => Some("unknown forgot pw error"))
  //     }
  //     Js.log2("forgot pw problem: ", err)
  //   },
  // }

  let dummyEmail = "success@simulator.amazonses.com"
  let dummyPassword = "lllLLL!!!111"

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
        | Some(msg) =>
          switch msg {
          | "PreSignUp failed with error user found." => {
              setCognitoError(_ => None)
              RescriptReactRouter.push(`/confirm?${url.search}`)
            }

          | _ => setCognitoError(_ => Some(msg))
          }
        | None => setCognitoError(_ => Some("unknown signup error"))
        }
        Js.log2("problem", ex)
      }
    }
  )

  React.useEffect1(() => {
    switch Js.Nullable.isNullable(cognitoUser) {
    | true => ()
    | false =>
      switch url.search {
      | "pw" => {
          let emailData: Types.attributeDataInput = {
            name: "email",
            value: dummyEmail,
          }
          let emailAttr = userAttributeConstructor(emailData)
          userpool->signUp(
            username,
            dummyPassword,
            Js.Nullable.return([emailAttr]),
            Js.Nullable.null,
            signupCallback,
            Js.Nullable.return({key: "fp"}),
          )
        }
      | _ => RescriptReactRouter.push(`/confirm?${url.search}`)
      }
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
