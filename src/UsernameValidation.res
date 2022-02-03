let reqs = " 3-10 characters; letters, numbers, and underscores only."

let useUsernameValidation = username => {
  let (unErr, setUnErr) = React.useState(_ => None)

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

  React.useEffect1(() => {
    checkUnLength(username)
    None
  }, [username])

  unErr
}
