import React from "react";
import { Cryptocurrency, cryptocurrencies } from "./model";

const CryptoTable: React.FC = () => {
  return (
    <div className="overflow-x-auto rounded">
      <table className="min-w-full bg-blue-900/10 rounded">
        <thead className="bg-blue-900/20">
          <tr>
            <th className="py-2 px-4">#</th>
            <th className="py-2 px-4">Name</th>
            <th className="py-2 px-4">Price</th>
            <th className="py-2 px-4">24h %</th>
            <th className="py-2 px-4">7d %</th>
            <th className="py-2 px-4">Market Cap</th>
            <th className="py-2 px-4">Volume</th>
            <th className="py-2 px-4">Chart</th>
          </tr>
        </thead>
        <tbody>
          {cryptocurrencies.map((crypto: Cryptocurrency) => (
            <tr key={crypto.id} className="border-t-[0.4px] border-t-blue-950">
              <td className="py-2 px-4">{crypto.id}</td>
              <td className="py-2 px-4 flex items-center">
                <img
                  src={crypto.logo}
                  alt={crypto.name}
                  className="w-6 h-6 mr-2 rounded-full"
                />
                {crypto.name}{" "}
                <span className="ml-1 text-gray-200 text-sm`">
                  {crypto.symbol}
                </span>
              </td>
              <td className="py-2 px-4">{crypto.price}</td>
              <td
                className={`py-2 px-4 ${
                  crypto.change24h.startsWith("-")
                    ? "text-red-500"
                    : "text-green-500"
                }`}
              >
                {crypto.change24h}
              </td>
              <td
                className={`py-2 px-4 ${
                  crypto.change7d.startsWith("-")
                    ? "text-red-500"
                    : "text-green-500"
                }`}
              >
                {crypto.change7d}
              </td>
              <td className="py-2 px-4">{crypto.marketCap}</td>
              <td className="py-2 px-4">{crypto.volume}</td>
              <td className="py-2 px-4">
                <img
                  src={crypto.chart}
                  alt={`Chart of ${crypto.name}`}
                  className="w-16 h-8"
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CryptoTable;
