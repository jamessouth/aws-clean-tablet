@react.component
let make = (~value, ~setFunc, ~setErrorFunc, ~funcList, ~propName) => {



  Validator.useValidator(value, setErrorFunc, funcList, propName ++ ": ")


  let onChange = e => setFunc(_ => ReactEvent.Form.target(e)["value"])

  <div className="w-full">
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor=propName>
      {React.string(propName)}
    </label>
    <input
      autoComplete=propName
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id=propName
      name=propName
      onChange
      spellCheck=false
      type_={switch propName == "username" {
      | true => "text"
      | false => propName
      }}
      value
    />
  </div>
}
