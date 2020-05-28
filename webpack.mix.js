const mix = require('laravel-mix');
require('laravel-mix-purgecss');

// ...
mix.js('resources/assets/js/app.js', 'public/js')
    .sass('resources/assets/sass/app.scss', 'public/css')
    .copy('resources/assets/images', 'public/images')
    .extract()
    .purgeCss()
    .browserSync('http://localhost:8080');