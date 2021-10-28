type revokeTokenCallback = Js.Exn.t => unit

@send
external signOut: (Js.Nullable.t<Signup.usr>, Js.Nullable.t<revokeTokenCallback>) => unit =
  "signOut"

@react.component
let make = (~cognitoUser, ~setToken, ~send, ~playerGame, ~close) => {
  let onClick = e => {
    Js.log2("btn clck", e)

    let pl: Game.routePayload = {
      action: "lobby",
      gameno: switch Js.String2.length(playerGame) == 0 {
      | true => "dc"
      | false => playerGame
      },
      tipe: "disconnect",
    }
    Js.Json.stringifyAny(pl)->send

    close(1000, "user sign-out")
    cognitoUser->signOut(Js.Nullable.null)
    setToken(_ => None)
  }

  <button
    className="absolute top-5px right-5px bg-transparent cursor-pointer" onClick type_="button">
    <img className="block" src="../assets/signout.png" />
  </button>
}
