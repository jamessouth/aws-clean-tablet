@react.component
let make: (
  ~players: array<Reducer.livePlayer>,
  ~oldWord: string,
  ~showAnswers: bool,
  ~winner: string,
  ~isWinner: bool,
  ~onEndClick: (. bool) => unit,
  ~playerName: string,
  ~noplrs: int,
) => React.element
