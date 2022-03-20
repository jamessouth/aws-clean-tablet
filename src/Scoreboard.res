@react.component
let make = (~players: array<Reducer.livePlayer>, ~previousWord, ~showAnswers, ~winner) => {
  Js.log2("score", players)

  let className = "block cursor-pointer text-warm-gray-800 font-perm m-auto px-8 py-2 text-2xl"

  let noplrs = Js.Array2.length(players)

  let onClick = _ => {
    // let pl: Game.lobbyPayload = {
    //   action: "lobby",
    //   gameno: "new",
    //   tipe: "join",
    // }
    // send(. Js.Json.stringifyAny(pl))

    RescriptReactRouter.push("/lobby")
  }

  <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
    <h2 className="text-center font-anon mb-5 text-warm-gray-100">
      {switch showAnswers {
      | true => React.string(previousWord)
      | false =>
        switch winner {
        | true => React.string("Winner!")
        | false => React.string("scores")
        }
      }}
    </h2>
    <ul
      className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
      {switch winner {
      | true => <>
          {players
          ->Js.Array2.mapi((p, i) => {
            switch i == 0 {
            | false => React.null
            | true =>
              <li
                className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl text-warm-gray-100"
                style={ReactDOM.Style.make(~backgroundColor=p.color, ())}
                key={j`${p.name}$i`}>
                <p> {React.string(p.name)} </p> <p> {React.string(p.score)} </p>
              </li>
            }
          })
          ->React.array}
          <Button textTrue="Return to lobby" textFalse="Return to lobby" onClick className />
        </>
      | false => players
        ->Js.Array2.mapi((p, i) => {
          <li
            className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl text-warm-gray-100"
            style={ReactDOM.Style.make(~backgroundColor=p.color, ())}
            key={j`${p.name}$i`}>
            <p
              className={switch p.hasAnswered {
              | true => "after:content-['\\22C5'] after:text-yellow-200 after:text-5xl after:absolute after:leading-25px"
              | false => ""
              }}>
              {React.string(p.name)}
            </p>
            {switch showAnswers {
            | true => <p> {React.string(p.answer)} </p>
            | false => <p> {React.string(p.score)} </p>
            }}
          </li>
        })
        ->React.array
      }}
    </ul>
  </div>
}
