import { logo } from "@/utils/Images";
import Image from "next/image";
import React from "react";

function Header() {
  return (
    <main className="flex  justify-between bg-black px-5  w-full sticky top-0 py-2 shadow-md z-50 border-b-[0.01px] border-tc">
      <div>
        <section>
          <Image src={logo} alt="logo" width={30} height={30} />
        </section>
      </div>

      <button className="bg-tc py-2 rounded-full px-4 font-oswald">
        Connect Wallet
      </button>
    </main>
  );
}

export default Header;
