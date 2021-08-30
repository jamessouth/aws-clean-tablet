
type authorized = Auth | Guest;
let wsContext = React.createContext(Guest)


let provider = React.Context.provider(wsContext)

@react.component
let make = (~value, ~children) => {
    React.createElement(provider, {"value": value, "children": children})
}

