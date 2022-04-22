type leaderPayload = {
  action: string,
  info: string,
}

@react.component
let make = (~send, ~rawLeaders: array<Reducer.stat>) => {
  React.useEffect0(() => {
    let pl = {
      action: "leaders",
      info: "hello",
    }
    send(. Js.Json.stringifyAny(pl))
    None
  })
  // let zzz:array<Reducer.stat> = [{name:"test",wins:2,totalPoints:15,games:3},{name:"test3",wins:1,totalPoints:12,games:3},{name:"test2",wins:0,totalPoints:9,games:3}]
  <div className="text-stone-100">
    <caption className="text-center"> {React.string("leaderboard")} </caption>
    <table className="mx-auto border-solid border-2 border-yellow-400 w-4/5">
      <thead>
        <tr>
          <th>
            <p> {React.string("name")} </p>
            // <Button textTrue="name" textFalse="name"/>
          </th>
          <th>
            <p> {React.string("wins")} </p>
            // <Button textTrue="wins" textFalse="wins"/>
          </th>
          <th>
            <p> {React.string("points")} </p>
            // <Button textTrue="points" textFalse="points"/>
          </th>
          <th>
            <p> {React.string("games")} </p>
            // <Button textTrue="games" textFalse="games"/>
          </th>
        </tr>
      </thead>
      <tbody>
        {rawLeaders
        ->Js.Array2.mapi((s, i) => {
          let {name, wins, totalPoints, games} = s
          <tr className="text-center odd:bg-blue-500" key={j`${name}$i`}>
            <td className=""> {React.string(name)} </td>
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
