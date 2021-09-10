type t

type poolData = {
    userPoolId: string,
    clientId: string
}

// @module("amazon-cognito-identity-js") @new external CognitoUserPool: poolData => t = "Foo"


// type t
@new @module("amazon-cognito-identity-js") external constructor : poolData => t = "CognitoUserPool"



let pool = {
  userPoolId: "",
  clientId: ""
}
let userpool = constructor(pool)

@react.component
let make = () => {

    let (disabled, setDisabled) = React.useState(_ => true)
    let (username, setUsername) = React.useState(_ => "")
    let (password, setPassword) = React.useState(_ => "")
    let (email, setEmail) = React.useState(_ => "")
    let onChange = (func, e) => {
        let value = ReactEvent.Form.target(e)["value"]
        (_ => value)->func
    }

    let handleSubmit = e => {
      e->ReactEvent.Form.preventDefault

    }

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
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="password">
            {"password:"->React.string}
          </label>
          <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 w-full text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            onChange=onChange(setPassword)
            // placeholder="Enter password"
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
            autoFocus=true
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
