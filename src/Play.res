type answerPayload = {
    action: string,
    gameno: string,
    answer: string,
    type_: string,
    playersCount: string,
    
}


let circ = Js.String2.fromCharCode(9862)
let ANSWER_MAX_LENGTH = 12

@react.component
let make = () => {

    let (answered, setAnswered) = React.useState(_ => false)
    let (inputText, setInputText) = React.useState(_ => "")
    let (currPrevWord, setCurrPrevWord) = React.useState(_ => false)

    let sendAnswer = _ => {
        let pl = {
            action: "play",
            gameno: j`${game.no}`,
            answer: inputText->Js.String2.slice(0, ANSWER_MAX_LENGTH),
            type_: "answer",
            playersCount: j`${game.players->Js.Array2.length}`
        }
        pl->send
        true->setAnswered
        ""->setInputText
    }

    let onAnimationEnd = _ => {
        sendAnswer()
    }

    let onEnter = _ => {
        sendAnswer()
    }


    React.useEffect1(() => {
    false->setAnswered
    
    }, [currentWord])


    React.useEffect2(() => {
        switch currentWord->Js.Nullable.toOption, previousWord->Js.Nullable.toOption {
        | Some(cw), _ | _, Some(pw) => true->setCurrPrevWord
        | None, None => false->setCurrPrevWord
        }
    
    }, [currentWord, previousWord])


    <div>
        <Scoreboard playerName=user players=game.players></Scoreboard>
        {
            switch game.playing, currPrevWord {
            | true, false => <p className="text-yellow-200 text-2xl font-bold">"Get Ready"->React.string<span className="animate-spin">{circ->React.string}</span></p>
            | false, _ | true, true => React.null
            }
        }

        <Word className=switch answered {
        | true => "animate-erase"
        | false => ""
        } onAnimationEnd playerColor currentWord></Word>

        <Form ANSWER_MAX_LENGTH answered inputText onEnter setInputText></Form>

        <Prompt></Prompt>

    </div>




}