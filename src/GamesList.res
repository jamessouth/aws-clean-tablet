@react.component
let make = (~games, ~playerGame, ~leadertoken: string, ~send) => {
  <ul className="mb-10 w-11/12 flex <lg:(max-w-lg flex-col flex-wrap justify-around)">
    {games->Js.Array2.mapi((game, i) => {
      let class = switch mod(i, 6) {
      | 0 => "game0"
      | 1 => "game1"
      | 2 => "game2"
      | 3 => "game3"
      | 4 => "game4"
      | _ => "game5"
      }
      <Game key=game.no game playerGame leadertoken send class/>
    })->React.array}
  </ul>
}
