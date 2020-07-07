const inViewport = (entries, _) =>
  entries.forEach(entry => {

    if (entry.intersectionRatio > 0) {
      const imgEl = entry.target
      imgEl.src = imgEl.dataset.src

      imgEl.addEventListener('load', () =>
        imgEl.classList.add('is-lazyloaded')
      );

      io.unobserve(entry.target)
    }
  })

const io = new IntersectionObserver(inViewport, {
  threshold: [0]
})

const load_image = () =>
  Array.from(document.querySelectorAll('.lazyload')).forEach(element => io.observe(element))

const check_login = (error, next) => {
  if (error == 'login') {
    Swal.fire({
      title: 'セッションが切れました',
      html: "ログイン画面に遷移します。<br>あらためてログインをお願いいたします。",
      icon: 'warning',
      allowOutsideClick: false,
      allowEscapeKey: false,
      confirmButtonColor: '#33',
      confirmButtonText: 'OK'
    }).then((_) => {
      location.href = '/login'
    })
  } else {
    next()
  }
}

const load_one_data = (date) => {
  Swal.queue([{
    title: 'Loading...',
    allowOutsideClick: false,
    allowEscapeKey: false,
    showConfirmButton: false,
    showCloseButton: false,
    showCancelButton: false,
    onOpen: () => {
      Swal.showLoading()
      return fetch('/detail/' + date)
        .then(response => response.json())
        .then(data => {
          Swal.hideLoading()
          Swal.close()
          check_login(data.error, () => {
            Swal.fire({
                type: 'success',
                title: data.title,
                html: marked(data.body),
                showConfirmButton: false,
                showCancelButton: true,
                cancelButtonText: '閉じる'
            })
            load_image()
          })
        })
        .catch(() => {
          Swal.hideLoading()
          Swal.insertQueueStep({
            title: '通信エラー',
            html: '通信エラーが発生しました。<br>繰り返し発生する場合は、お手数ですが管理者までお問い合わせください。。。',
            icon: 'warning',
            showConfirmButton: false,
            showCancelButton: true,
            cancelButtonText: '閉じる'
          })
        })
    }
  }])
}

const load_message = (minusMonth) => {
  set_meal_html(null)
  fetch('/issues/' + minusMonth, {
    method: 'GET',
    headers: {'Accept': 'application/json', 'Content-Type': 'application/json'}
  })
  .then(response => response.json())
  .then(json => {
    if (json.error) {
      check_login(json.error, () =>
        Swal.fire({
          title: '通信エラー',
          html: '通信エラーが発生しました。<br>繰り返し発生する場合は、お手数ですが管理者までお問い合わせください。。。',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#3085d6',
          cancelButtonColor: '#d33',
          confirmButtonText: 'Reload'
        }).then((_) => load_message(minusMonth))
      )
    } else {
      set_meal_html(json)
    }
  })
}

const set_meal_html = (json) => {
  if (json == null) {
    document.querySelector('#breakfast_progress').style.display = 'block'
    document.querySelector('#breakfast_area').innerHTML = ''
    document.querySelector('#lunch_progress').style.display = 'block'
    document.querySelector('#lunch_area').innerHTML = ''
    document.querySelector('#dinner_progress').style.display = 'block'
    document.querySelector('#dinner_area').innerHTML = ''
    document.querySelector('#impressionistic_progress').style.display = 'block'
    document.querySelector('#impressionistic_area').innerHTML = ''
  } else {
    document.querySelector('#breakfast_progress').style.display = 'none'
    document.querySelector('#breakfast_area').innerHTML = create_meal_html(json.Breakfasts)
    document.querySelector('#lunch_progress').style.display = 'none'
    document.querySelector('#lunch_area').innerHTML = create_meal_html(json.Lunchs)
    document.querySelector('#dinner_progress').style.display = 'none'
    document.querySelector('#dinner_area').innerHTML = create_meal_html(json.Dinners)
    document.querySelector('#impressionistic_progress').style.display = 'none'
    document.querySelector('#impressionistic_area').innerHTML = create_meal_html(json.Others, true)
  }
  load_image()
} 

const create_meal_html = (meal, isMark = false) => '<section class="mdl-grid mdl-grid--no-spacing">' +
  meal
    .map(element => {
      let title = ''
      if (isMark) {
        title = marked(element.Content)
      } else {
        title = '<pre>' + element.Content + '</pre>'
      }
      return '<div class="demo-card-wide mdl-card card-margin mdl-shadow--2dp">' +
      '<div class="mdl-card__title">' +
      '<h2 class="mdl-card__title-text">' + element.Date + '</h2>' +
      '</div>' +
      '<div class="flex-grow">' + title + '</div>' +
      element.Image +
      '</div>'
    })
    .join() + '</section>'

const logout = () =>
  Swal.fire({
    title: 'ログアウトしますか？',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'ログアウト'
  }).then((result) => {
    if (result.value) {
      location.href = "/logout"
    }
  })

load_message(0)