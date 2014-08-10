gulp = require 'gulp'
$ = require('gulp-load-plugins')()
nib = require 'nib'
browserify = require 'browserify'
source = require 'vinyl-source-stream'
stringify = require 'stringify'
coffeeify = require 'coffeeify'
path = require 'path'

filterMinifiedFiles = $.filter (file) ->
  extname = path.extname file.path
  src = file.path.substring 0, file.path.length - extname.length

  !/\.min$/.test src

gulp.task 'stylus', ->
  gulp.src 'public/styl/*.styl'
    .pipe $.filter (file) ->
      !/\/_/.test file.path
    .pipe $.stylus
      use: [nib()]
    .pipe gulp.dest 'public/css'

gulp.task 'browserify', ->
  browserify
    extensions: ['.coffee', '.html']
    debug: true
  .transform stringify ['.html']
  .transform coffeeify
  .require './public/coffee/app.coffee', entry : true
  .bundle()
  .pipe source 'app.js'
  .pipe gulp.dest 'public/js'

gulp.task 'minify-js', ['browserify'], ->
  gulp.src 'public/js/**/*.js'
    .pipe filterMinifiedFiles
    .pipe $.ngAnnotate()
    .pipe $.uglify()
    .pipe $.rename suffix: '.min'
    .pipe gulp.dest 'public/js'

gulp.task 'minify-css', ['stylus'], ->
  gulp.src 'public/css/*.css'
    .pipe filterMinifiedFiles
    .pipe $.minifyCss()
    .pipe $.rename suffix: '.min'
    .pipe gulp.dest 'public/css'

gulp.task 'watch', ->
  gulp.watch 'public/styl/**/*.styl', ['stylus']
  gulp.watch [ 'public/coffee/**/*.coffee', 'public/views/**/*.html' ], ['browserify']

gulp.task 'clean-css', ->
  gulp.src 'public/css/**/*.css', read: false
    .pipe $.rimraf()

gulp.task 'clean-js', ->
  gulp.src 'public/js/**/*.js', read: false
    .pipe $.rimraf()

gulp.task 'clean', ['clean-css', 'clean-js']
gulp.task 'build', ['minify-js', 'minify-css']
gulp.task 'default', ['build']
