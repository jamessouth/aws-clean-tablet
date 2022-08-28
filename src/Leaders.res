type propShape = {"leaderData": array<Reducer.stat>, "playerName": string}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

type sortDirection = Up | Down
let sortData = (field, dir, a: Reducer.stat, b: Reducer.stat) => {
  let res = switch field {
  | "name" => compare(a.name, b.name)
  | "wins" => compare(a.wins, b.wins)
  | "points" => compare(a.points, b.points)
  | "games" => compare(a.games, b.games)
  | "win %" => compare(a.winPct, b.winPct)
  | _ => compare(a.ppg, b.ppg)
  }

  switch dir {
  | Up => -res
  | Down => res
  }
}

@react.component
let make = (~leaderData: array<Reducer.stat>, ~playerName) => {
  let (nameDir, setNameDir) = React.Uncurried.useState(_ => Down)
  let (winDir, setWinDir) = React.Uncurried.useState(_ => Down)
  let (ptsDir, setPtsDir) = React.Uncurried.useState(_ => Up)
  let (gamesDir, setGamesDir) = React.Uncurried.useState(_ => Up)
  let (winPctDir, setWinPctDir) = React.Uncurried.useState(_ => Up)
  let (ppgDir, setPPGDir) = React.Uncurried.useState(_ => Up)

  let (sortedField, setSortedField) = React.Uncurried.useState(_ => "wins")
  let (arrow, setArrow) = React.Uncurried.useState(_ => "['\\2193']")

  let (data, setData) = React.Uncurried.useState(_ => [])

  let larrow = Js.String2.fromCharCode(8592)

  React.useEffect1(() => {
    Js.log("copyleader")
    setData(._ => leaderData->Js.Array2.copy)
    None
  }, [leaderData])

  Js.log(leaderData)
  Js.log(data)

  let onClick = (field, dir, func, _e) => {
    setSortedField(._ => field)
    switch dir {
    | Up => {
        setArrow(._ => "['\\2193']")
        func(._ => Down)
      }

    | Down => {
        setArrow(._ => "['\\2191']")
        func(._ => Up)
      }
    }

    setData(._ => data->Js.Array2.sortInPlaceWith(sortData(field, dir)))
  }

  let buttonBase = "bg-transparent text-dark-600 text-base font-anon font-bold w-full h-8"

  let arrowClass = ` relative after:content-${arrow} after:text-2xl after:font-over after:absolute`

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
          className="my-6 relative text-4xl md:my-12 md:text-5xl desk:my-18 desk:text-6xl font-fred font-bold text-shadow-lead">
          <Button
            textTrue=larrow
            textFalse=larrow
            onClick={_ => RescriptReactRouter.push("/auth/lobby")}
            className="cursor-pointer font-over text-5xl bg-transparent absolute left-10"
          />
          {React.string("Leaderboard")}
        </caption>
        <colgroup>
          <col
            className={switch sortedField == "name" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
          <col
            className={switch sortedField == "wins" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
          <col
            className={switch sortedField == "points" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
          <col
            className={switch sortedField == "games" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
          <col
            className={switch sortedField == "win %" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
          <col
            className={switch sortedField == "pts/gm" {
            | true => "bg-stone-100/12"
            | false => ""
            }}
          />
        </colgroup>
        <thead className="">
          <tr>
            <th className="sticky left-0 top-0 z-10 h-8 bg-amber-300 w-16.667vw min-w-104px">
              <Button
                textTrue="name"
                textFalse="name"
                onClick={onClick("name", nameDir, setNameDir)}
                className={switch sortedField == "name" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
            <th className="sticky top-0 h-8 bg-amber-300 w-16.667vw min-w-64px">
              <Button
                textTrue="wins"
                textFalse="wins"
                onClick={onClick("wins", winDir, setWinDir)}
                className={switch sortedField == "wins" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
            <th className="sticky top-0 h-8 bg-amber-300 w-16.667vw min-w-80px">
              <Button
                textTrue="points"
                textFalse="points"
                onClick={onClick("points", ptsDir, setPtsDir)}
                className={switch sortedField == "points" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
            <th className="sticky top-0 h-8 bg-amber-300 w-16.667vw min-w-72px">
              <Button
                textTrue="games"
                textFalse="games"
                onClick={onClick("games", gamesDir, setGamesDir)}
                className={switch sortedField == "games" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
            <th className="sticky top-0 h-8 bg-amber-300 w-16.667vw min-w-72px">
              <Button
                textTrue="win %"
                textFalse="win %"
                onClick={onClick("win %", winPctDir, setWinPctDir)}
                className={switch sortedField == "win %" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
            <th className="sticky top-0 h-8 bg-amber-300 w-16.667vw min-w-80px">
              <Button
                textTrue="pts/gm"
                textFalse="pts/gm"
                onClick={onClick("pts/gm", ppgDir, setPPGDir)}
                className={switch sortedField == "pts/gm" {
                | true => buttonBase ++ arrowClass
                | false => buttonBase
                }}
              />
            </th>
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

// let zzz: array<Reducer.stat> = [
//   {name: "mmmmmmmmmm", wins: 241, points: 5152, games: 83, winPct: 0.12, ppg: 11.0},
//   {name: "stu", wins: 121, points: 121, games: 32, winPct: 0.22, ppg: 11.1},
//   {name: "liz", wins: 50, points: 59, games: 363, winPct: 0.32, ppg: 11.2},
//   {name: "abner", wins: 42, points: 18, games: 173, winPct: 0.42, ppg: 11.3},
//   {name: "harold", wins: 40, points: 97, games: 333, winPct: 0.52, ppg: 11.4},
//   {name: "stacie", wins: 32, points: 17, games: 313, winPct: 0.62, ppg: 11.5},
//   {name: "marcy", wins: 30, points: 91, games: 332, winPct: 0.72, ppg: 11.6},
//   {name: "wes", wins: 22, points: 11, games: 213, winPct: 0.82, ppg: 11.7},
//   {name: "carl", wins: 21, points: 12, games: 23, winPct: 0.92, ppg: 11.8},
//   {name: "bill", wins: 10, points: 9, games: 323, winPct: 0.02, ppg: 11.9},
//   {name: "test", wins: 2, points: 1, games: 13, winPct: 0.17, ppg: 12.8},
// ]
