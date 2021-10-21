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

type state = {gamesList: Js.Nullable.t<array<game>>}

type action =
  | ListGames(Js.Nullable.t<array<game>>)
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
    gamesList: Js.Nullable.null,
  }
  let reducer = ({gamesList}, action) =>
    switch (Js.Nullable.toOption(gamesList), action) {
    | (None, ListGames(games)) => {gamesList: games}

    | (None, _) => {gamesList: gamesList}

    | (Some(gl), AddGame(game)) => {gamesList: Js.Nullable.return([game]->Js.Array2.concat(gl))}

    | (Some(gl), RemoveGame(game)) => {
        gamesList: Js.Nullable.return(gl->Js.Array2.filter(gm => gm.no !== game.no)),
      }

    | (Some(gl), UpdateGame(game)) => {
        gamesList: Js.Nullable.return(
          gl->Js.Array2.map(gm =>
            switch gm.no === game.no {
            | true => game
            | false => gm
            }
          ),
        ),
      }

    | (Some(_), ListGames(_)) => {gamesList: gamesList}
    }

  {
    initialState: initialState,
    reducer: reducer,
  }
}