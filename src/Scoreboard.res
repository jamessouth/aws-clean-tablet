@react.component
let make = (
  ~players: array<Reducer.livePlayer>,
  ~oldWord,
  ~showAnswers,
  ~winner,
  ~isWinner,
  ~onClickTrue,
  ~onClickFalse,
  ~playerName,
  ~noplrs,
) => {
  Js.log2("score", players)

  let hstyles = "text-center font-anon mb-5 text-stone-100 "
  let (count, setCount) = React.Uncurried.useState(_ => 30)
  let (bgimg, setBgimg) = React.Uncurried.useState(_ => "")

  React.useEffect3(() => {
    let bg = switch isWinner {
    | true =>
      switch winner == playerName {
      | true => "bg-win"
      | false =>
        switch Js.Math.unsafe_trunc(Js.Math.random() *. 4.) + 1 {
        | 1 => "bg-lose1"
        | 2 => "bg-lose2"
        | 3 => "bg-lose3"
        | _ => "bg-lose4"
        }
      }
    | false => ""
    }
    setBgimg(._ => bg)

    None
  }, (isWinner, winner, playerName))

  React.useEffect1(() => {
    Js.log("useff run")
    let id = if isWinner {
      Js.Global.setInterval(() => {
        setCount(. c => c - 1)
      }, 1000)
    } else {
      Js.Global.setInterval(() => (), 3600000)
    }

    Some(
      () => {
        Js.Global.clearInterval(id)
      },
    )
  }, [isWinner])

  React.useEffect1(() => {
    switch count == 0 {
    | true => onClickFalse()
    | false => ()
    }
    None
  }, [count])

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
          className={switch !isWinner {
          | true => hstyles
          | false => hstyles ++ "animate-blink"
          }}>
          {switch !isWinner {
          | false => {
              let hiScore = players->Js.Array2.unsafe_get(0)
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
            isWinner && i != 0
          ) {
            "filter brightness-30 contrast-60"
          } else if isWinner && i == 0 {
            "animate-rotate"
          } else {
            ""
          }}
          key={j`${p.name}$i`}
          style={ReactDOM.Style.make(~backgroundColor=p.color, ())}>
          <p
            className={switch p.hasAnswered {
            | true => "after:content-['\\22C5'] after:text-yellow-200 after:text-5xl after:absolute after:leading-25px"
            | false =>
              switch (isWinner, i == 0) {
              | (true, true) => "text-shadow-win"
              | _ => ""
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
              className={switch (isWinner, i == 0) {
              | (true, true) => "text-shadow-win"
              | _ => ""
              }}>
              {React.string(p.score)}
            </p>
          }}
        </li>
      })
      ->React.array}
    </ul>
    {switch count < 26 {
        | true =>
          <p className="font-perm absolute left-1/2 transform -translate-x-2/4 text-stone-100 text-xl">
            {React.string(j`Returning to lobby in: $count`)}
          </p>
        | false => React.null
        }}
    {switch isWinner {
    | true =>
      <>
        <div className={`w-64 h-96 bg-no-repeat opacity-0 m-auto animate-fadein ${bgimg}`} />
        <Button
          textTrue="Return to lobby"
          textFalse="Return to lobby"
          onClick=onClickTrue
          className="mt-1.5 mb-14 block cursor-pointer text-stone-800 font-perm mx-auto px-8 py-2 text-2xl"
        />
      </>
    | false => React.null
    }}
  </div>
}
