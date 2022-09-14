%%raw(`import 'virtual:windi-utilities.css'`)
// %%raw(`import 'virtual:windi-devtools'`)
%%raw(`import './css/index.css'`)

switch ReactDOM.querySelector("#root") {
| Some(root) => ReactDOM.render(<App />, root)
| None => ()
}


let loading = "./src/Loading.bs.js"

let ldhand = res => Promise.fetch(res)->Promise.then(ld => {
        let resps = [ld]
        resps->Js.Array2.forEach(r => Js.log2("Asset " ++ r.url ++ " fetched ok: ", r.ok))
        switch resps->Js.Array2.every(r => r.ok) {
        | true => Ok(resps->Js.Array2.map(r => r->Promise.blob))
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
      }->Promise.resolve
    )
    ->Promise.catch((. e) => {
      let msg = switch e {
      | Promise.JsError(err) =>
        switch Js.Exn.message(err) {
        | Some(msg) => msg
        | None => ""
        }
      | _ => "Unexpected error occurred"
      }
      Js.log2("Fetch error: ", msg)
      Error(msg)->Promise.resolve
    })
  ->ignore

let handler = (assets, _e) => ldhand(assets)
Web.addWindowEventListener("load", handler(loading))