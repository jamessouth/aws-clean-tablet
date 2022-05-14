type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
  index: string,
}

type scorePayload = {
  action: string,
  game: Reducer.liveGame,
}

@react.component
let make = (~game: Reducer.liveGame, ~playerColor, ~playerIndex, ~send, ~leader, ~playerName) => {
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  let (answered, setAnswered) = React.Uncurried.useState(_ => false)
  let (answer, setAnswer) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "ANSWER: 2-12 length; letters and spaces only; ",
  ))
  let {players, currentWord, previousWord, showAnswers, sk, winner} = game
  let answer_max_length = 12

  React.useEffect1(() => {
    ErrorHook.useError(answer, "ANSWER", setValidationError)
    None
  }, [answer])

  React.useEffect2(() => {
    Js.log("send start useeff")
    switch (leader, playerColor == "transparent") {
    | (_, true) | (false, _) => ()
    | (true, false) => {
        Js.log3("send start", leader, playerColor)
        let pl: Game.startPayload = {
          action: "start",
          gameno: sk,
        }
        send(. Js.Json.stringifyAny(pl))
      }
    }
    None
  }, (leader, playerColor))

  let hasRendered = React.useRef(false)
  Js.log3("play", game, hasRendered)

  React.useEffect3(() => {
    switch (leader, showAnswers, winner == "") {
    | (true, true, _) => Js.Global.setTimeout(() => {
        let pl: scorePayload = {
          action: "score",
          game: game,
        }
        send(. Js.Json.stringifyAny(pl))
      }, 8564)->ignore

    | (true, false, true) => {
        setAnswered(._ => false)
        switch hasRendered.current {
        | true => Js.Global.setTimeout(() => {
            let pl: Game.startPayload = {
              action: "start",
              gameno: sk,
            }
            send(. Js.Json.stringifyAny(pl))
          }, 2564)->ignore

        | false => hasRendered.current = true
        }
      }

    | (false, false, true) => setAnswered(._ => false)
    | (false, true, _) | (_, false, false) => ()
    }
    None
  }, (leader, showAnswers, winner))

  let sendAnswer = _ => {
    let pl = {
      action: "answer",
      gameno: sk,
      answer: answer
      ->Js.String2.slice(~from=0, ~to_=answer_max_length)
      ->Js.String2.replaceByRe(%re("/\d/g"), "")
      ->Js.String2.replaceByRe(%re("/[!-/:-@\[-`{-~]/g"), "")
      ->Js.String2.trim(_),
      index: playerIndex,
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

  let onClick = _ => {
    let pl: Game.startPayload = {
      action: "end",
      gameno: sk,
    }
    send(. Js.Json.stringifyAny(pl))

    RescriptReactRouter.push("/lobby")
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

  <div>
    <Scoreboard players currentWord previousWord showAnswers winner onClick playerName />
    {switch winner == "" {
    | false => React.null
    | true =>
      switch currentWord == "game over" {
      | true => <Word onAnimationEnd playerColor currentWord answered showTimer=false />
      | false => <>
          <Word onAnimationEnd playerColor currentWord answered showTimer={currentWord != ""} />
          {switch currentWord == "" {
          | true => React.null
          | false =>
            switch answered {
            | true => React.null
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

    // <Prompt></Prompt>
  </div>
}
