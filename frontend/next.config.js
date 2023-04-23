/** @type {import('next').NextConfig} */
let nextConfig = {
  distDir: 'build',
  experimental: {
    appDir: true,
  },

}

const env = process.env.NODE_ENV;

if (env === 'development') {
  nextConfig = {
    distDir: 'build',
    experimental: {
      appDir: true,
    },
     rewrites: async () => [
       { 
         source: '/api/v0/:path*', 
         destination: 'http://localhost/api/v0/:path*' },
     ]
  
  }
  }

module.exports = nextConfig
