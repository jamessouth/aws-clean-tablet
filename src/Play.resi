type answerPayload = {
  action: string,
  gameno: string,
  answer: string,
}

type endPayload = {
  action: string,
  gameno: string,
  token: string,
}

@react.component
let make: (
  ~players: array<Reducer.livePlayer>,
  ~sk: string,
  ~showAnswers: bool,
  ~winner: string,
  ~isWinner: bool,
  ~oldWord: string,
  ~word: string,
  ~playerColor: string,
  ~send: (. option<string>) => unit,
  ~playerName: string,
  ~endtoken: string,
  ~resetConnState: (. unit) => unit,
) => React.element
