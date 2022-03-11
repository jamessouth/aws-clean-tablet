let sectionClass = "relative flex flex-col justify-between items-center h-40 text-xl mb-12 "

let className = "font-anon text-xl text-true-gray-700 bg-true-gray-100 h-7 w-2/3 max-w-250px cursor-pointer border-none disabled:cursor-not-allowed disabled:contrast-[0.25]"

@react.component
let make = (~answer_max_length, ~answered, ~inputText, ~onEnter, ~setInputText, ~currentWord) => {
  let inputBox = React.useRef(Js.Nullable.null)
  let input_min_length = 2

  let (disableSubmit, setDisableSubmit) = React.Uncurried.useState(_ => true)
  let (isValidInput, setIsValidInput) = React.Uncurried.useState(_ => true)
  let (badChar: option<string>, setBadChar) = React.Uncurried.useState(_ => None)

  let onKeyPress = evt => {
    let key = ReactEvent.Keyboard.key(evt)
    switch (key, disableSubmit) {
    | ("Enter", false) => onEnter(. ignore())
    | (_, _) => ()
    }
  }

  let onChange = evt => {
    let value = ReactEvent.Form.currentTarget(evt)["value"]
    setInputText(._ => value)
  }

  React.useEffect1(() => {
    switch %re("/[^a-z '-]+/i")->Js.Re.test_(inputText) {
    | true => {
        setBadChar(._ => Some(inputText))
        setIsValidInput(._ => false)
      }
    | false => {
        setBadChar(._ => None)
        setIsValidInput(._ => true)
      }
    }
    None
  }, [inputText])

  React.useEffect2(() => {
    switch (answered, Js.Nullable.toOption(inputBox.current)) {
    | (true, Some(inp)) => Web.blur(inp)
    | (true, None) | (false, _) => ()
    }
    None
  }, (answered, inputBox.current))

  React.useEffect4(() => {
    setDisableSubmit(._ =>
      Js.String2.length(Js.String2.trim(inputText)) < input_min_length ||
      Js.String2.length(Js.String2.trim(inputText)) > answer_max_length ||
      answered ||
      !isValidInput
    )
    None
  }, (inputText, answer_max_length, answered, isValidInput))

  <section
    className={switch currentWord == "" {
    | true => sectionClass ++ "invisible"
    | false => sectionClass
    }}>
    {switch isValidInput {
    | true => React.null
    | false =>
      <p className="absolute text-smoke-100 bg-smoke-800 font-bold w-11/12 max-w-xl">
        {switch badChar {
        | Some(bc) => React.string(bc ++ " is not allowed")
        | None => React.string("That input is not allowed")
        }}
      </p>
    }}
    <label htmlFor="inputbox"> {React.string("Enter your answer:")} </label>
    <input
      className="h-7 w-3/5 text-xl pl-1 text-left bg-warm-gray-100 border-none text-smoke-700"
      id="inputbox"
      autoComplete="off"
      // autoFocus
      ref={ReactDOM.Ref.domRef(inputBox)}
      value={inputText}
      spellCheck=false
      onKeyPress
      onChange
      type_="text"
      placeholder={`2 - ${answer_max_length->Js.Int.toString} letters`}
      readOnly={switch answered {
      | true => true
      | false => false
      }}
    />
    <Button
      onClick={_ => onEnter(. ignore())}
      disabled=disableSubmit
      className
    />
  </section>
}
