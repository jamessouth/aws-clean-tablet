type revokeTokenCallback = Js.Exn.t => unit

@send
external signOut: (Js.Nullable.t<Signup.usr>, Js.Nullable.t<revokeTokenCallback>) => unit =
  "signOut"

@react.component
let make = (~cognitoUser, ~setToken) => {
  let onClick = e => {
    Js.log2("btn clck", e)
    cognitoUser->signOut(Js.Nullable.null)
    setToken(_ => None)
  }

  <button
    className="absolute top-5px right-5px bg-transparent cursor-pointer" onClick type_="button">
    <img className="block" src="../assets/signout.png" />
  </button>
}
