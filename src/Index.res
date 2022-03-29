%%raw(`import 'virtual:windi.css'`)
%%raw(`import 'virtual:windi-devtools'`)
%%raw(`import './css/windi.css'`)
@scope("window") @val
external addWindowEventListener: (string, unit => unit) => unit = "addEventListener"
@scope("window") @val
external removeWindowEventListener: (string, unit => unit) => unit = "removeEventListener"

type mediaQueryList = {
  matches: bool,
  media: string,
}


@scope("window") @val
external matchMedia: string => mediaQueryList = "matchMedia"


switch ReactDOM.querySelector("#root") {
| Some(root) => ReactDOM.render(<App />, root)
| None => ()
}

module Promise = {
  type t<+'a> = Js.Promise.t<'a>
  exception JsError(Js.Exn.t)
  external unsafeToJsExn: exn => Js.Exn.t = "%identity"
  @val @scope("Promise")
  external resolve: 'a => t<'a> = "resolve"
  @send external then: (t<'a>, @uncurry ('a => t<'b>)) => t<'b> = "then"
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
  @send external blob: resp => Promise.t<blob> = "blob"
  @val external fetch: string => Promise.t<Js.Nullable.t<resp>> = "fetch"
  let getPic = (asset) => {
    open Promise
    Js.log("load ev")
    fetch(asset)
    ->then(res => {
      Js.log2("res: ", res)
      switch Js.Nullable.toOption(res) {
      | None => Error("uh oh")
      | Some(r) =>
        switch r.ok {
        | true => Ok(r->blob)
        | false => {
            let stat = r.status
            Error(j`Fetch error: $stat - ${r.statusText}`)
          }
        }
      }->resolve
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
      Error(msg)->resolve
    })
  }->ignore
}

let handler = (asset, _e) => Prefetch.getPic(asset)

let bghand = handler("../../assets/chmob2x.webp")


addWindowEventListener("load", bghand)





let mob = matchMedia("(max-width: 767.9px)")
let tab = matchMedia("(min-width: 768px) and (max-width: 1439.9px)")
let big = matchMedia("(min-width: 1440px)")

// let getMedia = e => switch e.matches {
// | true => Js.log2("match", e)
// | false => Js.log2("no match", e)
// }
// @scope("mob") @val
// external addMobEventListener: (string, mediaQueryList => unit) => unit = "addEventListener"
// @scope("tab") @val
// external addTabEventListener: (string, mediaQueryList => unit) => unit = "addEventListener"
// @scope("big") @val
// external addBigEventListener: (string, mediaQueryList => unit) => unit = "addEventListener"


// addMobEventListener("change", getMedia)
// addTabEventListener("change", getMedia)
// addBigEventListener("change", getMedia)
Js.log2("mobb", mob.matches)
Js.log2("tabb", tab.matches)
Js.log2("bigg", big.matches)

// Js.Global.setTimeout(() => {
//     Js.log("7")
//   removeWindowEventListener("load", getPic)
// }, 7000)->ignore
