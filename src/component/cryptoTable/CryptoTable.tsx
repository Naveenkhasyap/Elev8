"use client";
import React from "react";
import { useRouter } from "next/navigation";

const CryptoTable = ({ tokens }: any) => {
  const router = useRouter();

  const handleRowClick = (id: number) => {
    router.push(`/crypto/${id}`);
  };

  const base64ToBlob = (base64: string, type: string) => {
    try {
      const byteCharacters = atob(base64);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);
      return new Blob([byteArray], { type });
    } catch (error) {
      console.error("Failed to convert base64 to Blob", error);
      return null;
    }
  };

  return (
    <div className="overflow-x-auto rounded">
      <table className="min-w-full bg-blue-900/10 rounded">
        <thead className="bg-blue-900/20">
          <tr>
            <th className="py-2 px-3 text-left">#</th>
            <th className="py-2 px-3 text-left">Name</th>
            <th className="py-2 px-3 text-left">Price</th>
            <th className="py-2 px-3 text-left">Market Cap</th>
            <th className="py-2 px-3 text-left">Chart</th>
          </tr>
        </thead>
        <tbody>
          {tokens.map((crypto: any, index: number) => {
            let imageUrl;
            const blob = base64ToBlob(crypto.image, "image/jpeg");
            if (blob) {
              imageUrl = URL.createObjectURL(blob);
            }

            return (
              <tr
                key={crypto.description}
                className="border-t-[0.4px] border-t-blue-950 cursor-pointer"
                onClick={() => handleRowClick(crypto.id)}
              >
                <td className="py-2 px-3">{index + 1}</td>
                <td className="py-2 px-3 flex items-center">
                  <img
                    src={imageUrl}
                    alt={crypto.name}
                    className="w-6 h-6 mr-2 rounded-full"
                  />
                  {crypto.name}
                  <span className="ml-1 text-gray-200 text-sm">
                    {crypto.symbol}
                  </span>
                </td>
                <td className="py-2 px-3">{crypto.price}</td>
                <td className="py-2 px-3">{crypto.marketCap}</td>
                <td className="py-2 px-3">
                  <img
                    src={"https://www.coingecko.com/coins/325/sparkline.svg"}
                    alt={`Chart of ${crypto.name}`}
                    className="w-16 h-8"
                  />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

export default CryptoTable;
