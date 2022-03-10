let className = "text-gray-700 mt-14 bg-warm-gray-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

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

  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (usernameError, setUsernameError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 characters; ",
  ))
  let (email, setEmail) = React.Uncurried.useState(_ => "")
  let (emailError, setEmailError) = React.Uncurried.useState(_ => Some(
    "EMAIL: 5-99 characters; enter a valid email address.",
  ))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  Js.log4("url", emailError, usernameError, cognitoError)

  React.useEffect1(() => {
    ErrorHook.useError(username, "USERNAME", setUsernameError)
    None
  }, [username])

  React.useEffect1(() => {
    ErrorHook.useError(email, "EMAIL", setEmailError)
    None
  }, [email])

  let dummyPassword = "lllLLL!!!111"
  let dummyUsername = "letmein"

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
      switch usernameError {
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
          switch url.search {
          | "un_em" =>
            switch (emailError, cognitoError) {
            | (Some(error), _) | (_, Some(error)) =>
              <span
                className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
                {React.string(error)}
              </span>
            | (None, None) => React.null
            }
          | _ =>
            switch (usernameError, cognitoError) {
            | (Some(error), _) | (_, Some(error)) =>
              <span
                className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
                {React.string(error)}
              </span>
            | (None, None) => React.null
            }
          }
        }}
        {switch url.search {
        | "un_em" => <Input value=email propName="email" inputMode="email" setFunc=setEmail />
        | _ => <Input value=username propName="username" setFunc=setUsername />
        }}
      </fieldset>
      <Button
        textTrue="submit"
        textFalse="submit"
        textProp=true
        onClick={onClick(url.search)}
        disabled=false
        className
      />
    </form>
  </main>
}
