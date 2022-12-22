@react.component
let make: (
  ~players: array<Reducer.player>,
  ~playerLiveGame: string,
  ~showAnswers: bool,
  ~winner: string,
  ~isGameOver: bool,
  ~oldWord: string,
  ~word: string,
  ~playerColor: string,
  ~send: (. option<string>) => unit,
  ~playerName: string,
  ~resetConnState: (. unit) => unit,
) => React.element
