import { logo } from "../../utils/Images";
import {
  DynamicConnectButton,
  useDynamicContext,
} from "@dynamic-labs/sdk-react-core";
import Image from "next/image";
import { useRouter } from "next/navigation";
import React from "react";

function Header() {
  const router = useRouter();
  const { primaryWallet } = useDynamicContext();
  const { handleLogOut } = useDynamicContext();

  return (
    <main className="flex  justify-between bg-black px-5  w-full sticky top-0 py-2 shadow-md z-40 border-b-[0.01px] border-tc">
      <div>
        <section className="flex gap-5">
          <Image
            src={logo}
            alt="logo"
            width={200}
            height={200}
            className="cursor-pointer"
            onClick={() => router.push("/")}
          />
        </section>
      </div>
      <section className="flex gap-5">
        <button
          className="bg-white text-black rounded-full px-4 font-oswald py-2"
          onClick={() => {
            router.push("/add");
          }}
        >
          Launch Token
        </button>
        {!primaryWallet?.address ? (
          <DynamicConnectButton>
            <button className="bg-tc py-2 rounded-full px-4 font-oswald">
              Connect Wallet
            </button>
          </DynamicConnectButton>
        ) : (
          <div className="flex items-center gap-3">
            <p>{`${primaryWallet?.address.slice(
              0,
              4
            )}...${primaryWallet?.address.slice(-4)}`}</p>
            <button
              className="bg-tc py-2 rounded-full px-4 font-oswald"
              onClick={() => handleLogOut()}
            >
              Disconnect
            </button>
          </div>
        )}
      </section>
    </main>
  );
}

export default Header;
