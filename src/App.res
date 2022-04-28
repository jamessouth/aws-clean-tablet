@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"

@react.component
let make = () => {
  Js.log("app")

  // let linkBase = "w-3/5 text-stone-100 block font-bold font-anon text-sm max-w-80 "

  // let linkBase2 = "w-3/5 border border-stone-100 block bg-stone-800/40 text-center text-stone-100 "

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
  // let (cognitoError, setCognitoError) = React.Uncurried.useState(_ => None)

  // let (token, setToken) = React.Uncurried.useState(_ => None)
  // let (showName, setShowName) = React.Uncurried.useState(_ => "")

  // let initialState: Reducer.state = {
  //   gamesList: Js.Nullable.null,
  //   game: {
  //     sk: "",
  //     players: [],
  //     currentWord: "",
  //     previousWord: "",
  //     showAnswers: false,
  //     winner: "",
  //   },
  // }

  // let (
  //   playerGame,
  //   playerName,
  //   playerColor,
  //   playerIndex,
  //   wsConnected,
  //   game,
  //   games,
  //   leader,
  //   leaderData,
  //   send,
  //   close,
  //   wsError,
  // ) = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, initialState)
  let zzz: array<Reducer.stat> = [
    {name: "mmmmmmmmmm", wins: 241, totalPoints: 5152, games: 83, winPct: 0.12, ppg: 11.0},
    {name: "stu", wins: 121, totalPoints: 121, games: 32, winPct: 0.22, ppg: 11.1},
    {name: "liz", wins: 50, totalPoints: 59, games: 363, winPct: 0.32, ppg: 11.2},
    {name: "abner", wins: 42, totalPoints: 18, games: 173, winPct: 0.42, ppg: 11.3},
    {name: "harold", wins: 40, totalPoints: 97, games: 333, winPct: 0.52, ppg: 11.4},
    {name: "stacie", wins: 32, totalPoints: 17, games: 313, winPct: 0.62, ppg: 11.5},
    {name: "marcy", wins: 30, totalPoints: 91, games: 332, winPct: 0.72, ppg: 11.6},
    {name: "wes", wins: 22, totalPoints: 11, games: 213, winPct: 0.82, ppg: 11.7},
    {name: "carl", wins: 21, totalPoints: 12, games: 23, winPct: 0.92, ppg: 11.8},
    {name: "bill", wins: 10, totalPoints: 9, games: 323, winPct: 0.02, ppg: 11.9},
    {name: "test", wins: 2, totalPoints: 1, games: 13, winPct: 0.17, ppg: 12.8},
  ]

  





  <Leaders leaderData=zzz  />
  // <>
  //   <header className="mb-10 newgmimg:mb-12">
  //     <p className="font-flow text-stone-100 text-4xl h-10 font-bold text-center">
  //       {React.string(playerName)}
  //     </p>
  //     <h1
  //       style={ReactDOM.Style.make(~backgroundColor={playerColor}, ())}
  //       className="text-6xl mt-11 mx-auto px-6 text-center font-arch decay-mask text-stone-100">
  //       {React.string("CLEAN TABLET")}
  //     </h1>
  //   </header>
  //   <main className="mb-8">
  //     {switch (path, token) {
  //     | (list{}, None) =>
  //       <nav className="flex flex-col items-center relative">
  //         <Link
  //           url="/signin"
  //           className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred mb-8 sm:mb-16"}
  //           content="SIGN IN"
  //         />
  //         <Link
  //           url="/signup"
  //           className={linkBase2 ++ "decay-mask text-3xl p-2 max-w-80 font-fred"}
  //           content="SIGN UP"
  //         />
  //         <Link url="/getinfo?cd_un" className={linkBase ++ "mt-10"} content="verification code?" />
  //         <Link url="/getinfo?pw_un" className={linkBase ++ "mt-6"} content="forgot password?" />
  //         <Link url="/getinfo?un_em" className={linkBase ++ "mt-6"} content="forgot username?" />
  //         {switch showName == "" {
  //         | true => React.null
  //         | false =>
  //           <p className="text-stone-100 absolute -top-20 w-4/5 bg-blue-gray-800 p-2 font-anon">
  //             {React.string("The username associated with the email you submitted is:" ++ showName)}
  //           </p>
  //         }}
  //       </nav>

  //     | (list{"signin"}, None) =>
  //       <Signin userpool setCognitoUser setToken cognitoUser cognitoError setCognitoError />

  //     | (list{"signup"}, None) => <Signup userpool setCognitoUser cognitoError setCognitoError />

  //     | (list{"getinfo"}, None) =>
  //       switch search {
  //       | "cd_un" | "pw_un" | "un_em" =>
  //         <GetInfo
  //           userpool cognitoUser setCognitoUser cognitoError setCognitoError setShowName search
  //         />
  //       | _ =>
  //         <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
  //       }

  //     | (list{"confirm"}, None) =>
  //       switch search {
  //       | "cd_un" | "pw_un" => <Confirm cognitoUser cognitoError setCognitoError search />
  //       | _ =>
  //         <div className="text-stone-100"> {React.string("unknown path, please try again")} </div>
  //       }

  //     | (list{"lobby"}, None) | (list{"game"}, None) | (list{"leaderboard"}, None) => {
  //         RescriptReactRouter.replace("/")
  //         React.null
  //       }

  //     | (list{}, Some(_))
  //     | (list{"signin"}, Some(_))
  //     | (list{"signup"}, Some(_))
  //     | (list{"getinfo"}, Some(_))
  //     | (list{"confirm"}, Some(_)) => {
  //         RescriptReactRouter.replace("/lobby")
  //         React.null
  //       }

  //     | (list{"lobby"}, Some(_)) =>
  //       switch wsConnected {
  //       | false => <Loading label="games..." />
  //       | true => <Lobby playerGame leader games send wsError close />
  //       }

  //     | (list{"game", gameno}, Some(_)) =>
  //       switch wsConnected {
  //       | true =>
  //         switch gameno == game.sk {
  //         | true => <Play game playerColor playerIndex send leader playerName />
  //         | false => <Loading label="game..." />
  //         }

  //       | false =>
  //         <p className="text-center text-stone-100 font-anon text-lg">
  //           {React.string("not connected...")}
  //         </p>
  //       }

  //     | (list{"leaderboard"}, Some(_)) => <Leaders send leaderData/>

  //     | (_, _) => <div> {React.string("other")} </div> // <PageNotFound/>
  //     }}
  //   </main>
  // </>
}
