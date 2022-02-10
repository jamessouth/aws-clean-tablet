let startBtnStyles = " mb-14 w-1/2 bg-warm-gray-100 text-gray-700 h-8 text-lg font-anon cursor-pointer"

@react.component
let make = (~wsConnected, ~playerGame, ~leadertoken, ~games, ~send, ~wsError, ~signOut) => {
  let onClick = _ => {
    let pl: Game.lobbyPayload = {
      action: "lobby",
      gameno: "new",
      tipe: "join",
    }
    send(. Js.Json.stringifyAny(pl))
  }

  <>
    {signOut}
    {switch wsError !== "" {
    | true =>
      <p className="text-center text-warm-gray-100 font-anon text-lg">
        {"not connected: connection error"->React.string}
      </p>
    | false =>
      switch (wsConnected, Js.Nullable.toOption(games)) {
      | (false, _) | (_, None) =>
        <p className="text-center text-warm-gray-100 font-anon text-lg">
          {React.string("loading games...")}
        </p>
      | (true, Some(gs)) =>
        <div className="flex flex-col items-center">
          <div className="relative m-auto <newgmimg:w-11/12 w-max">
            <img src="../../assets/ekko.png" alt="" className="block <newgmimg:max-w-full" width="421" height="80" />
            {switch playerGame === "" {
            | true =>
              <button
                className="h-full right-0 top-0 w-1/2 bg-transparent text-warm-gray-100 text-2xl font-flow cursor-pointer absolute border-l-2 border-gray-500/50"
                type_="button"
                onClick>
                {React.string("start a new game")}
              </button>
            | false => React.null
            }}
          </div>
          {switch gs->Js.Array2.length < 1 {
          | true =>
            <p className="text-warm-gray-100 font-anon text-lg">
              {React.string("no games found. start a new one!")}
            </p>
          | false => <GamesList games=gs playerGame leadertoken send />
          }}
        </div>
      }
    }}
  </>
}
