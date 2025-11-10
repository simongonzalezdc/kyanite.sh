/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Configure Turbopack for Next.js 16
  experimental: {
    turbo: {
      rules: {
        '*.mp3': ['file-loader'],
        '*.wav': ['file-loader'],
        '*.ogg': ['file-loader'],
      },
    },
  },
  // Keep webpack config for fallback
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(mp3|wav|ogg)$/,
      type: 'asset/resource'
    });
    return config;
  }
};

module.exports = nextConfig;

