let useUsernameValidation = (username, setUsernameError) => {
  let rec validate = (funcs, str, acc) =>
    switch funcs {
    | list{} => acc
    | list{h, ...t} => validate(t, str, acc ++ h(str))
    }

  let myList = list{
    s =>
      switch Js.String2.length(s) < 3 || Js.String2.length(s) > 10 {
      | false => ""
      | true => "3-10 characters; "
      },
    s =>
      switch Js.String2.match_(s, %re("/\W/")) {
      | None => ""
      | Some(_) => "letters, numbers, and underscores only; no whitespace."
      },
  }

  React.useEffect1(() => {
    let error = validate(myList, username, "")
    let final = switch error == "" {
    | true => None
    | false => Some("Username: " ++ error)
    }
    setUsernameError(_ => final)
    None
  }, [username])
}
