

@react.component
let make = () => {

    let url = RescriptReactRouter.useUrl()

    let {token} = AuthHook.useAuth()
    <>
        <h1 className="text-6xl mt-11 text-center font-arch decay-mask">{"CLEAN TABLET"->React.string}</h1>


        <div className="mt-10 sm:mt-20">

            {

                switch (url.path, token) {
                    | (list{}, Some(t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }
                    | (list{}, None) => <div className="flex flex-col items-center">

                    <Link url="/signin" className="w-3/5 border border-smoke-100 block font-fred decay-mask text-3xl p-2 mb-8 max-w-80 sm:mb-16" content="SIGN IN"/>

                    // <a onClick=onClick("/signin")  href="/signin">{->React.string}</a>
                    <a className="w-3/5 border border-smoke-100 block font-fred decay-mask text-3xl p-2 max-w-80" href="/login">{"SIGN UP"->React.string}</a>


                    <a className="w-3/5 border border-smoke-100 block text-xl mt-40 max-w-80" href="/leaderboards">{"Leaderboards"->React.string}</a>
                    
                    </div>

                    | (list{"leaderboards"}, _) => <div>{"leaderboard"->React.string}</div>

                    | (list{"signin"}, Some(t)) => {
                        RescriptReactRouter.replace("/lobby")
                        React.null
                        }

                    | (list{"signin"}, None) => <Signin/>


                    // | (list{"login"}, Some(t)) => {
                    //     RescriptReactRouter.replace("/lobby")
                    //     React.null
                    //     }

                    // | (list{"login"}, None) => <LoginPage/>

                    | (list{"lobby"}, Some(t)) => <Lobby/>

                    | (list{"lobby"}, None) => {
                        RescriptReactRouter.replace("/login")
                        React.null
                        }

                    | (list{"game", gameno}, Some(t)) => <Play/>

                    | (list{"game", gameno}, None) => {
                        RescriptReactRouter.replace("/login")
                        React.null
                        }


                    | (_, _) => <div>{"other"->React.string}</div>// <PageNotFound/>
                }
            }


        </div>
    </>
}