@react.component
let make = (~textTrue="submit", ~textFalse="submit", ~textProp=true, ~onClick, ~disabled=false, ~img=React.null, ~className) => {
  <button type_="button" className onClick disabled>
    {switch textProp {
    | true => React.string(textTrue)
    | false => React.string(textFalse)
    }}
    {img}
  </button>
}
