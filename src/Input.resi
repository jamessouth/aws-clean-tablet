@react.component
let make: (
  ~value: string,
  ~propName: string,
  ~autoComplete: string=?,
  ~inputMode: string=?,
  ~onKeyPress: ReactEvent.Keyboard.t => unit=?,
  ~setFunc: (. 'a => 'b) => unit,
) => React.element
