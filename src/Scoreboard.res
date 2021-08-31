



@react.component
let make = (~playerName, ~players) => {

    let noplyrs = players->Js.Array2.length

    <div className="w-full" style={ReactDOM.Style.make(~height=`calc(82px + (28px * ${noplyrs}))`, ())}>
        <h2 className="mb-5">{"score"->React.string}</h2>
        <ul className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
        {
            players->JS.Array2.map((p) => {
                switch p.color->Js.Nullable.toOption {
                | Some(c) => <li className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl" style={ReactDOM.Style.make(~backgroundColor=p.color, ())} key=p.connid>
                    <p>{p.name}</p>
                    <p>{p.score}</p>
                </li>
                | None => React.null
                }
            })
        }
        </ul>
    </div>

}