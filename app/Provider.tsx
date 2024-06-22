"use client";
import { useEffect, type ReactNode } from "react";
import { Provider } from "react-redux";
import { store } from "../store/store";
import { StarknetWalletConnectors } from "@dynamic-labs/starknet";
import {
  DynamicContextProvider,
  useUserWallets,
} from "@dynamic-labs/sdk-react-core";

export function Providers(props: { children: ReactNode }) {
  return (
    <Provider store={store}>
      <DynamicContextProvider
        settings={{
          environmentId: "fe350575-4c7c-4297-aff2-9a3de9e5c479",
          walletConnectors: [StarknetWalletConnectors],
        }}
      >
        {props.children}
      </DynamicContextProvider>
    </Provider>
  );
}
