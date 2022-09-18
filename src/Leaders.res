%%raw(`import './css/leader.css'`)

type propShape = {"leaderData": array<Reducer.stat>, "playerName": string}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

type field =
  | Name
  | Wins
  | Points
  | Games
  | WinPercentage
  | PointsPerGame

type sortDirection = Up | Down

let sortData = (field, dir, a: Reducer.stat, b: Reducer.stat) => {
  let res = switch field {
  | Name => compare(a.name, b.name)
  | Wins => compare(a.wins, b.wins)
  | Points => compare(a.points, b.points)
  | Games => compare(a.games, b.games)
  | WinPercentage => compare(a.winPct, b.winPct)
  | PointsPerGame => compare(a.ppg, b.ppg)
  }

  switch dir {
  | Up => -res
  | Down => res
  }
}

@react.component
let make = (~leaderData, ~playerName) => {
  let (nameDir, setNameDir) = React.Uncurried.useState(_ => Down)
  let (winDir, setWinDir) = React.Uncurried.useState(_ => Down)
  let (ptsDir, setPtsDir) = React.Uncurried.useState(_ => Up)
  let (gamesDir, setGamesDir) = React.Uncurried.useState(_ => Up)
  let (winPctDir, setWinPctDir) = React.Uncurried.useState(_ => Up)
  let (ppgDir, setPPGDir) = React.Uncurried.useState(_ => Up)

  let (sortedField, setSortedField) = React.Uncurried.useState(_ => Wins)
  let (arrowClass, setArrowClass) = React.Uncurried.useState(_ => "downArrow")

  let (data, setData) = React.Uncurried.useState(_ => [])

  React.useEffect1(() => {
    setData(._ => leaderData->Js.Array2.copy)
    None
  }, [leaderData])

  let onClick = (field, dir, func, _e) => {
    setSortedField(._ => field)
    switch dir {
    | Up => {
        setArrowClass(._ => "downArrow")
        func(._ => Down)
      }

    | Down => {
        setArrowClass(._ => "upArrow")
        func(._ => Up)
      }
    }

    setData(._ => data->Js.Array2.sortInPlaceWith(sortData(field, dir)))
  }

  <div className="leadermobbg leadertabbg leaderbigbg w-100vw h-100vh overflow-y-scroll leader">
    {switch data->Js.Array2.length == 0 {
    | true =>
      <>
        <div className="h-42vh" />
        <Loading label="data..." />
      </>
    | false =>
      <table
        className="border-collapse text-dark-600 font-anon table-fixed tablewidth:mx-8 lg:mx-16 desk:mx-32">
        <caption
          className="my-6 relative text-4xl md:(my-12 text-5xl) desk:(my-18 text-6xl) font-fred font-bold text-shadow-lead">
          <Button
            onClick={_ => Route.push(Lobby)}
            className="cursor-pointer font-over text-5xl bg-transparent absolute left-10">
            {React.string("←")}
          </Button>
          {React.string("Leaderboard")}
        </caption>
        <colgroup>
          {[Name, Wins, Points, Games, WinPercentage, PointsPerGame]
          ->Js.Array2.map(c =>
            <col
              key={j`$c`}
              className={switch sortedField == c {
              | true => "bg-stone-100/17"
              | false => ""
              }}
            />
          )
          ->React.array}
        </colgroup>
        <thead className="">
          <tr>
            {[
              ("min-w-104px left-0 z-10", "name", onClick(Name, nameDir, setNameDir), Name),
              ("min-w-64px", "wins", onClick(Wins, winDir, setWinDir), Wins),
              ("min-w-80px", "points", onClick(Points, ptsDir, setPtsDir), Points),
              ("min-w-72px", "games", onClick(Games, gamesDir, setGamesDir), Games),
              (
                "min-w-72px",
                "win %",
                onClick(WinPercentage, winPctDir, setWinPctDir),
                WinPercentage,
              ),
              ("min-w-80px", "pts/gm", onClick(PointsPerGame, ppgDir, setPPGDir), PointsPerGame),
            ]
            ->Js.Array2.map(c => {
              let (cn, btnText, oc, field) = c
              <th key=btnText className={"sticky top-0 h-8 bg-amber-300 w-16.667vw " ++ cn}>
                <Button
                  onClick=oc
                  className={"bg-transparent cursor-pointer text-dark-600 text-base font-anon font-bold w-full h-8" ++ if (
                    sortedField == field
                  ) {
                    ` relative ${arrowClass} after:(text-2xl font-over absolute)`
                  } else {
                    ""
                  }}>
                  {React.string(btnText)}
                </Button>
              </th>
            })
            ->React.array}
          </tr>
        </thead>
        <tbody>
          {data
          ->Js.Array2.mapi(({name, wins, points, games, winPct, ppg}, i) => {
            <tr
              className={switch name == playerName {
              | true => "text-center bg-blue-200/66 h-8 uppercase italic"
              | false => "text-center odd:bg-stone-100/16 h-8"
              }}
              key={j`${name}$i`}>
              <th className="sticky left-0 bg-amber-200"> {React.string(name)} </th>
              <td className=""> {React.string(j`$wins`)} </td>
              <td className=""> {React.string(j`$points`)} </td>
              <td className=""> {React.string(j`$games`)} </td>
              <td className="">
                {switch winPct == 0. || winPct == 1. {
                | true => React.string(j`$winPct.00`)
                | false => React.string(j`$winPct`)
                }}
              </td>
              <td className=""> {React.string(j`$ppg`)} </td>
            </tr>
          })
          ->React.array}
          <tr className="h-50vh" />
        </tbody>
      </table>
    }}
  </div>
}
