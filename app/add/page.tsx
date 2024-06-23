"use client";
import Button from "@/component/Button";
import Header from "@/component/header/Header";
import Success from "@/component/modal/success";
import { createToken } from "@/utils/apis";
import { useDynamicContext } from "@dynamic-labs/sdk-react-core";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import toast from "react-hot-toast";
import { IoMdArrowRoundBack } from "react-icons/io";

function Page() {
  const router = useRouter();
  const [name, setName] = useState<string>("");
  const [ticker, setTicker] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [image, setImage] = useState("");
  const { primaryWallet } = useDynamicContext();
  const [success, setSuccess] = useState<boolean>(false);
  const [loader, setLoader] = useState<boolean>(false);
  const inputStyle = `w-full px-4 py-2 mb-4 text-white bg-blue-800/30 border border-blue-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500`;

  const handleTokenCreate = async () => {
    setLoader(true);
    if (!primaryWallet?.address) {
      setLoader(false);
      return toast.error("Please Connect your Wallet");
    }

    if (!ticker) {
      setLoader(false);
      return toast.error("Please Enter Ticker");
    }
    if (!description) {
      setLoader(false);
      return toast.error("Please Enter Description");
    }
    if (!name) {
      setLoader(false);
      return toast.error("Please select a Name");
    }
    if (!image) {
      setLoader(false);
      return toast.error("Please select an Logo");
    }
    const response = await createToken({
      name,
      ticker,
      description,
      image,
      wallet: primaryWallet?.address,
    });
    if (response?.data.success == true) {
      setTimeout(() => {}, 1000);
    }
    setLoader(false);
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const base64String = reader.result?.toString().split(",")[1];
        setImage(base64String || "");
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div>
      <Header />
      {success && <Success onClose={setSuccess} router={router} />}
      <main className="flex justify-center items-center h-[93vh] w-full">
        <div
          className={` bg-blue-950/30 relative px-2 py-2 rounded mt-4 shadow-lg max-w-[50rem] w-[30rem] mx-auto z-40 flex flex-col `}
          onClick={(e) => e.stopPropagation()}
        >
          <div className={"flex items-center gap-[30%]"}>
            <IoMdArrowRoundBack
              className="text-3xl cursor-pointer"
              onClick={() => {
                router.back();
              }}
            />
            <p className=" text-center text-xl font-oswald">
              Create a new Coin
            </p>
          </div>

          <main className="py-4 px-3">
            <input
              type="text"
              placeholder="Coin Name"
              className={inputStyle}
              value={name}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setName(e.target.value)
              }
            />
            <input
              type="text"
              placeholder="Ticker"
              className={inputStyle}
              value={ticker}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setTicker(e.target.value)
              }
            />
            <textarea
              placeholder="Description"
              rows={4}
              className={inputStyle}
              value={description}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                setDescription(e.target.value)
              }
            />
            <input
              type="file"
              placeholder="Select logo"
              className={inputStyle}
              onChange={handleImageChange}
            />
            <Button
              text="Submit"
              onClick={() => handleTokenCreate()}
              isLoading={loader}
              loadingText="Processing...."
              disable={false}
            />
          </main>
        </div>
      </main>
    </div>
  );
}

export default Page;
