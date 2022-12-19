@react.component
let make: (
  ~players: array<Reducer.player>,
  ~oldWord: string,
  ~showAnswers: bool,
  ~winner: string,
  ~isGameOver: bool,
  ~onEndClick: (. bool) => unit,
  ~playerName: string,
  ~endtoken: string,
  ~noplrs: int,
) => React.element
