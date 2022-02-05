@react.component
let make = (~password, ~setPassword, ~setPasswordError, ~funcList) => {
  let (showPassword, setShowPassword) = React.useState(_ => false)


  Validator.useValidator(password, setPasswordError, funcList, "Password: ")

  let onClick = _ => {
    setShowPassword(prev => !prev)
  }

  let onChange = e => setPassword(_ => ReactEvent.Form.target(e)["value"])

  <div>
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="new-password">
      {React.string("password:")}
    </label>
    <input
      autoComplete="new-password"
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="new-password"
      name="password"
      onChange
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
