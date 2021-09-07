


@react.component
let make = (~url, ~className, ~content) => {
    let onClick = e => {
        e->ReactEvent.Mouse.preventDefault
        url->RescriptReactRouter.push
    }

    <a onClick className href={url}>{content->React.string}</a>
}