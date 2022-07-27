let className = "font-arch bg-transparent text-stone-100 text-2xl absolute right-0 top-0 cursor-pointer"

@react.component
let make = (
  ~value,
  ~propName,
  ~autoComplete=propName,
  ~inputMode="text",
  ~onKeyPress=_e => (),
  ~setFunc,
) => {
  let (showPassword, setShowPassword) = React.Uncurried.useState(_ => false)

  let onChange = e => setFunc(._ => ReactEvent.Form.target(e)["value"])

  <div
    className={switch propName == "password" {
    | true => "max-w-xs lg:max-w-sm w-full relative"
    | false => "max-w-xs lg:max-w-sm w-full"
    }}>
    <label className="text-2xl font-flow text-stone-100" htmlFor=autoComplete>
      {React.string(propName)}
    </label>
    <input
      autoComplete
      className="h-6 w-full text-xl outline-none font-anon bg-transparent border-b-1 text-stone-100 border-stone-100 indent-2px"
      id=autoComplete
      inputMode
      name=propName
      onChange
      onKeyPress
      spellCheck=false
      type_={switch propName == "username" || propName == "answer" || showPassword {
      | true => "text"
      | false => propName
      }}
      value
    />
    {switch propName == "password" {
    | true =>
      <Button
        textTrue="hide"
        textFalse="show"
        textProp=showPassword
        onClick={_ => setShowPassword(.prev => !prev)}
        className
      />
    | false => React.null
    }}
  </div>
}
