


@react.component
let make = (onAnimationEnd, playerColor, word) => {

    let blankPos = switch word->Js.String2.startsWith("_") {
    | true => "a blank then a word"
    | false => "a word then a blank"
    }


    <div className="bg-smoke-100 relative w-80 h-36 flex items-center justify-center">
        <svg className="overflow-visible absolute top-0 left-0 w-full h-full" preserveAspectRatio="none">
            <rect x="0" y="0" width="100%" height="100%" onAnimationEnd style={ReactDOM.Style.make(~stroke={playerColor}, ())} className="animate-change rect"></rect>
        </svg>
        <p ariaLabel={blankPos} role="alert" className="text-smoke-700 text-4xl py-0 px-6 font-perm">{word->React.string}</p>
    </div>
    
    


}