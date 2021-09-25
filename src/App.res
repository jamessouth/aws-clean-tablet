
@react.component
let make = () => {

// Js.Nullable.t<Signup.usr>

    let (cognitoUser, setCognitoUser) = React.useState(_ => None)

    let url = RescriptReactRouter.useUrl()

    let {token} = AuthHook.useAuth()
    <>
        <h1 className="text-6xl mt-11 text-center font-arch decay-mask text-warm-gray-100">{"CLEAN TABLET"->React.string}</h1>


        <div className="mt-10 sm:mt-20">

            {

                switch (url.path, token) {
                    | (list{}, Some(_t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }
                    | (list{}, None) => <div className="flex flex-col items-center">

                    <Link url="/signin" className="w-3/5 border border-warm-gray-100 block font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 mb-8 max-w-80 sm:mb-16" content="SIGN IN"/>


                    <Link url="/signup" className="w-3/5 border border-warm-gray-100 block font-fred text-center text-warm-gray-100 decay-mask text-3xl p-2 max-w-80" content="SIGN UP"/>


                    <Link url="/confirm" className="w-3/5 text-center text-warm-gray-100 block font-anon text-sm mt-4 max-w-80" content="Verification code?"/>


                    <Link url="/leaderboards" className="w-3/5 border border-warm-gray-100 text-center text-warm-gray-100 block font-anon text-xl mt-40 max-w-80" content="Leaderboards"/>
                    
                    </div>

                    | (list{"leaderboards"}, _) => <div>{"leaderboard"->React.string}</div>

                    | (list{"signin"}, Some(_t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }

                    | (list{"signin"}, None) => <Signin/>


                    | (list{"confirm"}, Some(_t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }

                    | (list{"confirm"}, None) => <Confirm cognitoUser/>




                    | (list{"signup"}, Some(_t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }

                    | (list{"signup"}, None) => <Signup setCognitoUser/>


                    | (list{"resetpwd"}, Some(_t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }

                    | (list{"resetpwd"}, None) => <Resetpwd/>

                    | (list{"lobby"}, Some(_t)) => <Lobby/>

                    | (list{"lobby"}, None) => {
                        RescriptReactRouter.replace("/login")
                        React.null
                        }

                    // | (list{"game", _gameno}, Some(_t)) => <Play/>

                    | (list{"game", _gameno}, None) => {
                        RescriptReactRouter.replace("/login")
                        React.null
                        }


                    | (_, _) => <div>{"other"->React.string}</div>// <PageNotFound/>
                }
            }


        </div>
    </>
}