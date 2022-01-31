@react.component
let make = (~games, ~playerGame, ~leadertoken: string, ~send) => {
  <ul className="mb-10 w-11/12 <md:(flex max-w-lg flex-col) md:(grid grid-cols-2 gap-8) lg:(gap-10 justify-items-center) xl:(grid-cols-3 gap-12 max-w-1688px)">
    {games->Js.Array2.mapi((game, i) => {
      let (class, readyColor) = switch mod(i, 6) {
      | 0 => ("game0", "#cc9e48")
      | 1 => ("game1", "#213e10")
      | 2 => ("game2", "#4e3942")
      | 3 => ("game3", "#4E4A2F")
      | 4 => ("game4", "#5f4500")
      | _ => ("game5", "#8d4f36")
      }
      <Game key=game.no game playerGame leadertoken send class readyColor/>
    })->React.array}
  </ul>
}
