@react.component
let make = (
  ~userpool,
  ~cognitoUser,
  ~setCognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~setShowName,
  ~search,
) => {
  let dummyPassword = "lllLLL!!!111"
  let dummyUsername = "letmein"
  let valErrInit = switch search {
  | "un_em" => "EMAIL: 5-99 length; enter a valid email address."
  | _ => "USERNAME: 3-10 length; "
  }
  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(valErrInit))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  Js.log("getinfo")
  let name_starts_index = 41
  let username_max_length = 10
  let email_max_length = 99

  React.useEffect2(() => {
    switch search {
    | "un_em" => ErrorHook.useError(email, "EMAIL", setValidationError)
    | _ => ErrorHook.useError(username, "USERNAME", setValidationError)
    }
    None
  }, (username, email))

  let signupCallback = (. err, res) => {
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
              | true => RescriptReactRouter.push(`/confirm?${search}`)
              | false => {
                  RescriptReactRouter.push("/")
                  setShowName(._ => Js.String2.sliceToEnd(msg, ~from=name_starts_index))
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
    | false => RescriptReactRouter.push(`/confirm?${search}`)
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
            username: username
            ->Js.String2.slice(~from=0, ~to_=username_max_length)
            ->Js.String2.replaceByRe(%re("/\W/g"), ""),
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
    | "un_em" =>
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
    | _ => ()
    }
  }

  <Form
    ht="h-52"
    onClick={onClick(search)}
    leg={switch search {
    | "un_em" => "Enter email"
    | _ => "Enter username"
    }}
    submitClicked
    validationError
    cognitoError>
    {switch search {
    | "un_em" => <Input value=email propName="email" inputMode="email" setFunc=setEmail />
    | _ => <Input value=username propName="username" setFunc=setUsername />
    }}
  </Form>
}
