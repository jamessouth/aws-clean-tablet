type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

let circ = Js.String2.fromCharCode(8635)
let answer_max_length = 12

// ~playerName,
@react.component
let make = (~wsConnected, ~game: Reducer.liveGame, ~playerColor, ~send, ~wsError) => {
  Js.log4("play", wsConnected, wsError, game)

  let (answered, setAnswered) = React.useState(_ => false)
  let (inputText, setInputText) = React.useState(_ => "")

  let (answersPhase, setAnswersPhase) = React.useState(_ => false)
  let (start, setStart) = React.useState(_ => false)

  let {players, currentWord, previousWord} = game

  React.useEffect2(() => {
    Js.log("wordz")
    switch (currentWord == "", previousWord == "") {
    | (_, true) | (false, false) => setAnswersPhase(_ => false)
    | (true, false) => setAnswersPhase(_ => true)
    }
    None
  }, (currentWord, previousWord))

  React.useEffect1(() => {
    Js.log("answersphase")
    switch answersPhase {
    | true => ()
    | false => setAnswered(_ => false)
    }
    None
  }, [answersPhase])

  React.useEffect1(() => {
    switch playerColor == "" {
    | true => Js.Global.setTimeout(() => {
        setStart(_ => false)
      }, 2)->ignore
    | false => Js.Global.setTimeout(() => {
        setStart(_ => true)
      }, 2000)->ignore
    }
    None
  }, [playerColor])

  let sendAnswer = _ => {
    let pl = {
      action: "answer",
      gameno: game.no,
      answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
    }
    send(. Js.Json.stringifyAny(pl))
    setAnswered(_ => true)
    setInputText(_ => "")
  }

  let onAnimationEnd = _ => {
    if !answered {
      sendAnswer()
    }
    Js.log("onanimend")
  }

  let onEnter = _ => {
    if !answered {
      sendAnswer()
    }
    Js.log("onenter")
  }

  <div>
    // playerName
    <Scoreboard players previousWord showAnswers=answersPhase />
    {switch start {
    | true => React.null
    | false =>
      <span className="animate-spin text-yellow-200 text-2xl font-bold absolute left-1/2">
        {React.string(circ)}
      </span>
    }}
    <Word onAnimationEnd playerColor currentWord answered showTimer={start && !answersPhase} />
    <Form answer_max_length answered inputText onEnter setInputText />

    // <Prompt></Prompt>
  </div>
}
