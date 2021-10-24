type routePayload = {
  action: string,
  gameno: string,
  tipe: string
}



let chk = Js.String2.fromCharCode(10003)

@react.component
let make = (~game: Reducer.game, ~leadertoken, ~playerGame, ~send) => {
  let (ready, setReady) = React.useState(_ => true)
  let (count, setCount) = React.useState(_ => 5)
  let (disabled1, setDisabled1) = React.useState(_ => false)
  let (disabled2, setDisabled2) = React.useState(_ => false)
  let (leader, setLeader) = React.useState(_ => "")
  //   let (prop1, setProp1) = React.useState(_ => false)
  //   let (inThisGame, setInThisGame) = React.useState(_ => false)

  let chkstyl = " text-2xl font-bold leading-3"

  let gameReady = switch game.leader {
  | Some(_) => true
  | None => false
  }

  let onClick1 = _e => {
    let pl = {
      action: "lobby",
      gameno: game.no,
      tipe: switch (playerGame == "", playerGame === game.no) {
      | (false, true) => "leave"
      | (_, _) => "join"
      }
    }
    Js.Json.stringifyAny(pl)->send
    switch (playerGame == "", playerGame === game.no) {
      | (false, true) => setReady(_ => true)
      | (_, _) => ()
      }
    
  }

  let onClick2 = _e => {
    let pl = {
      action: "lobby",
      gameno: game.no,
      tipe: switch ready {
      | true => "ready"
      | false => "unready"
      }
    }
    Js.Json.stringifyAny(pl)->send
    setReady(_ => !ready)
  }

  let ldr = switch game.leader {
  | Some(l) => l->Js.String2.split("_")
  | None => [""]
  }

  let leaderName = switch gameReady {
  | true => ldr[0]
  | false => ""
  }

  React.useEffect1(() => {
  switch game.leader {
  | None => setLeader(_ => "")
  | Some(l) => setLeader(_ => l)
  }
  None
}, [game.leader])

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
    let id = if gameReady && game.no === playerGame {
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
  }, (gameReady, game.no, playerGame))

  React.useEffect5(() => {
    switch (playerGame === game.no && count === 0, leader !== "" && leader === leadertoken) {
    | (true, true) => {
        
        RescriptReactRouter.push(`/game/${game.no}`)


        let pl = {
          action: "play",
          gameno: game.no,
          tipe: "start",
        }
        Js.Json.stringifyAny(pl)->send
      }
    | (true, false) => RescriptReactRouter.push(`/game/${game.no}`)
    



    | (false, _) => ()
    }
    None
  }, (count, game.no, playerGame, leadertoken, leader))

  <li className="mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8">
    <p className="text-center font-anon text-xs col-span-2"> {game.no->React.string} </p>
    <p className="text-center font-anon text-xs col-span-2"> {"players"->React.string} </p>
    {<>
      {game.players
      ->Js.Array2.map(p => {
        switch p.ready {
        | true =>
          <p className="text-center font-anon" key=p.connid>
            {p.name->React.string}
            <span
              className={switch leaderName === p.name {
              | true => `text-red-200${chkstyl}`
              | false => `text-green-200${chkstyl}`
              }}>
              {chk->React.string}
            </span>
          </p>
        | false => <p key=p.connid> {p.name->React.string} </p>
        }
      })
      ->React.array}
    </>}
    {switch (gameReady, playerGame !== game.no) {
    | (true, true) =>
        <p
          className="absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4">
          {"Starting soon..."->React.string}
        </p>

        
      | (true, false) =>
        <p
          className="absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4">
          {count->Js.Int.toString->React.string}
        </p>
      
    | (false, _) => React.null
    }}
    <button
        onClick=onClick1
      className="cursor-pointer font-anon w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70"
      disabled=disabled1>
      {switch (playerGame == "", playerGame === game.no) {
      | (false, true) => React.string("leave")
      | (_, _) => React.string("join")
      }}
    </button>
    <button
        onClick=onClick2
      className="cursor-pointer font-anon w-1/2 bottom-0 h-8 right-0 absolute pt-2 bg-smoke-700 bg-opacity-70"
      disabled=disabled2>
      {switch ready {
      | true => "ready"->React.string
      | false => "not ready"->React.string
      }}
    </button>
  </li>
}
