let sectionClass = "relative flex flex-col justify-between items-center h-40 text-xl mb-12 "

// let className = "text-gray-700 mt-14 bg-stone-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

@react.component
let make = (
  ~answered,
  ~answer,
  ~onEnter,
  ~setAnswer,
  ~currentWord,
  ~submitClicked,
  ~setSubmitClicked,
  ~validationError,
) => {
  // let inputBox = React.useRef(Js.Nullable.null)

  // let (disableSubmit, setDisableSubmit) = React.Uncurried.useState(_ => true)
  // let (isValidInput, setIsValidInput) = React.Uncurried.useState(_ => true)
  // let (badChar: option<string>, setBadChar) = React.Uncurried.useState(_ => None)

  // let onKeyPress = evt => {
  //   setSubmitClicked(._ => true)
  //   let key = ReactEvent.Keyboard.key(evt)
  //   switch (key, validationError) {
  //   | ("Enter", None) => onEnter(. ignore())
  //   | (_, _) => ()
  //   }
  // }

  let onClick = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None => onEnter(. ignore())
    | Some(_) => ()
    }
  }

  let onKeyPress = e => {
    let key = ReactEvent.Keyboard.key(e)
    switch key {
    | "Enter" => onClick()
    | _ => ()
    }
  }

  let onChange = e => setAnswer(._ => ReactEvent.Form.target(e)["value"])

  // React.useEffect1(() => {
  //   switch %re("/[^a-z '-]+/i")->Js.Re.test_(inputText) {
  //   | true => {
  //       setBadChar(._ => Some(inputText))
  //       setIsValidInput(._ => false)
  //     }
  //   | false => {
  //       setBadChar(._ => None)
  //       setIsValidInput(._ => true)
  //     }
  //   }
  //   None
  // }, [inputText])

  // React.useEffect2(() => {
  //   switch (answered, Js.Nullable.toOption(inputBox.current)) {
  //   | (true, Some(inp)) => Web.blur(inp)
  //   | (true, None) | (false, _) => ()
  //   }
  //   None
  // }, (answered, inputBox.current))

  // React.useEffect3(() => {
  //   setDisableSubmit(._ =>
  //     Js.String2.length(Js.String2.trim(inputText)) < min_answer_length ||
  //     Js.String2.length(Js.String2.trim(inputText)) > max_answer_length ||
  //     answered ||
  //     !isValidInput
  //   )
  //   None
  // }, (inputText, answered, isValidInput))

  <section
    className={switch currentWord == "" {
    | true => sectionClass ++ "invisible"
    | false => sectionClass
    }}>
    {
      // {switch isValidInput {
      // | true => React.null
      // | false =>
      //   <p className="absolute text-stone-100 bg-red-600 font-bold w-11/12 max-w-xl">
      //     {switch badChar {
      //     | Some(bc) => React.string(bc ++ " is not allowed")
      //     | None => React.string("That input is not allowed")
      //     }}
      //   </p>
      // }}

      switch submitClicked {
      | false => React.null
      | true =>
        switch validationError {
        | Some(error) =>
          <span
            className="absolute text-sm text-stone-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
            {React.string(error)}
          </span>
        | None => React.null
        }
      }
    }
    {switch answered {
    | true => React.null
    | false => <>
        <label className="text-stone-100 font-anon" htmlFor="inputbox">
          {React.string("Enter your answer:")}
        </label>
        <input
          className="h-7 w-3/5 text-xl pl-1 text-left bg-stone-100 border-none text-stone-800 max-w-xs"
          id="inputbox"
          autoComplete="off"
          // ref={ReactDOM.Ref.domRef(inputBox)}
          value={answer}
          spellCheck=false
          onKeyPress
          onChange
          type_="text"
        />
        // disabled=disableSubmit className
        <Button onClick  />
      </>
    }}
  </section>
}
