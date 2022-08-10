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

  React.useEffect0(() => {
    let tok = RescriptReactRouter.watchUrl(r => Js.log2("waa", r))
    Some(() => {RescriptReactRouter.unwatchUrl(tok)})
  })

  let initialState: Reducer.state = {
    gamesList: Js.Nullable.null,
    players: [],
    sk: "",
    oldWord: "",
    word: "",
    showAnswers: false,
    winner: "",
  }

  let (
    playerGame,
    playerName,
    playerColor,
    endtoken,
    count,
    wsConnected,
    players,
    sk,
    showAnswers,
    winner,
    oldWord,
    word,
    games,
    leaderData,
    send,
    resetConnState,
    close,
    wsError,
  ) = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, initialState)

  let loading1 = React.createElement(
    Loading.lazy_(() =>
      Loading.import_("./Loading.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Loading.makeProps(~label="games...", ()),
  )

  let loading2 = React.createElement(
    Loading.lazy_(() =>
      Loading.import_("./Loading.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Loading.makeProps(~label="game...", ()),
  )

 

  let signin = React.createElement(
    Signin.lazy_(() =>
      Signin.import_("./Signin.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Signin.makeProps(
  ~userpool=userpool,
  ~setCognitoUser=setCognitoUser,
  ~setToken=setToken,
  ~cognitoUser=cognitoUser,
  ~cognitoError=cognitoError,
  ~setCognitoError=setCognitoError,
  ()),
  )

  let signup = React.createElement(
    Signup.lazy_(() =>
      Signup.import_("./Signup.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
Signup.makeProps(
  ~cognitoError=cognitoError,
  ~setCognitoError=setCognitoError,
  ~setCognitoUser=setCognitoUser,
  ~userpool=userpool,
  ())


  )


// 66


  open Web
  <>
    {switch path {
    | list{"leaderboard"} => React.null
    | _ =>
      <header className="mb-10 newgmimg:mb-12">
        <p className="font-flow text-stone-100 text-4xl h-10 font-bold text-center">
          {React.string(playerName)}
        </p>
        <h1
          style={ReactDOM.Style.make(~backgroundColor={playerColor}, ())}
          className="text-6xl mt-11 mx-auto px-6 text-center font-arch decay-mask text-stone-100">
          {React.string("CLEAN TABLET")}
        </h1>
      </header>
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

      | (list{"signin"}, None) =>
        

        <React.Suspense fallback=React.null> signin </React.Suspense>

      | (list{"signup"}, None) => <React.Suspense fallback=React.null> signup </React.Suspense>

      | (list{"getinfo"}, None) =>
        switch search {
        | "cd_un" | "pw_un" | "un_em" =>
          <GetInfo
            userpool cognitoUser setCognitoUser cognitoError setCognitoError setShowName search
          />
        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (list{"confirm"}, None) =>
        switch search {
        | "cd_un" | "pw_un" => <Confirm cognitoUser cognitoError setCognitoError search />
        | _ =>
          <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
        }

      | (list{"lobby"}, None) | (list{"game"}, None) | (list{"leaderboard"}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      | (list{}, Some(_))
      | (list{"signin"}, Some(_))
      | (list{"signup"}, Some(_))
      | (list{"getinfo"}, Some(_))
      | (list{"confirm"}, Some(_)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"lobby"}, Some(_)) =>
        switch wsConnected {
        | false => {
            body(document)->setClassName("bodchmob bodchtab bodchbig")
            <React.Suspense fallback=React.null> loading1 </React.Suspense>
          }

        | true => {
            body(document)->classList->removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")
            <Lobby playerGame games send wsError close count />
          }
        }
      | (list{"game", gameno}, Some(_)) =>
        switch wsConnected {
        | true =>
          switch Js.Array2.length(players) > 0 && gameno == sk {
          | true =>
            <Play
              players
              sk
              showAnswers
              winner
              isWinner={winner != ""}
              oldWord
              word
              playerColor
              send
              playerName
              endtoken
              resetConnState
            />
          | false => <React.Suspense fallback=React.null> loading2 </React.Suspense>
          }

        | false =>
          <p className="text-center text-stone-100 font-anon text-lg">
            {React.string("not connected...")}
          </p>
        }

      | (list{"leaderboard"}, Some(_)) => {
          body(document)->classList->addClassList3("bodleadmob", "bodleadtab", "bodleadbig")
          <Leaders send leaderData playerName />
        }

      | (_, _) => <div> {React.string("other")} </div> // <PageNotFound/>
      }}
    </main>
  </>
}
