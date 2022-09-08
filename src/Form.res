@react.component
let make = (~ht="h-72", ~on_Click, ~leg, ~validationError, ~cognitoError, ~children) => {
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  let onClick = _ => {
    setSubmitClicked(._ => true)
    on_Click()
  }

  <form className="w-4/5 m-auto relative">
    <fieldset className={`flex flex-col items-center justify-around ${ht}`}>
      {switch Js.String2.length(leg) > 0 {
      | true =>
        <legend className="text-stone-100 m-auto mb-6 text-3xl font-fred">
          {React.string(leg)}
        </legend>
      | false => React.null
      }}
      {switch (submitClicked, validationError, cognitoError) {
      | (false, _, None) | (true, None, None) => React.null
      | (false, _, Some(error)) | (true, None, Some(error)) | (true, Some(error), _) =>
        <span
          className="absolute right-0 -top-24 text-sm text-stone-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
          {React.string(error)}
        </span>
      }}
      {children}
    </fieldset>
    <Button onClick> {React.string("submit")} </Button>
  </form>
}
