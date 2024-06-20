import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Head from "next/head";
import { DynamicContextProvider } from "@dynamic-labs/sdk-react-core";
import { StarknetWalletConnectors } from "@dynamic-labs/starknet";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Elev8",
  description: "A platform to launch coins on Starknet chain hassle free.",
};
const environmentId = process.env.DYNAMIC_API_KEY;
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
      <DynamicContextProvider
        settings={{
          environmentId: environmentId || "",
          walletConnectors: [StarknetWalletConnectors],
        }}
      >
        <body className={inter.className}>{children}</body>
      </DynamicContextProvider>
    </html>
  );
}
