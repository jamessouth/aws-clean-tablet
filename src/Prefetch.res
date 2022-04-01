// type blob
// type resp = {
//   ok: bool,
//   redirected: bool,
//   status: int,
//   statusText: string,
//   @as("type") _type: string,
//   url: string,
// }
//     open Promise
// @send external blob: resp => Promise.t<blob> = "blob"
// @val external fetchAll: string => Promise.t<Js.Nullable.t<resp>> = "fetch"
// @val external fetchAllSettled: string => Promise.t<Js.Nullable.t<outcome<resp>>> = "fetch"
// let getPics = assets =>
//   {
//     let (asset1, asset2, asset3) = assets
//     allSettled3((fetchAllSettled(asset1), fetchAllSettled(asset2), fetchAllSettled(asset3)))
//     ->then(((res1, res2, res3)) =>
//       Ok(
//         [res1, res2, res3]->Js.Array2.map(x =>
//           switch Js.Nullable.toOption(x) {
//           | None => Error("null result")
//           | Some(rs) => {
//             Js.log2("rs", rs)
//             switch rs {
//             | Rejected(str) => Error(`Fetch rejected: ${str}`)
//             | Fulfilled(r) => {
//                 Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok)
//                 switch r.ok {
//                 | true => Ok(r->blob)
//                 | false => {
//                     let stat = r.status
//                     Error(j`Fetch error: $stat - ${r.statusText}`)
//                   }
//                 }
//               }
//             }
//           }

//           }
//         ),
//       )->resolve
//     )
//     ->catch((. e) => {
//       let msg = switch e {
//       | JsError(err) =>
//         switch Js.Exn.message(err) {
//         | Some(msg) => msg
//         | None => ""
//         }
//       | _ => "Unexpected error occurred"
//       }
//       Js.log2("Fetch error: ", msg)
//       Error(msg)->resolve
//     })
//   }->ignore

// let handler = (assets, _e) => getPics(assets)

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
@val external fetch: string => Promise.t<resp> = "fetch"
let getPics = assets =>
  {
    open Promise
    let (asset1, asset2, asset3) = assets
    all3((fetch(asset1), fetch(asset2), fetch(asset3)))
    ->then(((res1, res2, res3)) =>
      {
        let resps = [res1, res2, res3]
        resps->Js.Array2.forEach(r => Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok))
        switch resps->Js.Array2.every(r => r.ok) {
        | true => Ok(resps->Js.Array2.map(r => Ok(r->blob)))
        | false => {
            let {status, statusText, url} = switch resps->Js.Array2.find(r => !r.ok) {
            | Some(r) => r
            | None => {
                ok: false,
                redirected: false,
                status: 0,
                statusText: "_",
                _type: "_",
                url: "_",
              }
            }
            let msg = j`Fetch error for asset ${url}: $status - ${statusText}`
            Js.log(msg)
            Error(msg)
          }
        }
      }->resolve
    )
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

let handler = (assets, _e) => getPics(assets)
