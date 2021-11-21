type answer = {
  playerid: string,
  answer: string,
}

type listPlayer = {
  name: string,
  connid: string,
  ready: bool,
}

type livePlayer = {
  name: string,
  connid: string,
  color: string,
  score: int,
  answer: answer,
}

type listGame = {
  no: string,
  ready: bool,
  players: array<listPlayer>,
}

type liveGame = {
  no: string,
  showAnswers: bool,
  currentWord: string,
  previousWord: string,
  players: array<livePlayer>,
}

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  game: liveGame,
}

type action =
  | ListGames(Js.Nullable.t<array<listGame>>)
  | AddGame(listGame)
  | RemoveGame(listGame)
  | UpdateListGame(listGame)
  | UpdateLiveGame(liveGame)

type return = {
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

    | (Some(gl), UpdateListGame(game)) => {
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

    | (Some(_), UpdateLiveGame(game)) => {
        ...state,
        game: game,
      }
    }

  {
    initialState: {
      gamesList: Js.Nullable.null,
      game: {
        no: "",
        players: [],
        showAnswers: false,
        currentWord: "",
        previousWord: "",
      },
    },
    reducer: reducer,
  }
}
