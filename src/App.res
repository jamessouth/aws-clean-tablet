
type player = {
    name: string,
    connid: string,
    ready: bool
}

type game = {
    leader: option<string>,
    players: array<player>,
    no: string
}

@react.component
let make = () => {


    <div className="mt-8"></div>
}