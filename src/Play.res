type answerPayload = {
    action: string,
    gameno: string,
    answer: string,
    type_: string,
    playersCount: string,
    
}


let circ = Js.String2.fromCharCode(9862)
let answer_max_length = 12

@react.component
let make = () => {
    let currentWord = Some("bill")
    let previousWord = Some("duck")

    let user = "mark"

    let (answered, setAnswered) = React.useState(_ => false)
    let (inputText, setInputText) = React.useState(_ => "")
    let (currPrevWord, setCurrPrevWord) = React.useState(_ => false)

    let sendAnswer = _ => {
        let pl = {
            action: "play",
            // gameno: j`${game.no}`,
            gameno: "555",
            answer: inputText->Js.String2.slice(~from=0, ~to_=answer_max_length),
            type_: "answer",
            // playersCount: j`${game.players->Js.Array2.length}`
            playersCount: "5"
        }
        // pl->send
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


    React.useEffect2(() => {
        switch (currentWord, previousWord) {
        | (Some(w), _) | (_, Some(w)) => (_ => true)->setCurrPrevWord
        | (None, None) => (_ => false)->setCurrPrevWord
        }
        None
    }, (currentWord, previousWord))


    <div>
        <Scoreboard playerName=user players=["a", "b"]></Scoreboard> //game.players
        {
            switch (true, currPrevWord) {//game.playing
            | (true, false) => <p className="text-yellow-200 text-2xl font-bold">{"Get Ready"->React.string}<span className="animate-spin">{circ->React.string}</span></p>
            | (false, _) | (true, true) => React.null
            }
        }

        // <Word className={switch answered {
        // | true => "animate-erase"
        // | false => ""
        // }} onAnimationEnd playerColor currentWord></Word>

        <Form answer_max_length answered inputText onEnter setInputText></Form>

        // <Prompt></Prompt>

    </div>




}