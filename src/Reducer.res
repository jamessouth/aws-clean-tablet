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
  connid: string,
}

type livePlayer = {
  // playerid: string,
  name: string,
  color: string,
  score: string, //sent as int
  answer: string,
  hasAnswered: bool,
  pointsThisRound: string,
}

type listGame = {
  no: string,
  timerCxld: bool,
  players: array<listPlayer>,
}

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  players: array<livePlayer>,
  sk: string, //game no
  oldWord: string,
  word: string,
  showAnswers: bool,
  winner: string,
  playerColor: string,
  playerGame: string,
  playerName: string,
  playerConnID: string,
}

type action =
  | ListGames(Js.Nullable.t<array<listGame>>, string)
  | AddGame(listGame)
  | RemoveGame(listGame)
  | UpdateListGame(listGame)
  | UpdatePlayers(array<livePlayer>, string, bool, string)
  | UpdateWord(string)
  | UpdatePlayerColor(string)
  | ResetPlayerState(state)

let init = clean => {
  gamesList: clean.gamesList,
  players: clean.players,
  sk: clean.sk,
  oldWord: clean.oldWord,
  word: clean.word,
  showAnswers: clean.showAnswers,
  winner: clean.winner,
  playerColor: clean.playerColor,
  playerName: clean.playerName,
  playerGame: clean.playerGame,
  playerConnID: clean.playerConnID,
}

let findName = (game: listGame, name: string, connid: string) => {
  switch game.players->Js.Array2.find(p => p.name == name && p.connid == connid) {
  | Some(_) => game.no
  | None => ""
  }
}

let reducer = (state, action) =>
  switch (Js.Nullable.toOption(state.gamesList), action) {
  | (_, ResetPlayerState(st)) => init(st)
  
  | (None, ListGames(games, name, connid)) => {...state, gamesList: games, playerName: name, playerConnID: connid}
  | (None, _) | (Some(_), ListGames(_)) => state




  | (Some(gl), AddGame(game)) => {
      ...state,
      gamesList: Js.Nullable.return([game]->Js.Array2.concat(gl)),
      playerGame: game->findName(state.playerName, state.playerConnID),
    }


  | (Some(gl), RemoveGame(game)) => {
      ...state,
      gamesList: Js.Nullable.return(gl->Js.Array2.filter(gm => gm.no !== game.no)),
      playerGame: switch game.no == state.playerGame {
      | true => ""
      | false => game.no
      }
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
      playerGame: game->findName(state.playerName, state.playerConnID),
    }



  | (Some(_), UpdatePlayers(players, sk, showAnswers, winner)) => {
      let ow = switch showAnswers {
      | true => state.word
      | false => ""
      }
      let nw = switch showAnswers {
      | true => ""
      | false => state.word
      }
      {
        ...state,
        players: players,
        sk: sk,
        showAnswers: showAnswers,
        winner: winner,
        oldWord: ow,
        word: nw,
      }
    }

  | (Some(_), UpdateWord(word)) => {
      ...state,
      word: word,
    }
  }
