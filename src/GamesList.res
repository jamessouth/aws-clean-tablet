@react.component
let make = (~games, ~playerGame, ~leadertoken: string, ~send) => {
  <ul className="mx-auto mb-10 w-11/12">
    {games->Js.Array2.mapi((game, i) => {
      let class = switch mod(i, 2) {
      | 0 => "game0"
      | _ => "game1"
      }
      <Game key=game.no game playerGame leadertoken send class/>
    })->React.array}
  </ul>
}
