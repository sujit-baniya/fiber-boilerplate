const mix = require('laravel-mix');
require('laravel-mix-purgecss');
const tailwindcss = require('tailwindcss');
// ...
mix.js('resources/assets/js/app.js', 'assets/js')
    .sass('resources/assets/sass/app.scss', 'assets/css')
    .sass('resources/assets/sass/vendor.scss', 'assets/css')
    .copy('resources/assets/images', 'assets/images')
    .options({
        processCssUrls: false,
        postCss: [tailwindcss('./tailwind.config.js')],
    })
    .extract()
    .purgeCss()
    .browserSync('http://localhost:8080');