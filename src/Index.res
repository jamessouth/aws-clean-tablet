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

let asset = switch mob.matches {
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

let bghand = Prefetch.handler(asset)

addWindowEventListener("load", bghand)

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
