let startBtnStyles = " mx-auto mb-8 w-1/2 bg-smoke-100 text-gray-700"

@react.component
let make = (~wsConnected, ~ingame, ~leadertoken, ~games, ~send, ~wsError) => {
  let onClick = _ => {
    let pl: Game.routePayload = {
      action: "lobby",
      gameno: "new",
      tipe: "join"
    }
    Js.Json.stringifyAny(pl)->send
  }

  // let sendfunc = val => val->send

  switch wsError !== "" {
  | true => <p> {"not connected: connection error"->React.string} </p>
  | false =>
    switch (wsConnected, games->Js.Array2.length === 0) {
    | (false, _) | (_, true) => <p> {"loading games..."->React.string} </p>
    | (true, false) =>
      <div className="flex flex-col mt-8">
        <button
          className={switch ingame === "" {
          | true => `visible${startBtnStyles}`
          | false => `invisible${startBtnStyles}`
          }}
          type_="button"
          onClick>
          {"start a new game"->React.string}
        </button>
        {switch games->Js.Array2.length < 1 {
        | true => <p> {"no games found. start a new one!"->React.string} </p>
        | false => <GamesList games ingame leadertoken send/>
        }}
      </div>
    }
  }
  // <div>{"lobby"->React.string}</div>
}
