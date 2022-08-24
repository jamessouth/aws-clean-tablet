@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"

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

  let {path, search} = RescriptReactRouter.useUrl()
  let (cognitoUser: Js.Nullable.t<usr>, setCognitoUser) = React.Uncurried.useState(_ =>
    Js.Nullable.null
  )
  let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)

  let (token, setToken) = React.Uncurried.useState(_ => None)
  let (showName, setShowName) = React.Uncurried.useState(_ => "")

  // 66
  // html - 1
  // css - 7
  // js - 46
  // xhr - 4
  // font - 4
  // img - 6
  // ws - 1

  let signin = React.createElement(
    Signin.lazy_(() =>
      Signin.import_("./Signin.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Signin.makeProps(
      ~userpool,
      ~setCognitoUser,
      ~setToken,
      ~cognitoUser,
      ~cognitoError,
      ~setCognitoError,
      (),
    ),
  )

  let signup = React.createElement(
    Signup.lazy_(() =>
      Signup.import_("./Signup.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Signup.makeProps(~cognitoError, ~setCognitoError, ~setCognitoUser, ~userpool, ()),
  )

  let getInfo = React.createElement(
    GetInfo.lazy_(() =>
      GetInfo.import_("./GetInfo.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    GetInfo.makeProps(
      ~userpool,
      ~cognitoUser,
      ~setCognitoUser,
      ~cognitoError,
      ~setCognitoError,
      ~setShowName,
      ~search,
      (),
    ),
  )

  let confirm = React.createElement(
    Confirm.lazy_(() =>
      Confirm.import_("./Confirm.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Confirm.makeProps(~cognitoUser, ~cognitoError, ~setCognitoError, ~search, ()),
  )

  open Web
  <>
    {switch path {
    | list{"leaderboard"} => React.null
    | _ =>
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
      className={switch path {
      | list{"leaderboard"} => ""
      | _ => "mb-8"
      }}>
      {switch (path, token) {
      | (list{}, None) => {
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

      | (list{"signin"}, None) => <React.Suspense fallback=React.null> signin </React.Suspense>

      | (list{"signup"}, None) => <React.Suspense fallback=React.null> signup </React.Suspense>

      | (list{"getinfo"}, None) =>
        switch search {
        | "cd_un" | "pw_un" | "un_em" =>
          <React.Suspense fallback=React.null> getInfo </React.Suspense>

        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (list{"confirm"}, None) =>
        switch search {
        | "cd_un" | "pw_un" => <React.Suspense fallback=React.null> confirm </React.Suspense>

        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (list{"auth", ..._}, None) | (list{"game"}, None) | (list{"leaderboard"}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      | (list{}, Some(_))
      | (list{"signin"}, Some(_))
      | (list{"signup"}, Some(_))
      | (list{"getinfo"}, Some(_))
      | (list{"confirm"}, Some(_)) => {
          RescriptReactRouter.replace("/auth/lobby")
          React.null
        }

      | (list{"auth", ...subpath}, Some(_)) =>
        <React.Suspense fallback=React.null>
          {React.createElement(
            Auth.lazy_(() =>
              Auth.import_("./Auth.bs")->Promise.then(comp => {
                Promise.resolve({"default": comp["make"]})
              })
            ),
            Auth.makeProps(~token, ~setToken, ~cognitoUser, ~setCognitoUser, ~subpath, ()),
          )}
        </React.Suspense>

      | (_, _) => <div> {React.string("other")} </div> // <PageNotFound/>
      }}
    </main>
  </>
}
