%%raw(`import 'virtual:windi.css'`)
%%raw(`import 'virtual:windi-devtools'`)
%%raw(`import './css/windi.css'`)
open Web

let mob = matchMedia("(max-width: 767.9px)")
let big = matchMedia("(min-width: 1440px)")

let retinaMQ = switch matchMedia("(resolution: 2x)").media != "not all" {
| true => "(min-resolution: 2x)"
| false => "(-webkit-min-device-pixel-ratio: 2)"
}

let retina = matchMedia(retinaMQ)

let chalkbg = switch mob.matches {
| true =>
  switch retina.matches {
  | true => "../../assets/chmob2x.webp"
  | false => "../../assets/chmob1x.webp"
  }
| false =>
  switch big.matches {
  | true =>
    switch retina.matches {
    | true => "../../assets/chbig2x.webp"
    | false => "../../assets/chbig1x.webp"
    }
  | false =>
    switch retina.matches {
    | true => "../../assets/chtab2x.webp"
    | false => "../../assets/chtab1x.webp"
    }
  }
}

let newgamebg = switch retina.matches {
| true => "../../assets/ekko2x.webp"
| false => "../../assets/ekko1x.webp"
}

let signout = "../../assets/signout.png"

let bghand = Prefetch.handlerAllSettled((chalkbg, newgamebg, signout))

addWindowEventListener("load", bghand)

// let loading = "./src/Loading.bs.js"

// let ldhand = jj => Prefetch.fetch(jj)->Promise.then(ld => {
//         let resps = [ld]
//         resps->Js.Array2.forEach(r => Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok))
//         switch resps->Js.Array2.every(r => r.ok) {
//         | true => Ok(resps->Js.Array2.map(r => r->Prefetch.blob))
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
//       }->Promise.resolve
//     )
//     ->Promise.catch((. e) => {
//       let msg = switch e {
//       | Promise.JsError(err) =>
//         switch Js.Exn.message(err) {
//         | Some(msg) => msg
//         | None => ""
//         }
//       | _ => "Unexpected error occurred"
//       }
//       Js.log2("Fetch error: ", msg)
//       Error(msg)->Promise.resolve
//     })
//   ->ignore

// let handler = (assets, _e) => ldhand(assets)
// let pp = handler(loading)
// addWindowEventListener("load", pp)

Js.log2("ret", retina.matches)
Js.log2("mobb", mob.matches)
Js.log2("bigg", big.matches)

// Js.Global.setTimeout(() => {
//     Js.log("7")
//   removeWindowEventListener("load", getPic)
// }, 7000)->ignore

switch ReactDOM.querySelector("#root") {
| Some(root) => ReactDOM.render(<App />, root)
| None => ()
}
