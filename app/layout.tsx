import type { Metadata } from "next";
import "./globals.css";
import Head from "next/head";
import { Providers } from "./Provider";
import { Toaster } from "react-hot-toast";

export const metadata: Metadata = {
  title: "Elev8",
  description: "A platform to launch coins on Starknet chain hassle free.",
};
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <Head>
        <link rel="icon" href="/favicon.ico" />
        <link
          href="https://fonts.googleapis.com/css2?family=Oswald:wght@200..700&display=swap"
          rel="stylesheet"
        />
      </Head>

      <body>
        <Toaster />
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
