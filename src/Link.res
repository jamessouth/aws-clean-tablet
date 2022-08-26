@react.component
let make = (~url, ~className, ~content="", ~image=React.null) => {
  let onClick = e => {
    ReactEvent.Mouse.preventDefault(e)
    RescriptReactRouter.push(url)
    switch url {
    | "/signin" =>
      ImageLoad.import_("./ImageLoad.bs")
      ->Promise.then(func => {
        Promise.resolve(func["bghand"](.))
      })
      ->ignore
    | _ => ()
    }
  }

  <a onClick className href={url}>
    {React.string(content)}
    {image}
  </a>
}
