type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

let circ = Js.String2.fromCharCode(9862)
let answer_max_length = 12

// ~playerName,
@react.component
let make = (~wsConnected, ~game: Reducer.liveGame, ~playerColor, ~send, ~wsError) => {
  Js.log4("play", wsConnected, wsError, game)

  let (answered, setAnswered) = React.useState(_ => false)
  let (inputText, setInputText) = React.useState(_ => "")

  let (answersPhase, setAnswersPhase) = React.useState(_ => false)

  let {players, currentWord, previousWord} = game

  React.useEffect2(() => {
    switch (currentWord != "", previousWord != "") {
    | (_, false) | (true, true) => setAnswersPhase(_ => false)
    | (false, true) => setAnswersPhase(_ => true)
    }
    None
  }, (currentWord, previousWord))

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

  React.useEffect1(() => {
    switch answersPhase {
    | true => ()
    | false => setAnswered(_ => false)
    }
    None
  }, [answersPhase])

  <div>
    // playerName
    <Scoreboard players previousWord showAnswers=answersPhase />
    // <p className="text-yellow-200 text-2xl font-bold">
    //   {"Get Ready"->React.string} <span className="animate-spin"> {React.string(circ)} </span>
    // </p>
    <Word onAnimationEnd playerColor currentWord answered showTimer={!answersPhase} />
    <Form answer_max_length answered inputText onEnter setInputText />

    // <Prompt></Prompt>
  </div>
}
