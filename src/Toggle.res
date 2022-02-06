@react.component
let make = (~toggleProp, ~toggleSetFunc) => {
  Js.log("render toggle")
  <button
    type_="button"
    className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 top-36 cursor-pointer"
    onClick={_ => toggleSetFunc(prev => !prev)}>
    {switch toggleProp {
    | true => "hide"->React.string
    | false => "show"->React.string
    }}
  </button>
}
