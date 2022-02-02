let reqs = " 8-98 characters; at least 1 each of uppercase and lowercase letters, numbers, and symbols required."

@react.component
let make = (
  ~pwVisited,
  ~setPwVisited,
  ~pwErr,
  ~setPwErr,
  ~setDisabled,
  ~password,
  ~setPassword,
) => {
  let (showPassword, setShowPassword) = React.useState(_ => false)

  let checkNoPwWhitespace = pw => {
    let r = %re("/\s/")
    switch Js.String2.match_(pw, r) {
    | Some(_) => setPwErr(_ => Some("No whitespace allowed." ++ reqs))
    | None => setPwErr(_ => None)
    }
  }

  let checkPwMaxLength = pw => {
    switch pw->Js.String2.length > 98 {
    | true => setPwErr(_ => Some("Password is too long." ++ reqs))
    | false => checkNoPwWhitespace(pw)
    }
  }

  let checkSymbol = pw => {
    let r = %re("/[!-*\[-`{-~./,:;<>?@]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a symbol." ++ reqs))
    | Some(_) => checkPwMaxLength(pw)
    }
  }

  let checkNumber = pw => {
    let r = %re("/\d/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a number." ++ reqs))
    | Some(_) => checkSymbol(pw)
    }
  }

  let checkUpper = pw => {
    let r = %re("/[A-Z]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add an uppercase letter." ++ reqs))
    | Some(_) => checkNumber(pw)
    }
  }

  let checkLower = pw => {
    let r = %re("/[a-z]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a lowercase letter." ++ reqs))
    | Some(_) => checkUpper(pw)
    }
  }

  let checkPwLength = pw => {
    switch pw->Js.String2.length < 8 {
    | true => setPwErr(_ => Some("Password is too short." ++ reqs))
    | false => checkLower(pw)
    }
  }

  let onClick = _ => {
    setShowPassword(prev => !prev)
  }

  let onBlur = _ => setPwVisited(_ => true)

  let onChange = e => setPassword(_ => ReactEvent.Form.target(e)["value"])

  React.useEffect2(() => {
    switch pwVisited {
    | true => checkPwLength(password)
    | false => setPwErr(_ => None)
    }
    None
  }, (password, pwVisited))

  React.useEffect1(() => {
    switch pwErr {
    | None => setDisabled(_ => false)
    | Some(_) => setDisabled(_ => true)
    }
    None
  }, [pwErr])

  <div>
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="new-password">
      {React.string("password:")}
    </label>
    <input
      autoComplete="new-password"
      autoFocus=false
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="new-password"
      maxLength=98
      minLength=8
      name="password"
      onBlur
      onChange
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
      className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 top-36 cursor-pointer"
      onClick>
      {switch showPassword {
      | true => "hide"->React.string
      | false => "show"->React.string
      }}
    </button>
  </div>
}
