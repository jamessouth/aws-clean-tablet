@react.component
let make = (~text, ~onClick) => {



      <button
        type_="button"
        className="text-gray-700 mt-14 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"
        onClick>
        {React.string(text)}
      </button>

}