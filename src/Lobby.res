
let startBtnStyles = " mx-auto mb-8 w-1/2 bg-smoke-100 text-gray-700"

@react.component
let make = () => {

    let onClick = _ => {
        let pl: Game.sendPayload = {
            action: "lobby",
            gameno: "new",
            type_: "join",
            value: false
        }
        pl->send
    }

switch wsError {
| true => <p>"not connected: connection error"->React.string</p>
| false => switch connectedWS, games->Js.Nullable.toOption {
| false, _ | _, None => <p>"loading games..."->React.string</p>
| true, Some(games) => <div
className="flex flex-col mt-8"
><button
className=switch ingame {
| true => `invisible${startBtnStyles}`
| false => `visible${startBtnStyles}`
}
type_="button"
onClick
>"start a new game"->React.string</button>
{
    switch games->Js.Array2.length < 1 {
    | true => <p>"no games found. start a new one!"->React.string</p>
    | false => <GamesList action games ingame leadertoken send=(val) => val->send user></GamesList>
    }
}
</div>
}
}


}