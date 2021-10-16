type player = {
  name: string,
  connid: string,
  ready: bool,
  color: option<string>,
  score: string,
}

type game = {
  leader: option<string>,
  no: string,
  starting: bool,
  players: array<player>,
}

type state = {games: Js.Nullable.t<array<game>>}

type action =
  | AddGame(game)
  | RemoveGame(game)
  | UpdateGame(game)

type return2 = {
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
    games: Js.Nullable.null,
  }
  let reducer = ({games}, action) =>
    switch (Js.Nullable.toOption(games), action) {
    | (None, AddGame(game)) => {games: Js.Nullable.return([game])}
    | (None, RemoveGame(_)) => {games: games}
    | (None, UpdateGame(_)) => {games: games}
    | (Some(gs), AddGame(game)) => {games: Js.Nullable.return([game]->Js.Array2.concat(gs))}
    | (Some(gs), RemoveGame(game)) => {games: Js.Nullable.return(gs->Js.Array2.filter(gm => gm.no !== game.no))}
    | (Some(gs), UpdateGame(game)) => {
        games: Js.Nullable.return(gs->Js.Array2.map(gm =>
          switch gm.no === game.no {
          | true => game
          | false => gm
          }
        )),
      }
    }

  {
    initialState: initialState,
    reducer: reducer,
  }
}
