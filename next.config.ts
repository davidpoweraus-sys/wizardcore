import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  reactCompiler: true,
  output: 'standalone',
  transpilePackages: ['pixelarticons'],
};

export default nextConfig;
