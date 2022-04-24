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
  ~sortData) => {

    let (winDir, setWinDir) = React.Uncurried.useState(_ =>
    Down
  )
    let (ptsDir, setPtsDir) = React.Uncurried.useState(_ =>
    Down
  )
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

let onClick = (field, dir, func, _e) => {
sortData(field, dir)
func(._ => switch dir {
| Up => Down
| Down => Up
})
}


  <div className="text-stone-100">
    <table className="mx-auto border-solid border-2 border-yellow-400 w-4/5 border-collapse">
      <caption className=""> {React.string("leaderboard")} </caption>
      <thead>
        <tr>
          <th>
            <p> {React.string("name")} </p>
            // <Button textTrue="name" textFalse="name" onClick/>
          </th>
          <th>
     
            <Button textTrue="wins" textFalse="wins" onClick={onClick("wins", winDir, setWinDir )} className="bg-transparent text-stone-100 after:content-['\\2193'] after:text-yellow-200 after:text-2xl after:absolute"/>
            
          </th>
          <th>
            <Button textTrue="points" textFalse="points" onClick={onClick("points", ptsDir, setPtsDir )} className="bg-transparent text-stone-100 after:content-['\\2193'] after:text-yellow-200 after:text-2xl after:absolute"/>
    
    
          </th>
          <th>
            <p> {React.string("games")} </p>
            // <Button textTrue="games" textFalse="games"/>
          </th>
        </tr>
      </thead>
      <tbody>
        {leaderData
        ->Js.Array2.mapi(({name, wins, totalPoints, games}, i) => {
          <tr className="text-center odd:bg-stone-300/20 h-8" key={j`${name}$i`}>
            <th className=""> {React.string(name)} </th>
            <td className=""> {React.string(j`$wins`)} </td>
            <td className=""> {React.string(j`$totalPoints`)} </td>
            <td className=""> {React.string(j`$games`)} </td>
          </tr>
        })
        ->React.array}
      </tbody>
    </table>
  </div>
}
