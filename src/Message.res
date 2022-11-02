type propShape = {"msg": string}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"

@react.component
let make = (~msg) => {
  <p
    className="text-stone-100 absolute -top-20 w-4/5 bg-red-800 p-2 left-1/2 transform -translate-x-2/4 text-center font-anon text-sm">
    {React.string(msg)}
  </p>
}
