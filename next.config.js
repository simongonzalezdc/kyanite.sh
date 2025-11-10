/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Webpack config for audio files
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(mp3|wav|ogg)$/,
      type: 'asset/resource'
    });
    return config;
  }
};

module.exports = nextConfig;

