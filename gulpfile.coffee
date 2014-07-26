gulp = require 'gulp'
$ = require( 'gulp-load-plugins' )()
nib = require 'nib'
browserify = require 'browserify'
source = require 'vinyl-source-stream'
coffeeify = require 'coffee-reactify'

filterMinifiedFiles = $.filter (file) ->
  !~file.path.indexOf '.min'

gulp.task 'stylus', ->
  gulp.src 'public/styl/*.styl'
    .pipe $.stylus
      use: [nib()]
    .pipe gulp.dest 'public/css'

gulp.task 'browserify', ->
  browserify
    extensions: ['.coffee', '.cjsx']
    debug: true
  .transform coffeeify
  .require './public/coffee/app.coffee', entry: true
  .bundle()
  .pipe source 'app.js'
  .pipe gulp.dest 'public/js'

gulp.task 'uglify', ['browserify'], ->
  gulp.src 'public/js/**/*.js'
    .pipe filterMinifiedFiles
    .pipe $.uglify()
    .pipe $.rename suffix: '.min'
    .pipe gulp.dest 'public/js'

gulp.task 'cssmin', ['stylus'], ->
  gulp.src 'public/css/*.css'
    .pipe filterMinifiedFiles
    .pipe $.minifyCss()
    .pipe $.rename suffix: '.min'
    .pipe gulp.dest 'public/css'

gulp.task 'watch', ->
  gulp.watch 'public/styl/**/*.styl', ['stylus']
  gulp.watch 'public/coffee/**/*.coffee', ['browserify']

gulp.task 'clean-css', ->
  gulp.src 'public/css/**/*.css', read: false
    .pipe $.rimraf()

gulp.task 'clean-js', ->
  gulp.src 'public/js/**/*.js', read: false
    .pipe $.rimraf()

gulp.task 'clean', ['clean-css', 'clean-js']
gulp.task 'build', ['uglify', 'cssmin']
gulp.task 'default', ['build']