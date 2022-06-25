type stat = {
  name: string,
  wins: int,
  points: int,
  games: int,
  winPct: float,
  ppg: float,
}

type listPlayer = {
  name: string,
  ready: bool,
}

type livePlayer = {
  playerid: string,
  name: string,
  color: string,
  score: string,//sent as int
  index: string,
  answer: string,
  hasAnswered: bool,
  pointsThisRound: string,
}

type listGame = {
  no: string,
  timerCxld: bool,
  players: array<listPlayer>,
}



// type liveGame = {
//   sk: string, 
//   currentWord: string,
//   previousWord: string,
//   players: array<livePlayer>,
//   showAnswers: bool,
//   winner: string,
// }

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  players: array<livePlayer>,
  sk: string,//game no
  word: string,
}

type action =
  | ListGames(Js.Nullable.t<array<listGame>>)
  | AddGame(listGame)
  | RemoveGame(listGame)
  | UpdateListGame(listGame)
  | UpdatePlayers(array<livePlayer>, string)
  | UpdateWord(string)
  | ResetPlayerState(state)

let init = clean => {
  gamesList: clean.gamesList,
  players: clean.players,
  sk: clean.sk,
  word: clean.word,
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
  | (Some(_), UpdatePlayers(players, sk)) => {
      ...state,
      players: players,
      sk: sk,
    }

  | (Some(_), UpdateWord(word)) => {
      ...state,
      word: word,
    }
  }
