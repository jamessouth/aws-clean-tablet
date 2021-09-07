



@react.component
let make = () => {

<main className="w-4/5 flex flex-col items-center h-60 justify-around m-auto">
    <form action="">
    <fieldset>
    <legend>{"Sign in"->React.string}</legend>
        <label className="text-xl self-start" htmlFor="username">
            {"Username:"->React.string}
        </label>
        <input
            autoComplete="username"
            autoFocus=true
            className="h-6 text-xl pl-1 text-left bg-transparent border-b-1 border-warm-gray-100"
            id="username"
            minLength=4
            name="username"
            placeholder="Enter username"
            required=true
            spellCheck=false
            type_="text"
        />
       
        <label className="text-xl self-start" htmlFor="password">
            {"Password:"->React.string}
        </label>
        <input
            autoComplete="current-password"
            autoFocus=false
            className="h-6 text-xl pl-1 text-left bg-transparent border-b-1 border-warm-gray-100"
            id="password"
            minLength=8
            name="password"
            placeholder="Enter password"
            required=true
            spellCheck=false
            type_="password"
        />
        </fieldset>
        
        
        
        <button className="text-gray-700">{"Submit"->React.string}</button>
    </form>
</main>


}