/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Analytics configuration
  env: {
    NEXT_PUBLIC_APP_VERSION: process.env.npm_package_version || '1.0.0',
    NEXT_PUBLIC_ANALYTICS_ENDPOINT: process.env.ANALYTICS_ENDPOINT,
    NEXT_PUBLIC_ANALYTICS_ENABLED: process.env.ANALYTICS_ENABLED !== 'false',
    NEXT_PUBLIC_ANALYTICS_DEBUG: process.env.NODE_ENV === 'development'
  },
  // Webpack config for audio files and Pixi.js
  webpack: (config, { isServer }) => {
    // Audio files
    config.module.rules.push({
      test: /\.(mp3|wav|ogg)$/,
      type: 'asset/resource'
    });

    // Fix for Pixi.js in Next.js - only on client side
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        path: false,
        crypto: false,
      };
    }

    return config;
  }
};

module.exports = nextConfig;

