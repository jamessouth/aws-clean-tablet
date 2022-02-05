let useValidator = (field, setErrorFunc, funcList, title) => {
  React.useEffect1(() => {
    let error = funcList->Js.Array2.reduce((acc, f) => acc ++ f(field), "")
    let final = switch error == "" {
    | true => None
    | false => Some(title ++ error)
    }
    setErrorFunc(_ => final)
    None
  }, [field])
}
