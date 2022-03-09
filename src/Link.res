@react.component
let make = (~url, ~className, ~content) => {
  let onClick = e => {
    ReactEvent.Mouse.preventDefault(e)
    RescriptReactRouter.push(url)
  }
  Js.log(url)
  <a onClick className href={url}> {React.string(content)} </a>
}
