type t

type poolData
type poolDataInput = {
    @as("UserPoolId") userPoolId: string,
    @as("ClientId") clientId: string,
    @as("AdvancedSecurityDataCollectionFlag") advancedSecurityDataCollectionFlag: bool
}

type attributeData
type attributeDataInput = {
    @as("Name") name: string,
    @as("Value") value: string
}

type userData
type userDataInput = {
    @as("Username") username: string,
    @as("Pool") pool: poolDataInput
}

@send external focus: Dom.element => unit = "focus"


@new @module("amazon-cognito-identity-js")
external userPoolConstructor : poolDataInput => poolData = "CognitoUserPool"
@new @module("amazon-cognito-identity-js")
external userAttributeConstructor : attributeDataInput => attributeData = "CognitoUserAttribute"
@new @module("amazon-cognito-identity-js")
external userConstructor : userDataInput => userData = "CognitoUser"



@val @scope(("import", "meta", "env"))
external upid: string = "VITE_UPID"
@val @scope(("import", "meta", "env"))
external cid: string = "VITE_CID"



// @val @scope("window")
// external alert: string => unit = "alert"


let cbToOption = (f) => (. err, result) =>
  switch (Js.Nullable.toOption(err), Js.Nullable.toOption(result)) {
  | (Some(err), _) => f(Error(err))
  | (_, Some(result)) => f(Ok(result))
  | _ => invalid_arg("invalid argument for cbToOption")
  }


let signupcallback = cbToOption(result =>
  switch result {
  | Ok(val) => {
      Js.log(val)
      ()
    }
  | Error(ex) => {
      Js.log(ex)
      ()
    }
  })




type clientMetadata = {
    key: string
}
type signUpCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<t>) => unit

@send external signUp: (poolData, string, string, Js.Nullable.t<array<attributeData>>, Js.Nullable.t<array<attributeData>>, signUpCB, Js.Nullable.t<clientMetadata>) => unit = "signUp"


let pool = {
    userPoolId: upid,
    clientId: cid,
    advancedSecurityDataCollectionFlag: false
}
let userpool = userPoolConstructor(pool)



