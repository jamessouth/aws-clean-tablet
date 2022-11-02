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

  let userpool = Cognito.userPoolConstructor({
    userPoolId: upid,
    clientId: cid,
    advancedSecurityDataCollectionFlag: false,
  })

  let route = Route.useRouter()
  let (cognitoUser: Js.Nullable.t<Cognito.usr>, setCognitoUser) = React.Uncurried.useState(_ =>
    Js.Nullable.null
  )

  let (token, setToken) = React.Uncurried.useState(_ => None)
  let (retrievedUsername, setRetrievedUsername) = React.Uncurried.useState(_ => "")
  let (wsError, setWsError) = React.Uncurried.useState(_ => "")

  Js.log2("wserr", wsError)

  let msg = React.createElement(
    Message.lazy_(() =>
      Message.import_("./Message.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Message.makeProps(~msg=wsError, ()),
  )

  let signin = React.createElement(
    Signin.lazy_(() =>
      Signin.import_("./Signin.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Signin.makeProps(~userpool, ~setCognitoUser, ~setToken, ~cognitoUser, ~retrievedUsername, ()),
  )

  let signup = React.createElement(
    Signup.lazy_(() =>
      Signup.import_("./Signup.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Signup.makeProps(~userpool, ~setCognitoUser, ()),
  )

  let auth = React.createElement(
    Auth.lazy_(() =>
      Auth.import_("./Auth.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Auth.makeProps(~token, ~setToken, ~cognitoUser, ~setCognitoUser, ~setWsError, ~route, ()),
  )

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
      | Home | SignIn | SignUp | GetInfo(_) | Confirm(_) | Lobby | Play(_) | Other => "mb-14"
      }}>
      {switch (route, token) {
      | (Home, None) => {
          open Route
          Web.body(Web.document)->Web.setClassName("bodmob bodtab bodbig")
          <nav
            className="relative font-anon text-sm grid grid-cols-2 grid-rows-[16fr,10fr,16fr,10fr,8fr,4fr,6fr]">
            <Link
              route=SignIn
              className={"w-3/5 border border-stone-100 bg-stone-800/40 text-center text-stone-100 " ++ "decay-mask text-3xl p-2 max-w-80 font-fred col-span-full justify-self-center"}
              content="SIGN IN"
            />
            <Link
              route=SignUp
              className={"w-3/5 border border-stone-100 bg-stone-800/40 text-center text-stone-100 " ++ "decay-mask text-3xl p-2 max-w-80 font-fred col-span-full justify-self-center row-start-3"}
              content="SIGN UP"
            />
            <Link
              route=GetInfo({search: ForgotPassword})
              className="text-stone-100 row-start-5 mr-6 justify-self-end"
              content="forgot password?"
            />
            <Link
              route=GetInfo({search: ForgotUsername})
              className="text-stone-100 row-start-5 ml-6"
              content="forgot username?"
            />
            <Link
              route=GetInfo({search: VerificationCode})
              className="text-stone-100 col-span-full justify-self-center row-start-7"
              content="have code?"
            />
            {switch wsError == "" {
            | true => React.null
            | false => <React.Suspense fallback=React.null> msg </React.Suspense>
            }}
          </nav>
        }

      | (SignIn, None) => <React.Suspense fallback=React.null> signin </React.Suspense>
      | (SignUp, None) => <React.Suspense fallback=React.null> signup </React.Suspense>

      | (GetInfo({search}), None) =>
        <React.Suspense fallback=React.null>
          {React.createElement(
            GetInfo.lazy_(() =>
              GetInfo.import_("./GetInfo.bs")->Promise.then(comp => {
                Promise.resolve({"default": comp["make"]})
              })
            ),
            GetInfo.makeProps(
              ~userpool,
              ~cognitoUser,
              ~setCognitoUser,
              ~setRetrievedUsername,
              ~search,
              (),
            ),
          )}
        </React.Suspense>

      | (Confirm({search}), None) =>
        <React.Suspense fallback=React.null>
          {React.createElement(
            Confirm.lazy_(() =>
              Confirm.import_("./Confirm.bs")->Promise.then(comp => {
                Promise.resolve({"default": comp["make"]})
              })
            ),
            Confirm.makeProps(~cognitoUser, ~search, ()),
          )}
        </React.Suspense>
      | (Lobby | Play(_) | Leaderboard, None) => {
          Route.replace(Home)
          React.null
        }

      | (Home | SignIn | SignUp | GetInfo(_) | Confirm(_), Some(_)) => {
          Route.replace(Lobby)
          React.null
        }

      | (Lobby | Play(_) | Leaderboard, Some(_)) =>
        <React.Suspense fallback=React.null> auth </React.Suspense>
      | (Other, _) =>
        <div className="text-center text-stone-100 text-4xl">
          {React.string("page not found")}
        </div>
      }}
    </main>
    {switch (route, token) {
    | (Home, None) =>
      <footer>
        <a
          href="https://github.com/jamessouth/aws-clean-tablet"
          className="w-7 h-7 block m-auto"
          rel="noopener noreferrer">
          <svg
            className="w-7 h-7 fill-stone-100 absolute"
            viewBox="0 0 32 32"
            xmlns="http://www.w3.org/2000/svg">
            <path
              d="M16 2a14 14 0 0 0-4.43 27.28c.7.13 1-.3 1-.67v-2.38c-3.89.84-4.71-1.88-4.71-1.88a3.71 3.71 0 0 0-1.62-2.05c-1.27-.86.1-.85.1-.85a2.94 2.94 0 0 1 2.14 1.45a3 3 0 0 0 4.08 1.16a2.93 2.93 0 0 1 .88-1.87c-3.1-.36-6.37-1.56-6.37-6.92a5.4 5.4 0 0 1 1.44-3.76a5 5 0 0 1 .14-3.7s1.17-.38 3.85 1.43a13.3 13.3 0 0 1 7 0c2.67-1.81 3.84-1.43 3.84-1.43a5 5 0 0 1 .14 3.7a5.4 5.4 0 0 1 1.44 3.76c0 5.38-3.27 6.56-6.39 6.91a3.33 3.33 0 0 1 .95 2.59v3.84c0 .46.25.81 1 .67A14 14 0 0 0 16 2Z"
            />
          </svg>
        </a>
      </footer>
    | (
      Leaderboard | SignIn | SignUp | GetInfo(_) | Confirm(_) | Lobby | Play(_) | Other,
      None | Some(_),
    )
    | (Home, Some(_)) => React.null
    }}
  </>
}
