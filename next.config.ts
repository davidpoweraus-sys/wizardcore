import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  reactCompiler: true,
  output: 'standalone',
  transpilePackages: ['pixelarticons'],
  
  async headers() {
    return [
      {
        // Apply these headers to all routes
        source: '/:path*',
        headers: [
          {
            key: 'Access-Control-Allow-Credentials',
            value: 'true',
          },
          {
            key: 'Access-Control-Allow-Origin',
            value: process.env.NODE_ENV === 'production'
              ? 'https://app.offensivewizard.com'
              : 'http://localhost:3000',
          },
          {
            key: 'Access-Control-Allow-Methods',
            value: 'GET,DELETE,PATCH,POST,PUT,OPTIONS',
          },
          {
            key: 'Access-Control-Allow-Headers',
            value: 'X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, Authorization, apikey, x-client-info, x-supabase-api-version, Next-Router-Prefetch, Next-Router-State-Tree, Next-Url, RSC',
          },
          {
            key: 'Access-Control-Expose-Headers',
            value: 'X-NextJS-Data, X-RSC, X-Action-Redirect, X-Middleware-Version',
          },
        ],
      },
    ]
  },
};

export default nextConfig;
