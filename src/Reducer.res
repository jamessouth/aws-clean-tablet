type player = {
  name: string,
  connid: string,
  ready: bool,
  color: option<string>,
  score: string,
  answer: answer,
}

type answer = {
  playerid: string,
  answer: string,
}

type listGame = {
  no: string,
  ready: bool,
  players: array<player>,
}

type liveGame = {
  no: string,
  players: array<player>,
}

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  game: liveGame,
  showAnswers: bool,
  currentWord: string,
  previousWord: string,
}

type action =
  | ListGames(Js.Nullable.t<array<listGame>>)
  | AddGame(listGame)
  | RemoveGame(listGame)
  | UpdateListGame(listGame)
  | UpdateLiveGame(liveGame)
  | Word(string)

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

    | (Some(gl), UpdateLiveGame(game)) =>
      switch game.showAnswers {
      | true => {
          ...state,
          previousWord: currentWord,
          game: game,
          showAnswers: true,
        }
      | false => state
      }

    | (Some(gl), Word(word)) => {
        ...state,
        currentWord: word,
        showAnswers: false,
      }
    }

  {
    initialState: {
      gamesList: Js.Nullable.null,
      game: {
        no: "",
        players: [],
      },
      showAnswers: false,
      currentWord: "",
      previousWord: "",
    },
    reducer: reducer,
  }
}
