type stat = {
  name: string,
  wins: int,
  points: int,
  games: int,
  winPct: float,
  ppg: float,
}

type player = {
  name: string,
  connid: string,
  color: string,
  score: string, //sent as int
  answer: string,
  hasAnswered: bool,
  pointsThisRound: string,
}

type listGame = {
  no: string,
  timerCxld: bool,
  players: array<player>,
}

type state = {
  gamesList: Js.Nullable.t<array<listGame>>,
  players: array<player>,
  playerLiveGame: string,
  oldWord: string,
  word: string,
  showAnswers: bool,
  winner: string,
  playerColor: string,
  playerConnID: string,
  playerListGame: string,
  playerName: string,
}

type action =
  | ListGames(Js.Nullable.t<array<listGame>>, string, string)
  | AddGame(listGame)
  | RemoveGame(listGame)
  | UpdateListGame(listGame)
  | UpdatePlayers(array<player>, string, bool, string)
  | UpdateWord(string)
  | ResetPlayerState(state)

let init = clean => {
  gamesList: clean.gamesList,
  players: clean.players,
  playerLiveGame: clean.playerLiveGame,
  oldWord: clean.oldWord,
  word: clean.word,
  showAnswers: clean.showAnswers,
  winner: clean.winner,
  playerColor: clean.playerColor,
  playerConnID: clean.playerConnID,
  playerListGame: clean.playerListGame,
  playerName: clean.playerName,
}

let reducer = (state, action) => {
  let predFunc = p => p.name == state.playerName && p.connid == state.playerConnID

  switch (Js.Nullable.toOption(state.gamesList), action) {
  | (_, ResetPlayerState(st)) => init(st)

  | (None, ListGames(games, name, connid)) => {
      ...state,
      gamesList: games,
      playerName: name,
      playerConnID: connid,
    }
  | (None, _) | (Some(_), ListGames(_)) => state

  | (Some(gl), AddGame(game)) => {
      ...state,
      gamesList: Js.Nullable.return([game]->Js.Array2.concat(gl)),
      playerListGame: switch game.players->Js.Array2.find(predFunc) {
      | Some(_) => game.no
      | None => ""
      },
    }

  | (Some(gl), RemoveGame(game)) => {
      ...state,
      gamesList: Js.Nullable.return(gl->Js.Array2.filter(gm => gm.no !== game.no)),
      playerListGame: switch game.no == state.playerListGame {
      | true => ""
      | false => game.no
      },
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
      playerListGame: switch game.players->Js.Array2.find(predFunc) {
      | Some(_) => game.no
      | None => ""
      },
    }

  | (Some(_), UpdatePlayers(players, playerLiveGame, showAnswers, winner)) => {
      let pc = switch state.playerColor == "transparent" {
      | true =>
        switch players->Js.Array2.find(predFunc) {
        | Some(p) => p.color
        | None => "black"
        }
      | false => state.playerColor
      }

      let (ow, nw) = switch showAnswers {
      | true => (state.word, "")
      | false => ("", state.word)
      }

      {
        ...state,
        players,
        playerLiveGame,
        showAnswers,
        winner,
        oldWord: ow,
        word: nw,
        playerColor: pc,
      }
    }

  | (Some(_), UpdateWord(word)) => {
      ...state,
      word,
    }
  }
}
