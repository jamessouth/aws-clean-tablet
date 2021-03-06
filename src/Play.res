type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

type endPayload = {
  action: string,
  gameno: string,
  token: string,
}

// type scorePayload = {
//   action: string,
//   game: Reducer.liveGame,
// }

@react.component
let make = (
  ~players: array<Reducer.livePlayer>,
  ~sk,
  ~showAnswers,
  ~winner,
  ~oldWord,
  ~word,
  ~playerColor,
  ~send,
  ~playerName,
  ~endtoken,
  ~resetConnState,
) => {
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  let (answered, setAnswered) = React.Uncurried.useState(_ => false)
  let (answer, setAnswer) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "ANSWER: 2-12 length; letters and spaces only; ",
  ))

  let answer_max_length = 12

  React.useEffect1(() => {
    ErrorHook.useError(answer, "ANSWER", setValidationError)
    None
  }, [answer])

  React.useEffect1(() => {
    setAnswered(._ => false)
    None
  }, [showAnswers])

  let hasRendered = React.useRef(false)
  Js.log3("play", players, hasRendered)

  let sendAnswer = _ => {
    let pl = {
      action: "answer",
      gameno: sk,
      answer: answer
      ->Js.String2.slice(~from=0, ~to_=answer_max_length)
      ->Js.String2.replaceByRe(%re("/\d/g"), "")
      ->Js.String2.replaceByRe(%re("/[!-/:-@\[-`{-~]/g"), "")
      ->Js.String2.trim(_),
    }
    send(. Js.Json.stringifyAny(pl))
    setAnswered(._ => true)
    setAnswer(._ => "")
    setSubmitClicked(._ => false)
  }

  let onAnimationEnd = _ => {
    if !answered {
      sendAnswer()
    }
    Js.log("onanimend")
  }

  let onEnter = (. _) => {
    if !answered {
      sendAnswer()
    }
    Js.log("onenter")
  }

  let reset = _ => {

    resetConnState()
    RescriptReactRouter.push("/lobby")
  }

  let onClick = (n, _) => {
   Js.log(n)
        let pl: endPayload = {
          action: "end",
          gameno: sk,
          token: endtoken,
        }
        send(. Js.Json.stringifyAny(pl))
   
    reset()

  }

  let onClick2 = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None => onEnter(. ignore())
    | Some(_) => ()
    }
  }

  let onKeyPress = e => {
    let key = ReactEvent.Keyboard.key(e)
    switch key {
    | "Enter" => onClick2()
    | _ => ()
    }
  }
  // open Web
  //   // addDocumentEventListener("popstate", e => Js.log2("vis", e))

  //   addWindowEventListener("popstate", e => {

  //     // preventDefault(e)
  //     // let c = confirm("you can't leave now!!!")
  //     Js.log2("bfu", e)

  // })

  <div>
    <Scoreboard players oldWord word showAnswers winner onClick reset playerName />
    {switch winner == "" {
    | false => React.null
    | true =>
      switch word == "game over" {
      | true => <Word onAnimationEnd playerColor word answered showTimer=false />
      | false => <>
          <Word onAnimationEnd playerColor word answered showTimer={word != ""} />
          {switch word == "" {
          | true => <div className="bg-transparent h-45 w-full" />
          | false =>
            switch answered {
            | true => <div className="bg-transparent h-45 w-full" />
            | false =>
              <Form
                ht="h-24" onClick=onClick2 leg="" submitClicked validationError cognitoError=None>
                <Input
                  value=answer
                  propName="answer"
                  autoComplete="off"
                  inputMode="text"
                  onKeyPress
                  setFunc=setAnswer
                />
              </Form>
            }
          }}
        </>
      }
    }}
  </div>
}
