// let dot = Js.String2.fromCharCode(8901)

@react.component
let make = (~players: array<Reducer.livePlayer>, ~previousWord, ~showAnswers) => {
  Js.log2("score", players)

  let noplrs = players->Js.Array2.length

  <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
    <h2 className="text-center font-anon mb-5 text-warm-gray-100">
      {switch showAnswers {
      | true => React.string(previousWord)
      | false => React.string("scores")
      }}
    </h2>
    <ul
      className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
      {players
      ->Js.Array2.map(p => {
        <li
          className="w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl text-warm-gray-100"
          style={ReactDOM.Style.make(~backgroundColor=p.color, ())}
          key=p.connid>
          <p
            className={switch p.hasAnswered {
            | true => "after:content-['\\22C5'] after:text-yellow-200 after:text-5xl after:absolute after:leading-25px"
            | false => ""
            }}>
            {React.string(p.name)}
          </p>
          {switch showAnswers {
          | true => <p> {React.string(p.answer.answer)} </p>
          | false => <p> {React.string(p.score)} </p>
          }}
        </li>
      })
      ->React.array}
    </ul>
  </div>
}
