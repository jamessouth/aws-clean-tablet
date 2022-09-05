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
external allSettled4: ((t<'a>, t<'b>, t<'c>, t<'d>)) => t<(outcome, outcome, outcome, outcome)> =
  "allSettled"
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
type blob
@send external blob: resp => t<blob> = "blob"
@val external fetch: string => t<resp> = "fetch"
@val external fetchAllSettled: string => t<outcome> = "fetch"
let getPicsAllSettled4 = assets =>
  {
    let (asset1, asset2, asset3, asset4) = assets
    allSettled4((
      fetchAllSettled(asset1),
      fetchAllSettled(asset2),
      fetchAllSettled(asset3),
      fetchAllSettled(asset4),
    ))
    ->then(((res1, res2, res3, res4)) => {
      [res1, res2, res3, res4]->Js.Array2.forEach(r =>
        switch r.status {
        | "fulfilled" =>
          switch Js.Nullable.toOption(r.value) {
          | Some(resp) => Js.log2("Asset " ++ resp.url ++ " fetched ok: ", resp.ok)
          | None => ()
          }
        | "rejected" =>
          switch Js.Nullable.toOption(r.reason) {
          | Some(msg) => Js.log("Asset fetch failed: " ++ msg)
          | None => ()
          }
        | _ => ()
        }
      )
      Ok(
        [res1, res2, res3, res4]
        ->Js.Array2.filter(r => r.status == "fulfilled")
        ->Js.Array2.map(r =>
          switch Js.Nullable.toOption(r.value) {
          | Some(resp) => resp->blob
          | None =>
            {
              ok: false,
              redirected: false,
              status: 0,
              statusText: "_",
              _type: "_",
              url: "_",
            }->blob
          }
        ),
      )->resolve
    })
    ->catch((. e) => {
      let msg = switch e {
      | JsError(err) =>
        switch Js.Exn.message(err) {
        | Some(msg) => msg
        | None => ""
        }
      | _ => "Unexpected error occurred"
      }
      Js.log2("Fetch error: ", msg)
      Error(msg)->resolve
    })
  }->ignore
