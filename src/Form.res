@send external blur: Dom.element => unit = "blur"

@react.component
let make = (~answer_max_length, ~answered, ~inputText, ~onEnter, ~setInputText) => {
  let inputBox = React.useRef(Js.Nullable.null)
  let input_min_length = 2

  let (disableSubmit, setDisableSubmit) = React.useState(_ => true)
  let (isValidInput, setIsValidInput) = React.useState(_ => true)
  let (badChar: option<string>, setBadChar) = React.useState(_ => None)

  let onKeyPress = evt => {
    let key = ReactEvent.Keyboard.key(evt)
    switch (key, disableSubmit) {
    | ("Enter", false) => onEnter()
    | (_, _) => ()
    }
  }

  let onChange = evt => {
    let value = ReactEvent.Form.currentTarget(evt)["value"]
    setInputText(value)
  }

  let onClick = _ => {
    onEnter()
  }

  React.useEffect1(() => {
    switch inputText->Js.String2.match_(%re("/[^a-z '-]+/i")) {
    | Some(arr) => {
        setBadChar(_ => Some(arr[0]))
        setIsValidInput(_ => false)
      }
    | None => {
        setBadChar(_ => None)
        setIsValidInput(_ => true)
      }
    }
    None
  }, [inputText])

  React.useEffect2(() => {
    switch (answered, Js.Nullable.toOption(inputBox.current)) {
    | (true, Some(inp)) => blur(inp)
    | (true, None) | (false, _) => ()
    }
    None
  }, (answered, inputBox.current))

  React.useEffect4(() => {
    setDisableSubmit(_ =>
      Js.String2.length(Js.String2.trim(inputText)) < input_min_length ||
      Js.String2.length(Js.String2.trim(inputText)) > answer_max_length ||
      answered ||
      !isValidInput
    )
    None
  }, (inputText, answer_max_length, answered, isValidInput))

  <section className="relative flex flex-col justify-between items-center h-40 text-xl mb-12">
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
    <button
      className="text-2xl text-smoke-700 bg-smoke-100 h-7 w-2/3 max-w-max cursor-pointer border-none disabled:cursor-not-allowed disabled:contrast-50"
      type_="button"
      onClick
      disabled={switch disableSubmit {
      | true => true
      | false => false
      }}>
      {"Submit"->React.string}
    </button>
  </section>
}
