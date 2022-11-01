@react.component
let make = (~ht="h-72", ~on_Click, ~leg, ~validationError, ~cognitoError, ~children) => {
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)

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
        <Message> {React.string(error)} </Message>
      }}
      {children}
    </fieldset>
    {switch (submitClicked, leg == "Sign in", validationError, cognitoError) {
    | (false, _, _, _)
    | (true, false, _, _)
    | (true, true, Some(_), _)
    | (true, true, _, Some(_)) => React.null
    | (true, true, None, None) =>
      <div className="absolute left-1/2 transform -translate-x-2/4 bottom-10">
        <Loading label="..." />
      </div>
    }}
    <Button
      onClick={_ => {
        setSubmitClicked(._ => true)
        on_Click(.)
      }}>
      {React.string("submit")}
    </Button>
  </form>
}
