





type cdd = {
    @as("AttributeName") attributeName: string,
    @as("DeliveryMedium") deliveryMedium: string,
    @as("Destination") destination: string
}

type clnt = {
    endpoint: string,
    fetchOptions: {}
}

type pl = {
    advancedSecurityDataCollectionFlag: bool,
    client: clnt,
    clientId: string,
    storage: {"length": float},
    userPoolId: string
}

type usr = {
    @as("Session") session: Js.Nullable.t<string>,
    authenticationFlowType: string,
    client: clnt,
    keyPrefix: string,
    pool: pl,
    signInUserSession: Js.Nullable.t<string>,
    storage: {"length": float},
    userDataKey: string,
    username: string
}

type signupOk = {
    codeDeliveryDetails: cdd,
    user: usr,
    userConfirmed: bool,
    userSub: string
}





@react.component
let make = (~setCognitoUser) => {
  // let pwInput = React.useRef(Js.Nullable.null)

  let (unVisited, setUnVisited) = React.useState(_ => false)
  
  
  let (unErr, setUnErr) = React.useState(_ => None)
  

  
  let (disabled, _setDisabled) = React.useState(_ => false)
  let (username, setUsername) = React.useState(_ => "")
  
  

  let (cognitoErr, setCognitoErr) = React.useState(_ => None)
  // let (cognitoResult, setCognitoResult) = React.useState(_ => false)





  let checkUnForbiddenChars = un => {
    let r = %re("/\W/")

    switch Js.String2.match_(un, r) {
    | Some(_) => (_ => Some("alphanumeric..."))->setUnErr
    | None => (_ => None)->setUnErr
    }
  }



  let checkUnMaxLength = un => {
    switch un->Js.String2.length > 10 {
    | true => (_ => Some("too long..."))->setUnErr
    | false => un->checkUnForbiddenChars
    }
  }



  let checkNoUnWhitespace = un => {
    let r = %re("/\s/")

    switch Js.String2.match_(un, r) {
    | Some(_) => (_ => Some("no whitespace..."))->setUnErr
    | None => un->checkUnMaxLength
    }
  }







  let checkUnLength = un => {
    switch un->Js.String2.length < 4 {
    | true => (_ => Some("too short..."))->setUnErr
    | false => un->checkNoUnWhitespace
    }
  }



  let onClick = _e => {
    (prev => !prev)->setShowPassword
  }

  let onChange = (func, e) => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->func
  }

  let onBlur = (input, _e) => {
    switch input {
    | "username" => (_ => true)->setUnVisited
    | "password" => (_ => true)->setPwVisited
    | _ => ()
    }
  }

  let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    let emailData = {
      name: "email",
      value: email,
    }
    let emailAttr = userAttributeConstructor(emailData)
    userpool->signUp(
      username,
      password,
      Js.Nullable.return([emailAttr]),
      Js.Nullable.null,
      signupCallback,
      Js.Nullable.null,
    )
  }

  React.useEffect2(() => {
    switch unVisited {
    | true => username->checkUnLength
    | false => (_ => None)->setUnErr
    }
    None
  }, (username, unVisited))



  // React.useEffect5(() => {
  //   switch (unErr, pwErr, username->Js.String2.length < 4, password->Js.String2.length < 8, email->Js.String2.length < 3) {
  //   | (None, None, false, false, false) => (_ => false)->setDisabled
  //   | (Some(_), _, _, _, _)
  //   | (_, Some(_), _, _, _)
  //   | (_, _, true, _, _)
  //   | (_, _, _, true, _)
  //   | (_, _, _, _, true) => (_ => true)->setDisabled
  //   }

  //   None
  // }, (unErr, pwErr, username, password, email))

  <main>
    <form className="w-4/5 m-auto" onSubmit={handleSubmit}>
      <fieldset className="flex flex-col items-center justify-around h-80">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred">
          {"Sign up"->React.string}
        </legend>
        <div className="relative">
          <label
            className={switch (unVisited, unErr) {
            | (true, Some(_)) => "text-2xl text-red-500 font-bold font-flow"
            | (false, _) | (true, None) => "text-2xl text-warm-gray-100 font-flow"
            }}
            htmlFor="username">
            {"username:"->React.string}
          </label>
          {switch (unVisited, unErr) {
          | (true, Some(err)) =>
            <span className="absolute right-0 text-2xl text-red-500 font-bold font-flow">
              {err->React.string}
            </span>
          | (false, _) | (true, None) => React.null
          }}
          <input
            autoComplete="username"
            autoFocus=true
            className={switch (unVisited, unErr) {
            | (
                true,
                Some(_),
              ) => "h-6 w-full text-xl pl-1 text-left outline-none text-red-500 bg-transparent border-b-1 border-red-500"
            | (false, _)
            | (
              true,
              None,
            ) => "h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            }}
            id="username"
            minLength=4
            name="username"
            onBlur={onBlur("username")}
            onChange={onChange(setUsername)}
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={username}
          />
        </div>
      </fieldset>
      {
        switch cognitoErr {
        | Some(msg) => <span className="text-sm text-warm-gray-100 absolute bg-red-500 text-center w-full left-1/2 transform max-w-lg -translate-x-1/2">{msg->React.string}</span>
        | None => React.null
        }
      }
      <button
        disabled
        className="text-gray-700 mt-16 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7">
        {"create"->React.string}
      </button>
    </form>
  </main>
}
