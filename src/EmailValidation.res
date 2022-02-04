let useEmailValidation = email => {
  let (err, setErr) = React.useState(_ => None)

  let checkEmail = em => {
    let r = %re(
      "/^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/"
    )
    switch Js.String2.match_(em, r) {
    | None => setErr(_ => Some("Please enter a valid email address."))
    | Some(_) => setErr(_ => None)
    }
  }

    let checkLength = em => {
    switch Js.String2.length(em) < 5 {
    | true => setErr(_ => Some("Email address is too short."))
    | false => checkEmail(em)
    }
  }

  React.useEffect1(() => {
    checkLength(email)
    None
  }, [email])

  err
}
