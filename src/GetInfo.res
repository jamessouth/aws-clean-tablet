let dummyPassword = "lllLLL!!!111"
let dummyUsername = "letmein"

@react.component
let make = (
  ~userpool,
  ~cognitoUser,
  ~setCognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~setShowName,
) => {
  let url = RescriptReactRouter.useUrl()

  let valErrInit = switch url.search {
  | "un_em" => "EMAIL: 5-99 characters; enter a valid email address."
  | _ => "USERNAME: 3-10 characters; "
  }
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErrInit))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  Js.log("getinfo")

  React.useEffect2(() => {
    switch url.search {
    | "un_em" => ErrorHook.useError(email, "EMAIL", setValidationError)
    | _ => ErrorHook.useError(username, "USERNAME", setValidationError)
    }

    None
  }, (username, email))

  let signupCallback = (. err, res) => {
    Js.log2("signup cb", url.search)
    switch (Js.Nullable.toOption(err), Js.Nullable.toOption(res)) {
    | (_, Some(_)) => ()
    | (Some(ex), _) => {
        switch Js.Exn.message(ex) {
        | Some(msg) =>
          switch Js.String2.startsWith(msg, "PreSignUp failed with error user found") {
          | true => {
              setCognitoError(._ => None)
              switch Js.String2.endsWith(msg, "error user found.") {
              | true => RescriptReactRouter.push(`/confirm?${url.search}`)
              | false => {
                  RescriptReactRouter.push("/")
                  setShowName(._ => Js.String2.sliceToEnd(msg, ~from=41))
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
    Js.log2("coguser useeff", url.search)
    switch Js.Nullable.isNullable(cognitoUser) {
    | true => ()
    | false => RescriptReactRouter.push(`/confirm?${url.search}`)
    }
    None
  }, [cognitoUser])

  let onClick = (tipe, _) => {
    open Cognito
    setSubmitClicked(._ => true)
    switch tipe {
    | "cd_un" =>
      switch validationError {
      | None => {
          let userdata: userDataInput = {
            username: username,
            pool: userpool,
          }
          setCognitoUser(._ => Js.Nullable.return(userConstructor(userdata)))
        }
      | Some(_) => ()
      }
    | "pw_un" =>
      switch validationError {
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
      switch validationError {
      | None => {
          let emailData: attributeDataInput = {
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

  <Form
    ht="h-52"
    onClick={onClick(url.search)}
    leg={switch url.search {
    | "un_em" => "Enter email"
    | _ => "Enter username"
    }}
    submitClicked
    validationError
    cognitoError>
    {switch url.search {
    | "un_em" => <Input value=email propName="email" inputMode="email" setFunc=setEmail />
    | _ => <Input value=username propName="username" setFunc=setUsername />
    }}
  </Form>
}
