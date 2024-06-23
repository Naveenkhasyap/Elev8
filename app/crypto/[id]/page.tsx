"use client";
import React from "react";
import { Chart, registerables } from "chart.js";
import { Line } from "react-chartjs-2";
import { useSelector } from "react-redux";
import { useParams } from "next/navigation";
import Header from "@/component/header/Header";
import ChartErrorBoundary from "@/component/ErrorBoundary";
// import ChartErrorBoundary from "@/component/ErrorBoundary";

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

  // if (!crypto) {
  //   return (
  //     <div className="flex justify-center items-center h-[90vh]">
  //       <img src="https://media.tenor.com/hB9OTbewrikAAAAi/work-work-in-progress.gif" />
  //     </div>
  //   );
  // }

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

  return (
    <>
      <Header />
      <main className="flex justify-between items-center py-10 px-10">
        <section className="w-1/2 text-blue-600">hi there</section>
        <ChartErrorBoundary>
          <div className="w-1/2">
            <Line key={chartKey} data={data} options={options} />
          </div>
        </ChartErrorBoundary>
      </main>
    </>
  );
};

export default CryptoDetail;
