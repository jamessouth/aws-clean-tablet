@react.component
let make: (
  ~players: array<Reducer.player>,
  ~oldWord: string,
  ~showAnswers: bool,
  ~winner: string,
  ~isGameOver: bool,
  ~resetConnState: (. unit) => unit,
  ~playerName: string,
  ~noplrs: int,
) => React.element
