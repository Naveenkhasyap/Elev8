"use client";
import Header from "../src/component/header/Header";
import { useEffect, useRef } from "react";
import Carts from "@/component/cart/Carts";
import CryptoTable from "@/component/cryptoTable/CryptoTable";
import { useDispatch, useSelector } from "react-redux";

export default function Home() {
  const bgRef = useRef<HTMLVideoElement>(null);
  const dispatch = useDispatch();
  const { modelOpen } = useSelector((state: any) => state.MainSlice);

  const largestGainer = [
    {
      price: 1.05,
      symbol: "Toshi",
      icon: "https://assets.coingecko.com/coins/images/31126/standard/2023-08-11_13.21.24.png?1697025145",
      currentState: true,
      currentStateVal: 1.05,
    },
    {
      price: 1.05,
      symbol: "ZKsync",
      icon: "https://assets.coingecko.com/coins/images/38043/standard/ZKTokenBlack.png?1718614502",
      currentState: true,
      currentStateVal: 1.05,
    },
    {
      price: 1.05,
      symbol: "LinqAI",
      icon: "https://assets.coingecko.com/coins/images/35645/standard/cmc-cg-linq-logo_%282%29.png?1709368877",
      currentState: false,
      currentStateVal: 1.05,
    },
  ];
  const treadingCoins = [
    {
      price: 1.05,
      symbol: "Toshi",
      icon: "https://assets.coingecko.com/coins/images/31126/standard/2023-08-11_13.21.24.png?1697025145",
      currentState: true,
      currentStateVal: 1.05,
    },
    {
      price: 1.05,
      symbol: "ZKsync",
      icon: "https://assets.coingecko.com/coins/images/38043/standard/ZKTokenBlack.png?1718614502",
      currentState: false,
      currentStateVal: 1.05,
    },
    {
      price: 1.05,
      symbol: "LinqAI",
      icon: "https://assets.coingecko.com/coins/images/35645/standard/cmc-cg-linq-logo_%282%29.png?1709368877",
      currentState: true,
      currentStateVal: 1.05,
    },
  ];

  useEffect(() => {
    if (bgRef.current) {
      bgRef.current.playbackRate = 0.5;
    }
  }, []);
  return (
    <main>
      <Header />


      <div className="px-10">
        <section className="flex flex-col py-15">
          <div className="flex  justify-between w-full gap-5 mt-24">
            <Carts title={"🔥 Trending"} coins={treadingCoins} />
            <Carts title={"🚀 Largest Gainers"} coins={largestGainer} />
          </div>
        </section>

        <section className="mt-10">
          <CryptoTable />
        </section>
      </div>
    </main>
  );
}
