


@react.component
let make = (~_action, ~games, ~_ingame, ~leadertoken: string, ~_send, ~_user) => {
  <ul className="mx-auto mb-10 w-10/12">
    {games->Js.Array2.map(game => {
      // <Game action key=game.no game ingame leadertoken send user />
      <Game key=game.no game leadertoken/>
    })->React.array}
  </ul>
}
