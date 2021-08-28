@send external blur: Dom.element => unit = "blur"

@react.component
let make = (ANSWER_MAX_LENGTH, answered, inputText, onEnter, setInputText) => {
  let inputBox = React.useRef(Js.Nullable.null)
  let INPUT_MIN_LENGTH = 2


  let (disableSubmit, setDisableSubmit) = React.useState(_ => true)
  let (isValidInput, setIsValidInput) = React.useState(_ => true)
  let (badChar, setBadChar) = React.useState(_ => Js.Nullable.null)

  let onKeyPress = evt => {
    let key = ReactEvent.Keyboard.key(evt)
    switch key, disableSubmit {
    | 'Enter', false => onEnter()
    | _, _ => ()
    }
  }

  let onChange = evt => {
    let value = ReactEvent.Form.currentTarget(evt)["value"]
    value->setInputText
  }

  let onClick = _ => {
    onEnter()
  }

  React.useEffect1(() => {
    switch Js.String.match_(%re("/[^a-z '-]+/i"), inputText) {
    | Some(arr) => {
        arr[0]->setBadChar
        false->setIsValidInput
      }
    | None => {
        Js.Nullable.null->setBadChar
        true->setIsValidInput
      }
    }
  }, [inputText])

  React.useEffect1(() => {
    switch answered, inputBox.current->Js.Nullable.toOption {
    | true, Some(inp) => inp->blur
    | true, None | false, _ => ()
    }
  }, [answered])

  React.useEffect4(() => {
    ((inputText->Js.String2.length < INPUT_MIN_LENGTH || inputText->Js.String2.length > ANSWER_MAX_LENGTH) || answered || !isValidInput)->setDisableSubmit
  }, [inputText, ANSWER_MAX_LENGTH, answered, isValidInput]])



  <section className="relative flex flex-col justify-between items-center h-40 text-xl mb-12">
    {
      !isValidInput &&
      <p className="absolute text-smoke-100 bg-smoke-800 font-bold w-11/12 max-w-xl" ariaLive="assertive">
      {switch badChar {
      | Some(bc) => bc ++ " is not allowed"
      | None => "That input is not allowed"
      }}
      </p>
    }
    <label ariaLive="assertive" for="inputbox">
      "Enter your answer:"->React.string
    </label>

    <input
      className="h-7 w-3/5 text-xl pl-1 text-left bg-transparent border-none text-smoke-700"
      id="inputbox"
      autoComplete="off"
      autoFocus
      ref={inputBox->ReactDOM.Ref.domRef}
      value={inputText}
      spellCheck="false"
      onKeyPress
      onChange
      type_="text"
      placeholder={`2 - ${ANSWER_MAX_LENGTH} letters`}
      readOnly={switch answered {
      | true => true
      | false => false
      }}
    ></input>
  
    <button
      className="text-2xl text-smoke-700 bg-smoke-100 h-7 w-2/3 max-w-max cursor-pointer border-none disabled:cursor-not-allowed disabled:contrast-50"
      type_="button"
      onClick
      disabled={switch disableSubmit {
      | true => true
      | false => false
      }}
    >
      "Submit"->React.string
    </button>
  
  </section>
}
