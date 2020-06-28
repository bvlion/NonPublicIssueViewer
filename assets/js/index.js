const $lazy = document.querySelectorAll('.lazyload')
const io = new IntersectionObserver(inViewport, {
  threshold: [0]
})

Array.from($lazy).forEach(element => {
  io.observe(element);
});

function inViewport(entries, _) {
  entries.forEach(entry => {

    if(entry.intersectionRatio > 0){
      const imgEl = entry.target
      imgEl.src = imgEl.dataset.src

      imgEl.addEventListener('load', () =>
        imgEl.classList.add('is-lazyloaded')
      );

      io.unobserve(entry.target)
    }
  });
}