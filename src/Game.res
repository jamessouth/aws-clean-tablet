type sendPayload = {
    action: string,
    gameno: string,
    type_: string,
    value: bool
}



let chk = Js.String2.fromCharCode(10003)

@react.component
let make = (~game: Types.game, ~leadertoken: string) => {//~ingame, ~send
  let gameReady = switch game.leader {
  | Some(_) => true
  | None => false
  }

//   let onClick1 = _ => {
//       let pl = {
//           action: "lobby",
//           gameno: j`${game.no}`,
//           type_: "leave or join",
//           value: true

//       }
//       pl->send
//       //if..... true->setReady
    
//   }

//   let onClick2 = _ => {
//       let pl = {
//           action: "lobby",
//           gameno: j`${game.no}`,
//           type_: "leave or join",
//           value: true
//       }
//       pl->send
//       //if..... !ready->setReady
    
//   }

    let ldr = switch game.leader {
    | Some(l) => l->Js.String2.split("_")
    | None => [""]
    }

//   let leaderName = switch gameReady {
//   | true => ldr[0]
//   | false => ""
//   }
  
  

  let (ready, _setReady) = React.useState(_ => true)
  let (count, _setCount) = React.useState(_ => 5)
  let (_leader, setLeader) = React.useState(_ => false)
//   let (disabled1, setDisabled1) = React.useState(_ => false)
//   let (disabled2, setDisabled2) = React.useState(_ => false)
//   let (startGame, setStartGame) = React.useState(_ => false)

//   let chkstyl = " text-2xl font-bold leading-3"

    React.useEffect2(() => {
    switch (game.leader, ldr[0] === leadertoken) {
    | (Some(_), true) => setLeader(_ => true)
    | (Some(_), false) | (None, _) => setLeader(_ => false)
    }
    None
  }, (game.leader, leadertoken))



    // React.useEffect3(() => {
    //     switch (ingame->Js.Nullable.toOption, ingame === game.no) {
    //     | (Some(g), false) => setDisabled1(_ => true)
    //     | None => setDisabled1(_ => false)
    //     }
    //     None
    // }, (ingame, game.no, game.players))



    // React.useEffect3(() => {
    //     switch (ingame->Js.Nullable.toOption, ingame === game.no) {
    //     | (Some(g), false) => setDisabled2(_ => true)
    //     | None => setDisabled2(_ => false)
    //     }
    //     None
    // }, (ingame, game.no, game.players))
   
   
   
//     React.useEffect3(() => {
//         let id = 0
//         switch (gameReady, game.no === ingame) {
//         | (true, true) => let id = setInterval(() => {
//             setCount(c => c - 1)
//         }, 1000)
//         | (_, _) => ()
//         }
//         None
//   }, (gameReady, game.no, ingame))



//     React.useEffect4(() => {
//         switch (ingame === game.no, count === 0, leader) {
//         | (true, true, true) => {
//             setStartGame(true)
//             send
//             ()
//         }
//         | _ => ()
//         }


//     None
//   }, (count, game.no, ingame, leader))



    <li
        className="mb-8 w-10/12 mx-auto grid grid-cols-2 grid-rows-gamebox relative pb-8"
    >
        <p className="text-xs col-span-2">
            {game.no->React.string}
        </p>

        <p className="text-xs col-span-2">
            {"players"->React.string}
        </p>

    // <>
    //     {
    //         game.players->Js.Array2.map((p) => {

    //         switch p.ready {
    //         | true => <p key=p.connid>{p.name->React.string}<span className={switch leaderName === p.name {
    //         | true => `text-red-200${chkstyl}`
    //         | false => `text-green-200${chkstyl}`
    //         }}>{chk->React.string}</span></p>
    //         | false => <p key=p.connid>{p.name->React.string}</p>
    //         }
    //         })
    //     }
    // </>



        {
            switch gameReady {
            | true => switch true { //ingame === game.no
                | true => <p className="absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4">{count->Js.Int.toString->React.string}</p>
                | false => <p className="absolute text-yellow-200 text-2xl font-bold left-1/2 bottom-1/4 transform -translate-x-2/4">{"Starting soon..."->React.string}</p>
                }
            | false => React.null
            }
        }


        <button
            className="w-1/2 bottom-0 h-8 left-0 absolute pt-2 bg-smoke-700 bg-opacity-70"
            // disabled=disabled1
            // onClick1
        >
            // {
                // switch value {
                // | pattern1 => expression
                // | pattern2 => expression
                // }
            // }
        </button>
        <button
            className="w-1/2 bottom-0 h-8 right-0 absolute pt-2 bg-smoke-700 bg-opacity-70"
            // disabled=disabled2
            // onClick2
        >{switch ready {
        | true => "ready"->React.string
        | false => "not ready"->React.string
        }}</button>
    </li>




}
