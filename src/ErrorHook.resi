type t =
  | Username
  | Password
  | Code
  | Email
  | Answer

let useMultiError: (
  Js.Array2.t<(Js.String2.t, t)>,
  (. option<string> => option<string>) => unit,
) => unit

let useError: (Js.String2.t, t, (. option<string> => option<string>) => unit) => unit
