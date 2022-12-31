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

let stringToGetInfo = s =>
  switch s {
  | "cd_un" => GetInfo({search: VerificationCode})
  | "pw_un" => GetInfo({search: ForgotPassword})
  | "un_em" => GetInfo({search: ForgotUsername})
  | _ => Other
  }

let stringToConfirm = s =>
  switch s {
  | "cd_un" => Confirm({search: VerificationCode})
  | "pw_un" => Confirm({search: ForgotPassword})
  | _ => Other
  }

let stringToAuthSubroute = l =>
  switch l {
  | list{"leaderboard"} => Auth({subroute: Leaderboard})
  | list{"lobby"} => Auth({subroute: Lobby})
  | list{"play", gameno} => Auth({subroute: Play({play: gameno})})
  | _ => Auth({subroute: Other})
  }

let queryToString = q =>
  switch q {
  | VerificationCode => "cd_un"
  | ForgotPassword => "pw_un"
  | ForgotUsername => "un_em"
  | Other => ""
  }

let authSubrouteToString = a =>
  switch a {
  | Leaderboard => "leaderboard"
  | Lobby => "lobby"
  | Play({play}) => `play/${play}`
  | Other => ""
  }

let urlStringToType = (url: RescriptReactRouter.url) =>
  switch url.path {
  | list{} => Home
  | list{"signin"} => SignIn
  | list{"signup"} => SignUp
  | list{"getinfo"} => stringToGetInfo(url.search)
  | list{"confirm"} => stringToConfirm(url.search)
  | list{"auth", ...subroutes} => stringToAuthSubroute(subroutes)
  | _ => Other
  }

let typeToUrlString = t =>
  switch t {
  | Home => "/"
  | SignIn => "/signin"
  | SignUp => "/signup"
  | GetInfo({search}) => `/getinfo?${queryToString(search)}`
  | Confirm({search}) => `/confirm?${queryToString(search)}`
  | Auth({subroute}) => `/auth/${authSubrouteToString(subroute)}`
  | Other => ""
  }

let useRouter = () => urlStringToType(RescriptReactRouter.useUrl())
let replace = route => route->typeToUrlString->RescriptReactRouter.replace
let push = route => route->typeToUrlString->RescriptReactRouter.push
