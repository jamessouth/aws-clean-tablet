type propShape = {
  "userpool": Cognito.poolData,
  "setCognitoUser": (. Js.Nullable.t<Cognito.usr> => Js.Nullable.t<Cognito.usr>) => unit,
  "setToken": (. option<string> => option<string>) => unit,
  "cognitoUser": Js.Nullable.t<Cognito.usr>,
  "cognitoError": option<string>,
  "setCognitoError": (. option<string> => option<string>) => unit,
  "token": option<string>,
  "path": list<string>,
}

@val
external import_: string => Promise.t<{"make": React.component<propShape>}> = "import"

@module("react")
external lazy_: (unit => Promise.t<{"default": React.component<propShape>}>) => React.component<
  propShape,
> = "lazy"


  let initialState: Reducer.state = {
    gamesList: Js.Nullable.null,
    players: [],
    sk: "",
    oldWord: "",
    word: "",
    showAnswers: false,
    winner: "",
  }



@react.component
let make = (
  ~userpool,
  ~setCognitoUser,
  ~setToken,
  ~cognitoUser,
  ~cognitoError,
  ~setCognitoError,
  ~token,
  ~path,
) => {





  let (
    playerGame,
    _,
    _,
    _,
    count,
    wsConnected,
    _,
    _,
    _,
    _,
    _,
    _,
    games,
    _,
    setLeaderData,
    send,
    _,
    close,
    wsError,
  ) = WsHook.useWs(token, setToken, cognitoUser, setCognitoUser, initialState)


  let lobby = React.createElement(
    Lobby.lazy_(() =>
      Lobby.import_("./Lobby.bs")->Promise.then(comp => {
        Promise.resolve({"default": comp["make"]})
      })
    ),
    Lobby.makeProps(~playerGame, ~games, ~send, ~wsError, ~close, ~count, ~setLeaderData, ~url=path, ()),
  )

  let loading1 = React.createElement(Loading.lazy_(() =>
    Loading.import_("./Loading.bs")->Promise.then(comp => {
      Promise.resolve({"default": comp["make"]})
    })
  ), Loading.makeProps(~label="games...", ()))





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

open Web

  {switch (path, token) {
    | (list{"signin"}, None) => {
        switch submitClicked {
        | true => <Loading label="lobby..." />
        | false =>
          <Form onClick leg="Sign in" submitClicked validationError cognitoError>
            <Input value=username propName="username" setFunc=setUsername />
            <Input
              value=password propName="password" autoComplete="current-password" setFunc=setPassword
            />
          </Form>
        }
      }

    | (list{"lobby"}, Some(_)) =>
        switch wsConnected {
        | false => {
            body(document)->setClassName("bodchmob bodchtab bodchbig")
            <React.Suspense fallback=React.null> loading1 </React.Suspense>
          }

        | true => {
            body(document)->classList->removeClassList3("bodleadmob", "bodleadtab", "bodleadbig")

            <React.Suspense fallback=React.null> lobby </React.Suspense>
          }
        }

    | (_, _) => <div> {React.string("other222")} </div> // <PageNotFound/>

  }}

  
}
