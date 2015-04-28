var gulp = require('gulp');
var sass = require('gulp-ruby-sass');
var react = require('gulp-react');
var uglify = require('gulp-uglify');
var colors = require('colors');
var mini = require('gulp-minify-css');

gulp.task('sass', function(){
    return sass('resources/scss/base.scss')
        .pipe(mini())
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

function getT( n ){
    return n > 9 ? n : '0' + n;
}

function info( path ){
    var s = (new Date).getMilliseconds();
    path = path.substring(path.indexOf('resources'), path.length);
    return function(){
        var d = new Date;
        var n = d.getMilliseconds() - s;
        var t = getT(d.getHours())+':'+ getT(d.getMinutes()) +':'+ getT(d.getSeconds());
        console.info('['+t.grey+']', 'Changed', path.cyan, 'after', (n+'ms').magenta);
    }
}

gulp.task('watch', function(){
    gulp.watch('resources/scss/**/*.scss', ['sass']);
    gulp.watch('resources/js/components/**/*.jsx', function(e){
        if( e.type === 'changed' ){
            var i = info(e.path);
            gulp.src(e.path, { base: 'resources/js/components/'})
                .pipe(react())
                .pipe(gulp.dest('public/js/components/'));
            i();
        }
    });
    gulp.watch('resources/js/**/*.js', function(e){
        if( e.type === 'changed' ){
            var i = info(e.path);
            gulp.src(e.path, { base: 'resources/js/'})
                //.pipe(uglify())
                .pipe(gulp.dest('public/js/'));
            i();
        }
    });
});

gulp.task('default', ['watch', 'sass', 'compressJs', 'jsx']);