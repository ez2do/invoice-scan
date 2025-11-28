/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class",
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        display: ['Inter', 'sans-serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      colors: {
        primary: '#137fec',
        'background-light': '#f6f7f8',
        'background-dark': '#101922',
        'text-light': '#0f172a',
        'text-dark': '#f8fafc',
        'text-light-secondary': '#64748b',
        'text-dark-secondary': '#94a3b8',
        'border-light': '#e2e8f0',
        'border-dark': '#334155',
        'surface-light': '#ffffff',
        'surface-dark': '#1e293b',
        'status-green': '#10b981',
        'status-yellow': '#f59e0b',
        'status-blue': '#3b82f6'
      },
      borderRadius: {
        DEFAULT: '0.25rem',
        lg: '0.5rem',
        xl: '0.75rem',
        full: '9999px'
      },
      spacing: {
        'safe-top': 'env(safe-area-inset-top)',
        'safe-bottom': 'env(safe-area-inset-bottom)',
      }
    },
  },
  plugins: [],
}