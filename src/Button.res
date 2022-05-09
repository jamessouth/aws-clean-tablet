@react.component
let make = (~textTrue="submit", ~textFalse="submit", ~textProp=true, ~onClick, ~disabled=false, ~img=React.null, ~className="text-stone-800 mt-14 bg-stone-100 hover:bg-stone-300 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7") => {
  <button type_="button" className onClick disabled>
    {switch textProp {
    | true => React.string(textTrue)
    | false => React.string(textFalse)
    }}
    {img}
  </button>
}
