@react.component
let make = (
  ~ht="h-72",
  ~onClick,
  ~leg,
  ~submitClicked,
  ~validationError,
  ~cognitoError,
  ~children,
) => {
  <form className="w-4/5 m-auto relative">
    <fieldset className={`flex flex-col items-center justify-around ${ht}`}>
      <legend className="text-stone-100 m-auto mb-6 text-3xl font-fred">
        {React.string(leg)}
      </legend>
      {switch submitClicked {
      | false => React.null
      | true =>
        switch (validationError, cognitoError) {
        | (Some(error), _) | (_, Some(error)) =>
          <span
            className="absolute right-0 -top-24 text-sm text-stone-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
            {React.string(error)}
          </span>
        | (None, None) => React.null
        }
      }}
      {children}
    </fieldset>
    <Button onClick />
  </form>
}
