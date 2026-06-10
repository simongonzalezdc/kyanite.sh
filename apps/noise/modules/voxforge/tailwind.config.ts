import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#EEF2FF',
          100: '#E0E7FF',
          500: '#3B82F6',
          600: '#2563EB',
          700: '#1D4ED8',
        },
        secondary: {
          500: '#8B5CF6',
          600: '#7C3AED',
        }
      },
      screens: {
        'xs': '320px',    // Extra small mobile
        'sm': '480px',    // Small mobile
        'md': '768px',    // Tablet
        'lg': '1024px',   // Desktop
        'xl': '1280px',   // Large desktop
        '2xl': '1536px',  // Extra large desktop
      },
      spacing: {
        'safe-t': 'var(--safe-area-inset-top)',
        'safe-b': 'var(--safe-area-inset-bottom)',
        'safe-l': 'var(--safe-area-inset-left)',
        'safe-r': 'var(--safe-area-inset-right)',
        'nav-h': 'var(--mobile-nav-height)',
      },
      minHeight: {
        'screen-safe': 'calc(100vh - var(--safe-area-inset-top) - var(--safe-area-inset-bottom))',
        'screen-nav': 'calc(100vh - var(--mobile-nav-height))',
        'screen-full-safe': 'calc(100vh - var(--safe-area-inset-top) - var(--safe-area-inset-bottom) - var(--mobile-nav-height))',
      },
      fontSize: {
        'fluid-xs': 'var(--font-size-xs)',
        'fluid-sm': 'var(--font-size-sm)',
        'fluid-base': 'var(--font-size-base)',
        'fluid-lg': 'var(--font-size-lg)',
        'fluid-xl': 'var(--font-size-xl)',
        'fluid-2xl': 'var(--font-size-2xl)',
        'fluid-3xl': 'var(--font-size-3xl)',
        'fluid-4xl': 'var(--font-size-4xl)',
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'fade-in': 'fadeIn 0.2s ease-in',
        'bounce-gentle': 'bounceGentle 0.6s ease-out',
      },
      keyframes: {
        slideUp: {
          '0%': { transform: 'translateY(100%)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideDown: {
          '0%': { transform: 'translateY(-100%)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        bounceGentle: {
          '0%, 20%, 53%, 80%, 100%': { transform: 'translate3d(0,0,0)' },
          '40%, 43%': { transform: 'translate3d(0, -8px, 0)' },
          '70%': { transform: 'translate3d(0, -4px, 0)' },
          '90%': { transform: 'translate3d(0, -2px, 0)' },
        },
      },
      touchAction: {
        'none': 'none',
        'auto': 'auto',
        'pan-x': 'pan-x',
        'pan-y': 'pan-y',
        'manipulation': 'manipulation',
      },
      willChange: {
        'auto': 'auto',
        'scroll': 'scroll-position',
        'contents': 'contents',
        'transform': 'transform',
      },
    },
  },
  plugins: [],
}
export default config

