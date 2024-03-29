let normalClose = 1000
let players_max_threshold = 7
let btnStyles = "absolute top-1 bg-transparent cursor-pointer "

type apigwAction =
  | Answer
  | Lobby
  | Logging
  | Query

type lobbyCommand =
  | Custom({cv: string})
  | Join
  | Leaders
  | Leave
  | ListGames

type lobbyPayload = {
  act: apigwAction,
  gn?: string,
  cmd: lobbyCommand,
}

type payloadOutput = {
  action: string,
  gameno?: string,
  aW5mb3Jt?: string,
  command?: string,
  log?: string,
}

let apigwActionToString = a =>
  switch a {
  | Answer => "answer"
  | Lobby => "lobby"
  | Logging => "logging"
  | Query => "query"
  }

let lobbyCommandToString = lc =>
  switch lc {
  | Custom({cv}) => cv
  | Join => "join"
  | Leaders => "leaders"
  | Leave => "leave"
  | ListGames => "listGames"
  }

let payloadToObj = pl => {
  switch pl.act {
  | Answer =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: switch pl.gn {
      | None => ""
      | Some(v) => v
      },
      aW5mb3Jt: lobbyCommandToString(pl.cmd),
    })
  | Query =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      command: lobbyCommandToString(pl.cmd),
    })
  | Lobby =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: switch pl.gn {
      | None => ""
      | Some(v) => v
      },
      command: lobbyCommandToString(pl.cmd),
    })
  | Logging =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: switch pl.gn {
      | None => ""
      | Some(v) => v
      },
      log: lobbyCommandToString(pl.cmd),
    })
  }
}

module Game = {
  let btnStyle = " cursor-pointer text-base font-bold text-stone-100 font-anon w-1/2 bottom-0 h-8 absolute bg-stone-700 bg-opacity-70 filter disabled:(cursor-not-allowed contrast-25)"

  @react.component
  let make = (~game, ~inThisGame, ~inAGame, ~count, ~send, ~class, ~isOnlyGame) => {
    let liStyle = `<md:mb-16 grid grid-cols-2 grid-rows-6 relative text-xl bg-bottom bg-no-repeat h-200px text-center font-bold text-dark-800 font-anon pb-8 ${class} lg:(max-w-lg w-full)`
    let (disabledJoin, setDisabledJoin) = React.Uncurried.useState(_ => false)
    let {no, timerCxld, players}: Reducer.listGame = game

    let onClickJoin = _ => {
      send(.
        payloadToObj({
          act: Lobby,
          gn: no,
          cmd: switch inThisGame {
          | true => Leave
          | false => Join
          },
        }),
      )
    }

    React.useEffect3(() => {
      let size = Js.Array2.length(players)
      switch (inThisGame, inAGame) {
      | (true, _) => setDisabledJoin(._ => false) //in this game
      | (false, true) => setDisabledJoin(._ => true) //in another game
      | (_, false) =>
        setDisabledJoin(._ =>
          if size > players_max_threshold {
            true
          } else {
            false
          }
        ) //not in a game
      }
      None
    }, (inThisGame, inAGame, players))

    React.useEffect3(() => {
      switch inThisGame && count == "start" {
      | true => Route.push(Auth({subroute: Play({play: no})}))
      | false => ()
      }
      None
    }, (inThisGame, count, no))

    <li
      className={switch (inThisGame, isOnlyGame) {
      | (true, false) => "shadow-lg shadow-stone-100 " ++ liStyle
      | (false, true) | (true, true) | (false, false) => liStyle
      }}>
      <p className="absolute text-stone-100 text-xs left-1/2 transform -translate-x-2/4 -top-3.5">
        {React.string(no)}
      </p>
      <p className="col-span-2" />
      {players
      ->Js.Array2.mapi((p, i) => {
        <p key={j`${p.name}$i`}> {React.string(p.name)} </p>
      })
      ->React.array}
      {switch (timerCxld, inThisGame) {
      | (false, false) =>
        <p
          className="absolute text-2xl animate-pulse font-perm left-1/2 top-2/3 transform -translate-x-2/4 w-full">
          {React.string("Starting soon...")}
        </p>
      | (false, true) =>
        <p
          className="absolute text-4xl animate-ping1 font-perm left-1/2 top-1/4 transform -translate-x-2/4">
          {React.string(count)}
        </p>
      | (true, false) | (true, true) => React.null
      }}
      <Button onClick=onClickJoin disabled=disabledJoin className={"left-0" ++ btnStyle}>
        {switch inThisGame {
        | true => React.string("leave")
        | false => React.string("join")
        }}
      </Button>
    </li>
  }
}

@react.component
let make = (~wsConnected, ~playerListGame, ~games, ~send, ~close, ~count) => {
  Js.log("lobbyyyyyy")
  React.useEffect2(() => {
    Js.log3("lobby useeff", wsConnected, games)

    switch (wsConnected, Js.Nullable.toOption(games)) {
    | (true, None) =>
      send(.
        payloadToObj({
          act: Query,
          cmd: ListGames,
        }),
      )

    | (false, _) | (true, Some(_)) => ()
    }

    None
  }, (wsConnected, games))

  let leaderboard = _ => {
    Route.push(Auth({subroute: Leaderboard}))
  }

  let signOut = _ => {
    Js.log("sign out click")

    let pl = switch playerListGame == "" {
    | true => None
    | false =>
      payloadToObj({
        act: Lobby,
        gn: playerListGame,
        cmd: Leave,
      })
    }
    send(. pl)
    close(. normalClose, "user sign-out")
  }

  <>
    <Button onClick=leaderboard className={btnStyles ++ "left-1"}>
      <img className="block" src="../../assets/leader.png" />
    </Button>
    <Button onClick=signOut className={btnStyles ++ "right-1"}>
      <img className="block" src="../../assets/signout.png" />
    </Button>
    {switch Js.Nullable.toOption(games) {
    | None => {
        Web.body(Web.document)->Web.setClassName("bodchmob bodchtab bodchbig")
        <Loading label="games..." />
      }

    | Some(gs: Js.Array2.t<Reducer.listGame>) => {
        Web.body(Web.document)
        ->Web.classList
        ->Web.removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")

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
              ->Js.Array2.map(game => {
                let class = "game" ++ Js.String2.sliceToEnd(game.no, ~from=18)
                <Game
                  key=game.no
                  game
                  inThisGame={playerListGame == game.no}
                  inAGame={playerListGame != ""}
                  count
                  send
                  class
                  isOnlyGame={Js.Array2.length(gs) == 1}
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
