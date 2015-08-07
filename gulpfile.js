var browserify = require('browserify');
var buffer     = require('vinyl-buffer');
var concat     = require('gulp-concat');
var gulp       = require('gulp');
var gutil      = require('gulp-util');
var minifyCSS  = require('gulp-minify-css');
var reactify   = require('reactify');
var rename     = require('gulp-rename');
var sass       = require('gulp-sass');
var source     = require('vinyl-source-stream');
var sourcemaps = require('gulp-sourcemaps');
var uglify     = require('gulp-uglify');

gulp.task('script', function () {
    // set up the browserify instance on a task basis
    var b = browserify({
        entries: 'script/site.jsx',
        debug: true,
        // defining transforms here will avoid crashing your stream
        transform: [reactify]
    });

    return b.bundle()
        .pipe(source('script/site.jsx'))
        .pipe(buffer())
        .pipe(sourcemaps.init({loadMaps: true}))
            // Add transformation tasks to the pipeline here.
            .pipe(uglify())
            .on('error', gutil.log)
        .pipe(rename('site.min.js'))
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest('static/'));
});

gulp.task('style', function () {
    gulp.src('./style/**/*.scss')
        .pipe(sass().on('error', sass.logError))
        .pipe(concat('style.css'))
        .pipe(minifyCSS())
        .pipe(rename('site.min.css'))
        .pipe(gulp.dest('./static/'));
});

gulp.task('watch', function () {
    gulp.watch('./style/**/*.scss', ['style']);
    gulp.watch('./script/**/*.jsx', ['script']);
});
