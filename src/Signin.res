@react.component
let make = () => {

    let (userInputText, setUserInputText) = React.useState(_ => "")
    let (pwInputText, setPwInputText) = React.useState(_ => "")
    let onChange = (func, e) => {
        let value = ReactEvent.Form.target(e)["value"]
        (_ => value)->func
    }

  <main>
    <form className="w-4/5 m-auto" action="">
      <fieldset className="flex flex-col items-center justify-around h-72">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred"> {"Sign in"->React.string} </legend>
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
            onChange=onChange(setUserInputText)
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
            value={userInputText}
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
            onChange=onChange(setPwInputText)
            // placeholder="Enter password"
            required=true
            spellCheck=false
            type_="password"
            value={pwInputText}
          />
        </div>
      </fieldset>
      <button className="text-gray-700 mt-20 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"> {"Submit"->React.string} </button>
    </form>
  </main>
}
