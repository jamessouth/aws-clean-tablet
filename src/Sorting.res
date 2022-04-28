type sortDirection =
  | Up
  | Down
let strSortName = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up =>
    switch n2.name > n1.name {
    | true => 1
    | false => -1
    }
  | Down =>
    switch n2.name <= n1.name {
    | true => 1
    | false => -1
    }
  }
let intSortWins = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up => n2.wins - n1.wins
  | Down => n1.wins - n2.wins
  }
let intSortPoints = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up => n2.totalPoints - n1.totalPoints
  | Down => n1.totalPoints - n2.totalPoints
  }
let intSortGames = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up => n2.games - n1.games
  | Down => n1.games - n2.games
  }
let fltSortWinPct = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up =>
    switch n2.winPct > n1.winPct {
    | true => 1
    | false => -1
    }
  | Down =>
    switch n2.winPct <= n1.winPct {
    | true => 1
    | false => -1
    }
  }
let fltSortPPG = (dir: sortDirection, n1: Reducer.stat, n2: Reducer.stat) =>
  switch dir {
  | Up =>
    switch n2.ppg > n1.ppg {
    | true => 1
    | false => -1
    }
  | Down =>
    switch n2.ppg <= n1.ppg {
    | true => 1
    | false => -1
    }
  }
let sortData = (input, dir, data, func) => {
  let arr = Js.Array2.copy(data)
  switch input {
  | "name" => arr->Js.Array2.sortInPlaceWith(strSortName(dir))
  | "wins" => arr->Js.Array2.sortInPlaceWith(intSortWins(dir))
  | "points" => arr->Js.Array2.sortInPlaceWith(intSortPoints(dir))
  | "games" => arr->Js.Array2.sortInPlaceWith(intSortGames(dir))
  | "win %" => arr->Js.Array2.sortInPlaceWith(fltSortWinPct(dir))
  | "pts/gm" => arr->Js.Array2.sortInPlaceWith(fltSortPPG(dir))
  | _ => []
  }->ignore
  func(._ => arr)
}
