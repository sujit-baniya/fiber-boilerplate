let mix = require('laravel-mix');
const del = require('del');
const path = require('path');
const tailwindcss = require('tailwindcss');

del('./public');
mix.webpackConfig({
    output: {chunkFilename: 'js/[name].js?id=[chunkhash]'},
    resolve: {
        alias: {
            vue$: 'vue/dist/vue.runtime.esm.js', // vue/dist/vue.esm-bundler.js
            '@': path.resolve('resources/src'),
        },
        extensions: ["*", ".js", ".jsx", ".vue", ".ts", ".tsx"]
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                loader: "ts-loader",
                options: {appendTsSuffixTo: [/\.vue$/]},
                exclude: /node_modules/
            }
        ]
    },
})
    .sourceMaps(false, 'source-map');

mix
    .setPublicPath('public')
    .webpackConfig({
        optimization: {
            providedExports: false,
            sideEffects: false,
            usedExports: false
        }
    })
    .vue({vue: 2})
    .sass('resources/src/sass/vendor.scss', 'css')

    .sass('resources/src/sass/landing.scss', 'css')

    .js('resources/src/js/app.js', 'public/js/')
    .sass('resources/src/sass/app.scss', 'css')

    .js('resources/src/js/payment.js', 'public/js/')

    .copyDirectory('resources/src/img', './public/img')
    .copyDirectory('resources/src/flags', './public/flags')
    .copyDirectory('resources/src/font', './public/font')
    .copy('resources/src/favicon.ico', './public/favicon.ico')
    .copy('resources/src/robots.txt', './public/robots.txt')
    .options({
        processCssUrls: false,
        postCss: [tailwindcss('./tailwind.config.js')]
    }).extract().browserSync('http://localhost:8080');
