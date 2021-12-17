

@react.component
let make = (~onAnimationEnd, ~playerColor, ~currentWord:string) => {

    let blankPos = switch currentWord->Js.String2.startsWith("_") {
    | true => "a blank then a word"
    | false => "a word then a blank"
    }




    <div className="mt-20 mb-10 mx-auto bg-smoke-100 relative w-80 h-36 flex items-center justify-center">
        <svg className="overflow-auto absolute top-0 left-0 w-full h-full" preserveAspectRatio="none">
            <rect x="0" y="0" width="100%" height="100%" onAnimationEnd style={ReactDOM.Style.make(~stroke={playerColor}, ())} className="animate-change stroke-offset-0 fill-transparent stroke-16 stroke-dash-1000"></rect>
        </svg>
        <p ariaLabel={blankPos} role="alert" className="text-smoke-700 text-4xl py-0 px-6 font-perm">{currentWord->React.string}</p>
    </div>
    
    


}