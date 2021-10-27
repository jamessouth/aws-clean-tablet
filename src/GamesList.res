@react.component
let make = (~games, ~playerGame, ~leadertoken: string, ~send) => {
  <ul className="mx-auto mb-10 w-11/12">
    {games->Js.Array2.mapi((game, i) => {
      switch value {
      | pattern1 => expression
      | pattern2 => expression
      }
      <Game key=game.no game playerGame leadertoken send />
    })->React.array}
  </ul>
}
