let reqs = " 3-10 characters; letters, numbers, and underscores only."

let useUsernameValidation = username => {
  let (err, setErr) = React.useState(_ => None)

  let checkUnForbiddenChars = un => {
    let r = %re("/\W/")
    switch Js.String2.match_(un, r) {
    | Some(_) => setErr(_ => Some("Alphanumeric characters only." ++ reqs))
    | None => setErr(_ => None)
    }
  }
  let checkUnMaxLength = un => {
    switch un->Js.String2.length > 10 {
    | true => setErr(_ => Some("Username is too long." ++ reqs))
    | false => checkUnForbiddenChars(un)
    }
  }
  let checkNoUnWhitespace = un => {
    let r = %re("/\s/")
    switch Js.String2.match_(un, r) {
    | Some(_) => setErr(_ => Some("No whitespace allowed." ++ reqs))
    | None => checkUnMaxLength(un)
    }
  }

  let checkUnLength = un => {
    switch un->Js.String2.length < 3 {
    | true => setErr(_ => Some("Username is too short." ++ reqs))
    | false => checkNoUnWhitespace(un)
    }
  }

  React.useEffect1(() => {
    checkUnLength(username)
    None
  }, [username])

  err
}
