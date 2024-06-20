import Slider from "react-slick";
import Image from "next/image"; // Assuming you're using Next.js Image component
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";

interface Data {
  name: string;
  logo: string;
  symbol: string;
  createdBy: string;
  price: number;
  description: string;
}

interface Props {
  data: Data[];
}

const Carousel = ({ data }: Props) => {
  const settings = {
    dots: true,
    infinite: true,
    speed: 500,
    slidesToShow: 4,
    slidesToScroll: 1,
    autoplay: true,
    autoplaySpeed: 3000,
  };

  return (
    <div className="w-[100%] px-20">
      <Slider {...settings}>
        {data.map((res: Data) => {
          const descArray = res.description.split(" ");
          const first10 = descArray.slice(0, 10);
          const desc = first10.join(" ");
          console.log(desc);
          return (
            <div
              key={res.logo}
              className="bg-white/95 w-1/3 h-[370px] m-2 rounded-xl text-black"
            >
              <Image
                src={res.logo}
                width={300}
                alt={"coin"}
                height={300}
                className=" rounded-t-md"
              />
              <div className="flex flex-col">
                <p className="text-center font-oswald font-semibold text-lg">
                  {res.symbol}
                </p>
                <section className="px-3">
                  <div className="flex justify-between">
                    <p>Owner </p>
                    <p>-</p>
                    <p>{res.createdBy}</p>
                  </div>

                  <div className="flex justify-between">
                    <p>Price </p>
                    <p>-</p>
                    <p>{res.price}</p>
                  </div>
                </section>
                <p className="px-3 pt-4">{desc} .....</p>
                <div className="flex justify-end px-4 items-end">
                  <button className="bg-tc text-white text-lg font-oswald px-4 py-1 rounded-lg">
                    Trade
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </Slider>
    </div>
  );
};

export default Carousel;
