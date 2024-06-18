import { logo } from "@/utils/Images";
import Image from "next/image";
import React from "react";

function Header() {
  return (
    <main className="flex  justify-between bg-black  w-full sticky top-0 py-2 shadow-md z-50">
      <Image src={logo} alt="logo" width={30} height={30} />

      <p>Connect Wallet</p>
    </main>
  );
}

export default Header;
