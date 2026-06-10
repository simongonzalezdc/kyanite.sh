import type { NextConfig } from "next";
import withSerwistInit from "@serwist/next";

const withSerwist = withSerwistInit({
  swSrc: "src/app/sw.ts",
  swDest: "public/sw.js",
  cacheOnNavigation: true,
  reloadOnOnline: true,
  // Disable during development since Serwist doesn't support Turbopack
  disable: process.env.NODE_ENV !== "production",
});

const nextConfig: NextConfig = {
  reactStrictMode: true,
  
  // Force webpack for builds (Serwist doesn't support Turbopack yet)
  // See: https://github.com/serwist/serwist/issues/54
  turbopack: {},
  
  // Static export for self-hosting
  // The PWA can be served from any static file host
  output: "export",
  
  // Trailing slash for static hosting compatibility
  trailingSlash: true,
  
  // Disable image optimization for static export
  images: {
    unoptimized: true,
  },
};

export default withSerwist(nextConfig);
