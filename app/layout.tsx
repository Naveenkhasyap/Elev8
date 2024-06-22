import type { Metadata } from "next";
import "./globals.css";
import Head from "next/head";
import { Providers } from "./Provider";

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
      {/* <DynamicContextProvider
        settings={{
          environmentId: "2762a57b-faa4-41ce-9f16-abff9300e2c9",
          recommendedWallets: [
            { walletKey: "phantomevm", label: "Popular" },
            { walletKey: "okxwallet" },
          ],
        }}
      >
        <body className={inter.className}>{children}</body>
      </DynamicContextProvider>  */}

      <body>
        {/* <Provider store={store}> */}
        <Providers>{children}</Providers>
        {/* </Provider> */}
      </body>
    </html>
  );
}
