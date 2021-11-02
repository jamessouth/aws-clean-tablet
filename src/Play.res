type answerPayload = {
    action: string,
    gameno: string,
    answer: string,
    tipe: string,
    playersCount: string,
    
}


let circ = Js.String2.fromCharCode(9862)
let answer_max_length = 12

@react.component
let make = (~playerName, ~wsConnected, ~game: Reducer.game, ~leadertoken, ~playerColor, ~send, ~wsError, ~currentWord, ~previousWord) => {
  
    Js.log4(wsConnected, leadertoken, wsError, previousWord)
    let (answered, setAnswered) = React.useState(_ => false)
    let (inputText, setInputText) = React.useState(_ => "")
    // let (currPrevWord, setCurrPrevWord) = React.useState(_ => false)

    let noplrs = Js.Array2.length(game.players)

    let sendAnswer = _ => {
        let pl = {
            action: "play",
            gameno: game.no,
            answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
            tipe: "answer",
            playersCount: j`$noplrs`
        }
        send(. Js.Json.stringifyAny(pl))
        (_ => true)->setAnswered
        (_ => "")->setInputText
    }

    let onAnimationEnd = _ => {
        sendAnswer()
    }

    let onEnter = _ => {
        sendAnswer()
    }


    React.useEffect1(() => {
    (_ => false)->setAnswered
    None
    }, [currentWord])


    // React.useEffect2(() => {
    //     switch (currentWord, previousWord) {
    //     | (Some(_w), _) | (_, Some(_w)) => (_ => true)->setCurrPrevWord
    //     | (None, None) => (_ => false)->setCurrPrevWord
    //     }
    //     None
    // }, (currentWord, previousWord))


    <div>
        <Scoreboard playerName players=game.players />
        // {
        //     switch (true, currPrevWord) {//game.playing
        //     | (true, false) => <p className="text-yellow-200 text-2xl font-bold">{"Get Ready"->React.string}<span className="animate-spin">{circ->React.string}</span></p>
        //     | (false, _) | (true, true) => React.null
        //     }
        // }
        //answered
        <Word onAnimationEnd playerColor currentWord />

        <Form answer_max_length answered inputText onEnter setInputText></Form>

        // <Prompt></Prompt>

    </div>




}