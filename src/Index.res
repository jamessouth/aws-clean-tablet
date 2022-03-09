%%raw(`import 'virtual:windi.css'`)
%%raw(`import 'virtual:windi-devtools'`)
%%raw(`import './css/windi.css'`)

// type s
// @send external hasOwnProperty: (s, string) => bool = "hasOwnProperty"
// @val external css: s = "CSS"

// type pw = {addModule: (. string) => Js.Promise.t<unit>}
// @val @scope("CSS")
// external paintWorklet: pw = "paintWorklet"

// switch css->hasOwnProperty("paintWorklet") {
// | true => paintWorklet.addModule(. "paint.js")->Js.Promise.catch(err => {
//     Js.log2("Error loading worklet:", err)
//     Js.Promise.resolve()
//   }, _)->ignore
// | false => Js.log("I am Firefox or Safari")
// }

switch ReactDOM.querySelector("#root") {
| Some(root) => ReactDOM.render(<App />, root)
| None => ()
}
