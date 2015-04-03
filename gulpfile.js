var gulp = require('gulp');
var sass = require('gulp-ruby-sass');
var react = require('gulp-react');
var uglify = require('gulp-uglify');
var mini = require('gulp-minify-css');

gulp.task('sass', function(){
    return sass('resources/scss/base.scss')
    //.pipe(mini())
    .pipe(gulp.dest('public/css/'))
});

gulp.task('jsx', function(){
    return gulp.src('resources/js/components/**/*.jsx')
    .pipe(react())
    .pipe(gulp.dest('public/js/components/'))
});

gulp.task('compressJs', function(){
    return gulp.src('resources/js/**/*.js')
    //.pipe(utlify())
    .pipe(gulp.dest('public/js/'))
});

gulp.task('bootstrap', function(){
    return sass('public/components/bootstrap-sass/assets/stylesheets/_bootstrap.scss')
    .pipe(gulp.dest('public/components/bootstrap-sass/assets/stylesheets'))
});

gulp.task('watch', function(){
    gulp.watch('resources/scss/**/*.scss', ['sass']);
    gulp.watch('resources/js/components/**/*.jsx', ['jsx']);
    gulp.watch('resources/js/**/*.js', ['compressJs'])
});

gulp.task('default', ['watch', 'sass', 'compressJs', 'jsx']);