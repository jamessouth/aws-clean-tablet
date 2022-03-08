@react.component
let make = (~error) => {
  <span
    className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
    {React.string(error)}
  </span>
}
