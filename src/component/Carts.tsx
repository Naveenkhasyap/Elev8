import Image from "next/image";
import React from "react";
import { IoMdArrowDropdown } from "react-icons/io";
import { MdArrowDropUp } from "react-icons/md";

interface CartProp {
  title?: string;
  coins?: {
    price: number;
    symbol: string;
    icon: string;
    currentState: boolean;
    currentStateVal: number;
  }[];
}
function Carts({ title, coins }: CartProp) {
  return (
    <div className=" shadow-xs bg-blue-950/30 shadow-neutral-100    w-1/2 text-white rounded py-4">
      <strong className="flex justify-center text-xl">{title}</strong>
      <section>
        {coins?.map((coin) => (
          <div
            className="flex justify-between items-center py-4 px-4"
            key={coin.symbol}
          >
            <div className="flex items-center">
              <Image
                src={coin.icon}
                alt="logo"
                width={30}
                height={30}
                className="rounded-full"
              />
              <p className="px-3">{coin.symbol}</p>
            </div>
            <div className="flex items-center">
              <p className="px-3">${coin.price}</p>
              <p
                className={`${
                  !coin.currentState ? "text-tc" : "text-green-500"
                } flex items-center`}
              >
                {coin.currentState ? <MdArrowDropUp /> : <IoMdArrowDropdown />}
                {coin.currentStateVal} %
              </p>
            </div>
          </div>
        ))}
      </section>
    </div>
  );
}

export default Carts;
