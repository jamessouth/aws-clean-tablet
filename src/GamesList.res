@react.component
let make = (action, games, ingame, leadertoken, send, user) => {
  <ul className="mx-auto mb-10 w-10/12">
    {games->Js.Array2.map(game => {
      <Game action key=game.no game ingame leadertoken send user />
    })}
  </ul>
}
