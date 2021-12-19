type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

let circ = Js.String2.fromCharCode(9862)
let answer_max_length = 12

// ~playerName,
@react.component
let make = (
  ~wsConnected,
  ~game: Reducer.liveGame,
  ~playerColor,
  ~send,
  ~wsError,
  ~currentWord,
  ~previousWord,
) => {
  Js.log3("play", wsConnected, wsError)
  Js.log3("play2", previousWord, game)
  let (answered, setAnswered) = React.useState(_ => false)
  let (inputText, setInputText) = React.useState(_ => "")

  let sendAnswer = _ => {
    let pl = {
      action: "answer",
      gameno: game.no,
      answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
    }
    send(. Js.Json.stringifyAny(pl))
    (_ => true)->setAnswered
    (_ => "")->setInputText
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
    <Scoreboard players=game.players currentWord previousWord />
    <p className="text-yellow-200 text-2xl font-bold">
      {"Get Ready"->React.string} <span className="animate-spin"> {React.string(circ)} </span>
    </p>
    <Word onAnimationEnd playerColor currentWord answered />
    <Form answer_max_length answered inputText onEnter setInputText />

    // <Prompt></Prompt>
  </div>
}
