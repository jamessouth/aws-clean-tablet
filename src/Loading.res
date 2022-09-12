@react.component
let make = (~fillColor="fill-stone-100", ~label="") => <>
  {switch label == "" {
  | true => React.null
  | false =>
    <p className="text-center text-stone-100 font-anon text-lg">
      {React.string(`loading ${label}`)}
    </p>
  }}
  <div className="h-8 m-auto w-8">
    <svg
      className={`w-8 h-8 animate-spin ${fillColor}`}
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg">
      <path
        d="M94 39a4 4 0 0 0 3-5A50 50 0 0 0 42 1a4 4 0 0 0-3 6 5 5 0 0 0 5 3 41 41 0 0 1 44 26 5 5 0 0 0 6 3Z"
        fill="currentFill"
      />
    </svg>
  </div>
</>
