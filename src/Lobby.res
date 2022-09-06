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

module Game = {
  type lobbyPayload = {
    action: string,
    gameno: string,
    tipe: string,
  }

  @react.component
  let make = (~game: Reducer.listGame, ~inThisGame, ~inAGame, ~count, ~send, ~class, ~onlyGame) => {
    let liStyle = `<md:mb-16 grid grid-cols-2 grid-rows-6 relative text-xl bg-bottom bg-no-repeat h-200px text-center font-bold text-dark-800 font-anon pb-8 ${class} lg:(max-w-lg w-full)`
    let btnStyle = " cursor-pointer text-base font-bold text-stone-100 font-anon w-1/2 bottom-0 h-8 absolute bg-stone-700 bg-opacity-70 filter disabled:cursor-not-allowed disabled:contrast-[0.25]"
    let (ready, setReady) = React.Uncurried.useState(_ => true)
    let (disabledJoin, setDisabledJoin) = React.Uncurried.useState(_ => false)
    let (disabledReady, setDisabledReady) = React.Uncurried.useState(_ => true)

    let onClickJoin = _ => {
      let pl = {
        action: "lobby",
        gameno: game.no,
        tipe: switch inThisGame {
        | true => "leave"
        | false => "join"
        },
      }
      send(. Js.Json.stringifyAny(pl))
      switch inThisGame {
      | true => setReady(._ => true)
      | false => ()
      }
    }

    let onClickReady = _ => {
      let pl = {
        action: "lobby",
        gameno: game.no,
        tipe: switch ready {
        | true => "ready"
        | false => "unready"
        },
      }
      send(. Js.Json.stringifyAny(pl))
      setReady(._ => !ready)
    }

    React.useEffect3(() => {
      let size = Js.Array2.length(game.players)
      switch (inThisGame, inAGame) {
      | (true, _) => {
          //in this game
          setDisabledJoin(._ => false)
          if size < 3 {
            setDisabledReady(._ => true)
          } else {
            setDisabledReady(._ => false)
          }
        }

      | (false, true) => {
          //in another game
          setDisabledJoin(._ => true)
          setDisabledReady(._ => true)
        }

      | (_, false) => {
          //not in a game
          setDisabledReady(._ => true)
          if size > 7 {
            setDisabledJoin(._ => true)
          } else {
            setDisabledJoin(._ => false)
          }
        }
      }
      None
    }, (inThisGame, inAGame, game.players))

    React.useEffect3(() => {
      switch inThisGame && count == "start" {
      | true => Route.push(Play({play: game.no}))
      | false => ()
      }
      None
    }, (inThisGame, count, game.no))

    <li
      className={switch (inThisGame, onlyGame) {
      | (true, false) => "shadow-lg shadow-stone-100 " ++ liStyle
      | (false, true) | (true, true) | (false, false) => liStyle
      }}>
      <p className="absolute text-stone-100 text-xs left-1/2 transform -translate-x-2/4 -top-3.5">
        {React.string(game.no)}
      </p>
      <p className="col-span-2" />
      {game.players
      ->Js.Array2.mapi((p, i) => {
        <p
          className={switch p.ready {
          | true => `underline decoration-stone-800 decoration-3 italic`
          | false => ""
          }}
          key={j`${p.name}$i`}>
          {React.string(p.name)}
        </p>
      })
      ->React.array}
      {switch (game.timerCxld, inThisGame) {
      | (false, false) =>
        <p
          className="absolute text-2xl animate-pulse font-perm left-1/2 top-2/3 transform -translate-x-2/4 w-full">
          {React.string("Starting soon...")}
        </p>
      | (false, true) =>
        switch count != "" {
        | true =>
          <p
            className="absolute text-4xl animate-ping1 font-perm left-1/2 top-1/4 transform -translate-x-2/4">
            {React.string(count)}
          </p>
        | false => React.null
        }
      | (true, false) | (true, true) => React.null
      }}
      <Button
        textTrue="leave"
        textFalse="join"
        textProp=inThisGame
        onClick=onClickJoin
        disabled=disabledJoin
        className={"left-0" ++ btnStyle}
      />
      <Button
        textTrue="ready"
        textFalse="not ready"
        textProp=ready
        onClick=onClickReady
        disabled=disabledReady
        className={"right-0" ++ btnStyle}
      />
    </li>
  }
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
    setLeaderData(._ => [])

    let pl = {
      action: "leaders",
      info: "hello",
    }
    send(. Js.Json.stringifyAny(pl))
    Route.push(Leaderboard)
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
