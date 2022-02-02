let reqs = " 3-10 characters; letters, numbers, and underscores only."

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
    | Some(_) => setUnErr(_ => Some("Alphanumeric characters only." ++ reqs))
    | None => setUnErr(_ => None)
    }
  }

  let checkUnMaxLength = un => {
    switch un->Js.String2.length > 10 {
    | true => setUnErr(_ => Some("Username is too long." ++ reqs))
    | false => checkUnForbiddenChars(un)
    }
  }

  let checkNoUnWhitespace = un => {
    let r = %re("/\s/")
    switch Js.String2.match_(un, r) {
    | Some(_) => setUnErr(_ => Some("No whitespace allowed." ++ reqs))
    | None => checkUnMaxLength(un)
    }
  }

  let checkUnLength = un => {
    switch un->Js.String2.length < 3 {
    | true => setUnErr(_ => Some("Username is too short." ++ reqs))
    | false => checkNoUnWhitespace(un)
    }
  }

  let onBlur = _ => setUnVisited(_ => true)
  //   let onFocus = _ => setUnVisited(_ => false)

  let onChange = e => setUsername(_ => ReactEvent.Form.target(e)["value"])

  React.useEffect2(() => {
    switch unVisited {
    | true => checkUnLength(username)
    | false => setUnErr(_ => None)
    }
    None
  }, (username, unVisited))

  React.useEffect1(() => {
    switch unErr {
    | None => setDisabled(_ => false)
    | Some(_) => setDisabled(_ => true)
    }
    None
  }, [unErr])

  <div>
    <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
      {React.string("username:")}
    </label>
    <input
      autoComplete="username"
      autoFocus=true
      className="h-6 w-full text-xl font-anon text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
      id="username"
      maxLength=10
      minLength=3
      name="username"
      onBlur
      onChange
      //   onFocus
      // placeholder="Enter username"
      required=true
      spellCheck=false
      type_="text"
      value={username}
    />
  </div>
}
