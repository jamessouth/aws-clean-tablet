type lobbyPayload = {
  action: string,
  gameno: string,
  tipe: string,
}

type startPayload = {
  action: string,
  gameno: string,
}

@react.component
let make = (
  ~game: Reducer.listGame,
  ~leader,
  ~inThisGame,
  ~inAGame,
  ~send,
  ~class,
  ~readyColor,
) => {
  let btnStyle = " cursor-pointer text-base font-bold text-warm-gray-100 font-anon w-1/2 bottom-0 h-8 absolute bg-smoke-700 bg-opacity-70"
  let (ready, setReady) = React.Uncurried.useState(_ => true)
  let (count, setCount) = React.Uncurried.useState(_ => 5)
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

  React.useEffect2(() => {
    let id = if game.ready && inThisGame {
      Js.Global.setInterval(() => {
        setCount(.c => c - 1)
      }, 1000)
    } else {
      Js.Global.setInterval(() => (), 3600000)
    }
    Some(
      () => {
        setCount(._ => 5)
        Js.Global.clearInterval(id)
      },
    )
  }, (inThisGame, game.ready))

  React.useEffect4(() => {
    switch (inThisGame && count == 0, leader) {
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
  }, (inThisGame, count, game.no, leader))

  <li
    className={`<md:mb-16 grid grid-cols-2 grid-rows-6 relative text-xl bg-bottom bg-no-repeat h-200px text-center font-bold text-dark-800 font-anon pb-8 ${class} lg:(max-w-lg w-full)`}>
    <p className="absolute text-warm-gray-100 text-xs left-1/2 transform -translate-x-2/4 -top-3.5">
      {React.string(game.no)}
    </p>
    <p className="col-span-2" />
    {game.players
    ->Js.Array2.mapi((p, i) => {
      <p
        className={switch p.ready {
        | true => `underline decoration-[${readyColor}] decoration-4`
        | false => ""
        }}
        key={j`${p.name}$i`}>
        {React.string(p.name)}
      </p>
    })
    ->React.array}
    {
      // <p key="2"> {React.string("z1")} </p>
      // <p key="23"> {React.string("z11")} </p>
      // <p key="24"> {React.string("z111")} </p>
      // <p key="25"> {React.string("z1111")} </p>
      // <p key="26"> {React.string("z11111")} </p>
      // <p key="27"> {React.string("z111111")} </p>
      // <p key="28"> {React.string("z1111111")} </p>
      switch (game.ready, inThisGame) {
      | (true, false) =>
        <p
          className="absolute text-3xl animate-pulse font-bold left-1/2 top-1/4 transform -translate-x-2/4">
          {React.string("Starting...")}
        </p>

      | (true, true) =>
        switch count > 0 {
        | true =>
          <p
            className="absolute text-3xl animate-ping font-bold left-1/2 top-1/4 transform -translate-x-2/4">
            {React.string(Js.Int.toString(count))}
          </p>
        | false => React.null
        }

      | (false, _) => React.null
      }
    }
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
