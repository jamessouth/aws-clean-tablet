type t

type confirmRegistrationCB = (. Js.Nullable.t<Js.Exn.t>, Js.Nullable.t<t>) => unit

@send
external confirmRegistration: (
  Js.Nullable.t<Signup.usr>,
  string,
  bool,
  confirmRegistrationCB,
  Js.Nullable.t<Signup.clientMetadata>,
) => unit = "confirmRegistration"



@react.component
let make = (~cognitoUser) => {
Js.log2("user", cognitoUser)
    let (showVerifCode, setShowVerifCode) = React.useState(_ => false)
    let (verifCode, setVerifCode) = React.useState(_ => "")

    let (disabled, setDisabled) = React.useState(_ => true)

    let onClick = _e => {
      (prev => !prev)->setShowVerifCode
    }

    let onChange = e => {
        let value = ReactEvent.Form.target(e)["value"]
        (_ => value)->setVerifCode
    }

    React.useEffect1(() => {
        switch verifCode->Js.String2.length != 6 {
        | true => (_ => true)->setDisabled
        | false => (_ => false)->setDisabled
        }
        None
    }, [verifCode])


  let confirmregistrationCallback = Signup.cbToOption(res =>
    switch res {
    | Ok(val) => {
        // (_ => None)->setCognitoErr
        // (_ => Some(val.user))->setCognitoUser
        // RescriptReactRouter.push("/confirm")


        Js.log2("conf rego res", val)
      }
    | Error(ex) => {
        // switch Js.Exn.message(ex) {
        // | Some(msg) => (_ => Some(msg))->setCognitoErr
        // | None => (_ => None)->setCognitoErr
        // }

        Js.log2("conf rego problem", ex)
      }
    }
  )






    let handleSubmit = e => {
      e->ReactEvent.Form.preventDefault
      cognitoUser->confirmRegistration(
        verifCode,
        false,
        confirmregistrationCallback,
        Js.Nullable.null,
      )

    }


  <main>
    <form className="w-5/6 m-auto" onSubmit={handleSubmit}>
      <fieldset className="h-40">
        <legend className="text-warm-gray-100 m-auto mb-8 text-3xl font-fred"> {"Confirm"->React.string} </legend>
        <div className="relative">
          <label
            className="block text-2xl text-warm-gray-100 font-flow"
            htmlFor="verif-code">
            {"enter code:"->React.string}
          </label>
          <input
            autoComplete="one-time-code"
            autoFocus=true
            className="h-8 text-xl outline-none text-warm-gray-100 bg-transparent border-b-1 border-warm-gray-100"
            id="verif-code"
            maxLength=6
            minLength=6
            inputMode="numeric"
            name="verifcode"
            onChange
            // placeholder="Enter password"
            // ref={pwInput->ReactDOM.Ref.domRef}
            pattern="^\d{6}$"
            required=true
            size=6
            spellCheck=false
            type_={switch showVerifCode {
            | true => "text"
            | false => "password"
            }}
            value={verifCode}
          />
          <button
            type_="button"
            className="font-arch bg-transparent text-warm-gray-100 text-2xl absolute right-0 cursor-pointer"
            onClick>
            {switch showVerifCode {
            | true => "hide"->React.string
            | false => "show"->React.string
            }}
          </button>
        </div>

  
      </fieldset>


      <button disabled className="text-gray-700 mt-10 bg-warm-gray-100 block font-flow text-2xl mx-auto cursor-pointer w-3/5 h-7"> {"confirm"->React.string} </button>
    </form>
  </main>


}