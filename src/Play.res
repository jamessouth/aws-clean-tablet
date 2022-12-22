let answer_max_length = 12

@react.component
let make = (
  ~players,
  ~playerLiveGame,
  ~showAnswers,
  ~winner,
  ~isGameOver,
  ~oldWord,
  ~word,
  ~playerColor,
  ~send,
  ~playerName,
  ~resetConnState,
) => {
  let (answered, setAnswered) = React.Uncurried.useState(_ => false)
  let (answer, setAnswer) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "ANSWER: 2-12 length; letters and spaces only; ",
  ))

  React.useEffect1(() => {
    ErrorHook.useError(answer, Answer, setValidationError)
    None
  }, [answer])

  React.useEffect1(() => {
    setAnswered(._ => false)
    None
  }, [showAnswers])

  let hasRendered = React.useRef(false)
  Js.log3("play", players, hasRendered)

  let sendAnswer = _ => {
    let ans =
      answer
      ->Js.String2.slice(~from=0, ~to_=answer_max_length)
      ->Js.String2.replaceByRe(%re("/\d/g"), "")
      ->Js.String2.replaceByRe(%re("/[!-/:-@\[-`{-~]/g"), "")
      ->Js.String2.trim(_)

    send(.
      Lobby.payloadToObj({
        act: Answer,
        gn: Gameno({no: playerLiveGame}),
        cmd: Custom({cv: ans}),
      }),
    )
    setAnswered(._ => true)
    setAnswer(._ => "")
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

  let onClick2 = (. ()) => {
    switch validationError {
    | None => onEnter(. ignore())
    | Some(_) => ()
    }
  }

  let onKeyPress = e => {
    let key = ReactEvent.Keyboard.key(e)
    switch key {
    | "Enter" => onClick2(.)
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
    <Scoreboard
      players
      oldWord
      showAnswers
      winner
      isGameOver
      resetConnState
      playerName
      noplrs={Js.Array2.length(players)}
    />
    {switch isGameOver {
    | true => React.null
    | false =>
      <>
        <Word onAnimationEnd playerColor word answered showTimer={word != ""} />
        {switch word == "" {
        | true => <div className="bg-transparent h-45 w-full" />
        | false =>
          switch answered {
          | true => <div className="bg-transparent h-45 w-full" />
          | false =>
            <Form ht="h-24" on_Click=onClick2 leg="" validationError cognitoError=None>
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
    }}
  </div>
}
