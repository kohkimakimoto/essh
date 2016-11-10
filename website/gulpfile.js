"use strict";

var gulp = require('gulp');
var shell = require('gulp-shell')
var concat = require("gulp-concat");
var sass  = require('gulp-sass');
var cleanCSS = require('gulp-clean-css');
var buffer = require('vinyl-buffer');
var source = require('vinyl-source-stream');
var browserify = require('browserify');
var webserver = require('gulp-webserver');

var destDir = "public";

// hugo
gulp.task('hugo', shell.task([
  'hugo',
]))

// css
gulp.task('css', function () {
    gulp.src('src/**/*.scss').
    pipe(sass().on('error', sass.logError)).
    pipe(gulp.dest(destDir))
    ;
});

// js
gulp.task('js', function () {
    browserify({
            entries: ['src/index.js'],
            extensions: ['.js', '.jsx'],
            debug: true,
        })
        .bundle()
        .on("error", function (err) {
            console.log("Error : " + err.message);
        })
        .pipe(source('bundle.js'))
        .pipe(buffer())
        .pipe(gulp.dest(destDir))
        ;
});

// html
gulp.task('html', () => {
  gulp.src('src/**/*.html')
    .pipe(gulp.dest(destDir));
});

gulp.task('webserver', function() {
  gulp.watch(['content/**/*', 'layout/**/*', 'src/**/*', 'static/**/*'], ['build']);
  gulp.src(destDir)
    .pipe(webserver({
      livereload: true,
      open: true,
    }));
});

gulp.task('build', ['hugo', 'css', 'js', 'html']);
gulp.task('serve', ['build', 'webserver']);
gulp.task('default', ['build']);
