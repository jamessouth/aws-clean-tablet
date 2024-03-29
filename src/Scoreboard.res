let hstyles = "text-center font-anon mb-5 text-stone-100 "

@react.component
let make = (
  ~players,
  ~oldWord,
  ~showAnswers,
  ~winner,
  ~isGameOver,
  ~resetConnState,
  ~playerName,
  ~noplrs,
) => {
  Js.log2("score", players)

  let (bgimg, setBgimg) = React.Uncurried.useState(_ => "")

  React.useEffect3(() => {
    let bg = switch isGameOver {
    | true =>
      switch winner == playerName {
      | true => "bg-win"
      | false => {
          let i = Js.Math.unsafe_trunc(Js.Math.random() *. 4.) + 1
          j`bg-lose$i`
        }
      }
    | false => ""
    }
    setBgimg(._ => bg)

    None
  }, (isGameOver, winner, playerName))

  <div className="w-full" style={ReactDOM.Style.make(~height=j`calc(82px + (28px * $noplrs))`, ())}>
    {switch showAnswers {
    | true =>
      <>
        <p className="text-center font-anon font-bold text-stone-100 text-xl mb-2">
          {React.string("Answers for:")}
        </p>
        <h2 className=hstyles> {React.string(oldWord)} </h2>
      </>
    | false =>
      <>
        <p className="h-7 mb-2" />
        <h2
          className={switch !isGameOver {
          | true => hstyles
          | false => hstyles ++ "animate-blink"
          }}>
          {switch !isGameOver {
          | false => {
              let hiScore: Reducer.player = players->Js.Array2.unsafe_get(0)
              React.string(winner ++ " wins with " ++ hiScore.score ++ " points!")
            }

          | true => React.string("Scores:")
          }}
        </h2>
      </>
    }}
    <ul
      className="bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center">
      {players
      ->Js.Array2.mapi((p, i) => {
        <li
          className={"w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl text-stone-100 " ++ if (
            isGameOver && i != 0
          ) {
            "filter brightness-30 contrast-60"
          } else if isGameOver && i == 0 {
            "animate-rotate"
          } else {
            ""
          }}
          key={j`${p.name}$i`}
          style={ReactDOM.Style.make(~backgroundColor=p.color, ())}>
          <p
            className={switch p.hasAnswered {
            | true => "dot after:(text-yellow-200 text-5xl absolute leading-25px)"
            | false =>
              switch (isGameOver, i == 0) {
              | (true, true) => "text-shadow-win"
              | (false, true) | (true, false) | (false, false) => ""
              }
            }}>
            {React.string(p.name)}
          </p>
          {switch showAnswers {
          | true =>
            <>
              <p className="animate-pulse font-luck"> {React.string("+" ++ p.pointsThisRound)} </p>
              <p> {React.string(p.answer)} </p>
            </>
          | false =>
            <p
              className={switch (isGameOver, i == 0) {
              | (true, true) => "text-shadow-win"
              | (false, true) | (true, false) | (false, false) => ""
              }}>
              {React.string(p.score)}
            </p>
          }}
        </li>
      })
      ->React.array}
    </ul>
    {switch isGameOver {
    | true =>
      <div className={`w-64 h-96 bg-no-repeat opacity-0 m-auto animate-fadein ${bgimg}`}>
        <Button
          onClick={_ => {
            resetConnState(.)
            Route.push(Auth({subroute: Lobby}))
          }}
          className="mt-1.5 mb-4 block cursor-pointer text-stone-800 font-perm mx-auto px-8 py-2 text-2xl">
          {React.string("Return to lobby")}
        </Button>
      </div>

    | false => React.null
    }}
  </div>
}
