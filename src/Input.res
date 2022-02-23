@react.component
let make = (
  ~value,
  ~propName,
  ~autoComplete=propName,
  ~toggleProp=false,
  ~toggleButton=React.null,
  ~inputMode="text",
  ~setFunc,
) => {
  let onChange = e => setFunc(_ => ReactEvent.Form.target(e)["value"])

  <div
    className={switch Js.Nullable.isNullable(Js.Nullable.return(toggleButton)) {
    | true => "w-full"
    | false => "w-full relative"
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
      type_={switch propName == "username" || toggleProp {
      | true => "text"
      | false => propName
      }}
      value
    />
    {toggleButton}
  </div>
}
