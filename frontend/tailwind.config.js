/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        bg: '#f4f2ec',
        panel: '#fffdf7',
        ink: {
          DEFAULT: '#1f2933',
          soft: '#52606d',
        },
        line: '#e3dccf',
        brand: {
          DEFAULT: '#0f766e',
          light: '#14b8a6',
        },
        ok: {
          DEFAULT: '#166534',
          bg: '#dcfce7',
        },
        bad: {
          DEFAULT: '#991b1b',
          bg: '#fee2e2',
        },
      },
      fontFamily: {
        heading: ['"Space Grotesk"', 'sans-serif'],
        body: ['"IBM Plex Sans"', 'sans-serif'],
      },
      borderRadius: {
        panel: '14px',
        card: '12px',
        btn: '10px',
      },
      boxShadow: {
        panel: '0 10px 30px rgba(37, 46, 53, 0.08)',
      },
    },
  },
  plugins: [],
}
