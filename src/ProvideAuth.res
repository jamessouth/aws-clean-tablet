
type authorized = Auth | Guest;
let authContext = React.createContext(Guest)


let provider = React.Context.provider(authContext)

@react.component
let make = (~value, ~children) => {
    React.createElement(provider, {"value": value, "children": children})
}

