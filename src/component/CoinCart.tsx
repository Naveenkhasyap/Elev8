import Image from "next/image";
import React from "react";

type CoinCart = {
  title: string;
  description: string;
  image: string;
  price: number;
  capital: number;
  createdBy: string;
  symbol: string;
};
function CoinCart({
  title,
  symbol,
  description,
  image,
  price,
  capital,
  createdBy,
}: CoinCart) {
  ``;
  return (
    <div className="flex shadow gap-3">
      <Image
        src={image}
        alt={title}
        width={100}
        height={100}
        className="rounded"
      />

      <section className="flex flex-col ">
        <div className="flex gap-3">
          <p>Creted by</p>
          <p>{createdBy}</p>
        </div>

        <div className="flex gap-3">
          <p>price</p>
          <p>{price}</p>
        </div>
        <section className="flex gap-3">
          <p>Market capital</p>
          <p>{capital}</p>
        </section>
      </section>
    </div>
  );
}

export default CoinCart;
