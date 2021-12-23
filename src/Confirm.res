type t

type confirmRegistrationCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<t>) => unit

@send
external confirmRegistration: (
  Js.Nullable.t<Signup.usr>,
  string,
  bool,
  confirmRegistrationCB,
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "confirmRegistration"

@send
external confirmPassword: (
  Js.Nullable.t<Signup.usr>, //user object
  string, //conf code
  string, //new pw
  GetUsername.passwordPWCB, //cb obj
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "confirmPassword"

@react.component
let make = (~cognitoUser, ~cognitoErr, ~setCognitoErr) => {
  let url = RescriptReactRouter.useUrl()
  Js.log3("user", cognitoUser, url)
  let (password, setPassword) = React.useState(_ => "")
  let (showPassword, setShowPassword) = React.useState(_ => false)
  let (pwVisited, setPwVisited) = React.useState(_ => false)
  let (pwErr, setPwErr) = React.useState(_ => None)
  let (showVerifCode, setShowVerifCode) = React.useState(_ => false)
  let (verifCode, setVerifCode) = React.useState(_ => "")
  let (disabled, setDisabled) = React.useState(_ => true)

  

  let onClick = _e => {
    (prev => !prev)->setShowVerifCode
  }

  let onClick2 = _e => {
    (prev => !prev)->setShowPassword
  }

  let onBlur = (input, _e) => {
    switch input {
    // | "username" => (_ => true)->setUnVisited
    | "password" => (_ => true)->setPwVisited
    | _ => ()
    }
  }

  let onChange = (func, e) => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->func
  }

  React.useEffect1(() => {
    switch verifCode->Js.String2.length != 6 {
    | true => (_ => true)->setDisabled
    | false => (_ => false)->setDisabled
    }
    None
  }, [verifCode])

  let checkPwForbiddenChars = pw => {
    let r = %re("/[-=+]/")

    switch Js.String2.match_(pw, r) {
    | Some(_) => (_ => Some("no +, -, or = ..."))->setPwErr
    | None => (_ => None)->setPwErr
    }
  }

  let checkPwMaxLength = pw => {
    switch pw->Js.String2.length > 98 {
    | true => (_ => Some("too long..."))->setPwErr
    | false => pw->checkPwForbiddenChars
    }
  }

  let checkNoPwWhitespace = pw => {
    let r = %re("/\s/")

    switch Js.String2.match_(pw, r) {
    | Some(_) => (_ => Some("no whitespace..."))->setPwErr
    | None => pw->checkPwMaxLength
    }
  }

  let checkSymbol = pw => {
    let r = %re("/[!-*\[-`{-~./,:;<>?@]/")

    switch Js.String2.match_(pw, r) {
    | None => (_ => Some("add symbol..."))->setPwErr
    | Some(_) => pw->checkNoPwWhitespace
    }
  }

  let checkNumber = pw => {
    let r = %re("/\d/")

    switch Js.String2.match_(pw, r) {
    | None => (_ => Some("add number..."))->setPwErr
    | Some(_) => pw->checkSymbol
    }
  }

  let checkUpper = pw => {
    let r = %re("/[A-Z]/")

    switch Js.String2.match_(pw, r) {
    | None => (_ => Some("add uppercase..."))->setPwErr
    | Some(_) => pw->checkNumber
    }
  }

  let checkLower = pw => {
    let r = %re("/[a-z]/")

    switch Js.String2.match_(pw, r) {
    | None => (_ => Some("add lowercase..."))->setPwErr
    | Some(_) => pw->checkUpper
    }
  }

  let checkPwLength = pw => {
    switch pw->Js.String2.length < 8 {
    | true => (_ => Some("too short..."))->setPwErr
    | false => pw->checkLower
    }
  }

  React.useEffect2(() => {
    switch pwVisited {
    | true => password->checkPwLength
    | false => (_ => None)->setPwErr
    }
    None
  }, (password, pwVisited))

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

  let confirmregistrationCallback = Signup.cbToOption(res =>
    switch res {
    | Ok(val) => {
        setCognitoErr(_ => None)
        RescriptReactRouter.push("/signin")
        Js.log2("conf rego res", val)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => setCognitoErr(_ => Some(msg))
        | None => setCognitoErr(_ => Some("unknown confirm rego error"))
        }

        Js.log2("conf rego problem", ex)
      }
    }
  )

  let confirmpasswordCallback: GetUsername.passwordPWCB = {
    onSuccess: str => {
        setCognitoErr(_ => None)
        RescriptReactRouter.push("/signin")
      Js.log2("pw confirmed: ", str)
    },
    onFailure: err => {
      switch Js.Exn.message(err) {
      | Some(msg) => setCognitoErr(_ => Some(msg))
      | None => setCognitoErr(_ => Some("unknown confirm pw error"))
      }
      Js.log2("confirm pw problem: ", err)
    },
  }

  let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    switch Js.Nullable.isNullable(cognitoUser) {
    | false =>
      switch url.search {
      | "code" =>
        cognitoUser->confirmRegistration(
          verifCode,
          false,
          confirmregistrationCallback,
          Js.Nullable.null,
        )
      | "pw" => cognitoUser->confirmPassword(verifCode, password, confirmpasswordCallback, Js.Nullable.null)
      | _ => (_ => Some("unknown method - not submitting"))->setCognitoErr
      }

    | true => (_ => Some("null user - not submitting"))->setCognitoErr
    }
  }

  switch url.search {
  | "code" | "pw" =>
    <main>
      <form className="w-4/5 m-auto" onSubmit={handleSubmit}>
        <fieldset className="flex flex-col justify-between h-52">
          <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
            {switch url.search {
            | "code" => React.string("Confirm code")
            | "pw" => React.string("Change password")
            | _ => React.string("other")
            }}
          </legend>
          <div className="relative">
            <label className="block text-2xl text-warm-gray-100 font-flow" htmlFor="verif-code">
              {"enter code:"->React.string}
            </label>
            <input
              autoComplete="one-time-code"
              autoFocus=true
              className="h-8 text-xl outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
              id="verif-code"
              maxLength=6
              minLength=6
              inputMode="numeric"
              name="verifcode"
              onChange={onChange(setVerifCode)}
              // placeholder="Enter password"
              // ref={pwInput->ReactDOM.Ref.domRef}
              pattern="^\d{6}$"
              required=true
              size=6
              spellCheck=false
              type_={switch showVerifCode {
              | true => "text"
              | false => "password"
              }}
              value={verifCode}
            />
            <button
              type_="button"
              className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 cursor-pointer"
              onClick>
              {switch showVerifCode {
              | true => "hide"->React.string
              | false => "show"->React.string
              }}
            </button>
          </div>
          {switch url.search {
          | "pw" =>
            <div className="relative">
              <label
                className={switch (pwVisited, pwErr) {
                | (true, Some(_)) => "text-2xl text-red-500 font-bold font-flow"
                | (false, _) | (true, None) => "text-2xl text-warm-gray-100 font-flow"
                }}
                htmlFor="new-password">
                {"password:"->React.string}
              </label>
              {switch (pwVisited, pwErr) {
              | (true, Some(err)) =>
                <span className="absolute right-0 text-2xl text-red-500 font-bold font-flow">
                  {err->React.string}
                </span>
              | (false, _) | (true, None) => React.null
              }}
              <input
                autoComplete="new-password"
                autoFocus=false
                className={switch (pwVisited, pwErr) {
                | (
                    true,
                    Some(_),
                  ) => "h-8 w-3/4 text-xl outline-none text-red-500 bg-transparent border-b-1 border-red-500"
                | (false, _)
                | (
                  true,
                  None,
                ) => "h-8 w-3/4 text-xl outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
                }}
                // className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
                id="new-password"
                // minLength=8
                name="password"
                onBlur={onBlur("password")}
                onChange={onChange(setPassword)}
                // placeholder="Enter password"
                // ref={pwInput->ReactDOM.Ref.domRef}
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
                onClick=onClick2>
                {switch showPassword {
                | true => "hide"->React.string
                | false => "show"->React.string
                }}
              </button>
            </div>
          | _ => React.null
          }}
        </fieldset>
        {switch cognitoErr {
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
          {"confirm"->React.string}
        </button>
      </form>
    </main>
  | _ => <div> {"other"->React.string} </div>
  }
}
