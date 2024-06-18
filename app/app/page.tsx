import CoinCart from "@/component/CoinCart";
import Header from "../src/component/Header";

export default function Home() {
  return (
    <main className="flex flex-col py-2 px-10">
      <Header />

      <section className="flex justify-center items-center flex-col">
        <strong> Start a new coin </strong>

        <strong>Top coins</strong>

        <div className="shadow-lg">
          <CoinCart
            capital={100}
            createdBy="cdoc21"
            description="this token is crated for fun"
            image={
              "https://pump.mypinata.cloud/ipfs/QmdWPWmGjoDVpacNw51dtetcHnJitWPkdAWDq3CrAiP4nt?img-width=128&img-dpr=2&img-onerror=redirect"
            }
            price={0.5}
            title="Pepe coin"
            symbol="PEPE"
          />
        </div>
      </section>
    </main>
  );
}
