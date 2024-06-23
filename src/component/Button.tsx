import React from "react";
import Loading from "./Loading";

function Button({
  isLoading,
  onClick,
  styles,
  disable,
  text,
  loadingText,
}: any) {
  return (
    <button
      className={` ${
        disable && "bg-green-500 cursor-not-allowed"
      } flex items-center rounded-lg w-full justify-center   px-4 py-2 text-white ${
        !isLoading &&
        !disable &&
        " hover:bg-green-400  active:scale-95 transition-all duration-500 ease-out "
      }  ${
        isLoading && !disable ? "bg-green-500 " : "bg-green-500 "
      } ${styles}`}
      disabled={disable || isLoading}
      onClick={onClick}
    >
      {isLoading ? (
        <>
          <Loading
            text={loadingText ? loadingText : "Verifying.... "}
            width={"!w-12"}
          />
        </>
      ) : (
        <span className="font-medium">{text ? text : "Submit"}</span>
      )}
    </button>
  );
}

export default Button;
