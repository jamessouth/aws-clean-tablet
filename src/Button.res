@react.component
let make = (~textTrue, ~textFalse, ~textProp, ~onClick, ~disabled, ~img=React.null, ~className) => {
  <button type_="button" className onClick disabled>
    {switch textProp {
    | true => React.string(textTrue)
    | false => React.string(textFalse)
    }}
    {img}
  </button>
}
