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

type pwError = 
  | Good(string)
  | Invalid(string)

let checkLength = (~pw, ~pw2): pwError => {
    switch pw->Js.String2.length < 8 {
    | true => Invalid("8+ chars")
    | false => Good(pw)
    }
    switch pw2->Js.String2.length < 8 {
    | true => Invalid("8+ chars")
    | false => Good(pw)
    }
}

@react.component
let make = () => {

    // let pwInput = React.useRef(Js.Nullable.null)

    let (pwVisited, setPwVisited) = React.useState(_ => false)
    let (pwErrs, setPwErrs) = React.useState(_ => [])


    let (disabled, setDisabled) = React.useState(_ => true)
    let (username, setUsername) = React.useState(_ => "")
    let (password, setPassword) = React.useState(_ => "")
    let (email, setEmail) = React.useState(_ => "")
    let onChange = (func, e) => {
        let value = ReactEvent.Form.target(e)["value"]
        (_ => value)->func
    }

    let onBlur = _e => {
      (_ => true)->setPwVisited
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

    React.useEffect1(() => {
      
      switch password->Js.String2.length < 8 {
      | true => (prev => Js.Array2.concat(prev, ["must be longer"]))->setPwErrs
      | false => (_ => [])->setPwErrs
      }

      None
    }, [password])


    // React.useEffect1(() => {
      
    //   switch pwVisited {
    //   | false => expression
    //   | true => expression
    //   }

    //   None
    // }, [pwVisited])

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
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
            {"username:"->React.string}
          </label>
          <input
            autoComplete="username"
            autoFocus=true
            className="h-6 w-full text-xl pl-1 text-left font-anon outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="username"
            minLength=4
            name="username"
            onChange=onChange(setUsername)
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={username}
          />
        </div>
        <div>
          <label className={switch (pwVisited, pwErrs->Js.Array2.length > 0) {
          | (true, true) => "text-2xl text-red-500 font-bold font-flow" 
          | (false, _) | (true, false) => "text-2xl text-warm-gray-100 font-flow" 
          }} htmlFor="password">
            {"password:"->React.string}
          </label>
          <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            onBlur
            onChange=onChange(setPassword)
            // placeholder="Enter password"
            // ref={pwInput->ReactDOM.Ref.domRef}
            required=true
            spellCheck=false
            type_="password"
            value={password}
          />
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
