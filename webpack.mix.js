const mix = require('laravel-mix');
const tailwindcss = require('tailwindcss');
// ...
mix.js('resources/src/js/app.js', 'resources/assets/js')
    .js('resources/src/js/payment.js', 'resources/assets/js')
    .sass('resources/src/sass/app.scss', 'resources/assets/css')
    .sass('resources/src/sass/vendor.scss', 'resources/assets/css')
    .copy('resources/src/images', 'resources/assets/images')
    .options({
        processCssUrls: false,
        postCss: [tailwindcss('./tailwind.config.js')],
    })
    .extract()
    .browserSync('http://localhost:8080');