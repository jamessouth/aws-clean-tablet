@react.component
let make = (~url, ~className, ~content="", ~image=React.null) => {
  let onClick = e => {
    ReactEvent.Mouse.preventDefault(e)
    RescriptReactRouter.push(url)
  }

  <a onClick className href={url}> {React.string(content)} {image} </a>
}
