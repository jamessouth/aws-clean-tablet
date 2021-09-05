

@react.component
let make = () => {

    let url = RescriptReactRouter.useUrl()

    let {token} = AuthHook.useAuth()

    <div className="mt-8">

        {

            switch (url.path, token) {
                | (list{}, Some(t)) => {
                    RescriptReactRouter.replace("/lobby")
                    React.null
                    }
                | (list{}, None) => <div className="flex flex-col items-center">
                <a className="w-3/5 border border-smoke-100 block font-fred decay-mask text-5xl leading-12rem sm:mt-16 sm:text-8xl sm:leading-12rem" href="/login">{"ENTER"->React.string}</a>
                <a className="w-3/5 border border-smoke-100 mb-28 mt-10 block text-xl sm:mt-16 sm:text-2xl" href="/leaderboards">{"Leaderboards"->React.string}</a>
                
                </div>

                | (list{"leaderboards"}, _) => <div>{"leaderboard"->React.string}</div>

                | (list{"login"}, Some(t)) => {
                    RescriptReactRouter.replace("/lobby")
                    React.null
                    }

                | (list{"login"}, None) => <LoginPage/>

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

}