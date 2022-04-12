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

  let initialState: Reducer.state = {
    gamesList: Js.Nullable.null,
    game: {
      sk: "",
      players: [],
      currentWord: "",
      previousWord: "",
      showAnswers: false,
      winner: "",
    },
  }

  let (
    playerGame,
    playerName,
    playerColor,
    playerIndex,
    wsConnected,
    game,
    games,
    leader,
    send,
    close,
    wsError,
  ) = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, initialState)

  <>
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
    <main className="mb-8">
      {switch (path, token) {
      | (list{}, None) =>
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
          <Link url="/getinfo?cd_un" className={linkBase ++ "mt-8"} content="verification code?" />
          <Link url="/getinfo?pw_un" className={linkBase ++ "mt-4"} content="forgot password?" />
          <Link url="/getinfo?un_em" className={linkBase ++ "mt-4"} content="forgot username?" />
          <Link
            url="/leaderboards"
            className={linkBase2 ++ "font-anon text-xl mt-20 max-w-80"}
            content="Leaderboards"
          />
          {switch showName == "" {
          | true => React.null
          | false =>
            <p className="text-stone-100 absolute -top-20 w-4/5 bg-blue-gray-800 p-2 font-anon">
              {React.string("The username associated with the email you submitted is:" ++ showName)}
            </p>
          }}
        </nav>

      | (list{"signin"}, None) =>
        <Signin userpool setCognitoUser setToken cognitoUser cognitoError setCognitoError />

      | (list{"signup"}, None) => <Signup userpool setCognitoUser cognitoError setCognitoError />

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

      | (list{"lobby"}, None) | (list{"game"}, None) => {
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
        | false => <Loading label="games..." />
        | true => <Lobby playerGame leader games send wsError close />
        }

      | (list{"game", gameno}, Some(_)) =>
        switch wsConnected {
        | true =>
          switch gameno == game.sk {
          | true => <Play game playerColor playerIndex send leader playerName />
          | false => <Loading label="game..." />
          }

        | false =>
          <p className="text-center text-stone-100 font-anon text-lg">
            {React.string("not connected...")}
          </p>
        }

      | (list{"leaderboards"}, _) => <div> {React.string("leaderboard")} </div>

      | (_, _) => <div> {React.string("other")} </div> // <PageNotFound/>
      }}
    </main>
  </>
}
