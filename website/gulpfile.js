"use strict";

var gulp = require('gulp');
var shell = require('gulp-shell')
var concat = require("gulp-concat");
var sass  = require('gulp-sass');
var cleanCSS = require('gulp-clean-css');
var uglify = require('gulp-uglify');
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
    gulp.src('src/scss/style.scss')
    .pipe(sass().on('error', sass.logError))
    .pipe(cleanCSS())
    .pipe(gulp.dest(destDir))
    ;
});

// js
gulp.task('js', function () {
    // bundle.js
    browserify({
            entries: ['src/js/index.js'],
            extensions: ['.js', '.jsx'],
            debug: true,
        })
        .bundle()
        .on("error", function (err) {
            console.log("Error : " + err.message);
        })
        .pipe(source('bundle.js'))
        .pipe(buffer())
        .pipe(uglify({preserveComments: 'some'}))
        .pipe(gulp.dest(destDir))
        ;
});

// html
gulp.task('html', () => {
  gulp.src('src/**/*.html')
    .pipe(gulp.dest(destDir));
});

// font
gulp.task('font', () => {
  gulp.src('node_modules/font-awesome/fonts/*')
    .pipe(gulp.dest(destDir + "/fonts"));
  gulp.src('node_modules/simple-line-icons/fonts/*')
    .pipe(gulp.dest(destDir + "/fonts"));
});


gulp.task('webserver', function() {
  gulp.watch(['content/**/*', 'layouts/**/*', 'src/**/*', 'static/**/*'], ['build']);
  gulp.src(destDir)
    .pipe(webserver({
      livereload: true,
      open: true,
    }));
});

gulp.task('build', ['hugo', 'css', 'js', 'html', 'font']);
gulp.task('serve', ['build', 'webserver']);
gulp.task('default', ['build']);
