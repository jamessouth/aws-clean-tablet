type t
type openEventHandler = unit => unit
type errorEventHandler = Dom.errorEvent => unit
type messageEvent = {data: string}
type messageEventHandler = messageEvent => unit
type closeEvent = {
  code: int,
  reason: string,
  wasClean: bool,
}
type closeEventHandler = closeEvent => unit
@new external newWs: string => t = "WebSocket"
@set external onOpen: (Js.Nullable.t<t>, openEventHandler) => unit = "onopen"
@set external onError: (Js.Nullable.t<t>, errorEventHandler) => unit = "onerror"
@set external onMessage: (Js.Nullable.t<t>, messageEventHandler) => unit = "onmessage"
@set external onClose: (Js.Nullable.t<t>, closeEventHandler) => unit = "onclose"
@send external closeCode: (Js.Nullable.t<t>, int) => unit = "close"
@send external closeCodeReason: (Js.Nullable.t<t>, int, string) => unit = "close"
@send external sendString: (Js.Nullable.t<t>, string) => unit = "send"
@val external document: Dom.document = "document"
@get external body: Dom.document => Dom.htmlBodyElement = "body"
@set external setClassName: (Dom.htmlBodyElement, string) => unit = "className"
@send external blur: Dom.element => unit = "blur"
@scope("window") @val
external addWindowEventListener: (string, unit => unit) => unit = "addEventListener"
@scope("window") @val
external removeWindowEventListener: (string, unit => unit) => unit = "removeEventListener"
type mediaQueryList = {
  matches: bool,
  media: string,
}
@scope("window") @val
external matchMedia: string => mediaQueryList = "matchMedia"
