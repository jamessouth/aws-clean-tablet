type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
  index: int,
}

type scorePayload = {
  action: string,
  game: Reducer.liveGame,
}

@react.component
let make = (~game: Reducer.liveGame, ~playerColor, ~send, ~leader) => {
  let answer_max_length = 12

  let (answered, setAnswered) = React.Uncurried.useState(_ => false)
  let (inputText, setInputText) = React.Uncurried.useState(_ => "")
  let {players, currentWord, previousWord, showAnswers, sk} = game

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
    switch (leader, showAnswers, game.winner == "") {
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
  }, (leader, showAnswers, game.winner))

  let sendAnswer = _ => {
    let index = switch players->Js.Array2.find(p => p.color == playerColor) {
    | Some(v) => v.index
    | None => -1
    }
    let pl = {
      action: "answer",
      gameno: sk,
      answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
      index: index,
    }
    send(. Js.Json.stringifyAny(pl))
    setAnswered(._ => true)
    setInputText(._ => "")
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
    // reset conns, delete game

    let pl: Game.lobbyPayload = {
      action: "lobby",
      gameno: sk,
      tipe: "gameover",
    }
    send(. Js.Json.stringifyAny(pl))

    RescriptReactRouter.push("/lobby")
  }

  <div>
    // playerName
    <Scoreboard players previousWord showAnswers winner=game.winner onClick />
    {switch game.winner == "" {
    | false => React.null
    | true => <>
        <Word onAnimationEnd playerColor currentWord answered showTimer={currentWord != ""} />
        <Answer answer_max_length answered inputText onEnter setInputText currentWord />
      </>
    }}

    // <Prompt></Prompt>
  </div>
}
