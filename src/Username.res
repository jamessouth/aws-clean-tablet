let reqs = "Usernames must: contain only letters, numbers, and underscores; be 3-10 characters long"

@react.component
let make = (
  ~unVisited,
  ~setUnVisited,
  ~unErr,
  ~setUnErr,
  ~setDisabled,
  ~username,
  ~setUsername,
) => {
  let checkUnForbiddenChars = un => {
    let r = %re("/\W/")
    switch Js.String2.match_(un, r) {
    | Some(_) => setUnErr(_ => Some("Alphanumeric characters only.\n" ++ reqs))
    | None => setUnErr(_ => None)
    }
  }

  let checkUnMaxLength = un => {
    switch un->Js.String2.length > 10 {
    | true => setUnErr(_ => Some("Username too long.\n" ++ reqs))
    | false => checkUnForbiddenChars(un)
    }
  }

  let checkNoUnWhitespace = un => {
    let r = %re("/\s/")
    switch Js.String2.match_(un, r) {
    | Some(_) => setUnErr(_ => Some("No whitespace allowed.\n" ++ reqs))
    | None => checkUnMaxLength(un)
    }
  }

  let checkUnLength = un => {
    switch un->Js.String2.length < 3 {
    | true => setUnErr(_ => Some("Username too short.\n" ++ reqs))
    | false => checkNoUnWhitespace(un)
    }
  }

  let onBlur = _ => setUnVisited(_ => true)

  let onChange = e => setUsername(_ => ReactEvent.Form.target(e)["value"])

  React.useEffect2(() => {
    switch unVisited {
    | true => checkUnLength(username)
    | false => setUnErr(_ => None)
    }
    None
  }, (username, unVisited))

  React.useEffect2(() => {
    switch (unErr, username->Js.String2.length < 4) {
    | (None, false) => setDisabled(_ => false)
    | (Some(_), _) | (_, true) => setDisabled(_ => true)
    }
    None
  }, (unErr, username))

  <div className="relative">
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
      {React.string("username:")}
    </label>
    {switch (unVisited, unErr) {
    | (true, Some(err)) =>
      <span
        className="absolute right-0 text-lg text-warm-gray-100 bg-red-500 font-anon font-flow h-30 w-2/3 z-10">
        {React.string(err)}
      </span>
    | (false, _) | (true, None) => React.null
    }}
    <input
      autoComplete="username"
      autoFocus=true
      className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="username"
      minLength=4
      name="username"
      onBlur
      onChange
      // placeholder="Enter username"
      required=true
      spellCheck=false
      type_="text"
      value={username}
    />
  </div>
}
