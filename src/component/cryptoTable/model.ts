export interface Cryptocurrency {
  id: number;
  name: string;
  symbol: string;
  price: string;
  change24h: string;
  change7d: string;
  marketCap: string;
  volume: string;
  chart: string;
  logo: string;
}

export const cryptocurrencies: Cryptocurrency[] = [
  {
    id: 1,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 2,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 3,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 4,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 5,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 6,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 7,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 8,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 9,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 10,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
  {
    id: 11,
    name: "Bitcoin",
    symbol: "BTC",
    price: "$66,317.63",
    change24h: "1.5%",
    change7d: "2.8%",
    marketCap: "$1,296,701,047,451",
    volume: "$19,451,094,162",
    chart: "https://www.coingecko.com/coins/1/sparkline.svg",
    logo: "https://assets.coingecko.com/coins/images/1/standard/bitcoin.png?1696501400",
  },
];

export type Token = {
  name: string;
  ticker: string;
  description: string;
  image: string;
  userAccountAddress: string;
  status: string;
  txnHash: string;
  created_at: string;
  updated_at: string;
  change24hr: string;
  change7day: string;
  price: string;
  marketCap: string;
};

export type buy = {
  amount: number;
  address: string;
  ticker: string;
  type: string;
};
