type query =
  | VerificationCode
  | ForgotPassword
  | ForgotUsername
  | Other

type authSubroute =
  | Leaderboard
  | Lobby
  | Play({play: string})
  | Other

type t =
  | Home
  | SignIn
  | SignUp
  | GetInfo({search: query})
  | Confirm({search: query})
  | Auth({subroute: authSubroute})
  | Other

let typeToUrlString: t => string

let useRouter: unit => t
let replace: t => unit
let push: t => unit
