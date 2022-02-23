let checkLength = (min, max, str) =>
  switch Js.String2.length(str) < min || Js.String2.length(str) > max {
  | false => ""
  | true => j`$min-$max characters; `
  }
let checkInclusion = (re, msg, str) =>
  switch Js.String2.match_(str, re) {
  | None => msg
  | Some(_) => ""
  }
let checkExclusion = (re, msg, str) =>
  switch Js.String2.match_(str, re) {
  | None => ""
  | Some(_) => msg
  }

let getFuncs = input =>
  switch input {
  | "username" => [
      s => checkLength(3, 10, s),
      s =>
        checkExclusion(
          %re("/\W/"),
          "letters, numbers, and underscores only; no whitespace or symbols.",
          s,
        ),
    ]
  | "password" => [
      s => checkLength(8, 98, s),
      s => checkInclusion(%re("/[!-/:-@\[-`{-~]/"), "at least 1 symbol; ", s),
      s => checkInclusion(%re("/\d/"), "at least 1 number; ", s),
      s => checkInclusion(%re("/[A-Z]/"), "at least 1 uppercase letter; ", s),
      s => checkInclusion(%re("/[a-z]/"), "at least 1 lowercase letter; ", s),
      s => checkExclusion(%re("/\s/"), "no whitespace.", s),
    ]
  | "code" => [s => checkInclusion(%re("/^\d{6}$/"), "6-digit number only.", s)]
  | "email" => [
      s => checkLength(5, 99, s),
      s =>
        checkInclusion(
          %re(
            "/^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/"
          ),
          "enter a valid email address.",
          s,
        ),
    ]
  | _ => []
  }

let useError = (value, propName, setErrorFunc) => {
  Js.log("Errorhook")

  React.useEffect1(() => {
    let error = getFuncs(propName)->Js.Array2.reduce((acc, f) => acc ++ f(value), "")
    let final = switch error == "" {
    | true => None
    | false => Some(propName ++ ": " ++ error)
    }
    setErrorFunc(_ => final)
    None
  }, [value])
}
