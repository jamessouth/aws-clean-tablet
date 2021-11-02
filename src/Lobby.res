let startBtnStyles = " mb-8 w-1/2 bg-warm-gray-100 text-gray-700 h-8 text-lg font-anon cursor-pointer"

@react.component
let make = (~wsConnected, ~playerGame, ~leadertoken, ~games, ~send, ~wsError) => {
  let onClick = _ => {
    let pl: Game.routePayload = {
      action: "lobby",
      gameno: "new",
      tipe: "join"
    }
    send(. Js.Json.stringifyAny(pl))
  }



  switch wsError !== "" {
  | true => <p className="text-center text-warm-gray-100 font-anon text-lg"> {"not connected: connection error"->React.string} </p>
  | false =>
    switch (wsConnected, Js.Nullable.toOption(games)) {
    | (false, _) | (_, None) => <p className="text-center text-warm-gray-100 font-anon text-lg"> {"loading games..."->React.string} </p>
    | (true, Some(gs)) =>
      <div className="flex flex-col mt-8 items-center">
        <button
          className={switch playerGame === "" {
          | true => `visible${startBtnStyles}`
          | false => `invisible${startBtnStyles}`
          }}
          type_="button"
          onClick>
          {"start a new game"->React.string}
        </button>
        {switch gs->Js.Array2.length < 1 {
        | true => <p className="text-warm-gray-100 font-anon text-lg"> {"no games found. start a new one!"->React.string} </p>
        | false => <GamesList games=gs playerGame leadertoken send/>
        }}
      </div>
    }
  }
  // <div>{"lobby"->React.string}</div>
}
