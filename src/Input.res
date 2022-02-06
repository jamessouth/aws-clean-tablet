@react.component
let make = (~value, ~setFunc, ~setErrorFunc, ~funcList, ~propName, ~autoComplete=propName, ~toggleProp=false, ~toggleSetFunc) => {

  Validator.useValidator(value, setErrorFunc, funcList, propName ++ ": ")

  let onChange = e => setFunc(_ => ReactEvent.Form.target(e)["value"])

  <div className="w-full">
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor=autoComplete>
      {React.string(propName)}
    </label>
    <input
      autoComplete
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id=autoComplete
      name=propName
      onChange
      spellCheck=false
      type_={switch propName == "username" || toggleProp {
      | true => "text"
      | false => propName
      }}
      value
    />
    {switch propName == "password" {
    | true => <Toggle toggleProp toggleSetFunc/>
    | false => React.null
    }}
  </div>
}
