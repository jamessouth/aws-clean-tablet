type leaderPayload = {
  action: string,
  info: string,
}

@react.component
let make = (
  // ~send,
  ~leaderData: array<Reducer.stat>,
) => {
  open Sorting
  let (nameDir, setNameDir) = React.Uncurried.useState(_ => Down)
  let (winDir, setWinDir) = React.Uncurried.useState(_ => Down)
  let (ptsDir, setPtsDir) = React.Uncurried.useState(_ => Up)
  let (gamesDir, setGamesDir) = React.Uncurried.useState(_ => Up)
  let (winPctDir, setWinPctDir) = React.Uncurried.useState(_ => Up)
  let (ppgDir, setPPGDir) = React.Uncurried.useState(_ => Up)

  let (sortedField, setSortedField) = React.Uncurried.useState(_ => "wins")
  let (arrow, setArrow) = React.Uncurried.useState(_ => "['\\2193']")

  let (data, setData) = React.Uncurried.useState(_ => leaderData)

  // React.useEffect0(() => {
  //   let pl = {
  //     action: "leaders",
  //     info: "hello",
  //   }
  //   send(. Js.Json.stringifyAny(pl))
  //   None
  // })

  Js.log(leaderData)

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
    sortData(field, dir, data, setData)
  }

  let buttonBase = "bg-transparent text-dark-600 text-base font-anon font-bold w-full h-8"

  let arrowClass = ` relative after:content-${arrow} after:text-2xl after:font-over after:absolute`

  <div className="leaderbg overflow-y-scroll">
    <table className="w-full border-collapse text-dark-600 font-anon table-fixed">
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
      <thead className="sticky top-0 h-8 bg-amber-200">
        <tr>
          <th className="first:w-16.7vw first:min-w-104px">
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
          <th>
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
          <th>
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
        ->Js.Array2.mapi(({name, wins, totalPoints, games, winPct, ppg}, i) => {
          <tr className="text-center odd:bg-stone-100/16 h-8" key={j`${name}$i`}>
            <th className=""> {React.string(name)} </th>
            <td className=""> {React.string(j`$wins`)} </td>
            <td className=""> {React.string(j`$totalPoints`)} </td>
            <td className=""> {React.string(j`$games`)} </td>
            <td className=""> {React.string(j`$winPct`)} </td>
            <td className=""> {React.string(j`$ppg`)} </td>
          </tr>
        })
        ->React.array}
        <tr className="h-50vh" />
      </tbody>
    </table>
  </div>
}
