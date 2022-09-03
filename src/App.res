@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"

module Link = {
  @react.component
  let make = (~url, ~className, ~content="") => {
    let onClick = e => {
      ReactEvent.Mouse.preventDefault(e)
      RescriptReactRouter.push(url)
      switch url {
      | "/signin" =>
        ImageLoad.import_("./ImageLoad.bs")
        ->Promise.then(func => {
          Promise.resolve(func["bghand"](.))
        })
        ->ignore
      | _ => ()
      }
    }

    <a onClick className href={url}> {React.string(content)} </a>
  }
}


module Route = {
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

let urlStringToType = (url: ReasonReactRouter.url) =>
  switch url.path {
  | list{} => Home
  | list{"signin"} => Signin
  | list{"signup"} => Signup
  | list{"getinfo"} => GetInfo({search: url.search})
  | list{"confirm"} => Confirm({search: url.search})
  | list{"auth", "lobby"} => Lobby
  | list{"auth", "leaderboard"} => Leaderboard
  | list{"auth", "play", gameno} => Play({play: gameno})
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
}








@react.component
let make = () => {
  Js.log("app")

  let linkBase = "w-3/5 text-stone-100 block font-bold font-anon text-sm max-w-80 "

  let linkBase2 = "w-3/5 border border-stone-100 block bg-stone-800/40 text-center text-stone-100 "

  open Cognito
  let userpool = userPoolConstructor({
    userPoolId: upid,
    clientId: cid,
    advancedSecurityDataCollectionFlag: false,
  })

  let route = Route.urlStringToType(RescriptReactRouter.useUrl())
  let (cognitoUser: Js.Nullable.t<usr>, setCognitoUser) = React.Uncurried.useState(_ =>
    Js.Nullable.null
  )
  

  let (token, setToken) = React.Uncurried.useState(_ => None)
  let (showName, setShowName) = React.Uncurried.useState(_ => "")



  let auth = React.useMemo4(_ =>
    React.createElement(
      Auth.lazy_(() =>
        Auth.import_("./Auth.bs")->Promise.then(
          comp => {
            Promise.resolve({"default": comp["make"]})
          },
        )
      ),
      Auth.makeProps(~token, ~setToken, ~cognitoUser, ~setCognitoUser, ()),
    )
  , (token, setToken, cognitoUser, setCognitoUser))

  open Web
  <>
    {switch route {
    | Leaderboard => React.null
    | Home  | SignIn | SignUp | GetInfo | Confirm | Lobby | Play | Other =>
      switch token {
      | None =>
        <header className="mb-10 newgmimg:mb-12">
          <h1
            className="text-6xl mt-21 mx-auto px-6 text-center font-arch decay-mask text-stone-100">
            {React.string("CLEAN TABLET")}
          </h1>
        </header>
      | Some(_) => React.null
      }
    }}
    <main
      className={switch route {
      | Leaderboard => ""
      | Home  | SignIn | SignUp | GetInfo | Confirm | Lobby | Play | Other => "mb-8"
      }}>
      {switch (route, token) {
      | (Home, None) => {
          body(document)->setClassName("bodmob bodtab bodbig")
          <nav className="flex flex-col items-center relative">
            <Link
              url="/signin"
              className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred mb-8 sm:mb-16"}
              content="SIGN IN"
            />
            <Link
              url="/signup"
              className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred"}
              content="SIGN UP"
            />
            <Link
              url="/getinfo?cd_un" className={linkBase ++ "mt-10"} content="verification code?"
            />
            <Link url="/getinfo?pw_un" className={linkBase ++ "mt-6"} content="forgot password?" />
            <Link url="/getinfo?un_em" className={linkBase ++ "mt-6"} content="forgot username?" />
            {switch showName == "" {
            | true => React.null
            | false =>
              <p className="text-stone-100 absolute -top-20 w-4/5 bg-blue-gray-800 p-2 font-anon">
                {React.string(
                  "The username associated with the email you submitted is:" ++ showName,
                )}
              </p>
            }}
          </nav>
        }

      | (Signin, None) =>
      <Signin userpool setCognitoUser setToken cognitoUser 
      // cognitoError setCognitoError 
      />
      //  <React.Suspense fallback=React.null> signin </React.Suspense>

      | (Signup, None) => 
      <Signup userpool setCognitoUser 
      // cognitoError setCognitoError 
      />
      // <React.Suspense fallback=React.null> signup </React.Suspense>

      | (GetInfo, None) =>
        switch search {
        | "cd_un" | "pw_un" | "un_em" =>
          <GetInfo
            userpool cognitoUser setCognitoUser 
            // cognitoError setCognitoError
             setShowName search
          />
          // <React.Suspense fallback=React.null> getInfo </React.Suspense>

        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (Confirm, None) =>
        switch search {
        | "cd_un" | "pw_un" => 
        <Confirm cognitoUser 
        // cognitoError setCognitoError
         search />
        // <React.Suspense fallback=React.null> confirm </React.Suspense>

        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (Lobby | Play | Leaderboard, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      | (Home | Signin | Signup | GetInfo | Confirm, Some(_)) => {
          RescriptReactRouter.replace(Route.typeToUrlString(Lobby))
          React.null
        }

      | (Lobby | Play | Leaderboard, Some(_)) => <React.Suspense fallback=React.null> auth </React.Suspense>

      | (Other, _) => <div> {React.string("other")} </div> // <PageNotFound/>
      }}
    </main>
  </>
}
