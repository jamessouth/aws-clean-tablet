

module Lazy = {
 
  type dynamicImport = Promise.t<{"make": React.component<{"fillColor": option<string>, "label": option<string>}>}>

  @module("react")
  external lazy_: (unit => Promise.t<React.component<{"fillColor": option<string>, "label": option<string>}>>) => React.component<{"fillColor": option<string>, "label": option<string>}> = "lazy"

  let load: ((unit => dynamicImport), {"fillColor": option<string>, "label": option<string>}) => React.element = (func, props) => {
    let propsPromise = Promise.promiseConstructor((resolve, _) => resolve(. props))

    lazy_(() => Promise.all2((func(), propsPromise))->Promise.then(((comp, prps)) => {
      let ccc = comp["make"]

      React.createElement(ccc, prps)
      }->Promise.resolve))
  }
  
}

@val
external import_: string => Lazy.dynamicImport = "import"



@react.component
let make = (
  ~userpool,
  ~setCognitoUser,
  ~setToken,
  ~cognitoUser,
  ~cognitoError,
  ~setCognitoError,
) => {

let comp = Lazy.load(() => import_("./Spin.bs"), Spin.makeProps(
  ~fillColor= "fill-blue-700",
  // ~label="jj",
  // ~key="",
  (),
))
// Js.log(comp)
// let comp = ->cmp

  let (username, setUsername) = React.Uncurried.useState(_ => "")
  let (password, setPassword) = React.Uncurried.useState(_ => "")
  let (validationError, setValidationError) = React.Uncurried.useState(_ => Some(
    "USERNAME: 3-10 length; PASSWORD: 8-98 length; at least 1 symbol; at least 1 number; at least 1 uppercase letter; at least 1 lowercase letter; ",
  ))
  let (submitClicked, setSubmitClicked) = React.Uncurried.useState(_ => false)
  let username_max_length = 10
  let password_max_length = 98

  React.useEffect2(() => {
    ErrorHook.useMultiError([(username, "USERNAME"), (password, "PASSWORD")], setValidationError)
    None
  }, (username, password))

  open Cognito
  let onClick = _ => {
    setSubmitClicked(._ => true)
    switch validationError {
    | None => {
        let cbs = {
          onSuccess: res => {
            setCognitoError(._ => None)
            Js.log2("signin result:", res)
            setToken(._ => Some(res.accessToken.jwtToken))
          },
          onFailure: ex => {
            switch Js.Exn.message(ex) {
            | Some(msg) => setCognitoError(._ => Some(msg))
            | None => setCognitoError(._ => Some("unknown signin error"))
            }

            setCognitoUser(._ => Js.Nullable.null)
            Js.log2("problem", ex)
          },
          newPasswordRequired: Js.Nullable.null,
          mfaRequired: Js.Nullable.null,
          customChallenge: Js.Nullable.null,
        }
        let authnData = {
          username: username
          ->Js.String2.slice(~from=0, ~to_=username_max_length)
          ->Js.String2.replaceByRe(%re("/\W/g"), ""),
          password: password
          ->Js.String2.slice(~from=0, ~to_=password_max_length)
          ->Js.String2.replaceByRe(%re("/\s/g"), ""),
          validationData: Js.Nullable.null,
          authParameters: Js.Nullable.null,
          clientMetadata: Js.Nullable.null,
        }
        let authnDetails = authenticationDetailsConstructor(authnData)

        switch Js.Nullable.isNullable(cognitoUser) {
        | true => {
            let userdata = {
              username: username
              ->Js.String2.slice(~from=0, ~to_=username_max_length)
              ->Js.String2.replaceByRe(%re("/\W/g"), ""),
              pool: userpool,
            }
            let user = Js.Nullable.return(userConstructor(userdata))
            user->authenticateUser(authnDetails, cbs)
            setCognitoUser(._ => user)
          }
        | false => cognitoUser->authenticateUser(authnDetails, cbs)
        }
      }
    | Some(_) => ()
    }
  }

//   let f = Spin.makeProps(
//   ~fillColor= "fill-blue-700",
//   // ~label="jj",
//   // ~key="",
//   (),
// )->Spin.make
// Js.log(f)

  {
    switch submitClicked {
    // | true => <Loading label="lobby..." />
    | true => React.null
    | false =>
    <>
    
    // f
    <React.Suspense fallback=React.null>
      comp
    </React.Suspense>
      <Form onClick leg="Sign in" submitClicked validationError cognitoError>
        <Input value=username propName="username" setFunc=setUsername />
        <Input
          value=password propName="password" autoComplete="current-password" setFunc=setPassword
        />
      </Form>
      </>
    }
  }
}
