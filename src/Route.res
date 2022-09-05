  type t =
    | Home
    | SignIn
    | SignUp
    | GetInfo({search: string})
    | Confirm({search: string})
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
      switch url.search {
      | "cd_un" => GetInfo({search: "cd_un"})
      | "pw_un" => GetInfo({search: "pw_un"})
      | "un_em" => GetInfo({search: "un_em"})
      | _ => Other
      }
    | list{"confirm"} =>
      switch url.search {
      | "cd_un" => Confirm({search: "cd_un"})
      | "pw_un" => Confirm({search: "pw_un"})
      | _ => Other
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
    | GetInfo({search}) => `/getinfo?${search}`
    | Confirm({search}) => `/confirm?${search}`
    | Lobby => "/auth/lobby"
    | Leaderboard => "/auth/leaderboard"
    | Play({play}) => `/auth/play/${play}`
    | Other => ""
    }

  let useRouter = () => urlStringToType(RescriptReactRouter.useUrl())
  let replace = route => route->typeToUrlString->RescriptReactRouter.replace
  let push = route => route->typeToUrlString->RescriptReactRouter.push
