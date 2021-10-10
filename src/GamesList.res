


@react.component
let make = (~games, ~ingame, ~leadertoken: string, ~send) => {
  <ul className="mx-auto mb-10 w-10/12">
    {games->Js.Array2.map(game => {
      <Game key=game.no game ingame leadertoken send />
    })->React.array}
  </ul>
}
