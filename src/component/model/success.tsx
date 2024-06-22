import React, { useState } from "react";
import { RxCross2 } from "react-icons/rx";

function Success({ dispatch, onClose }: any) {
  const [isClosing, setIsClosing] = useState(false);
  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      dispatch(onClose(false));
    }, 500);
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
      <div
        className={`bg-[#edeef0] text-black relative px-2 py-2 rounded mt-4 shadow-lg max-w-[50rem] w-[30rem] mx-auto z-40 flex flex-col ${
          isClosing
            ? "motion-safe:animate-fade-out"
            : "motion-safe:animate-fade-in"
        }`}
        onClick={(e) => e.stopPropagation()}
      >
        <section className="py-4 px-3 relative">
          <RxCross2
            className=" w-10  h-7 text-tc font-bold cursor-pointer z-[1000] absolute -top-2 -right-4"
            onClick={() => handleClose()}
          />
          <img src="https://ik.imagekit.io/z3vwasz5xwz/kyc-success_cDeeKNCbjV.gif?updatedAt=1714994891805" />
          <p>Your new crypto coin has been launched successfully! 🚀🎉</p>
        </section>
      </div>
    </div>
  );
}

export default Success;
