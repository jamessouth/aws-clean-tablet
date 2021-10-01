@val external localStorage: Dom.Storage2.t = "localStorage"
external getItem: (Dom.Storage2.t, string) => option<string> = "getItem"
external setItem: (Dom.Storage2.t, string, string) => unit = "setItem"

type returnVal = {
    setToken: string => unit,
    token: option<string>
}

let useAuth = () => {
Js.log("authhook")
    let storedToken = localStorage->Dom.Storage2.getItem("token")
    
    
    
    // let userinit = switch storedToken {
    // | Some(t) => t
    // | None => ""
    // }

    let (token, setToken) = React.useState(_ => storedToken)


    let saveToken = token => {
        localStorage->Dom.Storage2.setItem("token", token)
        (_ => Some(token))->setToken
    }

    let return = {
        setToken: saveToken,
        token
    }

    return


}

