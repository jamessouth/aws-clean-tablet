type blob
type resp = {
  ok: bool,
  redirected: bool,
  status: int,
  statusText: string,
  @as("type") _type: string,
  url: string,
}
@send external blob: resp => Promise.t<blob> = "blob"
@val external fetch: string => Promise.t<Js.Nullable.t<resp>> = "fetch"
let getPics = assets =>
  {
    open Promise
    let (asset1, asset2, asset3) = assets
    all3((fetch(asset1), fetch(asset2), fetch(asset3)))
    ->then(((res1, res2, res3)) => {
      // switch Js.Nullable.toOption(res) {
      // | None => Error("uh oh")
      // | Some(r) => {
      //     Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok)
      //     switch r.ok {
      //     | true => Ok(r->blob)
      //     | false => {
      //         let stat = r.status
      //         Error(j`Fetch error: $stat - ${r.statusText}`)
      //       }
      //     }
      //   }
      // }
      Js.log4("rezz", res1, res2, res3)
      Ok((res1, res2, res3))
      ->resolve
    })
    // ->catch((. e) => {
    //   let msg = switch e {
    //   | JsError(err) =>
    //     switch Js.Exn.message(err) {
    //     | Some(msg) => msg
    //     | None => ""
    //     }
    //   | _ => "Unexpected error occurred"
    //   }
    //   Js.log2("Fetch error: ", msg)
    //   Error(msg)->resolve
    // })
  }->ignore

let handler = (assets, _e) => getPics(assets)
