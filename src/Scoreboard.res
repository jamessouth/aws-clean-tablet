



@react.component
let make = (~playerName, ~players: array<Reducer.player>) => {

    let noplrs = players->Js.Array2.length

    <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
        <h2 className="text-center mb-5">{React.string("score")}</h2>

        <p>{React.string(playerName)}</p>

        <ul className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
        {
            players->Js.Array2.map((p) => {
                switch p.color {
                | Some(c) => <li className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl" style={ReactDOM.Style.make(~backgroundColor=c, ())} key=p.connid>
                    <p>{p.name->React.string}</p>
                    <p>{p.score->React.string}</p>
                </li>
                | None => React.null
                }
            })->React.array
        }
        </ul>
    </div>

}