type player = {
  name: string,
  connid: string,
  ready: bool,
  color: option<string>,
  score: string,
}

type answer = {
  playerid: string,
  answer: string,
}

type game = {
  no: string,
  ready: bool,
  starting: bool,
  loading: bool,
  playing: bool,
  players: array<player>,
  answers: array<answer>,
}

type state = {
  gamesList: Js.Nullable.t<array<game>>,
  game: game,
  currentWord: string,
  previousWord: string,
}

type action =
  | ListGames(Js.Nullable.t<array<game>>)
  | AddGame(game)
  | RemoveGame(game)
  | UpdateGame(game)
  | Word(string)

type return2 = {
  initialState: state,
  reducer: (state, action) => state,
}

let appState = () => {
  Js.log("appState")
  let reducer = (state, action) =>
    switch (Js.Nullable.toOption(state.gamesList), action) {
    | (None, ListGames(games)) => {...state, gamesList: games}

    | (None, _) | (Some(_), ListGames(_)) => state

    | (Some(gl), AddGame(game)) => {
        ...state,
        gamesList: Js.Nullable.return([game]->Js.Array2.concat(gl)),
      }

    | (Some(gl), RemoveGame(game)) => {
        ...state,
        gamesList: Js.Nullable.return(gl->Js.Array2.filter(gm => gm.no !== game.no)),
      }

    | (Some(gl), UpdateGame(game)) =>
      switch game.loading {
      | true => {
          ...state,
          game: game,
        }
      | false => {
          ...state,
          gamesList: Js.Nullable.return(
            gl->Js.Array2.map(gm =>
              switch gm.no === game.no {
              | true => game
              | false => gm
              }
            ),
          ),
        }
      }
    | (Some(gl), Word(word)) =>
      switch game.loading {
      | true => {
          ...state,
          game: game,
        }
      | false => {
          ...state,
          gamesList: Js.Nullable.return(
            gl->Js.Array2.map(gm =>
              switch gm.no === game.no {
              | true => game
              | false => gm
              }
            ),
          ),
        }
      }
    }

  {
    initialState: {
      gamesList: Js.Nullable.null,
      game: {
        no: "",
        ready: false,
        starting: false,
        loading: false,
        playing: false,
        players: [],
        answers: [],
      },
      currentWord: "",
      previousWord: "",
    },
    reducer: reducer,
  }
}
