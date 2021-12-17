

let dot = Js.String2.fromCharCode(8901)

@react.component
let make = (~players: array<Reducer.livePlayer>, ~currentWord, ~previousWord) => {

    Js.log2("score", players)
    let noplrs = players->Js.Array2.length

    <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
        <h2 className="text-center mb-5">{React.string("scores")}</h2>

        // <p className="absolute">{React.string(playerName)}</p>

        <ul className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
        {
            players->Js.Array2.map((p) => {
                <li className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl" style={ReactDOM.Style.make(~backgroundColor=p.color, ())} key=p.connid>
                    <p>{p.name->React.string}</p>
                    {switch p.hasAnswered {
                    | true => <span className="text-yellow-200">{React.string(dot)}</span>
                    | false => React.null
                    }}
                    {switch (currentWord != "", previousWord != "") {
                    | (true, false) | (false, false) => <p>{React.string(p.score)}</p>
                    | (false, true) => <p>{React.string(p.answer.answer)}</p>
                    | (true, true) => React.null
                    }}
                </li>
              
            })->React.array
        }
        </ul>
    </div>

}