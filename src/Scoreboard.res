
type player = {
    name: string,
    connid: string,
    score: string,
    color: option<string>,
    ready: bool
}


@react.component
let make = (~_playerName, ~players) => {

    let noplyrs = players->Js.Array2.length->Js.Int.toString

    <div className="w-full" style={ReactDOM.Style.make(~height=`calc(82px + (28px * ${noplyrs}))`, ())}>
        <h2 className="mb-5">{"score"->React.string}</h2>
        <ul className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
        // {
        //     players->Js.Array2.map((p) => {
        //         switch p.color {
        //         | Some(c) => <li className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl" style={ReactDOM.Style.make(~backgroundColor=c, ())} key=p.connid>
        //             <p>{p.name->React.string}</p>
        //             <p>{p.score->React.string}</p>
        //         </li>
        //         | None => React.null
        //         }
        //     })
        // }
        </ul>
    </div>

}