"use client";
import { createToken } from "@/utils/apis";
import React, { useState } from "react";
import { RxCross2 } from "react-icons/rx";

const Modal = ({ onClose, setSuccess }: any) => {
  const [isClosing, setIsClosing] = useState(false);
  const [name, setName] = useState<string>("");
  const [ticker, setTicker] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [image, setImage] = useState("");
  const [wallet, setWallet] = useState<string>("");

  const inputStyle = `w-full px-4 py-2 mb-4 text-white bg-blue-800/30 border border-blue-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500`;

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      onClose(false);
    }, 500);
  };

  const handleTokenCreate = async () => {
    const response = await createToken({
      name,
      ticker,
      description,
      image,
      wallet,
    });
    if (response?.data.success == true) {
      setIsClosing(true);
      setTimeout(() => {
        onClose(false);
        setSuccess(true);
      });
    }
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const base64String = reader.result?.toString().split(",")[1];
        setImage(base64String || ""); // Use empty string as fallback
      };
      reader.readAsDataURL(file);
    }
  };
  return (
    <div
      className={`fixed inset-0 bg-black bg-opacity-75 flex justify-center items-center z-30  ${
        isClosing
          ? "motion-safe:animate-fade-out"
          : "motion-safe:animate-fade-in"
      }`}
      onClick={handleClose}
    >
      {/* {image && <RenderBase64Image />} */}
      <div
        className={`bg-blue-900/50 relative px-2 py-2 rounded mt-4 shadow-lg max-w-[50rem] w-[30rem] mx-auto z-40 flex flex-col ${
          isClosing
            ? "motion-safe:animate-fade-out"
            : "motion-safe:animate-fade-in"
        }`}
        onClick={(e) => e.stopPropagation()}
      >
        <section className="flex justify-between items-center">
          <p />
          <p className="text-center self-center font-oswald">
            Create a new Coin
          </p>
          <RxCross2
            className=" w-10  h-7 text-tc font-bold cursor-pointer z-[1000]"
            onClick={() => handleClose()}
          />
        </section>

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
            onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
              setDescription(e.target.value);
            }}
          />
          <input
            type="file"
            placeholder="Select logo"
            className={inputStyle}
            onChange={handleImageChange}
          />
          <input
            type="text"
            placeholder="Select logo"
            className={inputStyle}
            value={wallet}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setWallet(e.target.value)
            }
          />
          <button
            className="w-full px-4 py-2 mb-4 text-white bg-green-600 rounded hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500"
            onClick={() => {
              handleTokenCreate();
            }}
          >
            Submit
          </button>
        </main>
        <div className="flex justify-end mr-3">
          <button
            onClick={handleClose}
            className="px-4 py-2 bg-red-500  ml-5 text-white rounded"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default Modal;
