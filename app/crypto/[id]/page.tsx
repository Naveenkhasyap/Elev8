"use client";
import React, { useEffect, useState } from "react";
import { Chart, registerables } from "chart.js";
import { Line } from "react-chartjs-2";
import { useParams } from "next/navigation";
import Header from "@/component/header/Header";
import ChartErrorBoundary from "@/component/ErrorBoundary";
import Button from "@/component/Button";
import { Token } from "@/component/cryptoTable/model";
import toast from "react-hot-toast";
import BuySell from "@/component/buySell/BuySell";
import { buyToken, fetchTokens, sellToken } from "@/utils/apis";
import { useDynamicContext } from "@dynamic-labs/sdk-react-core";

Chart.register(...registerables);

const CryptoDetail: React.FC = () => {
  const params = useParams<{ id: string }>();
  const [isBuyLoading, setIsBuyLoading] = useState(false);
  const [buySell, setBuySell] = useState(false);
  const [text, setText] = useState("");
  const [style, setStyle] = useState("");
  const { primaryWallet } = useDynamicContext();
  const [cryptoCoin, setCryptoCoin] = useState([]);
  const [currentToken, setCurrentToken] = useState<Token | null>(null);

  const fetchData = async () => {
    try {
      const response = await fetchTokens();
      if (response && response?.data !== undefined) {
        setCryptoCoin(response?.data?.data);
      } else {
        return;
      }
    } catch (error) {
      toast.error("something went wrong");
      console.error("Error fetching tokens:", error);
    }
  };

  if (!params || !params.id) {
    return <div>No ID available</div>;
  }

  const { id } = params;
  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    const crypto = cryptoCoin.find((c: Token) => c?.name === id);
    if (crypto) {
      setCurrentToken(crypto);
    }
  }, [cryptoCoin, id]);

  const options = {
    scales: {
      y: {
        beginAtZero: false,
      },
    },
  };

  if (!currentToken) {
    return (
      <>
        <Header />
        <div className="flex justify-center items-center h-[90vh]">
          <p>No Crypto Coin found</p>
        </div>
      </>
    );
  }

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

  const chartKey = `${currentToken?.name}-chart`;
  let imageUrl = `data:image/png;base64,${currentToken?.image}`;

  const BuyHandler = async (type: string, amount: number) => {
    try {
    } catch (err) {
      console.log(err);
    }
    const address = primaryWallet?.address || "";
    if (amount === 0) {
      return toast.error(`Amount must be greater than zero`);
    }

    const ticker = currentToken.ticker;

    if (!ticker) {
      toast.error("Something went wrong Please try again");
    }
    if (!address) {
      toast.error("Please Connect your Wallet");
    }
    const lowerCaseTypeValue = type.toLowerCase();
    let response;

    if (lowerCaseTypeValue === "buy") {
      response = await buyToken({
        amount,
        address,
        ticker,
        type: lowerCaseTypeValue,
      });
    } else if (lowerCaseTypeValue === "sell") {
      response = await sellToken({
        amount,
        address,
        ticker,
        type: lowerCaseTypeValue,
      });
    }
    if (response?.data.success) {
      setBuySell(false);
      toast.success("Transaction Successful");
    } else {
      setBuySell(false);
      toast.error("Something went wrong Please try again");
    }
  };

  return (
    <main>
      <Header />
      <div className="flex justify-between items-center gap-4 py-10 px-10">
        {buySell && (
          <BuySell
            text={text}
            clickHandler={BuyHandler}
            onClose={setBuySell}
            style={style}
          />
        )}
        <ChartErrorBoundary>
          <div className="flex-1 h-full">
            <Line key={chartKey} data={data} options={options} />
          </div>
        </ChartErrorBoundary>
        <section className="w-1/3  flex-col">
          <div className="flex items-center gap-2">
            <img
              src={imageUrl}
              width={40}
              height={40}
              alt="coin_logo"
              className="rounded-full"
            />
            <p className=" text-3xl">
              <strong>{currentToken?.name}</strong>
              <span className=""> Price #1</span>
            </p>
          </div>
          <p className="text-xl font-oswald text-green-500 mt-2 ml-2 ">{`$ ${currentToken?.price}`}</p>
          <div className="flex gap-4 mt-4">
            <Button
              text="Buy"
              styles={"w-10 bg-green-600 hover:bg-green-700"}
              isLoading={isBuyLoading}
              onClick={() => {
                setStyle("");
                setText("Buy");
                setBuySell(true);
              }}
            />
            <Button
              text="Sell"
              styles={"w-10 bg-tc hover:bg-red-700"}
              isLoading={isBuyLoading}
              onClick={() => {
                setStyle("bg-red-700 hover:bg-red-800");
                setText("Sell");
                setBuySell(true);
              }}
            />
          </div>
        </section>
      </div>
      <section className="px-5 flex flex-col gap-6">
        <p className="text-2xl font-oswald px-4">
          {" "}
          {`${currentToken.name} Statistics`}
        </p>
        <div className="bg-blue-800/30 px-5 py-5 rounded-md flex flex-col gap-6">
          <div className="flex items-center justify-between gap-2 ">
            <p className="text-lg font-oswald">Price</p>
            <p className="text-lg font-oswald">${currentToken?.price}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Market Cap</p>
            <p className="text-lg font-oswald">${currentToken?.marketCap}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">Volume</p>
            <p className="text-lg font-oswald">${currentToken?.marketCap}</p>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-lg font-oswald">24 Hour Tranding Vol</p>
            <p className="text-lg font-oswald">${currentToken.change24hr}</p>
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
    </main>
  );
};

export default CryptoDetail;
