"use client";
import CoinCart from "@/component/CoinCart";
import Header from "../src/component/Header";
import Head from "next/head";
import { useEffect, useRef } from "react";
import Carousel from "@/component/Carousel";

export default function Home() {
  const bgRef = useRef<HTMLVideoElement>(null);

  const dummyCoins = [
    {
      name: "Hulk",
      symbol: "HULK",
      logo: "https://w0.peakpx.com/wallpaper/495/822/HD-wallpaper-hulk-marvel.jpg",
      price: 0.5,

      description: "This token represents the Marvel Hulk ",
      createdBy: "cdoc21",
    },
    {
      name: "RedHulk",
      symbol: "RULK",
      logo: "https://images2.alphacoders.com/134/thumb-1920-1348343.jpeg",
      price: 0.5,

      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {

      name: "LOKI",
      symbol: "LOKI",
      logo: "https://c4.wallpaperflare.com/wallpaper/84/141/133/marvel-cinematic-universe-loki-tom-hiddleston-thor-ragnarok-wallpaper-preview.jpg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {
      name: "Captain America",
      symbol: "CAPTAIN",
      logo: "https://c4.wallpaperflare.com/wallpaper/713/13/242/movies-marvel-cinematic-universe-marvel-comics-avengers-endgame-captain-america-hd-wallpaper-preview.jpg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {
      name: "Hulk",
      symbol: "HULK",
      logo: "https://w0.peakpx.com/wallpaper/495/822/HD-wallpaper-hulk-marvel.jpg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {
      name: "RedHulk",
      symbol: "RULK",
      logo: "https://images2.alphacoders.com/134/thumb-1920-1348343.jpeg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {
      name: "LOKI",
      symbol: "LOKI",
      logo: "https://c4.wallpaperflare.com/wallpaper/84/141/133/marvel-cinematic-universe-loki-tom-hiddleston-thor-ragnarok-wallpaper-preview.jpg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
    {
      name: "Captain America",
      symbol: "CAPTAIN",
      logo: "https://c4.wallpaperflare.com/wallpaper/713/13/242/movies-marvel-cinematic-universe-marvel-comics-avengers-endgame-captain-america-hd-wallpaper-preview.jpg",
      price: 0.5,
      description: "This token represents the Marvel Hulk",
      createdBy: "cdoc21",
    },
  ];

  useEffect(() => {
    if (bgRef.current) {
      bgRef.current.playbackRate = 0.5; // Set the playback rate to slow down the video (0.5 means half the normal speed)
    }
  }, []);
  return (

    // <>
    //     // <main className="">
    // //   <Header />

    // //   <section className="flex justify-center items-center flex-col">
    // //     <strong> Start a new coin </strong>

    // //     <strong>Top coins</strong>

    // //     <div className="shadow-lg">
    // //       <CoinCart
    // //         capital={100}
    // //         createdBy="cdoc21"
    // //         description="this token is crated for fun"
    // //         image={
    // //           "https://pump.mypinata.cloud/ipfs/QmdWPWmGjoDVpacNw51dtetcHnJitWPkdAWDq3CrAiP4nt?img-width=128&img-dpr=2&img-onerror=redirect"
    // //         }
    // //         price={0.5}
    // //         title="Pepe coin"
    // //         symbol="PEPE"
    // //       />
    // //     </div>
    // //   </section>
    // // </main>
    // </>

    <main className="flex flex-col ">
      <Header />
      <section className="flex flex-col ">
        <strong className="font-oswald text-4xl text-center m-16">
          TOP COINS
        </strong>
        <Carousel data={dummyCoins} />
        {/* <section className="flex justify-center items-center flex-col">
        {" "}
        hello there
      </section> */}

        {/* <video
        ref={bgRef}
        autoPlay
        loop
        muted
        className="absolute top-0 left-0 w-full h-full object-cover z-0"
      >
        <source
          src="https://ik.imagekit.io/zt5dmrbrg/Elev8/Elev8.mp4"
          type="video/mp4"
        />
      </video>
      <div className="relative z-10 text-white text-center p-4">
      </div> */}
      </section>
    </main>
  );
}
