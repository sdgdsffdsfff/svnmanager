var gulp = require('gulp');
var sass = require('gulp-ruby-sass'),
    react = require('gulp-react'),
    coffee = require('gulp-coffee'),
    cjsx = require('gulp-cjsx'),
    mini = require('gulp-minify-css');

gulp.task('sass', function(){
    return sass('scss/base.scss')
    //.pipe(mini())
    .pipe(gulp.dest('css/'))
});

gulp.task('jsx', function(){
    return gulp.src('./js/jsx/**/*.jsx')
    .pipe(react())
    .pipe(gulp.dest('./js/components/'))
});

gulp.task('bootstrap', function(){
    return sass('./components/bootstrap-sass/assets/stylesheets/_bootstrap.scss')
    .pipe(gulp.dest('./components/bootstrap-sass/assets/stylesheets'))
});

gulp.task('watch', function(){
    gulp.watch('./scss/**/*.scss', ['sass']);
    gulp.watch('./components/bootstrap-sass/**/*.scss', ['bootstrap']);
    gulp.watch('./js/jsx/**/*.jsx', ['jsx']);
});

gulp.task('default', ['watch', 'jsx']);