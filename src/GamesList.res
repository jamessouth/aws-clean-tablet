


@react.component
let make = (~action, ~games, ~ingame, ~leadertoken: string, ~send, ~user) => {
  <ul className="mx-auto mb-10 w-10/12">
    {games->Js.Array2.map(game => {
      // <Game action key=game.no game ingame leadertoken send user />
      <Game key=game.no game leadertoken/>
    })->React.array}
  </ul>
}
