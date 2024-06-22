import { logo } from "../../utils/Images";
import { DynamicConnectButton } from "@dynamic-labs/sdk-react-core";
import Image from "next/image";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import Modal from "../model/Modal";
import Success from "../model/success";

function Header() {
  const [success, setSuccess] = useState(false);
  const [coinModel, setCoinModel] = useState(false);
  const router = useRouter();
  return (
    <main className="flex  justify-between bg-black px-5  w-full sticky top-0 py-2 shadow-md z-40 border-b-[0.01px] border-tc">
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
        <button
          className="bg-white text-black rounded-full px-4 font-oswald py-2"
          onClick={() => {
            setCoinModel(true);
          }}
        >
          Launch Token
        </button>
        <DynamicConnectButton>
          <button className="bg-tc py-2 rounded-full px-4 font-oswald">
            Connect Wallet
          </button>
        </DynamicConnectButton>
        {success && <Success onClose={setSuccess} />}
        {coinModel && <Modal onClose={setCoinModel} setSuccess={setSuccess} />}
      </section>
    </main>
  );
}

export default Header;
