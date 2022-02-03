let reqs = " 3-10 characters; letters, numbers, and underscores only."

@react.component
let make = (~email, ~setEmail) => {
  let onChange = e => setEmail(_ => ReactEvent.Form.target(e)["value"])

  <div className="w-full">
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="email">
      {React.string("email:")}
    </label>
    <input
      autoComplete="email"
      autoFocus=false
      className="h-6 w-full text-base font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="email"
    //   maxLength=50
    //   minLength=7
      name="email"
      onChange
      // placeholder="Enter username"
    //   required=true
      spellCheck=false
      type_="email"
      value={email}
    />
  </div>
}
