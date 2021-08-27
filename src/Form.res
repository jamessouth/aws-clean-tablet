@send external blur: Dom.element => unit = "blur"

@react.component
let make = (ANSWER_MAX_LENGTH, answered, inputText, onEnter, setInputText) => {
  let inputBox = React.useRef(Js.Nullable.null)
  let INPUT_MIN_LENGTH = 2

  let (disableSubmit, setDisableSubmit) = React.useState(_ => true)
  let (isValidInput, setIsValidInput) = React.useState(_ => true)
  let (badChar, setBadChar) = React.useState(_ => Js.Nullable.null)

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
    switch answered {
    | true =>
      switch inputBox.current->Js.Nullable.toOption {
      | Some(inp) => inp->blur
      | None => ()
      }
    | false => ()
    }
  }, [answered])

  <div> {React.string("Hello ReScripters!")} </div>
}
