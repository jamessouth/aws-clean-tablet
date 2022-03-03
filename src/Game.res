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
let make = (~game: Reducer.listGame, ~leader, ~playerGame, ~send, ~class, ~readyColor) => {
  let (ready, setReady) = React.Uncurried.useState(_ => true)
  let (count, setCount) = React.useState(_ => 5)
  let (disabled1, setDisabled1) = React.useState(_ => false)
  let (disabled2, setDisabled2) = React.useState(_ => true)

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

      // if Js.Array2.length(game.players) > 7 {} else {}
    switch (playerGame === game.no, playerGame == "") {
    | (true, _) => {
      setDisabled1(_ => false)
      setDisabled2(_ => false)
    }
    | (false, false) => {
      setDisabled1(_ => true)
      setDisabled2(_ => true)
      }


    | (_, true) => 
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

  <li className={`<md:mb-16 grid grid-cols-2 grid-rows-6 relative text-xl bg-bottom bg-no-repeat text-center font-bold text-dark-800 font-anon pb-8 ${class} lg:(max-w-lg w-full)`}>
    <p className="absolute text-warm-gray-100 text-xs left-1/2 transform -translate-x-2/4 -top-3.5">
      {React.string(game.no)}
    </p>
    <p className="col-span-2"></p>

    {game.players
    ->Js.Array2.mapi((p, i) => {
      <p className={switch p.ready {
        | true => `underline decoration-[${readyColor}] decoration-4`
        | false => ""
        }} key=j`${p.name}$i`>
        {React.string(p.name)}
      </p>
    })
    ->React.array}
<p key="2">
        {React.string("z1")}
      </p><p key="23">
        {React.string("z11")}
      </p><p key="24">
        {React.string("z111")}
      </p><p key="25">
        {React.string("z1111")}
      </p><p key="26">
        {React.string("z11111")}
      </p><p key="27">
        {React.string("z111111")}
      </p><p key="28">
        {React.string("z1111111")}
      </p>




    {switch (game.ready, playerGame !== game.no) {
    | (true, true) =>
      <p
        className={"absolute text-3xl animate-pulse font-bold left-1/2 bottom-1/4 transform -translate-x-2/4"}>
        {React.string("Starting soon...")}
      </p>

    | (true, false) => switch count > 0 {
      | true =>
        <p
          className={"absolute text-3xl animate-ping font-bold left-1/2 bottom-1/4 transform -translate-x-2/4"}>
          {count->Js.Int.toString->React.string}
        </p>
      | false => React.null
      }

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
      | true => React.string("ready")
      | false => React.string("not ready")
      }}
    </button>
  </li>
}
