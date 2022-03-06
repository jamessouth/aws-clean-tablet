@react.component
let make = (~toggleProp, ~toggleSetFunc) => {
  <button
    type_="button"
    className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 top-0 cursor-pointer"
    onClick={_ => toggleSetFunc(. prev => !prev)}>
    {switch toggleProp {
    | true => React.string("hide")
    | false => React.string("show")
    }}
  </button>
}
