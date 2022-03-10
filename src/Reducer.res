type listPlayer = {
  name: string,
  ready: bool,
}

type livePlayer = {
  name: string,
  color: string,
  score: string,
  answer: string,
  hasAnswered: bool,
  index: int,
}

type listGame = {
  no: string,
  ready: bool,
  players: array<listPlayer>,
}

type liveGame = {
  sk: string, //game no
  currentWord: string,
  previousWord: string,
  players: array<livePlayer>,
  hiScore: int,
  gameTied: bool,
  showAnswers: bool,
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
  | ResetPlayerState(state)

let init = clean => {
  gamesList: clean.gamesList,
  game: clean.game,
}

let reducer = (state, action) =>
  switch (Js.Nullable.toOption(state.gamesList), action) {
  | (_, ResetPlayerState(st)) => init(st)
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
          switch gm.no == game.no {
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
