const login = () => {
  const passphrase = document.querySelector('#passphrase').value
  if (!passphrase) {
    Swal.fire('合言葉を入力してください。', '', 'error')
    return
  }
  document.querySelector('.loading-overlay').style.display = 'block'
  document.querySelector('.loading-message').style.display = 'block'
  document.querySelector('.loader').style.display = 'block'
  fetch('/login', {
    credentials: 'same-origin',
    method: 'POST',
    body: JSON.stringify({passphrase: passphrase}),
    headers: {'Accept': 'application/json', 'Content-Type': 'application/json'}
  }).then((response) => {
    return response.json()
  }).then((json) => {
    if (json.error) {
      document.querySelector('.loading-overlay').style.display = 'none'
      document.querySelector('.loading-message').style.display = 'none'
      document.querySelector('.loader').style.display = 'none'
      Swal.fire(json.error, '', 'error')
    } else {
      window.location = '/'
    }
  })
}

const enter = () => {
  if (window.event.keyCode === 13) {
    login()
  }
}

let view = false
const showChange = () => {
  const passphrase = document.querySelector('#passphrase')
  if (view) {
    passphrase.setAttribute("type", "password")
    view = false
  } else {
    passphrase.setAttribute("type", "text")
    view = true
  }
}