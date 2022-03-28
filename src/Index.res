%%raw(`import 'virtual:windi.css'`)
%%raw(`import 'virtual:windi-devtools'`)
%%raw(`import './css/windi.css'`)
@scope("window") @val
external addWindowEventListener: (string, unit => unit) => unit = "addEventListener"
@scope("window") @val
external removeWindowEventListener: (string, unit => unit) => unit = "removeEventListener"

switch ReactDOM.querySelector("#root") {
| Some(root) => ReactDOM.render(<App />, root)
| None => ()
}

// module Response = {
//   type t<'data>
// }

module Promise = {
  type t<+'a> = Js.Promise.t<'a>

  exception JsError(Js.Exn.t)
  external unsafeToJsExn: exn => Js.Exn.t = "%identity"

  @val @scope("Promise")
  external resolve: 'a => t<'a> = "resolve"

  @send external then: (t<'a>, @uncurry ('a => t<'b>)) => t<'b> = "then"

  @send external thenResolve: (t<'a>, @uncurry ('a => 'b)) => t<'b> = "then"

  @send external _catch: (t<'a>, @uncurry (exn => t<'a>)) => t<'a> = "catch"

  let catch = (promise, callback) => {
    _catch(promise, err => {
      let v = if Js.Exn.isCamlExceptionOrOpenVariant(err) {
        err
      } else {
        JsError(unsafeToJsExn(err))
      }
      callback(v)
    })
  }
}

module Prefetch = {
  type blob
  type resp = {
    ok: bool,
    redirected: bool,
    status: int,
    statusText: string,
    @as("type") _type: string,
    url: string,
  }
  // type response = {"response": Js.Nullable.t<resp>, "error": Js.Nullable.t<string>}

  @send external blob: resp => Promise.t<blob> = "blob"

  @val external fetch: string => Promise.t<Js.Nullable.t<resp>> = "fetch"
}

let getPic = () => {
  open Promise
  open Prefetch
  Js.log("load ev")
  fetch("../../assets/chmob2x.webp")
  ->then(res => {
    Js.log2("res: ", res)
    switch Js.Nullable.isNullable(res) {
    | true => Error("uh oh")
    | false =>
      switch res.ok {
      // | true => Ok(res->blob)
      | false => Error(`Fetch error: ${res.status} - ${res.statusText}`)
      }
    }
  })
  ->thenResolve(res => {
    Js.log2("res2", res)
  })
  ->catch(e => {
    let msg = switch e {
    | JsError(err) =>
      switch Js.Exn.message(err) {
      | Some(msg) => msg
      | None => ""
      }
    | _ => "Unexpected error occurred"
    }
    Error(msg)->resolve
  })
}
addWindowEventListener("load", getPic)

// Js.Global.setTimeout(() => {
//     Js.log("7")
//   removeWindowEventListener("load", getPic)
// }, 7000)->ignore
