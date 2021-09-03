
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