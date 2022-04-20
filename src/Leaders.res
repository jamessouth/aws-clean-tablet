type leaderPayload = {
  action: string,
  info: string,
}

@react.component
let make = (~send) => {
  let pl = {
    action: "leaders",
    info: "hello",
  }
  send(. Js.Json.stringifyAny(pl))

  <div className="text-stone-100"> {React.string("leaderboard coming soon!")} </div>
}
