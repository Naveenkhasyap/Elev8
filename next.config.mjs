/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: [
      "pump.mypinata.cloud",
      "w0.peakpx.com",
      "images2.alphacoders.com",
      "c4.wallpaperflare.com",
      "assets.coingecko.com",
      "ik.imagekit.io",
    ],
  },
  async headers() {
    return [
      {
        // matching all API routes
        source: "/api/:path*",
        headers: [
          { key: "Access-Control-Allow-Credentials", value: "true" },
          { key: "Access-Control-Allow-Origin", value: "*" },
          { key: "Access-Control-Allow-Methods", value: "GET,OPTIONS,PATCH,DELETE,POST,PUT" },
          { key: "Access-Control-Allow-Headers", value: "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version" },
        ]
      }
    ]
  },
  env: {
    BASE_URL: process.env.BASE_URL,
    NEXT_URL: process.env.NEXT_URL,
  },
};

export default nextConfig;
