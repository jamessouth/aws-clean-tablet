@react.component
let make = (
  ~players: array<Reducer.livePlayer>,
  ~previousWord,
  ~showAnswers,
  ~winner,
  ~onClick,
) => {
  Js.log2("score", players)

  let className = "mt-10 block cursor-pointer text-stone-800 font-perm m-auto px-8 py-2 text-2xl"

  let noplrs = Js.Array2.length(players)

  <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
    {switch showAnswers {
    | true => <>
        <p className="font-anon font-bold text-stone-100 text-xl">
          {React.string("Answers for:")}
        </p>
        <h2 className="text-center font-anon mb-5 text-stone-100">
          {React.string(previousWord)}
        </h2>
      </>
    | false => <>
        <p className="h-7" />
        <h2 className="text-center font-anon mb-5 text-stone-100">
          {switch winner == "" {
          | false => React.string(winner ++ " wins!")
          | true => React.string("Scores:")
          }}
        </h2>
      </>
    }}
    <ul
      className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
      {players
      ->Js.Array2.mapi((p, i) => {
        let pts_ans = p.answer->Js.String2.split("_")

        let (points, answer) = (pts_ans->Js.Array2.unsafe_get(0), pts_ans->Js.Array2.unsafe_get(1))

        <li
          className={"w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl text-stone-100 " ++ if (
            winner != "" && i != 0
          ) {
            "filter brightness-25"
          } else if winner != "" && i == 0 {
            "animate-rotate"
          } else {
            ""
          }}
          key={j`${p.name}$i`}
          style={ReactDOM.Style.make(~backgroundColor=p.color, ())}>
          <p
            className={switch p.hasAnswered {
            | true => "after:content-['\\22C5'] after:text-yellow-200 after:text-5xl after:absolute after:leading-25px"
            | false => ""
            }}>
            {React.string(p.name)}
          </p>
          {switch showAnswers {
          | true => <> <p className="animate-pulse font-luck"> {React.string("+" ++ points)} </p> <p> {React.string(answer)} </p> </>
          | false => <p> {React.string(p.score)} </p>
          }}
        </li>
      })
      ->React.array}
    </ul>
    {switch winner == "" {
    | false => <Button textTrue="Return to lobby" textFalse="Return to lobby" onClick className />
    | true => React.null
    }}
  </div>
}
