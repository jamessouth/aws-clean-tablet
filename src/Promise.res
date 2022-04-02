type t<+'a> = Js.Promise.t<'a>
exception JsError(Js.Exn.t)
external unsafeToJsExn: exn => Js.Exn.t = "%identity"
@val @scope("Promise")
external resolve: 'a => t<'a> = "resolve"
@send external then: (t<'a>, @uncurry ('a => t<'b>)) => t<'b> = "then"
type resp = {
  ok: bool,
  redirected: bool,
  status: int,
  statusText: string,
  @as("type") _type: string,
  url: string,
}
type outcome = {
  status: string,
  value: Js.Nullable.t<resp>,
  reason: Js.Nullable.t<string>,
}
@val @scope("Promise")
external allSettled3: ((t<'a>, t<'b>, t<'c>)) => t<(outcome, outcome, outcome)> = "allSettled"
@val @scope("Promise")
external all3: ((t<'a>, t<'b>, t<'c>)) => t<('a, 'b, 'c)> = "all"
@send external _catch: (t<'a>, @uncurry (exn => t<'a>)) => t<'a> = "catch"
let catch = (promise, callback) => {
  _catch(promise, err => {
    let v = if Js.Exn.isCamlExceptionOrOpenVariant(err) {
      err
    } else {
      JsError(unsafeToJsExn(err))
    }
    callback(. v)
  })
}
