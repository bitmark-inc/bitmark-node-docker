var gulp = require("gulp"),
    browserify = require("browserify"),
    vueify = require('vueify'),
    source = require('vinyl-source-stream'),
    hmr = require('browserify-hmr'),
    babelify = require('babelify'),

    webserver = require('gulp-webserver'),
    notify = require('gulp-notify'),
    useref = require('gulp-useref'),
    clean = require('gulp-clean'),
    gulpif = require('gulp-if'),
    rev = require("gulp-rev"),
    rename = require('gulp-rename'),
    revReplace = require("gulp-rev-replace");

gulp.task("dev", ["clean", "html"], function () {
  const b = browserify('./src/main.js')
    .plugin(hmr)
    .transform(vueify)
    .transform(babelify, {presets: ["es2015"]})

  function bundle() {
    b.bundle()
    .on('error', function(err){
      console.error('' + err);
      gulp.src('').pipe(notify('✖ Bunlde Failed ✖'));
      this.emit('end');
    })
    .pipe(source("bundle.js"))
    .pipe(gulp.dest("./public/assets/js/"));
  }

  b.bundle()
    .pipe(source("bundle.js"))
    .pipe(gulp.dest("./public/assets/js/"))
    .on('end', function () {
      gulp.src('public')
      .pipe(webserver({
        fallback: 'index.html',
        livereload: false,
        directoryListing: false,
        open: false
      }));
    })
  gulp.watch(["src/index.html", "src/**/*.vue", "src/**/*.js", "!src/dist/**"],
              function(event) { bundle(); });
})

gulp.task('clean', function () {
  return gulp.src('public/assets/js/bundle-*', {read: false})
    .pipe(clean());
});

gulp.task("watch", ["clean", "fonts", "bundle-js", "html"], function() {
  gulp.watch(["src/index.html", "src/**/*.vue", "src/**/*.js",
              "!src/dist/**"],
             ["clean", "html", function(event) {
  }]);
});

gulp.task("html", ["bundle-js"], function(){
  return gulp.src("./src/index.html")
    .pipe(useref())
    .pipe(gulpif('*.css', rev()))
    .pipe(revReplace())
    .pipe(gulp.dest("./public"));
});

gulp.task("html-bundle", ["bundle-js"], function(){
  return gulp.src("./src/index.html")
    .pipe(useref())
    .pipe(gulpif('*.css', rev()))
    .pipe(gulpif('*.js', rev()))
    .pipe(revReplace())
    .pipe(gulp.dest("./public"));
});

gulp.task("bundle-js", function() {
  return browserify('./src/main.js')
    .transform(vueify)
    .transform(babelify, {presets: ["es2015"]})
    .bundle()
    .on('error', function(err){
      console.error('' + err);
      gulp.src('').pipe(notify('✖ Bunlde Failed ✖'));
      this.emit('end');
    })
    .pipe(source("bundle.js"))
    .pipe(gulp.dest("./src/dist"));
});

gulp.task("bundle", ["clean", "bundle-js", "html-bundle"]);
gulp.task("default", ["clean", "bundle-js", "html"]);
