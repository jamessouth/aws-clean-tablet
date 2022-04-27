type leaderPayload = {
  action: string,
  info: string,
}

type sortDirection =
  | Up
  | Down

@react.component
let make = (
  // ~send,
  ~leaderData: array<Reducer.stat>,
) => {
  let (nameDir, setNameDir) = React.Uncurried.useState(_ => Down)
  let (winDir, setWinDir) = React.Uncurried.useState(_ => Up)
  let (ptsDir, setPtsDir) = React.Uncurried.useState(_ => Up)
  let (gamesDir, setGamesDir) = React.Uncurried.useState(_ => Up)

  let (sortedField, setSortedField) = React.Uncurried.useState(_ => "")
  let (arrow, setArrow) = React.Uncurried.useState(_ => "")

  let (dt, setDt) = React.Uncurried.useState(_ => leaderData)

  let strSortName = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
    switch dir {
    | Up =>
      switch n2.name > n1.name {
      | true => 1
      | false => -1
      }
    | Down =>
      switch n2.name <= n1.name {
      | true => 1
      | false => -1
      }
    }
  let numSortWins = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
    switch dir {
    | Up => n2.wins - n1.wins
    | Down => n1.wins - n2.wins
    }
  let numSortPoints = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
    switch dir {
    | Up => n2.totalPoints - n1.totalPoints
    | Down => n1.totalPoints - n2.totalPoints
    }
  let numSortGames = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
    switch dir {
    | Up => n2.games - n1.games
    | Down => n1.games - n2.games
    }

  // React.useEffect0(() => {
  //   let pl = {
  //     action: "leaders",
  //     info: "hello",
  //   }
  //   send(. Js.Json.stringifyAny(pl))
  //   None
  // })

  Js.log(leaderData)

  let sortData = (input, dir) => {
    let arr = Js.Array2.copy(dt)
    switch input {
    | "name" => arr->Js.Array2.sortInPlaceWith(strSortName(dir))
    | "wins" => arr->Js.Array2.sortInPlaceWith(numSortWins(dir))
    | "points" => arr->Js.Array2.sortInPlaceWith(numSortPoints(dir))
    | "games" => arr->Js.Array2.sortInPlaceWith(numSortGames(dir))
    | _ => []
    }->ignore
    setDt(._ => arr)
  }

  let onClick = (field, dir, func, _e) => {
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
    setSortedField(._ => field)
    sortData(field, dir)
  }

  let buttonBase = "bg-transparent text-dark-600 text-base font-anon font-bold"

  let arrowClass = ` relative after:content-${arrow} after:text-2xl after:font-over after:absolute`

  <div className="leaderbg overflow-y-scroll">
    <table className="w-full border-collapse text-dark-600 font-anon">
      <caption className="my-4 text-4xl font-fred font-bold text-shadow-lead">
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
      </colgroup>
      <thead className="sticky top-0 h-8 bg-amber-200">
        <tr>
          <th>
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
          <th>
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
          <th>
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
          <th>
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
        </tr>
      </thead>
      <tbody>
        {dt
        ->Js.Array2.mapi(({name, wins, totalPoints, games}, i) => {
          <tr className="text-center odd:bg-stone-100/16 h-8" key={j`${name}$i`}>
            <th className=""> {React.string(name)} </th>
            <td className=""> {React.string(j`$wins`)} </td>
            <td className=""> {React.string(j`$totalPoints`)} </td>
            <td className=""> {React.string(j`$games`)} </td>
          </tr>
        })
        ->React.array}
        <tr className="h-50vh" />
      </tbody>
    </table>
  </div>
}
