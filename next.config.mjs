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
  env: {
    BASE_URL: process.env.BASE_URL,
    NEXT_URL: process.env.NEXT_URL,
  },
};

export default nextConfig;
