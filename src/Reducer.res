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
  playing: bool,
  players: array<player>,
}

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  game: liveGame,
  currentWord: string,
  previousWord: string,
  showAnswers: bool,
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

    | (Some(gl), UpdateLiveGame(game)) => switch value {
    | pattern1 => expression
    | pattern2 => expression
    }
    
    
    
    
    
    
    {
        ...state,
        previousWord: currentWord,
        game: game,
        showAnswers: true,
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
        playing: false,
        players: [],
      },
      currentWord: "",
      previousWord: "",
      showAnswers: false,
    },
    reducer: reducer,
  }
}
