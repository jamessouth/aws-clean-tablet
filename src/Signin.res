









@react.component
let make = (~userpool, ~setCognitoUser, ~setToken) => {
  let (_cognitoErr, setCognitoErr) = React.useState(_ => None)
  let (showPassword, setShowPassword) = React.useState(_ => false)
  let (disabled, setDisabled) = React.useState(_ => true)
  let (username, setUsername) = React.useState(_ => "")
  let (password, setPassword) = React.useState(_ => "")
  let onChange = (func, e) => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->func
  }

  let {handleSubmit: onSubmit} = SigninHook.useSignin(username, password, userpool, setCognitoErr, setToken, setCognitoUser)

  let onClick = _e => {
    (prev => !prev)->setShowPassword
  }

  React.useEffect2(() => {
    switch (username->Js.String2.length > 3, password->Js.String2.length > 7) {
    | (true, true) => (_ => false)->setDisabled
    | (false, true) | (true, false) | (false, false) => (_ => true)->setDisabled
    }

    None
  }, (username, password))

  <main>
    <form className="w-4/5 m-auto" onSubmit>
      <fieldset className="flex flex-col items-center justify-around h-80">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {"Sign in"->React.string}
        </legend>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
            {"username:"->React.string}
          </label>
          <input
            autoComplete="username"
            autoFocus=true
            className="h-6 w-full text-xl pl-1 text-left font-anon outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="username"
            minLength=4
            name="username"
            onChange={onChange(setUsername)}
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={username}
          />
        </div>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="password">
            {"password:"->React.string}
          </label>
          <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            onChange={onChange(setPassword)}
            // placeholder="Enter password"
            required=true
            spellCheck=false
            type_={switch showPassword {
            | true => "text"
            | false => "password"
            }}
            value={password}
          />
          <button
            type_="button"
            className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 cursor-pointer"
            onClick>
            {switch showPassword {
            | true => "hide"->React.string
            | false => "show"->React.string
            }}
          </button>
        </div>
        <Link
          url="/resetpwd"
          className="self-end text-sm cursor-pointer font-anon text-warm-gray-100"
          content="forgot password?"
        />
      </fieldset>
      <button
        disabled
        className="text-gray-700 mt-16 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7">
        {"submit"->React.string}
      </button>
    </form>
  </main>
}
