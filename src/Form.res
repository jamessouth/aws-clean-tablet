@react.component
let make = (~ht="72", ~btn, ~leg, ~children) => {
  <form className="w-4/5 m-auto relative">
    <fieldset className={`flex flex-col items-center justify-around h-${ht}`}>
      <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
        {React.string(leg)}
      </legend>
      {children}
    </fieldset>
    {btn}
  </form>
}
