type lobbyPayload = {
  action: string,
  gameno: string,
  tipe: string,
}

type startPayload = {
  action: string,
  gameno: string,
}

let chk = Js.String2.fromCharCode(10003)

@react.component
let make = (~game: Reducer.listGame, ~leadertoken, ~playerGame, ~send, ~class, ~textcolor) => {
  let (ready, setReady) = React.Uncurried.useState(_ => true)
  let (count, setCount) = React.useState(_ => 5)
  let (disabled1, setDisabled1) = React.useState(_ => false)
  let (disabled2, setDisabled2) = React.useState(_ => false)
  let (leader, setLeader) = React.useState(_ => false)

  let onClick1 = _ => {
    let pl = {
      action: "lobby",
      gameno: game.no,
      tipe: switch playerGame === game.no {
      | true => "leave"
      | false => "join"
      },
    }
    send(. Js.Json.stringifyAny(pl))
    switch playerGame === game.no {
    | true => setReady(._ => true)
    | false => ()
    }
  }

  let onClick2 = _ => {
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
    switch (playerGame == "", playerGame === game.no, Js.Array2.length(game.players) > 7) {
    | (false, false, _) | (true, _, true) => setDisabled1(_ => true)
    | (_, _, _) => setDisabled1(_ => false)
    }
    None
  }, (playerGame, game.no, game.players))

  React.useEffect3(() => {
    switch (playerGame == "", playerGame === game.no, Js.Array2.length(game.players) < 3) {
    | (false, false, _) | (true, _, _) | (_, _, true) => setDisabled2(_ => true)
    | (_, _, _) => setDisabled2(_ => false)
    }
    None
  }, (playerGame, game.no, game.players))

  React.useEffect3(() => {
    let id = if game.ready && game.no === playerGame {
      Js.Global.setInterval(() => {
        setCount(c => c - 1)
      }, 1000)
    } else {
      Js.Global.setInterval(() => (), 3600000)
    }
    Some(
      () => {
        setCount(_ => 5)
        Js.Global.clearInterval(id)
      },
    )
  }, (game.ready, game.no, playerGame))

  React.useEffect2(() => {
    switch Js.Array2.length(game.players) > 0 {
    | true =>
      switch game.players[0].name ++ game.players[0].connid == leadertoken {
      | true => setLeader(_ => true)
      | false => setLeader(_ => false)
      }
    | false => setLeader(_ => false)
    }
    None
  }, (game.players, leadertoken))

  React.useEffect4(() => {
    switch (playerGame === game.no && count === 0, leader) {
    | (true, true) => {
        RescriptReactRouter.push(`/game/${game.no}`)

        let pl = {
          action: "prep",
          gameno: game.no,
        }
        send(. Js.Json.stringifyAny(pl))
      }
    | (true, false) => RescriptReactRouter.push(`/game/${game.no}`)

    | (false, _) => ()
    }
    None
  }, (count, game.no, playerGame, leader))

  <li className={`mb-8 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8 ${class}`}>
    <p className={`text-center font-bold ${textcolor} font-anon text-sm col-span-2`}>
      {game.no->React.string}
    </p>
    <p className={`text-center font-bold ${textcolor} font-anon text-sm col-span-2`}>
      {"players"->React.string}
    </p>
    {game.players
    ->Js.Array2.map(p => {
      <p className={`text-center font-bold ${textcolor} font-anon`} key=p.connid>
        {p.name->React.string}
        {switch p.ready {
        | true =>
          <span className="text-yellow-400 text-2xl leading-3"> {React.string(chk)} </span>
        | false => React.null
        }}
      </p>
    })
    ->React.array}
    {switch (game.ready, playerGame !== game.no) {
    | (true, true) =>
      <p
        className={`absolute ${textcolor} text-3xl animate-pulse font-bold left-1/2 bottom-1/4 transform -translate-x-2/4`}>
        {React.string("Starting soon...")}
      </p>

    | (true, false) =>
      <p
        className={`absolute ${textcolor} text-3xl animate-ping font-bold left-1/2 bottom-1/4 transform -translate-x-2/4`}>
        {count->Js.Int.toString->React.string}
      </p>

    | (false, _) => React.null
    }}
    <button
      onClick=onClick1
      className="cursor-pointer text-base text-warm-gray-100 font-anon w-1/2 bottom-0 h-8 left-0 absolute bg-smoke-700 bg-opacity-70"
      disabled=disabled1>
      {switch playerGame === game.no {
      | true => React.string("leave")
      | false => React.string("join")
      }}
    </button>
    <button
      onClick=onClick2
      className="cursor-pointer text-base text-warm-gray-100 font-anon w-1/2 bottom-0 h-8 right-0 absolute bg-smoke-700 bg-opacity-70"
      disabled=disabled2>
      {switch ready {
      | true => "ready"->React.string
      | false => "not ready"->React.string
      }}
    </button>
  </li>
}
