@react.component
let make = (~validationError=None, ~cognitoError=None) => {
  switch (validationError, cognitoError) {
  | (Some(err), _) | (_, Some(err)) =>
    <span
      className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
      {React.string(err)}
    </span>
  | (None, None) => React.null
  }
}
