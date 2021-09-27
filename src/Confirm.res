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

@react.component
let make = (~cognitoUser) => {
  Js.log2("user", cognitoUser)
  let (showVerifCode, setShowVerifCode) = React.useState(_ => false)
  let (verifCode, setVerifCode) = React.useState(_ => "")

  let (disabled, setDisabled) = React.useState(_ => true)

  let (cognitoErr, setCognitoErr) = React.useState(_ => None)

  let onClick = _e => {
    (prev => !prev)->setShowVerifCode
  }

  let onChange = e => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->setVerifCode
  }

  React.useEffect1(() => {
    switch verifCode->Js.String2.length != 6 {
    | true => (_ => true)->setDisabled
    | false => (_ => false)->setDisabled
    }
    None
  }, [verifCode])

  // testtest1!T
  let confirmregistrationCallback = Signup.cbToOption(res =>
    switch res {
    | Ok(val) => {
        (_ => None)->setCognitoErr
        // (_ => Some(val.user))->setCognitoUser
        // RescriptReactRouter.push("/confirm")

        Js.log2("conf rego res", val)
      }
    | Error(ex) => {
        switch Js.Exn.message(ex) {
        | Some(msg) => (_ => Some(msg))->setCognitoErr
        | None => (_ => None)->setCognitoErr
        }

        Js.log2("conf rego problem", ex)
      }
    }
  )

  let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    switch Js.Nullable.isNullable(cognitoUser) {
    | false =>
      cognitoUser->confirmRegistration(
        verifCode,
        false,
        confirmregistrationCallback,
        Js.Nullable.null,
      )
    | true => (_ => Some("null user - not submitting"))->setCognitoErr
    }
  }

  <main>
    <form className="w-5/6 m-auto" onSubmit={handleSubmit}>
      <fieldset className="h-40">
        <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
          {"Confirm"->React.string}
        </legend>
        {switch Js.Nullable.isNullable(cognitoUser) {
        | true =>
          <div className="relative">
            <label
              className={switch (unVisited, unErr) {
              | (true, Some(_)) => "text-2xl text-red-500 font-bold font-flow"
              | (false, _) | (true, None) => "text-2xl text-warm-gray-100 font-flow"
              }}
              htmlFor="username">
              {"username:"->React.string}
            </label>
            {switch (unVisited, unErr) {
            | (true, Some(err)) =>
              <span className="absolute right-0 text-2xl text-red-500 font-bold font-flow">
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
                ) => "h-6 w-full text-xl pl-1 text-left outline-none text-red-500 bg-transparent border-b-1 border-red-500"
              | (false, _)
              | (
                true,
                None,
              ) => "h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
              }}
              id="username"
              minLength=4
              name="username"
              onBlur={onBlur("username")}
              onChange={onChange(setUsername)}
              // placeholder="Enter username"
              required=true
              spellCheck=false
              type_="text"
              value={username}
            />
          </div>
        | false => React.null
        }}
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
            onChange
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
}
