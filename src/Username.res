@react.component
let make = (~username, ~setUsername, ~setUsernameError) => {

  UsernameValidation.useUsernameValidation(username, setUsernameError)


  let onChange = e => setUsername(_ => ReactEvent.Form.target(e)["value"])

  <div>
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
      {React.string("username:")}
    </label>
    <input
      autoComplete="username"
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="username"
      name="username"
      onChange
      spellCheck=false
      type_="text"
      value={username}
    />
  </div>
}
