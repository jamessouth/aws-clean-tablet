let className = "h-full right-0 top-0 w-1/2 bg-transparent text-warm-gray-100 text-2xl font-flow cursor-pointer absolute border-l-2 border-gray-500/50"

@react.component
let make = (~playerGame, ~leader, ~games, ~send, ~wsError, ~close) => {
  let onClick = _ => {
    let pl: Game.lobbyPayload = {
      action: "lobby",
      gameno: "new",
      tipe: "join",
    }
    send(. Js.Json.stringifyAny(pl))
  }

  let signOut = _ => {
    Js.log("sign out click")

    let pl: Game.lobbyPayload = {
      action: "lobby",
      gameno: switch playerGame == "" {
      | true => "dc"
      | false => playerGame
      },
      tipe: "disconnect",
    }
    send(. Js.Json.stringifyAny(pl))
    close(. 1000, "user sign-out")
  }

  <>
    <Button
      textTrue=""
      textFalse=""
      onClick=signOut
      className="absolute top-1 right-1 bg-transparent cursor-pointer"
      img={<img className="block" src="../assets/signout.png" />}
    />
    {switch wsError !== "" {
    | true =>
      <p className="text-center text-warm-gray-100 font-anon text-lg">
        {React.string("not connected: connection error")}
      </p>
    | false =>
      switch Js.Nullable.toOption(games) {
      | None =>
        <p className="text-center text-warm-gray-100 font-anon text-lg">
          {React.string("loading games...")}
        </p>
      | Some(gs) =>
        <div className="flex flex-col items-center">
          <div className="relative m-auto <newgmimg:w-11/12 w-max">
            <img
              srcSet="../../assets/ekko2x.webp 2x"
              src="../../assets/ekko1x.webp"
              alt=""
              className="block <newgmimg:max-w-full"
              width="421"
              height="80"
            />
            {switch playerGame == "" {
            | true =>
              <Button textTrue="start a new game" textFalse="start a new game" onClick className />
            | false => React.null
            }}
          </div>
          {switch Js.Array2.length(gs) < 1 {
          | true =>
            <p className="text-warm-gray-100 font-anon text-lg">
              {React.string("no games found. start a new one!")}
            </p>
          | false =>
            <ul
              className="m-12 w-11/12 <md:(flex max-w-lg flex-col) md:(grid grid-cols-2 gap-8) lg:(gap-10 justify-items-center) xl:(grid-cols-3 gap-12 max-w-1688px)">
              {gs
              ->Js.Array2.mapi((game: Reducer.listGame, i) => {
                let (class, readyColor) = switch mod(i, 6) {
                | 0 => ("game0", "#cc9e48")
                | 1 => ("game1", "#213e10")
                | 2 => ("game2", "#4e3942")
                | 3 => ("game3", "#4E4A2F")
                | 4 => ("game4", "#5f4500")
                | _ => ("game5", "#8d4f36")
                }
                <Game
                  key=game.no
                  game
                  leader
                  inThisGame={playerGame == game.no}
                  inAGame={playerGame != ""}
                  send
                  class
                  readyColor
                />
              })
              ->React.array}
            </ul>
          }}
        </div>
      }
    }}
  </>
}
