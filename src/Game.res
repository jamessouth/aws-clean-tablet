let chk = Js.String2.fromCharCode(10003)

@react.component
let make = (game, ingame, leadertoken, send) => {
    let gameReady = switch game.leader->Js.Nullable.toOption {
    | Some(_) => true
    | None => false
    }

    let leaderName = gameReady && game.leader->Js.String2.split("_")[0]

    





}