type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

type scorePayload = {
  action: string,
  game: Reducer.liveGame,
}

let circ = Js.String2.fromCharCode(8635)
let answer_max_length = 12

// ~playerName,
@react.component
let make = (~wsConnected, ~game: Reducer.liveGame, ~playerColor, ~send, ~wsError, ~leadertoken) => {
  Js.log4("play", wsConnected, wsError, game)

  let (answered, setAnswered) = React.useState(_ => false)
  let (inputText, setInputText) = React.useState(_ => "")

  let (answersPhase, setAnswersPhase) = React.useState(_ => false)
  // let (start, setStart) = React.useState(_ => false)
  let (leader, setLeader) = React.useState(_ => false)

  let {players, currentWord, previousWord} = game

  React.useEffect2(() => {
    switch Js.Array2.length(game.players) > 0 {
    | true =>
      switch game.players[0].name ++ game.players[0].connid == leadertoken {
      | true => setLeader(_ => true)
      | false => setLeader(_ => false)
      }
    | false => setLeader(_ => false)
    }
    None
  }, (game.players, leadertoken))

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

  React.useEffect2(() => {
    switch (leader, answersPhase) {
    | (true, true) => Js.Global.setTimeout(() => {
        let pl: scorePayload = {
          action: "score",
          game: game,
        }
        send(. Js.Json.stringifyAny(pl))
      }, 8564)->ignore
    | _ => ()
    }
    None
  }, (leader, answersPhase))

  React.useEffect2(() => {
    switch (leader, playerColor == "") {
    | (_, true) | (false, false) => ()
    | (true, false) => {
        let pl: Game.startPayload = {
          action: "start",
          gameno: game.no,
        }
        send(. Js.Json.stringifyAny(pl))
      }
    }
    None
  }, (leader, playerColor))

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
    {switch (playerColor == "", currentWord == "") {
    | (true, true) =>
      <span className="animate-spin text-yellow-200 text-2xl font-bold absolute left-1/2">
        {React.string(circ)}
      </span>
    | _ => React.null
    }}
    <Word onAnimationEnd playerColor currentWord answered showTimer={currentWord != ""} />
    {switch currentWord == "" {
    | true => React.null
    | false => <Form answer_max_length answered inputText onEnter setInputText />
    }}

    // <Prompt></Prompt>
  </div>
}
