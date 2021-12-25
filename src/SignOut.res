@react.component
let make = (~send, ~playerGame, ~close) => {
  let onClick = e => {
    Js.log2("btn clck", e)

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

  <button
    className="absolute top-5px right-5px bg-transparent cursor-pointer" onClick type_="button">
    <img className="block" src="../assets/signout.png" />
  </button>
}
