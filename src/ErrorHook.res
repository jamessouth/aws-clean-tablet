let checkLength = (. min, max, str) =>
  switch Js.String2.length(str) < min || Js.String2.length(str) > max {
  | false => ""
  | true => j`$min-$max length; `
  }
let checkInclusion = (. re, msg, str) =>
  switch Js.String2.match_(str, re) {
  | None => msg
  | Some(_) => ""
  }
let checkExclusion = (. re, msg, str) =>
  switch Js.String2.match_(str, re) {
  | None => ""
  | Some(_) => msg
  }

let getFuncs = input =>
  switch input {
  | "USERNAME" => [
      (. s) => checkLength(. 3, 10, s),
      (. s) =>
        checkExclusion(.
          %re("/\W/"),
          "letters, numbers, and underscores only; no whitespace or symbols; ",
          s,
        ),
    ]
  | "PASSWORD" => [
      (. s) => checkLength(. 8, 98, s),
      (. s) => checkInclusion(. %re("/[!-/:-@\[-`{-~]/"), "at least 1 symbol; ", s),
      (. s) => checkInclusion(. %re("/\d/"), "at least 1 number; ", s),
      (. s) => checkInclusion(. %re("/[A-Z]/"), "at least 1 uppercase letter; ", s),
      (. s) => checkInclusion(. %re("/[a-z]/"), "at least 1 lowercase letter; ", s),
      (. s) => checkExclusion(. %re("/\s/"), "no whitespace; ", s),
    ]
  | "CODE" => [(. s) => checkInclusion(. %re("/^\d{6}$/"), "6-digit number only; ", s)]
  | "EMAIL" => [
      (. s) => checkLength(. 5, 99, s),
      (. s) =>
        checkInclusion(.
          %re(
            "/^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/"
          ),
          "enter a valid email address.",
          s,
        ),
    ]
  | "ANSWER" => [
      (. s) => checkLength(. 2, 12, s),
      (. s) => checkInclusion(. %re("/[a-z ]/i"), "letters and spaces only; ", s),
      (. s) => checkExclusion(. %re("/\d/"), "no numbers; ", s),
      (. s) => checkExclusion(. %re("/[!-/:-@\[-`{-~]/"), "no symbols; ", s),
      (. s) => checkExclusion(. %re("/^\s|\s$/"), "must begin and end with letters; ", s),
    ]
  | _ => []
  }

let useMultiError = (fields, setErrorFunc) => {
  let errs = fields->Js.Array2.map(fld => {
    let (val, prop) = fld
    let error = getFuncs(prop)->Js.Array2.reduce((acc, f) => acc ++ f(. val), "")
    switch error == "" {
    | true => ""
    | false => prop ++ ": " ++ error
    }
  })
  let total = errs->Js.Array2.joinWith("")
  let final = switch total == "" {
  | true => None
  | false => Some(total)
  }
  setErrorFunc(._ => final)
}

let useError = (value, propName, setErrorFunc) => {
  Js.log("Errorhook2")

  let error = getFuncs(propName)->Js.Array2.reduce((acc, f) => acc ++ f(. value), "")
  let final = switch error == "" {
  | true => None
  | false => Some(propName ++ ": " ++ error)
  }
  setErrorFunc(._ => final)
}
