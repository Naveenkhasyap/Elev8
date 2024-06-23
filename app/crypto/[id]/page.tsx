"use client";
import React from "react";
import { Chart, registerables } from "chart.js";
import { Line } from "react-chartjs-2";
import { useSelector } from "react-redux";
import { useParams } from "next/navigation";
import Header from "@/component/header/Header";
import ChartErrorBoundary from "@/component/ErrorBoundary";
import Button from "@/component/Button";

Chart.register(...registerables);
interface CryptoCoin {
  name: string;
  price: number;
  change24h: number;
  change7d: number;
  marketCap: number;
  volume: number;
  chart: string;
  priceRange: string;
  fullyDilutedValuation: number;
  volume24h: number;
  circulatingSupply: number;
  totalSupply: number;
  maxSupply: number;
  symbol: string;
}

const CryptoDetail: React.FC = () => {
  const { cryptoCoins } = useSelector((state: any) => state.MainSlice);
  const params = useParams<{ id: string }>();

  if (!params || !params.id) {
    return <div>No ID available</div>;
  }

  const { id } = params;
  const crypto = cryptoCoins.find((c: CryptoCoin) => c?.name === id);

  if (!crypto) {
    return (
      <div className="flex justify-center items-center h-[90vh]">
        <img src="https://media.tenor.com/hB9OTbewrikAAAAi/work-work-in-progress.gif" />
      </div>
    );
  }

  const options = {
    scales: {
      y: {
        beginAtZero: false,
      },
    },
  };

  const data = {
    labels: ["Red", "Blue", "Yellow", "Green", "Purple", "Orange"],
    datasets: [
      {
        label: "Price",
        data: [100, 150, 120, 180, 200, 190, 210, 220],
        borderColor: "rgba(75, 192, 192, 1)",
        tension: 0.1,
      },
    ],
  };

  const chartKey = `${crypto?.name}-chart`;
  let imageUrl = `data:image/png;base64,${crypto?.image}`;

  return (
    <>
      <Header />
      <main className="flex justify-between items-center gap-8 py-10 px-10">
        <section className="w-1/2  flex-col">
          <div className="flex items-center gap-2">
            <img
              src={imageUrl}
              width={20}
              height={20}
              alt="coin_logo"
              className="rounded-full"
            />
            <p className=" text-lg ">
              <strong>{crypto?.name}</strong>
              <span className=""> Price #1</span>
            </p>
          </div>
          <p className="text-xl font-oswald text-green-500">{`$ ${crypto?.price}`}</p>
          <div className="flex gap-4">
            <Button text="Buy" styles={"w-10"} />
            <Button text="Sell" styles={"w-10 bg-tc"} />
          </div>
        </section>
        <ChartErrorBoundary>
          <div className="w-1/2 h-full">
            <Line key={chartKey} data={data} options={options} />
          </div>
        </ChartErrorBoundary>
      </main>
      <section className="px-5 flex flex-col gap-6">
        <p className="text-2xl font-oswald px-4">
          {" "}
          {`${crypto.name} Statistics`}
        </p>
        <div className="bg-blue-800/30 px-5 py-5 rounded-md flex flex-col gap-6">
          <div className="flex items-center justify-between gap-2 ">
            <p className="text-lg font-oswald">Price</p>
            <p className="text-lg font-oswald">${crypto?.price}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Market Cap</p>
            <p className="text-lg font-oswald">${crypto?.marketCap}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Volume</p>
            <p className="text-lg font-oswald">${crypto?.marketCap}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">24 Hour Tranding Vol</p>
            <p className="text-lg font-oswald">${crypto.change24hr}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Circulation Supply</p>
            <p className="text-lg font-oswald">21,000,000</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Max Supply</p>
            <p className="text-lg font-oswald">21,000,000</p>
          </div>
        </div>
      </section>
    </>
  );
};

export default CryptoDetail;
