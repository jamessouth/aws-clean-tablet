@val external localStorage: Dom.Storage2.t = "localStorage"
external getItem: (Dom.Storage2.t, string) => option<string> = "getItem"
external setItem: (Dom.Storage2.t, string, string) => unit = "setItem"


let useAuth = () => {

    let userinit = switch localStorage->Dom.Storage2.getItem("token") {
    | Some(t) => t
    | None => ""
    }

    let (token, setToken) = React.useState(_ => userinit)

    let authinit = switch user {
        | "" => "signin"
        | _ => "signedin"
        }


    let (authState, setAuthState) = React.useState(_ => authinit)

    


}