/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  darkMode: 'media',
  theme: {
    extend: {
      fontFamily: {
        sans: [
          'SF Pro Display',
          'SF Pro Text',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'Noto Sans',
          'sans-serif',
          'Apple Color Emoji',
          'Segoe UI Emoji',
        ],
      },
      keyframes: {
        blob: {
          '0%': { transform: 'translate(0px, 0px) scale(1)' },
          '33%': { transform: 'translate(20px, -30px) scale(1.1)' },
          '66%': { transform: 'translate(-25px, 10px) scale(0.95)' },
          '100%': { transform: 'translate(0px, 0px) scale(1)' },
        },
        glow: {
          '0%, 100%': { opacity: '0.0' },
          '50%': { opacity: '0.15' },
        },
      },
      animation: {
        blob: 'blob 8s ease-in-out infinite',
        glow: 'glow 3s ease-in-out infinite',
        pulse4: 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        pulse5: 'pulse 5s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        pulse6: 'pulse 6s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
    },
  },
  safelist: [
    // Dynamic gradient classes that might be purged
    'from-purple-500',
    'to-fuchsia-600',
    'from-orange-500',
    'to-amber-600',
    'from-cyan-500',
    'to-sky-600',
    'bg-gradient-to-br',
    'bg-gradient-to-r',
  ],
  plugins: [],
}

