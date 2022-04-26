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
  let (winDir, setWinDir) = React.Uncurried.useState(_ => Down)
  let (ptsDir, setPtsDir) = React.Uncurried.useState(_ => Down)

  let (sortedField, setSortedField) = React.Uncurried.useState(_ => "")
  let (arrow, setArrow) = React.Uncurried.useState(_ => "")

  let (dt, setDt) = React.Uncurried.useState(_ => leaderData)

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

  // React.useEffect0(() => {
  //   let pl = {
  //     action: "leaders",
  //     info: "hello",
  //   }
  //   send(. Js.Json.stringifyAny(pl))
  //   None
  // })
  // let zzz:array<Reducer.stat> = [{name:"test",wins:2,totalPoints:15,games:3},{name:"test3",wins:1,totalPoints:12,games:3},{name:"test2",wins:0,totalPoints:9,games:3}]
  Js.log(leaderData)

    let sortData = (input, dir) => {
        let arr = Js.Array2.copy(dt)
    switch input {
    | "wins" => {
        Js.Array.sortInPlaceWith(numSortWins(dir), arr)->ignore
        setDt(._ => arr)
      }
    | "points" => {
        Js.Array.sortInPlaceWith(numSortPoints(dir), arr)->ignore
        setDt(._ => arr)
      }
    | _ => ()
    }
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

  let arrowClass = ` relative after:-top-8px after:content-${arrow} after:text-stone-300 after:text-2xl after:absolute`

  <div className="text-stone-800 leaderbg overflow-y-scroll">
    <table className="w-full shadow-lead border-collapse">
      // <caption className="mb-8 text-4xl font-fred"> {React.string("Leaderboard")} </caption>
      <colgroup>
        <col />
        <col
          className={switch sortedField == "wins" {
          | true => "bg-stone-800/14"
          | false => ""
          }}
        />
        <col
          className={switch sortedField == "points" {
          | true => "bg-stone-800/14"
          | false => ""
          }}
        />
        <col />
      </colgroup>
      <thead className="sticky top-0">
        <tr>
          <th>
            <p> {React.string("name")} </p>
            // <Button textTrue="name" textFalse="name" onClick/>
          </th>
          <th
            className={switch sortedField == "wins" {
            | true => "text-left"
            | false => ""
            }}>
            <Button
              textTrue="wins"
              textFalse="wins"
              onClick={onClick("wins", winDir, setWinDir)}
              className={switch sortedField == "wins" {
              | true => "bg-transparent text-stone-100" ++ arrowClass
              | false => "bg-transparent text-stone-100"
              }}
            />
          </th>
          <th
            className={switch sortedField == "points" {
            | true => "text-left"
            | false => ""
            }}>
            <Button
              textTrue="points"
              textFalse="points"
              onClick={onClick("points", ptsDir, setPtsDir)}
              className={switch sortedField == "points" {
              | true => "bg-transparent text-stone-100" ++ arrowClass
              | false => "bg-transparent text-stone-100"
              }}
            />
          </th>
          <th>
            <p> {React.string("games")} </p>
            // <Button textTrue="games" textFalse="games"/>
          </th>
        </tr>
      </thead>
      <tbody>
        {dt
        ->Js.Array2.mapi(({name, wins, totalPoints, games}, i) => {
          <tr className="text-center odd:bg-stone-800/14 h-8" key={j`${name}$i`}>
            <th className=""> {React.string(name)} </th>
            <td className=""> {React.string(j`$wins`)} </td>
            <td className=""> {React.string(j`$totalPoints`)} </td>
            <td className=""> {React.string(j`$games`)} </td>
          </tr>
        })
        ->React.array}
        <tr className="text-center h-8" key={"900"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
        <tr className="text-center h-8" key={"901"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
        <tr className="text-center h-8" key={"902"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
        <tr className="text-center h-8" key={"903"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
        <tr className="text-center h-8" key={"904"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
        <tr className="text-center h-8" key={"905"}>
            <th className=""> {React.string("")} </th>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
            <td className=""> {React.string("")} </td>
          </tr>
      </tbody>
    </table>
  </div>
}
