

type passwordPWCB = {
  onFailure: Js.Exn.t => unit,
  onSuccess: string => unit,
}

@send
external forgotPassword: (
  Js.Nullable.t<Signup.usr>, //user object
  passwordPWCB, //cb obj
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "forgotPassword"



@react.component
let make = () => {
  let (disabled, setDisabled) = React.useState(_ => true)
  let (email, setEmail) = React.useState(_ => "")

  let onChange = e => {
    let value = ReactEvent.Form.target(e)["value"]
    (_ => value)->setEmail
  }

  React.useEffect1(() => {
    switch email->Js.String2.length > 0 {
    | true => (_ => false)->setDisabled
    | false => (_ => true)->setDisabled
    }

    None
  }, [email])

  let handleSubmit = e => {
    e->ReactEvent.Form.preventDefault
    let forgotPWcb = {
      onSuccess: str => {
        Js.log2("forgot pw initiated: ", str)
      RescriptReactRouter.push("/confirm")
      
      },
      onFailure: err => {
        switch Js.Exn.message(err) {
        | Some(msg) => (_ => Some(msg))->setCognitoErr
        | None => (_ => Some("unknown forgot pw error"))->setCognitoErr
        }
        Js.log2("forgot pw problem: ", err)
      },
    }





  }

  <main>
    <form className="w-5/6 m-auto" onSubmit={handleSubmit}>
      <fieldset className="h-40">
        <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred">
          {"Reset password"->React.string}
        </legend>
        <div>
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
            onChange
            // placeholder="Enter username"
            required=true
            spellCheck=false
            type_="email"
            value={email}
          />
        </div>
      </fieldset>
      <button
        disabled
        className="text-gray-700 mt-10 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7">
        {"send code"->React.string}
      </button>
    </form>
  </main>
}
