"use client";
import { type ReactNode } from "react";
import { Provider } from "react-redux";
import { store } from "../store/store";
import { StarknetWalletConnectors } from "@dynamic-labs/starknet";
import { DynamicContextProvider } from "@dynamic-labs/sdk-react-core";

export function Providers(props: { children: ReactNode }) {
  const Dynamic_KEY = process.env.DYNAMIC_API_KEY;

  return (
    <Provider store={store}>
      <DynamicContextProvider
        settings={{
          environmentId: Dynamic_KEY || "",
          walletConnectors: [StarknetWalletConnectors],
        }}
      >
        {props.children}
      </DynamicContextProvider>
    </Provider>
  );
}
