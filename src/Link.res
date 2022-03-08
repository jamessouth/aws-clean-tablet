@react.component
let make = (~url, ~className, ~content) => {
  let onClick = e => {
    e->ReactEvent.Mouse.preventDefault
    url->RescriptReactRouter.push
  }
  Js.log(url)
  <a onClick className href={url}> {React.string(content)} </a>
}
