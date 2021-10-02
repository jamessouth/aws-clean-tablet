@val external localStorage: Dom.Storage2.t = "localStorage"
external getItem: (Dom.Storage2.t, string) => option<string> = "getItem"
external setItem: (Dom.Storage2.t, string, string) => unit = "setItem"
external removeItem: (Dom.Storage2.t, string) => unit = "removeItem"
external key: (Dom.Storage2.t, int) => option<string> = "key"
external length: Dom.Storage2.t => int = "length"

type returnVal = {
  setToken: string => unit,
  clearCognitoKeys: (. unit) => unit,
  token: option<string>,
}

let useAuth = () => {
  Js.log("authhook")
  let storedToken = localStorage->Dom.Storage2.getItem("token")

  let (token, setToken) = React.Uncurried.useState(_ => storedToken)

  let saveToken = token => {
    localStorage->Dom.Storage2.setItem("token", token)
    setToken(. _ => Some(token))
  }

  let clearCognitoKeys = (. ()) => {
    let keys = []
    for i in 0 to localStorage->Dom.Storage2.length - 1 {
      switch localStorage->Dom.Storage2.key(i) {
      | Some(k) => keys->Js.Array2.push(k)->ignore
      | None => ()
      }
    }
    for i in 0 to Js.Array2.length(keys) - 1 {
      localStorage->Dom.Storage2.removeItem(keys[i])
    }
  }

  let return = {
    setToken: saveToken,
    clearCognitoKeys: clearCognitoKeys,
    token: token,
  }

  return
}
