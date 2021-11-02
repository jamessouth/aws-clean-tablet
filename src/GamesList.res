@react.component
let make = (~games, ~playerGame, ~leadertoken: string, ~send) => {
  <ul className="mx-auto mb-10 w-11/12">
    {games->Js.Array2.mapi((game, i) => {
      let (class, textcolor) = switch mod(i, 2) {
      | 0 => ("game0", "text-dark-800")
      | _ => ("game1", "text-warm-gray-100")
      }
      <Game key=game.no game playerGame leadertoken send class textcolor/>
    })->React.array}
  </ul>
}
