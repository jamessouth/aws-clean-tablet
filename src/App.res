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

@react.component
let make = () => {
  Js.log("app")
  let url = RescriptReactRouter.useUrl()
  

  let (cognitoUser: Js.Nullable.t<Signup.usr>, setCognitoUser) = React.Uncurried.useState(_ =>
    Js.Nullable.null
  )

  let (playerName, setPlayerName) = React.Uncurried.useState(_ => "")
  
  
  let (token, setToken) = React.Uncurried.useState(_ => None)


  React.useEffect1(() => {
    switch Js.Nullable.toOption(cognitoUser) {
    | None => setPlayerName(._ => "")
    | Some(user) => setPlayerName(._ => user.username)
    }
  None
}, [cognitoUser])

  let {
    playerGame,
    playerColor,
    wsConnected,
    game,
    games,
    currentWord,
    previousWord,
    connID,
    send,
    close,
    wsError,
    setWs
  } = WsHook.useWs(token, setToken)
let a: array<Reducer.livePlayer> = [{
  name: "bill",
  connid: "111",
  color: "red",
  score: "12",
  answer: {playerid: "1", answer: ""},
  hasAnswered: true,
}, {
  name: "carl",
  connid: "222",
  color: "blue",
  score: "13",
  answer: {playerid: "2", answer: ""},
  hasAnswered: true,
}, {
  name: "wes",
  connid: "333",
  color: "green",
  score: "1",
  answer: {playerid: "3", answer: ""},
  hasAnswered: false,
}]
  <>
<Scoreboard players=a currentWord="fret" previousWord="" />
    <SignOut cognitoUser setToken send playerGame close setCognitoUser setPlayerName setWs/>
    <h1 className="text-6xl mt-11 text-center font-arch decay-mask text-warm-gray-100">
      {"CLEAN TABLET"->React.string}
    </h1>
    <div className="mt-10 sm:mt-20">
      {switch (url.path, token) {
      | (list{}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }
      | (list{}, None) =>
        <div className="flex flex-col items-center">
          <Link
            url="/signin"
            className="w-3/5 border border-warm-gray-100 block font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 mb-8 max-w-80 sm:mb-16"
            content="SIGN IN"
          />
          <Link
            url="/signup"
            className="w-3/5 border border-warm-gray-100 block font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 max-w-80"
            content="SIGN UP"
          />
          <Link
            url="/getusername"
            className="w-3/5 text-center text-warm-gray-100 block font-anon text-sm mt-4 max-w-80"
            content="verification code?"
          />
          <Link
            url="/leaderboards"
            className="w-3/5 border border-warm-gray-100 text-center text-warm-gray-100 block font-anon text-xl mt-40 max-w-80"
            content="Leaderboards"
          />
        </div>

      | (list{"leaderboards"}, _) => <div> {"leaderboard"->React.string} </div>

      | (list{"signin"}, Some(_)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"signin"}, None) => <Signin userpool setCognitoUser setToken cognitoUser/>

      | (list{"confirm"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"confirm"}, None) => <Confirm cognitoUser />

      | (list{"getusername"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"getusername"}, None) => <GetUsername userpool setCognitoUser />

      | (list{"signup"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"signup"}, None) => <Signup userpool setCognitoUser />

      | (list{"resetpwd"}, Some(_t)) => {
          RescriptReactRouter.replace("/lobby")
          React.null
        }

      | (list{"resetpwd"}, None) => <Resetpwd />

      | (list{"lobby"}, Some(_)) => <Lobby wsConnected playerGame leadertoken=(playerName ++ connID) games send wsError/>

      | (list{"lobby"}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

      | (list{"game", _gameno}, None) => {
          RescriptReactRouter.replace("/")
          React.null
        }

// playerName
      | (list{"game", _gameno}, Some(_)) => <Play wsConnected game playerColor send wsError currentWord previousWord/>

      | (_, _) => <div> {"other"->React.string} </div> // <PageNotFound/>
      }}
    </div>
  </>
}
