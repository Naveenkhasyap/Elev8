"use client";
import { useParams } from "next/navigation";
import Header from "@/component/header/Header";
import TradingViewWidget from "@/component/cryptoDetail/TradingViewWidget";
import { useSelector } from "react-redux";

const CryptoDetail = () => {
  const params = useParams<{ id: string }>();
  const { cryptoCoins } = useSelector((state: any) => state.MainSlice);
  if (!params || !params.id) {
    return <div>No ID available</div>;
  }

  const { id } = params;

  const crypto = cryptoCoins.find((c: any) => c.id.toString() === id);

  if (!crypto) {
    return (
      <div className="flex justify-center items-center h-[90vh]">
        <img src="https://media.tenor.com/hB9OTbewrikAAAAi/work-work-in-progress.gif" />
      </div>
    );
  }

  return (
    <>
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold">{crypto.name}</h1>
        <p>Price: {crypto.price}</p>
        <p>24h Change: {crypto.change24h}</p>
        <p>7d Change: {crypto.change7d}</p>
        <p>Market Cap: {crypto.marketCap}</p>
        <p>Volume: {crypto.volume}</p>
        <img src={`/charts/${crypto.chart}`} alt={`Chart of ${crypto.name}`} />
        <TradingViewWidget />
      </div>
    </>
  );
};

export default CryptoDetail;
