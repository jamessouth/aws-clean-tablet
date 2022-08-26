
type propShape = {"bghand": (.unit) => unit}

@val
external import_: string => Promise.t<propShape> = "import"



open Web

let mob = matchMedia("(max-width: 767.9px)")
let big = matchMedia("(min-width: 1440px)")

let retinaMQ = switch matchMedia("(resolution: 2x)").media != "not all" {
| true => "(min-resolution: 2x)"
| false => "(-webkit-min-device-pixel-ratio: 2)"
}

let retina = matchMedia(retinaMQ)

let chalkbg = switch mob.matches {
| true =>
  switch retina.matches {
  | true => "../../assets/chmob2x.webp"
  | false => "../../assets/chmob1x.webp"
  }
| false =>
  switch big.matches {
  | true =>
    switch retina.matches {
    | true => "../../assets/chbig2x.webp"
    | false => "../../assets/chbig1x.webp"
    }
  | false =>
    switch retina.matches {
    | true => "../../assets/chtab2x.webp"
    | false => "../../assets/chtab1x.webp"
    }
  }
}

let newgamebg = switch retina.matches {
| true => "../../assets/ekko2x.webp"
| false => "../../assets/ekko1x.webp"
}

let signout = "../../assets/signout.png"
let leader = "../../assets/leader.png"

let bghand = (._) => Prefetch.getPicsAllSettled4((chalkbg, newgamebg, signout, leader))