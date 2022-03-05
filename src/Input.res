@react.component
let make = (~value, ~propName, ~autoComplete=propName, ~inputMode="text", ~setFunc) => {
  let (showPassword, setShowPassword) = React.useState(_ => false)

  let onChange = e => setFunc(_ => ReactEvent.Form.target(e)["value"])

  <div
    className={switch propName == "password" {
    | true => "w-full relative"
    | false => "w-full"
    }}>
    <label className="text-2xl font-flow text-warm-gray-100" htmlFor=autoComplete>
      {React.string(propName)}
    </label>
    <input
      autoComplete
      className="h-6 w-full text-xl font-anon bg-transparent border-b-1 text-warm-gray-100 border-warm-gray-100"
      id=autoComplete
      inputMode
      name=propName
      onChange
      spellCheck=false
      type_={switch propName == "username" || showPassword {
      | true => "text"
      | false => propName
      }}
      value
    />
    {switch propName == "password" {
    | true =>
      React.useMemo1(
        _ => {<Toggle toggleProp=showPassword toggleSetFunc=setShowPassword />},
        [showPassword],
      )
    | false => React.null
    }}
  </div>
}
