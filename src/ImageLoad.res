
type propShape = {"ppp": (int,int) => int}

@val
external import_: string => Promise.t<propShape> = "import"


let ppp = i => j => i + j