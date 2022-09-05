type query =
  | VerificationCode
  | ForgotPassword
  | ForgotUsername
  | Other

let fromStringToQuery = s =>
  switch s {
  | "cd_un" => VerificationCode
  | "pw_un" => ForgotPassword
  | "un_em" => ForgotUsername
  | _ => Other
  }

let fromQueryToString = q =>
  switch q {
  | VerificationCode => "cd_un"
  | ForgotPassword => "pw_un"
  | ForgotUsername => "un_em"
  | Other => ""
  }

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

let urlStringToType = (url: RescriptReactRouter.url) =>
  switch url.path {
  | list{} => Home
  | list{"signin"} => SignIn
  | list{"signup"} => SignUp
  | list{"getinfo"} =>
    switch fromStringToQuery(url.search) {
    | VerificationCode => GetInfo({search: VerificationCode})
    | ForgotPassword => GetInfo({search: ForgotPassword})
    | ForgotUsername => GetInfo({search: ForgotUsername})
    | Other => Other
    }
  | list{"confirm"} =>
    switch fromStringToQuery(url.search) {
    | VerificationCode => Confirm({search: VerificationCode})
    | ForgotPassword => Confirm({search: ForgotPassword})
    | ForgotUsername | Other => Other
    }
  | list{"auth", ...subroutes} =>
    switch subroutes {
    | list{"lobby"} => Lobby
    | list{"leaderboard"} => Leaderboard
    | list{"play", gameno} => Play({play: gameno})
    | _ => Other
    }
  | _ => Other
  }

let typeToUrlString = t =>
  switch t {
  | Home => "/"
  | SignIn => "/signin"
  | SignUp => "/signup"
  | GetInfo({search}) => `/getinfo?${fromQueryToString(search)}`
  | Confirm({search}) => `/confirm?${fromQueryToString(search)}`
  | Lobby => "/auth/lobby"
  | Leaderboard => "/auth/leaderboard"
  | Play({play}) => `/auth/play/${play}`
  | Other => ""
  }

let useRouter = () => urlStringToType(RescriptReactRouter.useUrl())
let replace = route => route->typeToUrlString->RescriptReactRouter.replace
let push = route => route->typeToUrlString->RescriptReactRouter.push
