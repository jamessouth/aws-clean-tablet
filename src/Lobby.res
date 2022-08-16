type propShape = {
  "close": (. int, string) => unit,
  "count": string,
  "games": Js.Nullable.t<Js.Array2.t<Reducer.listGame>>,
  "playerGame": string,
  "send": (. option<string>) => unit,
  "wsError": string,
  "setLeaderData": (. array<Reducer.stat> => array<Reducer.stat>) => unit,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"


type leaderPayload = {
  action: string,
  info: string,
}


@react.component
let make = (~playerGame, ~games, ~send, ~wsError, ~close, ~count, ~setLeaderData) => {
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


  let leaderboard = _ => {
    Js.log("leaderboard click")
    setLeaderData(._ => [])

    let pl = {
      action: "leaders",
      info: "hello",
    }
    send(. Js.Json.stringifyAny(pl))
    RescriptReactRouter.push("/leaderboard")
  }

  <>
    <Button
      textTrue=""
      textFalse=""
      onClick=leaderboard
      className="absolute top-1 left-1 bg-transparent cursor-pointer"
      img={<img className="block" src="../../assets/leader.png" />}
    />
    <Button
      textTrue=""
      textFalse=""
      onClick=signOut
      className="absolute top-1 right-1 bg-transparent cursor-pointer"
      img={<img className="block" src="../../assets/signout.png" />}
    />
    {switch wsError !== "" {
    | true =>
      <p className="text-center text-stone-100 font-anon text-lg">
        {React.string("not connected: connection error")}
      </p>
    | false =>
      switch Js.Nullable.toOption(games) {
      | None => <Loading label="games..." />
      | Some(gs) =>
        <div className="flex flex-col items-center">
          <div className="relative m-auto <newgmimg:w-11/12 w-max">
            <img
              srcSet="../../assets/ekko1x.webp, ../../assets/ekko2x.webp 2x"
              src="../../assets/ekko1x.webp"
              alt=""
              className="block <newgmimg:max-w-full"
              width="421"
              height="80"
            />
            {switch playerGame == "" {
            | true =>
              <Button
                textTrue="start a new game"
                textFalse="start a new game"
                onClick
                className="h-full right-0 top-0 w-1/2 bg-transparent text-stone-100 text-2xl font-flow cursor-pointer absolute border-l-2 border-gray-500/50"
              />
            | false => React.null
            }}
          </div>
          {switch Js.Array2.length(gs) < 1 {
          | true =>
            <p className="text-stone-100 font-anon text-lg mt-8">
              {React.string("no games found.")}
            </p>
          | false =>
            <ul
              className="m-12 newgmimg:mt-14 w-11/12 <md:(flex max-w-lg flex-col) md:(grid grid-cols-2 gap-8) lg:(gap-10 justify-items-center) xl:(grid-cols-3 gap-12 max-w-1688px)">
              {gs
              ->Js.Array2.map((game: Reducer.listGame) => {
                let class = "game" ++ Js.String2.sliceToEnd(game.no, ~from=18)
                <Game
                  key=game.no
                  game
                  inThisGame={playerGame == game.no}
                  inAGame={playerGame != ""}
                  count
                  send
                  class
                  onlyGame={Js.Array2.length(gs) == 1}
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
