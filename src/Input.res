@react.component
let make = (
  ~value,
  ~setFunc,
  ~setErrorFunc,
  ~funcList,
  ~propName,
  ~autoComplete=propName,
  ~toggleProp=false,
  ~toggleButton=React.null,
  ~validationError,
) => {
  let (class, setClass) = React.useState(_ => "warm-gray-100")

  React.useEffect1(() => {
    switch validationError {
    | None => setClass(_ => "warm-gray-100")
    | Some(err) =>
      switch Js.String2.startsWith(err, propName) {
      | false => setClass(_ => "warm-gray-100")
      | true => setClass(_ => "red-500")
      }
    }
    None
  }, [validationError])

  Validator.useValidator(value, setErrorFunc, funcList, propName ++ ": ")

  let onChange = e => setFunc(_ => ReactEvent.Form.target(e)["value"])

  <div className="w-full">
    <label className={`text-2xl font-flow text-${class}`} htmlFor=autoComplete>
      {React.string(propName)}
    </label>
    <input
      autoComplete
      className={`h-6 w-full text-xl font-anon bg-transparent border-b-1 text-warm-gray-100 border-${class}`}
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
    {toggleButton}
  </div>
}
