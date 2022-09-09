type propShape = {"bghand": (. unit) => unit}

@val
external import_: string => Promise.t<propShape> = "import"

let bghand: (. 'a) => unit
