open Promise
type blob
@send external blob: resp => Promise.t<blob> = "blob"
@val external fetch: string => Promise.t<resp> = "fetch"
@val external fetchAllSettled: string => Promise.t<outcome> = "fetch"
let getPicsAllSettled4 = assets =>
  {
    let (asset1, asset2, asset3, asset4) = assets

    allSettled4((fetchAllSettled(asset1), fetchAllSettled(asset2), fetchAllSettled(asset3), fetchAllSettled(asset4)))
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
// let getPics = assets =>
//   {
//     let (asset1, asset2, asset3) = assets
//     all3((fetch(asset1), fetch(asset2), fetch(asset3)))
//     ->then(((res1, res2, res3)) =>
//       {
//         let resps = [res1, res2, res3]
//         resps->Js.Array2.forEach(r => Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok))
//         switch resps->Js.Array2.every(r => r.ok) {
//         | true => Ok(resps->Js.Array2.map(r => r->blob))
//         | false => {
//             let {status, statusText, url} = switch resps->Js.Array2.find(r => !r.ok) {
//             | Some(r) => r
//             | None => {
//                 ok: false,
//                 redirected: false,
//                 status: 0,
//                 statusText: "_",
//                 _type: "_",
//                 url: "_",
//               }
//             }
//             let msg = j`Fetch error for asset ${url}: $status - ${statusText}`
//             Js.log(msg)
//             Error(msg)
//           }
//         }
//       }->resolve
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
// let handlerAllSettled = assets => getPicsAllSettled(assets)
// let handler = (assets, _e) => getPics(assets)
