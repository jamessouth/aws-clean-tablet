

@react.component
let make = (~email, ~setEmail, ~setEmailError) => {


  EmailValidation.useEmailValidation(email, setEmailError)

  let onChange = e => setEmail(_ => ReactEvent.Form.target(e)["value"])

  <div className="w-full">
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="email">
      {React.string("email:")}
    </label>
    <input
      autoComplete="email"
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="email"
      name="email"
      onChange
      spellCheck=false
      type_="email"
      value={email}
    />
  </div>
}
