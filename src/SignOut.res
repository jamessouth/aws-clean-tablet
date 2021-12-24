type revokeTokenCallback = Js.Exn.t => unit

@send
external signOut: (Js.Nullable.t<Signup.usr>, Js.Nullable.t<revokeTokenCallback>) => unit =
  "signOut"

@react.component
let make = (~cognitoUser, ~setToken, ~send, ~playerGame, ~close, ~setCognitoUser, ~setPlayerName, ~setWs, ~setPlayerGame, ~setConnID) => {
  let onClick = e => {
    Js.log2("btn clck", e)

    let pl: Game.lobbyPayload = {
      action: "lobby",
      gameno: switch Js.String2.length(playerGame) == 0 {
      | true => "dc"
      | false => playerGame
      },
      tipe: "disconnect",
    }
    send(. Js.Json.stringifyAny(pl))

    close(. 1000, "user sign-out")
    cognitoUser->signOut(Js.Nullable.null)
    setCognitoUser(._ => Js.Nullable.null)
    setPlayerName(._ => "")
    setPlayerGame(._ => "")
    setConnID(._ => "")
    setToken(._ => None)
    setWs(._ => Js.Nullable.null)
  }

  <button
    className="absolute top-5px right-5px bg-transparent cursor-pointer" onClick type_="button">
    <img className="block" src="../assets/signout.png" />
  </button>
}
