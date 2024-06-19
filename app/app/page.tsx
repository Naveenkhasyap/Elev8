"use client";
import CoinCart from "@/component/CoinCart";
import Header from "../src/component/Header";
import Head from "next/head";
import { useEffect, useRef } from "react";

export default function Home() {
  const bgRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (bgRef.current) {
      bgRef.current.playbackRate = 0.5; // Set the playback rate to slow down the video (0.5 means half the normal speed)
    }
  }, []);
  return (
    // <main className="">
    //   <Header />

    //   <section className="flex justify-center items-center flex-col">
    //     <strong> Start a new coin </strong>

    //     <strong>Top coins</strong>

    //     <div className="shadow-lg">
    //       <CoinCart
    //         capital={100}
    //         createdBy="cdoc21"
    //         description="this token is crated for fun"
    //         image={
    //           "https://pump.mypinata.cloud/ipfs/QmdWPWmGjoDVpacNw51dtetcHnJitWPkdAWDq3CrAiP4nt?img-width=128&img-dpr=2&img-onerror=redirect"
    //         }
    //         price={0.5}
    //         title="Pepe coin"
    //         symbol="PEPE"
    //       />
    //     </div>
    //   </section>
    // </main>

    <main className="flex flex-col py-2 px-10">
      <Header />
      <section className="flex justify-center items-center flex-col">
        {" "}
        hello there
      </section>
      <video
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
        <strong className="font-oswald text-5xl">TOP COINS</strong>
      </div>
    </main>
  );
}
