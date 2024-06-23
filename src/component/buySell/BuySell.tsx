"use client";
import React, { useState } from "react";
import { RxCross2 } from "react-icons/rx";

type BuySellProps = {
  onClose: (value: boolean) => void;
  setSuccess?: (value: boolean) => void;
  text: string;
  style: string;
  clickHandler: (type: string, amount: number) => void;
};

const BuySell = ({
  onClose,
  setSuccess,
  text,
  clickHandler,
  style,
}: BuySellProps) => {
  const [isClosing, setIsClosing] = useState(false);
  const [amount, setAmount] = useState<number>();

  const inputStyle = `w-full px-4 py-2 mb-4 text-white bg-blue-800/30 border border-blue-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500`;

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      onClose(false);
    }, 500);
  };

  return (
    <div
      className={`fixed inset-0 bg-black bg-opacity-75 flex justify-center items-center z-30  ${
        isClosing
          ? "motion-safe:animate-fade-out"
          : "motion-safe:animate-fade-in"
      }`}
      onClick={handleClose}
    >
      <div
        className={`bg-blue-900/50 relative px-2 py-2 rounded mt-4 shadow-lg max-w-[50rem] w-[30rem] mx-auto z-40 flex flex-col ${
          isClosing
            ? "motion-safe:animate-fade-out"
            : "motion-safe:animate-fade-in"
        }`}
        onClick={(e) => e.stopPropagation()}
      >
        <section className="flex justify-between items-center">
          <p />
          <p className="text-center self-center font-oswald">
            Create a new Coin
          </p>
          <RxCross2
            className=" w-10  h-7 text-tc font-bold cursor-pointer z-[1000]"
            onClick={() => handleClose()}
          />
        </section>

        <main className="py-4 px-3">
          <input
            type="number"
            placeholder="Amount"
            className={inputStyle}
            value={amount}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              setAmount(parseInt(e.target.value));
            }}
          />

          <button
            className={`${style} w-full px-4 py-2 mb-4 text-white bg-green-600 rounded hover:bg-green-700 focus:outline-none focus:ring-2`}
            onClick={() => {
              clickHandler(text, amount || 0);
            }}
          >
            {text ? text : "Submit"}
          </button>
        </main>
        <div className="flex justify-end mr-3">
          <button
            onClick={handleClose}
            className="px-4 py-2 bg-red-500  ml-5 text-white rounded"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default BuySell;
