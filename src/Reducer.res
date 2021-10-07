type player = {
  name: string,
  connid: string,
  ready: bool,
  color: string,
  score: int,
}

type game = {
  leader: string,
  no: string,
  starting: bool,
  players: array<player>,
}

type state = {games: array<game>}

type action =
  | AddGame(game)
  | RemoveGame(game)
  | UpdateGame(game)

type return = {
  initialState: state,
  reducer: (state, action) => state,
}

// let mergeGame = (arr, ni) => {
//     let list = Js.Array2.copy(arr)
//     for i in 0 to Js.Array2.length(arr) - 1 {
//       if arr[i].no == ni.no {
//         switch ni.starting {
//         | true => list->Js.Array2.spliceInPlace(~pos=i, ~remove=1, ~add=[])
//         | false => {
//           list[i] = ni
//             list
//         }
//         }
//       }
//     }
//     Js.Array2.concat([ni], list)
// }

let appState = () => {
  Js.log("appState")
  let initialState = {
    games: [],
  }
  let reducer = (state, action) =>
    switch action {
    | AddGame(game) => {games: [game]->Js.Array2.concat(state.games)}
    | RemoveGame(game) => {games: state.games->Js.Array2.filter(gm => gm.no !== game.no)}
    | UpdateGame(game) => {
        games: state.games->Js.Array2.map(gm =>
          switch gm.no === game.no {
          | true => game
          | false => gm
          }
        ),
      }
    }

  {
    initialState: initialState,
    reducer: reducer,
  }
}
