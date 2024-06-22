"use client";
import { useParams } from "next/navigation";
import { cryptocurrencies } from "../../../src/component/cryptoTable/model";
import Header from "@/component/header/Header";

const CryptoDetail = () => {
  const params = useParams();
  const { id } = params;

  const crypto = cryptocurrencies.find((c) => c.id.toString() === id);

  if (!crypto) {
    return <div>Cryptocurrency not found</div>;
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
      </div>
    </>
  );
};

export default CryptoDetail;
