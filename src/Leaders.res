type leaderPayload = {
  action: string,
  info: string,
}

@react.component
let make = (~send, ~rawLeaders) => {

  React.useEffect0(() => {
    let pl = {
      action: "leaders",
      info: "hello",
    }
    send(. Js.Json.stringifyAny(pl))
    None
  })


  <div className="">
    <h2 className="text-stone-100">{React.string("leaderboard coming soon!")}</h2>
  
  <ul
      className="">
      {rawLeaders
      ->Js.Array2.mapi((s, i) => {
        <li
          className="flex flex-row"
          key={j`${s.name}$i`}>
          <p
            className="">
            {React.string(s.name)}
          </p>
          <p
            className="">
            {React.string(s.name)}
          </p>
          <p
            className="">
            {React.string(s.name)}
          </p>
          <p
            className="">
            {React.string(s.name)}
          </p>
          
        </li>
      })
      ->React.array}
    </ul>
  
  
  
  
  </div>

  
}
