@react.component
let make = () => {
  <main>
    <form className="w-4/5 m-auto" action="">
      <fieldset className="flex flex-col items-center justify-around h-60">
        <legend className="text-warm-gray-100 m-auto mb-6 text-3xl font-fred"> {"Sign in"->React.string} </legend>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="username">
            {"username:"->React.string}
          </label>
          <input
            autoComplete="username"
            autoFocus=true
            className="h-6 w-64 text-xl pl-1 text-left font-anon outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="username"
            minLength=4
            name="username"
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
          />
        </div>
        <div>
          <label className="text-2xl text-warm-gray-100 font-flow" htmlFor="password">
            {"password:"->React.string}
          </label>
          <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 w-64 text-xl pl-1 text-left outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            // placeholder="Enter password"
            required=true
            spellCheck=false
            type_="password"
          />
        </div>
      </fieldset>
      <button className="text-gray-700 mt-20 block font-flow mx-auto w-1/2 h-6"> {"Submit"->React.string} </button>
    </form>
  </main>
}
