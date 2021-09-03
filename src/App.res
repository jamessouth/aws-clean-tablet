

@react.component
let make = () => {

    let url = RescriptReactRouter.useUrl()
    <div className="mt-8">

        {

            switch url.path {
                | list{"/"} => {
                    RescriptReactRouter.replace("/lobby")
                    React.null
                    }
                | list{"/else"} => <div className="flex flex-col items-center">
                <a className="w-3/5 border border-smoke-100 block font-fred decay-mask text-5xl leading-12rem sm:mt-16 sm:text-8xl sm:leading-12rem" href="/lobby">{"ENTER"->React.string}</a>
                <a className="w-3/5 border border-smoke-100 mb-28 mt-10 block text-xl sm:mt-16 sm:text-2xl" href="/leaderboards">{"Leaderboards"->React.string}</a>
                
                </div>
                | list{"/leaderboards"} => <div>{"leaderboard"->React.string}</div>
                | list{"/login"} => <LoginPage/>





                | list{"/lobby"} => <Lobby/>
                | list{"/game/:gameno"} => <Play/>
                | _ => <div>{"other"->React.string}</div>// <PageNotFound/>
            }
        }


    </div>
}