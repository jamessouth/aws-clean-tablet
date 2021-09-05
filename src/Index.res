%%raw(`import 'virtual:windi.css'`)
%%raw(`import './css/windi.css'`)

switch (ReactDOM.querySelector("#root")) {
| Some(root) => ReactDOM.render(<App/>, root)
| None => ()
}