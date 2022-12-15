let players_min_threshold = 3
let players_max_threshold = 7

type apigwAction =
  | Answer
  | End
  | Leaders
  | Lobby
  | LobbyNonApigw
  | Logging

type lobbyGameno =
  | Gameno({no: string})
  | Newgame

type lobbyCommand =
  | Custom({cv: string})
  | Join
  | Leave

type lobbyPayload = {
  act: apigwAction,
  gn: lobbyGameno,
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
  | End => "end"
  | Leaders => "leaders"
  | Lobby => "lobby"
  | LobbyNonApigw => "lobbyNonApigw"
  | Logging => "logging"
  }

let lobbyGamenoToString = gn =>
  switch gn {
  | Gameno({no}) => no
  | Newgame => "newgame"
  }

let lobbyCommandToString = lc =>
  switch lc {
  | Custom({cv}) => cv
  | Join => "join"
  | Leave => "leave"
  }

let payloadToObj = pl => {
  switch pl.act {
  | Answer =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: lobbyGamenoToString(pl.gn),
      aW5mb3Jt: lobbyCommandToString(pl.cmd),
    })
  | End | Leaders =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
    })
  | Lobby | LobbyNonApigw =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: lobbyGamenoToString(pl.gn),
      command: lobbyCommandToString(pl.cmd),
    })
  | Logging =>
    Js.Json.stringifyAny({
      action: apigwActionToString(pl.act),
      gameno: lobbyGamenoToString(pl.gn),
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
      let pl = switch inThisGame {
      | true =>
        payloadToObj({
          act: Lobby,
          gn: Gameno({no: no}),
          cmd: Leave,
        })
      | false =>
        payloadToObj({
          act: LobbyNonApigw,
          gn: Gameno({no: no}),
          cmd: Join,
        })
      }
      send(. pl)
    }

    React.useEffect3(() => {
      let size = Js.Array2.length(players)
      switch (inThisGame, inAGame) {
      | (true, _) => {
          //in this game
          setDisabledJoin(._ => false)
        }

      | (false, true) => {
          //in another game
          setDisabledJoin(._ => true)
        }

      | (_, false) => {
          //not in a game
          setDisabledJoin(._ =>
            if size > players_max_threshold {
              true
            } else {
              false
            }
          )
        }
      }
      None
    }, (inThisGame, inAGame, players))

    React.useEffect3(() => {
      switch inThisGame && count == "start" {
      | true => Route.push(Play({play: no}))
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
        <p
          key={j`${p.name}$i`}>
          {React.string(p.name)}
        </p>
      })
      ->React.array}
      {switch (timerCxld, inThisGame) {
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
let make = (~playerGame, ~games, ~send, ~close, ~count, ~setLeaderData) => {
  let onClick = _ =>
    send(.
      payloadToObj({
        act: LobbyNonApigw,
        gn: Newgame,
        cmd: Join,
      }),
    )

  let signOut = _ => {
    Js.log("sign out click")

    let pl = switch playerGame == "" {
    | true => None
    | false =>
      payloadToObj({
        act: Lobby,
        gn: Gameno({no: playerGame}),
        cmd: Leave,
      })
    }
    send(. pl)
    close(. 1000, "user sign-out")
  }

  let leaderboard = _ => {
    setLeaderData(._ => [])
    send(.
      payloadToObj({
        act: Leaders,
        gn: Newgame, //placeholder
        cmd: Join, //placeholder
      }),
    )
    Route.push(Leaderboard)
  }

  <>
    <Button onClick=leaderboard className="absolute top-1 left-1 bg-transparent cursor-pointer">
      <img className="block" src="../../assets/leader.png" />
    </Button>
    <Button onClick=signOut className="absolute top-1 right-1 bg-transparent cursor-pointer">
      <img className="block" src="../../assets/signout.png" />
    </Button>
    {switch Js.Nullable.toOption(games) {
    | None => <Loading label="games..." />
    | Some(gs: Js.Array2.t<Reducer.listGame>) =>
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
              onClick
              className="h-full right-0 top-0 w-1/2 bg-transparent text-stone-100 text-2xl font-flow cursor-pointer absolute border-l-2 border-gray-500/50">
              {React.string("start a new game")}
            </Button>
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
            ->Js.Array2.map(game => {
              let class = "game" ++ Js.String2.sliceToEnd(game.no, ~from=18)
              <Game
                key=game.no
                game
                inThisGame={playerGame == game.no}
                inAGame={playerGame != ""}
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
    }}
  </>
}
