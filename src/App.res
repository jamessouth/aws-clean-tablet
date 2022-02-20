@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"

@val external localStorage: Dom.Storage2.t = "localStorage"
external getItem: (Dom.Storage2.t, string) => option<string> = "getItem"

// @module external signoutimg: string =

@new @module("amazon-cognito-identity-js")
external userPoolConstructor: Types.poolDataInput => Types.poolData = "CognitoUserPool"

let pool: Types.poolDataInput = {
  userPoolId: upid,
  clientId: cid,
  advancedSecurityDataCollectionFlag: false,
}
let userpool = userPoolConstructor(pool)







    let checkLength = (min, max, str) =>
    switch Js.String2.length(str) < min || Js.String2.length(str) > max {
    | false => ""
    | true => j`$min-$max characters; `
    }
  let checkInclusion = (re, msg, str) =>
    switch Js.String2.match_(str, re) {
    | None => msg
    | Some(_) => ""
    }
  let checkExclusion = (re, msg, str) =>
    switch Js.String2.match_(str, re) {
    | None => ""
    | Some(_) => msg
    }


  let usernameFuncList = [
    s => checkLength(3, 10, s),
    s =>
      checkExclusion(
        %re("/\W/"),
        "letters, numbers, and underscores only; no whitespace or symbols.",
        s,
      ),
  ]

  let passwordFuncList = [
    s => checkLength(8, 98, s),
    s => checkInclusion(%re("/[!-/:-@\[-`{-~]/"), "at least 1 symbol; ", s),
    s => checkInclusion(%re("/\d/"), "at least 1 number; ", s),
    s => checkInclusion(%re("/[A-Z]/"), "at least 1 uppercase letter; ", s),
    s => checkInclusion(%re("/[a-z]/"), "at least 1 lowercase letter; ", s),
    s => checkExclusion(%re("/\s/"), "no whitespace.", s),
  ]

  let emailFuncList = [
    s => checkLength(5, 99, s),
    s =>
      checkInclusion(
        %re(
          "/^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/"
        ),
        "enter a valid email address.",
        s,
      ),
  ]


@react.component
let make = () => {
  Js.log("app")
  let url = RescriptReactRouter.useUrl()

  let (cognitoUser: Js.Nullable.t<Signup.usr>, setCognitoUser) = React.Uncurried.useState(_ =>
    Js.Nullable.null
  )
  let (cognitoError, setCognitoError) = React.useState(_ => None)

  let (playerName, setPlayerName) = React.Uncurried.useState(_ => "")

  let (token, setToken) = React.Uncurried.useState(_ => None)
  let (showName, setShowName) = React.Uncurried.useState(_ => "")


  React.useEffect1(() => {
    switch Js.Nullable.toOption(cognitoUser) {
    | None => setPlayerName(._ => "")
    | Some(user) => setPlayerName(._ => user.username)
    }
    None
  }, [cognitoUser])

  let {
    playerGame,
    // setPlayerGame,
    playerColor,
    wsConnected,
    game,
    games,
    connID,
    // setConnID,
    send,
    close,
    wsError,
    // setWs,
    // dispatch
  } = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, setPlayerName)

  let signOut = <SignOut send playerGame close />

  <>
    <p className="font-flow text-warm-gray-100 text-4xl h-10 font-bold text-center">{React.string(playerName)}</p>
    <h1 style={ReactDOM.Style.make(~backgroundColor={playerColor}, ())} className="text-6xl mt-11 mx-auto w-11/12 text-center font-arch decay-mask text-warm-gray-100">
      {"CLEAN TABLET"->React.string}
    </h1>
    <div className="mt-10">
      {switch (url.path, token) {
      | (list{}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }
      | (list{}, None) =>
        <div className="flex flex-col items-center relative">
          <Link
            url="/signin"
            className="w-3/5 border border-warm-gray-100 block bg-warm-gray-800/40 font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 mb-8 max-w-80 sm:mb-16"
            content="SIGN IN"
          />
          <Link
            url="/signup"
            className="w-3/5 border border-warm-gray-100 block bg-warm-gray-800/40 font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 max-w-80"
            content="SIGN UP"
          />
          <Link
            url="/getinfo?cd_un"
            className="w-3/5 text-warm-gray-100 block font-bold font-anon text-sm mt-4 max-w-80"
            content="verification code?"
          />
          <Link
            url="/getinfo?pw_un"
            className="w-3/5 text-warm-gray-100 block font-bold font-anon text-sm mt-4 max-w-80"
            content="forgot password?"
          />
          <Link
            url="/getinfo?un_em"
            className="w-3/5 text-warm-gray-100 block font-bold font-anon text-sm mt-4 max-w-80"
            content="forgot username?"
          />
          <Link
            url="/leaderboards"
            className="w-3/5 border border-warm-gray-100 text-center text-warm-gray-100 block bg-warm-gray-800/40 font-anon text-xl mt-24 max-w-80"
            content="Leaderboards"
          />
          {switch showName == "" {
          | true => React.null
          | false => <p className="text-warm-gray-100 absolute -top-20 w-4/5 bg-blue-gray-800 p-2 font-anon">{React.string("The username associated with the email you submitted is:" ++ showName)}</p>
          }}
        </div>

      | (list{"leaderboards"}, _) => <div> {"leaderboard"->React.string} </div>

      | (list{"signin"}, Some(_)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"signin"}, None) =>
        <Signin userpool setCognitoUser setToken cognitoUser cognitoError setCognitoError />

      | (list{"confirm"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"confirm"}, None) => <Confirm cognitoUser cognitoError setCognitoError />

      | (list{"getinfo"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"getinfo"}, None) =>
        <GetInfo userpool cognitoUser setCognitoUser cognitoError setCognitoError usernameFuncList emailFuncList setShowName/>

      | (list{"signup"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"signup"}, None) => <Signup userpool setCognitoUser cognitoError setCognitoError usernameFuncList passwordFuncList emailFuncList/>

      | (list{"lobby"}, Some(_)) =>
        <Lobby
          wsConnected playerGame leadertoken={playerName ++ connID} games send wsError signOut
        />

      | (list{"lobby"}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      | (list{"game", _gameno}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      // playerName
      | (list{"game", _gameno}, Some(_)) => <Play wsConnected game playerColor send wsError leadertoken={playerName ++ connID}/>

      | (_, _) => <div> {"other"->React.string} </div> // <PageNotFound/>
      }}
    </div>
  </>
}
