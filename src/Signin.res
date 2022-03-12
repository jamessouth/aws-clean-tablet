let className = "text-gray-700 mt-14 bg-warm-gray-100 block max-w-xs lg:max-w-sm font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"

@react.component
let make = (
  ~signinOnClick,
  ~cognitoError,
  ~submitClicked,
  ~validationError,
  ~username,
  ~password,
  ~setUsername,
  ~setPassword,
) => {
  Js.log("signin")

  let error = switch submitClicked {
  | false => React.null
  | true =>
    switch (validationError, cognitoError) {
    | (Some(error), _) | (_, Some(error)) =>
      <span
        className="absolute right-0 -top-24 text-sm text-warm-gray-100 bg-red-600 font-anon w-4/5 leading-4 p-1">
        {React.string(error)}
      </span>
    | (None, None) => React.null
    }
  }

  let btn = <Button onClick=signinOnClick className />

  <Form btn leg="Sign in">
    {error}
    <Input value=username propName="username" setFunc=setUsername />
    <Input value=password propName="password" autoComplete="current-password" setFunc=setPassword />
  </Form>
}
