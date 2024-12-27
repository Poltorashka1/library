/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./templates/**/*.{html,js}", // Шаблоны HTML
        "./static/**/*.{html,js}",   // Другие статические файлы
        "./**/*.go",                 // Go-файлы
    ],
    theme: {
    },
    plugins: [],
}

