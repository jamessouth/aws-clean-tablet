let pStyle = " text-stone-700 py-0 px-6 font-perm"
let maxWordLength = 12

@react.component
let make = (~onAnimationEnd, ~playerColor, ~word, ~answered, ~showTimer) => {
  let (alpha, setAlpha) = React.Uncurried.useState(_ => "")

  let blankPos = switch word->Js.String2.startsWith("_") {
  | true => "a blank then a word"
  | false => "a word then a blank"
  }

  React.useEffect1(() => {
    let alph = switch answered {
    | true => "70"
    | false => ""
    }
    setAlpha(._ => alph)
    None
  }, [answered])

  <div
    className="mt-20 mb-10 mx-auto bg-stone-100 relative w-80 h-36 flex flex-col items-center justify-center">
    {switch (playerColor == "transparent", word == "") {
    | (true, true) => <Loading fillColor="fill-stone-800" />
    | (false, true) | (true, false) | (false, false) => React.null
    }}
    {switch showTimer {
    | true =>
      <svg className="overflow-auto absolute top-0 left-0 w-full h-full" preserveAspectRatio="none">
        <rect
          x="0"
          y="0"
          width="100%"
          height="100%"
          onAnimationEnd
          style={ReactDOM.Style.make(~stroke={playerColor ++ alpha}, ())}
          className="animate-change stroke-offset-0 fill-transparent stroke-16 stroke-dash-1000"
        />
      </svg>
    | false => React.null
    }}
    <p
      ariaLabel={blankPos}
      role="alert"
      className={switch Js.String2.length(word) > maxWordLength {
      | true => "text-3xl" ++ pStyle
      | false => "text-4xl" ++ pStyle
      }}>
      {React.string(word)}
    </p>
  </div>
}
