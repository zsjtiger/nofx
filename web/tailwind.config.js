/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Color values live in src/app/index.css (:root). Channel variables allow
        // opacity modifiers such as bg-crux-gold/20.
        'crux-gold': {
          DEFAULT: 'rgb(var(--crux-accent-rgb) / <alpha-value>)',
          dim: 'rgba(224, 72, 59, 0.10)',
          glow: 'rgba(224, 72, 59, 0.30)',
          highlight: 'var(--crux-accent-dark)',
        },
        'crux-bg': {
          DEFAULT: 'var(--crux-bg)',
          deeper: 'var(--crux-bg-deeper)',
          lighter: 'var(--panel-bg)',
        },
        'crux-accent': {
          DEFAULT: 'rgb(var(--crux-accent-rgb) / <alpha-value>)',
          text: 'var(--crux-accent-text)',
        },
        'crux-warning': {
          DEFAULT: 'rgb(var(--crux-warning-rgb) / <alpha-value>)',
          text: 'var(--crux-warning-text)',
        },
        'crux-text': {
          DEFAULT: 'var(--text-primary)',
          main: 'var(--text-primary)',
          muted: 'var(--text-tertiary)',
        },
        'crux-success': {
          DEFAULT: 'rgb(var(--crux-green-rgb) / <alpha-value>)',
          text: 'var(--crux-green-text)',
        },
        'crux-danger': {
          DEFAULT: 'rgb(var(--crux-red-rgb) / <alpha-value>)',
          text: 'var(--crux-red-text)',
        },
      },
      fontFamily: {
        sans: ['IBM Plex Mono', 'ui-monospace', 'Menlo', 'monospace'],
        mono: ['IBM Plex Mono', 'Menlo', 'Monaco', 'Courier New', 'monospace'],
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(circle at center, var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
        'scanlines': "url(\"data:image/svg+xml,%3Csvg width='4' height='4' viewBox='0 0 4 4' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M0 0H4V2H0V0Z' fill='rgba(0,0,0,0.4)'/%3E%3C/svg%3E\")",
        'grid-pattern': "linear-gradient(to right, #1f2937 1px, transparent 1px), linear-gradient(to bottom, #1f2937 1px, transparent 1px)",
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'scan': 'scan 8s linear infinite',
        'scan-fast': 'scan 2s linear infinite',
        'float': 'float 6s ease-in-out infinite',
        'glitch': 'glitch 0.3s cubic-bezier(.25, .46, .45, .94) both infinite',
        'shimmer': 'shimmer 2s linear infinite',
      },
      keyframes: {
        scan: {
          '0%': { backgroundPosition: '0 0' },
          '100%': { backgroundPosition: '0 100%' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        glitch: {
          '0%': { transform: 'translate(0)' },
          '20%': { transform: 'translate(-2px, 2px)' },
          '40%': { transform: 'translate(-2px, -2px)' },
          '60%': { transform: 'translate(2px, 2px)' },
          '80%': { transform: 'translate(2px, -2px)' },
          '100%': { transform: 'translate(0)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
      },
      boxShadow: {
        'neon': '0 0 5px theme("colors.crux-gold.DEFAULT"), 0 0 20px theme("colors.crux-gold.dim")',
        'neon-blue': '0 0 5px theme("colors.crux-accent"), 0 0 20px rgba(0, 240, 255, 0.2)',
      },
    },
  },
  plugins: [],
}
