type propShape = {
  "userpool": Cognito.poolData,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setRetrievedUsername": (. Js.String2.t => Js.String2.t) => unit,
  "search": Route.query,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

let dummyUsername = "letmein"
let dummyPassword = "lllLLL!!!111"
let name_starts_index = 41
let username_max_length = 10
let email_max_length = 99

@react.component
let make = (~userpool, ~cognitoUser, ~setCognitoUser, ~setRetrievedUsername, ~search) => {
  let valErrInit = {
    open Route
    switch search {
    | ForgotUsername => "EMAIL: 5-99 length; enter a valid email address."
    | VerificationCode | ForgotPassword | Other => "USERNAME: 3-10 length; "
    }
  }
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErrInit))
  Js.log("getinfo")

  let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)
  React.useEffect2(() => {
    open Route
    switch search {
    | ForgotUsername => ErrorHook.useError(email, Email, setValidationError)
    | VerificationCode | ForgotPassword | Other =>
      ErrorHook.useError(username, Username, setValidationError)
    }
    None
  }, (username, email))

  let signupCallback = (. err, res) => {
    open Route
    Js.log2("signup cb", search)
    switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
    | (_, Some(_)) => ()
    | (Some(ex), _) => {
        switch Js.Exn.message(ex) {
        | Some(msg) =>
          switch Js.String2.startsWith(msg, "PreSignUp failed with error user found") {
          | true => {
              setCognitoError(._ => None)
              switch Js.String2.endsWith(msg, "error user found.") {
              | true => push(Confirm({search: search}))
              | false => {
                  push(SignIn)
                  setRetrievedUsername(._ => Js.String2.sliceToEnd(msg, ~from=name_starts_index))
                }
              }
            }

          | false => setCognitoError(._ => Some(msg))
          }
        | None => setCognitoError(._ => Some("unknown signup error"))
        }
        Js.log2("problem", ex)
      }

    | _ => Js.Exn.raiseError("invalid cb argument")
    }
  }

  React.useEffect1(() => {
    Js.log2("coguser useeff", search)
    switch Js.Nullable.isNullable(cognitoUser) {
    | true => ()
    | false => Route.push(Confirm({search: search}))
    }
    None
  }, [cognitoUser])

  let onClick = (. tipe, ()) => {
    open Cognito
    open Route
    switch tipe {
    | VerificationCode =>
      switch validationError {
      | None => {
          let userdata: userDataInput = {
            username: username
            ->Js.String2.slice(~from=0, ~to_=username_max_length)
            ->Js.String2.replaceByRe(%re("/\W/g"), ""),
            pool: userpool,
          }
          setCognitoUser(._ => Js.Nullable.return(userConstructor(userdata)))
        }

      | Some(_) => ()
      }
    | ForgotPassword =>
      switch validationError {
      | None =>
        userpool->signUp(
          username
          ->Js.String2.slice(~from=0, ~to_=username_max_length)
          ->Js.String2.replaceByRe(%re("/\W/g"), ""),
          dummyPassword,
          Js.Nullable.null,
          Js.Nullable.null,
          signupCallback,
          Js.Nullable.return({key: "forgotpassword"}),
        )
      | Some(_) => ()
      }
    | ForgotUsername =>
      switch validationError {
      | None => {
          let emailData: attributeDataInput = {
            name: "email",
            value: email->Js.String2.slice(~from=0, ~to_=email_max_length),
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
    | Other => ()
    }
  }

  <Form
    ht="h-52"
    on_Click={(. ()) => onClick(. search, ())}
    leg={switch search {
    | ForgotUsername => "Enter email"
    | VerificationCode | ForgotPassword | Other => "Enter username"
    }}
    validationError
    cognitoError>
    {switch search {
    | ForgotUsername => <Input value=email propName="email" inputMode="email" setFunc=setEmail />
    | VerificationCode | ForgotPassword | Other =>
      <Input value=username propName="username" setFunc=setUsername />
    }}
  </Form>
}
