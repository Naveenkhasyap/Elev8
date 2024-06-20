import { logo } from "../utils/Images";
import { DynamicWidget } from "@dynamic-labs/sdk-react-core";
import Image from "next/image";
import { useRouter } from "next/navigation";
import React from "react";

function Header() {
  const router = useRouter();
  return (
    <main className="flex  justify-between bg-black px-5  w-full sticky top-0 py-2 shadow-md z-50 border-b-[0.01px] border-tc">
      <div>
        <section className="flex gap-5">
          <Image
            src={logo}
            alt="logo"
            width={30}
            height={30}
            className="cursor-pointer"
            onClick={() => router.push("/")}
          />
        </section>
      </div>
      <section className="flex gap-5">
        <button className="bg-white text-black rounded-full px-4 font-oswald py-2">
          Launch Token
        </button>
        <button className="bg-tc py-2 rounded-full px-4 font-oswald">
          Connect Wallet
        </button>
        <DynamicWidget />
      </section>
    </main>
  );
}

export default Header;
