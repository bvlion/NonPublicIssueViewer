const gulp = require('gulp')
const uglify = require('gulp-uglify')
const babel = require('gulp-babel')
const cleanCSS = require('gulp-clean-css')
const plumber = require('gulp-plumber')
const notify  = require('gulp-notify')
const rename = require('gulp-rename')
const postcss = require("gulp-postcss")
const autoprefixer = require("autoprefixer")

const cssSrc = 'css/*.css'
const cssTask = 'minify-css'
const jsSrc = 'js/*.js'
const jsTask = 'minify-js'

gulp.task(jsTask, () =>
  gulp.src(jsSrc)
    .pipe(babel({
      'presets': ['@babel/preset-env']
    }))
    .pipe(uglify())
    .pipe(rename({
      extname: '.min.js'
    }))
    .pipe(plumber({
      errorHandler: notify.onError('Error: <%= error.message %>')
    }))
    .pipe(gulp.dest('../source/public/js/'))
)

gulp.task(cssTask, () =>
  gulp.src(cssSrc)
    .pipe(postcss([
      autoprefixer({cascade: false})
    ]))
    .pipe(cleanCSS())
    .pipe(rename({
      extname: '.min.css'
    }))
    .pipe(plumber({
      errorHandler: notify.onError('Error: <%= error.message %>') 
    }))
    .pipe(gulp.dest('../source/public/css/'))
)

gulp.task('watch', () => {
  gulp.watch(jsSrc, gulp.series(jsTask))
  gulp.watch(cssSrc, gulp.series(cssTask))
})