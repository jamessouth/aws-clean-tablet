let useValidator = (field, setErrorFunc, funcList, title) => {
  let rec validate = (funcs, str, acc) =>
    switch funcs {
    | list{} => acc
    | list{h, ...t} => validate(t, str, acc ++ h(str))
    }

  React.useEffect1(() => {
    let error = validate(funcList, field, "")
    let final = switch error == "" {
    | true => None
    | false => Some(title ++ error)
    }
    setErrorFunc(_ => final)
    None
  }, [field])
}
