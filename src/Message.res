@react.component
let make = (~children) => {
  <p
    className="text-stone-100 absolute -top-20 w-4/5 bg-red-800 p-2 left-1/2 transform -translate-x-2/4 text-center font-anon text-sm">
    {children}
  </p>
}
