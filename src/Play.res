type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
  tipe: string,
  playersCount: string,
}

let circ = Js.String2.fromCharCode(9862)
let answer_max_length = 12

// ~playerName,
@react.component
let make = (
  ~wsConnected,
  ~game: Reducer.liveGame,
  ~leadertoken,
  ~playerColor,
  ~send,
  ~wsError,
  ~currentWord,
  ~previousWord,
) => {
  Js.log4("play", wsConnected, leadertoken, wsError)
  Js.log3("play2", previousWord, game)
  let (answered, setAnswered) = React.useState(_ => false)
  let (inputText, setInputText) = React.useState(_ => "")
  let (gameStarted, setGameStarted) = React.useState(_ => false)

  let noplrs = Js.Array2.length(game.players)

  let _sendAnswer = _ => {
    let pl = {
      action: "answer",
      gameno: game.no,
      answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
      playersCount: j`$noplrs`,
    }
    send(. Js.Json.stringifyAny(pl))
    (_ => true)->setAnswered
    (_ => "")->setInputText
  }

  let onAnimationEnd = _ => {
    // sendAnswer()
    Js.log("onanimend")
  }

  let onEnter = _ => {
    // sendAnswer()
    Js.log("onenter")
  }

  React.useEffect1(() => {
    (_ => false)->setAnswered
    None
  }, [currentWord])

  React.useEffect1(() => {
    (_ => true)->setGameStarted
    None
  }, [currentWord])

  <div>
    // playerName
    <Scoreboard players=game.players />
    {switch gameStarted {
    | false =>
      <p className="text-yellow-200 text-2xl font-bold">
        {"Get Ready"->React.string} <span className="animate-spin"> {circ->React.string} </span>
      </p>
    | true => React.null
    }}
    <Word onAnimationEnd playerColor currentWord />
    <Form answer_max_length answered inputText onEnter setInputText />

    // <Prompt></Prompt>
  </div>
}
