@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"

module Link = {
  open Route
  @react.component
  let make = (~route, ~className, ~content="") => {
    let onClick = e => {
      ReactEvent.Mouse.preventDefault(e)
      push(route)
      switch route {
      | SignIn =>
        ImageLoad.import_("./ImageLoad.bs")
        ->Promise.then(func => {
          Promise.resolve(func["bghand"](.))
        })
        ->ignore
      | Home | SignUp | GetInfo(_) | Confirm(_) | Lobby | Leaderboard | Play(_) | Other => ()
      }
    }

    <a onClick className href={typeToUrlString(route)}> {React.string(content)} </a>
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

  let route = Route.useRouter()
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
  open Route
  <>
    {switch route {
    | Leaderboard => React.null
    | Home | SignIn | SignUp | GetInfo(_) | Confirm(_) | Lobby | Play(_) | Other =>
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
      | Home | SignIn | SignUp | GetInfo(_) | Confirm(_) | Lobby | Play(_) | Other => "mb-8"
      }}>
      {switch (route, token) {
      | (Home, None) => {
          body(document)->setClassName("bodmob bodtab bodbig")
          <nav className="flex flex-col items-center relative">
            <Link
              route=SignIn
              className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred mb-8 sm:mb-16"}
              content="SIGN IN"
            />
            <Link
              route=SignUp
              className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred"}
              content="SIGN UP"
            />
            <Link
              route=GetInfo({search: VerificationCode})
              className={linkBase ++ "mt-10"}
              content="verification code?"
            />
            <Link
              route=GetInfo({search: ForgotPassword})
              className={linkBase ++ "mt-6"}
              content="forgot password?"
            />
            <Link
              route=GetInfo({search: ForgotUsername})
              className={linkBase ++ "mt-6"}
              content="forgot username?"
            />
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

      | (SignIn, None) => <Signin userpool setCognitoUser setToken cognitoUser />
      | (SignUp, None) => <Signup userpool setCognitoUser />
      | (GetInfo({search}), None) =>
        <GetInfo userpool cognitoUser setCognitoUser setShowName search />
      | (Confirm({search}), None) => <Confirm cognitoUser search />
      | (Lobby | Play(_) | Leaderboard, None) => {
          replace(Home)
          React.null
        }

      | (Home | SignIn | SignUp | GetInfo(_) | Confirm(_), Some(_)) => {
          replace(Lobby)
          React.null
        }

      | (Lobby | Play(_) | Leaderboard, Some(_)) =>
        <React.Suspense fallback=React.null> auth </React.Suspense>
      | (Other, _) => <div> {React.string("page not found")} </div> // <PageNotFound/>
      }}
    </main>
  </>
}
