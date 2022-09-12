@react.component
let make: (
  ~onAnimationEnd: ReactEvent.Animation.t => unit,
  ~playerColor: string,
  ~word: string,
  ~answered: bool,
  ~showTimer: bool,
) => React.element
