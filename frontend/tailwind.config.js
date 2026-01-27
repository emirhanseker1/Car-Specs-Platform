/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: {
                    DEFAULT: '#ff5a1f', // Distinct orange from the design
                    hover: '#e04e1b',
                },
                background: '#f8f8f8', // Very light gray background
                surface: '#ffffff',
                text: {
                    main: '#1a1a1a',
                    muted: '#858585'
                },
                border: '#e5e5e5'
            },
            fontFamily: {
                sans: ['Inter', 'sans-serif'],
            },
            boxShadow: {
                'card': '0 2px 8px rgba(0, 0, 0, 0.04)',
            }
        },
    },
    plugins: [],
}
