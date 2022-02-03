let reqs = " 8-98 characters; at least 1 each of uppercase and lowercase letters, numbers, and symbols required."

let usePasswordValidation = password => {
  let (pwErr, setPwErr) = React.useState(_ => None)

  let checkNoPwWhitespace = pw => {
    let r = %re("/\s/")
    switch Js.String2.match_(pw, r) {
    | Some(_) => setPwErr(_ => Some("No whitespace allowed." ++ reqs))
    | None => setPwErr(_ => None)
    }
  }

  let checkPwMaxLength = pw => {
    switch pw->Js.String2.length > 98 {
    | true => setPwErr(_ => Some("Password is too long." ++ reqs))
    | false => checkNoPwWhitespace(pw)
    }
  }

  let checkSymbol = pw => {
    let r = %re("/[!-*\[-`{-~./,:;<>?@]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a symbol." ++ reqs))
    | Some(_) => checkPwMaxLength(pw)
    }
  }

  let checkNumber = pw => {
    let r = %re("/\d/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a number." ++ reqs))
    | Some(_) => checkSymbol(pw)
    }
  }

  let checkUpper = pw => {
    let r = %re("/[A-Z]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add an uppercase letter." ++ reqs))
    | Some(_) => checkNumber(pw)
    }
  }

  let checkLower = pw => {
    let r = %re("/[a-z]/")
    switch Js.String2.match_(pw, r) {
    | None => setPwErr(_ => Some("Add a lowercase letter." ++ reqs))
    | Some(_) => checkUpper(pw)
    }
  }

  let checkPwLength = pw => {
    switch pw->Js.String2.length < 8 {
    | true => setPwErr(_ => Some("Password is too short." ++ reqs))
    | false => checkLower(pw)
    }
  }

  React.useEffect1(() => {
    checkPwLength(password)
    None
  }, [password])

  pwErr
}
