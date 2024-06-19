import React from "react";

export interface data {
  name: string;
  symbol: string;
  logo: string;
  price: number;
  description: string;
  createdBy: string;
}

interface Props {
  data: data[];
}
function Carousel({ data }: Props) {
  return (
    <div className="w-3/4">
      <div className="mt-20">
        {data.map((res: data) => {
          return <div key={res.logo}>{res.symbol}</div>;
        })}
      </div>
    </div>
  );
}

export default Carousel;