@react.component
let make = () => {

    // let pwInput = React.useRef(Js.Nullable.null)

    let (unVisited, setUnVisited) = React.useState(_ => false)
    let (pwVisited, setPwVisited) = React.useState(_ => false)
    let (unErr, setUnErr) = React.useState(_ => None)
    let (pwErr, setPwErr) = React.useState(_ => None)


    let (showPassword, setShowPassword) = React.useState(_ => false)
    let (disabled, setDisabled) = React.useState(_ => true)
    let (username, setUsername) = React.useState(_ => "")
    let (password, setPassword) = React.useState(_ => "")
    let (email, setEmail) = React.useState(_ => "")





    let checkUnForbiddenChars = un => {
      let r = %re("/\W/")

      switch Js.String2.match_(un, r) {
      | Some(_) => (_ => Some("alphanumeric..."))->setUnErr
      | None => (_ => None)->setUnErr
      }
    }

    let checkPwForbiddenChars = pw => {
      let r = %re("/[-=+]/")

      switch Js.String2.match_(pw, r) {
      | Some(_) => (_ => Some("no +, -, or = ..."))->setPwErr
      | None => (_ => None)->setPwErr
      }
    }

    let checkUnMaxLength = un => {
      switch un->Js.String2.length > 10 {
      | true => (_ => Some("too long..."))->setUnErr
      | false => un->checkUnForbiddenChars
      }
    }

    let checkPwMaxLength = pw => {
      switch pw->Js.String2.length > 98 {
      | true => (_ => Some("too long..."))->setPwErr
      | false => pw->checkPwForbiddenChars
      }
    }

    let checkNoUnWhitespace = un => {
      let r = %re("/\s/")

      switch Js.String2.match_(un, r) {
      | Some(_) => (_ => Some("no whitespace..."))->setUnErr
      | None => un->checkUnMaxLength
      }
    }

    let checkNoPwWhitespace = pw => {
      let r = %re("/\s/")

      switch Js.String2.match_(pw, r) {
      | Some(_) => (_ => Some("no whitespace..."))->setPwErr
      | None => pw->checkPwMaxLength
      }
    }

    let checkSymbol = pw => {
      let r = %re("/[!-*\[-`{-~./,:;<>?@]/")

      switch Js.String2.match_(pw, r) {
      | None => (_ => Some("add symbol..."))->setPwErr
      | Some(_) => pw->checkNoPwWhitespace
      }
    }

    let checkNumber = pw => {
      let r = %re("/\d/")

      switch Js.String2.match_(pw, r) {
      | None => (_ => Some("add number..."))->setPwErr
      | Some(_) => pw->checkSymbol
      }
    }

    let checkUpper = pw => {
      let r = %re("/[A-Z]/")

      switch Js.String2.match_(pw, r) {
      | None => (_ => Some("add uppercase..."))->setPwErr
      | Some(_) => pw->checkNumber
      }
    }

    let checkLower = pw => {
      let r = %re("/[a-z]/")

      switch Js.String2.match_(pw, r) {
      | None => (_ => Some("add lowercase..."))->setPwErr
      | Some(_) => pw->checkUpper
      }
    }

    let checkUnLength = un => {
        switch un->Js.String2.length < 4 {
        | true => (_ => Some("too short..."))->setUnErr
        | false => un->checkNoUnWhitespace
        }
    }

    let checkPwLength = pw => {
        switch pw->Js.String2.length < 8 {
        | true => (_ => Some("too short..."))->setPwErr
        | false => pw->checkLower
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
        value: email
      }
      let emailAttr = userAttributeConstructor(emailData)
      userpool->signUp(username, password, Js.Nullable.return([emailAttr]), Js.Nullable.null, signupcallback, Js.Nullable.null)
    }

    React.useEffect2(() => {
      switch unVisited {
      | true => username->checkUnLength
      | false => (_ => None)->setUnErr
      }
      None
    }, (username, unVisited))

    React.useEffect2(() => {
      switch pwVisited {
      | true => password->checkPwLength
      | false => (_ => None)->setPwErr
      }
      None
    }, (password, pwVisited))



    React.useEffect3(() => {
      switch (username->Js.String2.length > 3, password->Js.String2.length > 7, email->Js.String2.length > 0) {
      | (true, true, true) => (_ => false)->setDisabled
      | (false, _, _) | (_, false, _) | (_, _, false) => (_ => true)->setDisabled
      }

      None
    }, (username, password, email))


  <main>
    <form className="w-4/5 m-auto" onSubmit={handleSubmit}>
      <fieldset className="flex flex-col items-center justify-around h-80">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred"> {"Sign up"->React.string} </legend>
        <div className="relative">
          <label className={switch (unVisited, unErr) {
          | (true, Some(_)) => "text-2xl text-red-500 font-bold font-flow"
          | (false, _) | (true, None) => "text-2xl text-warm-gray-100 font-flow"
          }} htmlFor="username">
            {"username:"->React.string}
          </label>
          {
              switch (unVisited, unErr) {
          | (true, Some(err)) => <span className="absolute right-0 text-2xl text-red-500 font-bold font-flow">{err->React.string}</span>
          | (false, _) | (true, None) => React.null
          }
            }
          <input
            autoComplete="username"
            autoFocus=true
            className={switch (unVisited, unErr) {
              | (true, Some(_)) => "h-6 w-full text-xl pl-1 text-left outline-none text-red-500 bg-transparent border-b-1 border-red-500"
              | (false, _) | (true, None) => "h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
              }}
            id="username"
            minLength=4
            name="username"
            onBlur=onBlur("username")
            onChange=onChange(setUsername)
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={username}
          />
        </div>
        <div className="relative">
          <label className={switch (pwVisited, pwErr) {
          | (true, Some(_)) => "text-2xl text-red-500 font-bold font-flow"
          | (false, _) | (true, None) => "text-2xl text-warm-gray-100 font-flow"
          }} htmlFor="new-password">
            {"password:"->React.string}
          </label>
          {
              switch (pwVisited, pwErr) {
          | (true, Some(err)) => <span className="absolute right-0 text-2xl text-red-500 font-bold font-flow">{err->React.string}</span>
          | (false, _) | (true, None) => React.null
          }
            }
          <input
            autoComplete="new-password"
            autoFocus=false

            className={switch (pwVisited, pwErr) {
              | (true, Some(_)) => "h-6 w-3/4 text-xl pl-1 text-left outline-none text-red-500 bg-transparent border-b-1 border-red-500"
              | (false, _) | (true, None) => "h-6 w-3/4 text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
              }}



            // className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="new-password"
            minLength=8
            name="password"
            onBlur=onBlur("password")
            onChange=onChange(setPassword)
            // placeholder="Enter password"
            // ref={pwInput->ReactDOM.Ref.domRef}
            required=true
            spellCheck=false
            type_={switch showPassword {
            | true => "text"
            | false => "password"
            }}
            value={password}
          />
          
          <button type_="button" className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 cursor-pointer" onClick>{switch showPassword {
          | true => "hide"->React.string
          | false => "show"->React.string
          }}</button>
        </div>


        <div className="w-full">
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="email">
            {"email:"->React.string}
          </label>
          <input
            autoComplete="email"
            autoFocus=false
            className="h-6 w-full text-base pl-1 text-left font-anon outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="email"
            // minLength=4
            name="email"
            onChange=onChange(setEmail)
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="email"
            value={email}
          />
        </div>

      </fieldset>


      <button disabled className="text-gray-700 mt-16 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"> {"create"->React.string} </button>
    </form>
  </main>
}