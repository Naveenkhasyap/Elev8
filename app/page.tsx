"use client";
import Header from "../src/component/header/Header";
import { useEffect, useRef, useState } from "react";
import Carts from "@/component/cart/Carts";
import CryptoTable from "@/component/cryptoTable/CryptoTable";
import { fetchTokens } from "@/utils/apis";

export default function Home() {
  const bgRef = useRef<HTMLVideoElement>(null);
  const [tokens, setTokens] = useState([]);
  const [largestGainer, setLargestGainer] = useState([]);
  const [treadingCoins, setTreadingCoins] = useState([]);

  const fetchData = async () => {
    try {
      const response = await fetchTokens();
      if (response && response.data !== undefined) {
        setTokens(response.data.data);
        setLargestGainer(response.data.data.slice(0, 3));
        setTreadingCoins(response.data.data.slice(3, 6));
      }
    } catch (error) {
      console.error("Error fetching tokens:", error);
    }
  };

  useEffect(() => {
    fetchData();
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
          <CryptoTable tokens={tokens} />
        </section>
      </div>
    </main>
  );
}
