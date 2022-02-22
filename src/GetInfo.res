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

// @send
// external forgotPassword: (
//   Js.Nullable.t<Signup.usr>, //user object
//   passwordPWCB, //cb obj
//   Js.Nullable.t<Signup.clientMetadata>,
// ) => unit = "forgotPassword"

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
  ~emailFuncList,
  ~setShowName,
) => {
  // let pwInput = React.useRef(Js.Nullable.null)

  let url = RescriptReactRouter.useUrl()
  Js.log2("url", url)

  let (username, setUsername) = React.useState(_ => "")
  let (usernameError, setUsernameError) = React.useState(_ => Some("username: 3-10 characters; "))

  let (email, setEmail) = React.useState(_ => "")
  let (emailError, setEmailError) = React.useState(_ => Some(
    "email: 5-99 characters; enter a valid email address.",
  ))

  let (submitClicked, setSubmitClicked) = React.useState(_ => false)

  // let dummyEmail = "success@simulator.amazonses.com"
  let dummyPassword = "lllLLL!!!111"
  let dummyUsername = "letmein"

  let signupCallback = cbToOption(res =>
    Js.log2("signup cb", url.search)
    switch res {
    | Ok(val) => ()
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) =>
          switch Js.String2.startsWith(msg, "PreSignUp failed with error user found") {
          | true => {
              setCognitoError(_ => None)

              switch Js.String2.endsWith(msg, "error user found.") {
              | true => RescriptReactRouter.push(`/confirm?${url.search}`)
              | false => {
                  RescriptReactRouter.push("/")
                  setShowName(._ => Js.String2.sliceToEnd(msg, ~from=41))
                }
              }
            }

          | false => setCognitoError(_ => Some(msg))
          }
        | None => setCognitoError(_ => Some("unknown signup error"))
        }
        Js.log2("problem", ex)
      }
    }
  )

  React.useEffect1(() => {
    Js.log2("coguser useeff", url.search)
    switch Js.Nullable.isNullable(cognitoUser) {
    | true => ()
    | false => RescriptReactRouter.push(`/confirm?${url.search}`)
    }
    None
  }, [cognitoUser])

  let onClick = (tipe, _) => {
    setSubmitClicked(_ => true)
    switch tipe {
    | "cd_un" =>
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
    | "pw_un" =>
      switch usernameError {
      | None =>
        userpool->signUp(
          username,
          dummyPassword,
          Js.Nullable.null,
          Js.Nullable.null,
          signupCallback,
          Js.Nullable.return({key: "forgotpassword"}),
        )
      | Some(_) => ()
      }
    | "un_em" =>
      switch emailError {
      | None => {
          let emailData: Types.attributeDataInput = {
            name: "email",
            value: email,
          }
          let emailAttr = userAttributeConstructor(emailData)
          userpool->signUp(
            dummyUsername,
            dummyPassword,
            Js.Nullable.return([emailAttr]),
            Js.Nullable.null,
            signupCallback,
            Js.Nullable.return({key: "forgotusername"}),
          )
        }
      | Some(_) => ()
      }
    | _ => ()
    }
  }

  switch url.search {
  | "cd_un" | "pw_un" | "un_em" =>
    <main>
      <form className="w-4/5 m-auto relative">
        <fieldset className="flex flex-col items-center justify-around h-52">
          <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
            {switch url.search {
            | "un_em" => React.string("Enter email")
            | _ => React.string("Enter username")
            }}
          </legend>
          {switch submitClicked {
          | false => React.null
          | true =>
            <Error
              validationError={switch url.search {
              | "un_em" => emailError
              | _ => usernameError
              }}
              cognitoError
            />
          }}
          {switch url.search {
          | "un_em" =>
            <Input
              submitClicked
              value=email
              setFunc=setEmail
              setErrorFunc=setEmailError
              funcList=emailFuncList
              propName="email"
              inputMode="email"
              validationError=emailError
            />
          | _ =>
            <Input
              submitClicked
              value=username
              setFunc=setUsername
              setErrorFunc=setUsernameError
              funcList=usernameFuncList
              propName="username"
              validationError=usernameError
            />
          }}
        </fieldset>
        <Button text="submit" onClick={onClick(url.search)} />
      </form>
    </main>
  | _ =>
    <div className="text-warm-gray-100"> {React.string("unknown path, please try again")} </div>
  }
}
