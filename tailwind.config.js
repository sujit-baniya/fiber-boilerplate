module.exports = {
    purge: [
        './resources/**/*.html',
        './resources/**/*.js',
        './resources/**/*.vue',
        './resources/**/*.scss',
        './resources/**/*.css',
    ],
    theme: {
        extend: {
            colors: {
                black: '#0f1c33',
            },
            margin: {
                '96': '24rem',
                '128': '32rem',
            },
        }
    },
    variants: {
        tableLayout: ['responsive', 'hover', 'focus'],
    },
    plugins: [],
}
