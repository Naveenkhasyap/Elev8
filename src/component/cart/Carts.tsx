import React from "react";
import { IoMdArrowDropdown } from "react-icons/io";
import { MdArrowDropUp } from "react-icons/md";
import { Token } from "../cryptoTable/model";

interface CartProp {
  title?: string;
  coins?: Token[];
}

function Carts({ title, coins }: CartProp) {
  return (
    <div className=" shadow-xs bg-blue-950/30 shadow-neutral-100    w-1/2 text-white rounded py-4">
      <strong className="flex justify-center text-xl">{title}</strong>
      <section>
        {coins?.map((coin) => {
          let imageUrl = `data:image/png;base64,${coin.image}`;

          return (
            <div
              className="flex justify-between items-center py-4 px-4"
              key={coin.name}
            >
              <div className="flex items-center">
                <img
                  src={imageUrl}
                  alt="logo"
                  width={30}
                  height={30}
                  className="rounded-full"
                />
                <p className="px-3">{coin.name}</p>
              </div>
              <div className="flex items-center">
                <p className="px-3">${coin.price}</p>
                <p
                  className={`${
                    !coin.change24hr ? "text-tc" : "text-green-500"
                  } flex items-center`}
                >
                  {coin.change24hr ? <MdArrowDropUp /> : <IoMdArrowDropdown />}
                  {coin.change24hr} %
                </p>
              </div>
            </div>
          );
        })}
      </section>
    </div>
  );
}

export default Carts;
