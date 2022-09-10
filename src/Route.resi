type query =
  | VerificationCode
  | ForgotPassword
  | ForgotUsername
  | Other

type t =
  | Home
  | SignIn
  | SignUp
  | GetInfo({search: query})
  | Confirm({search: query})
  | Lobby
  | Leaderboard
  | Play({play: string})
  | Other

let typeToUrlString: t => string

let useRouter: unit => t
let replace: t => unit
let push: t => unit
