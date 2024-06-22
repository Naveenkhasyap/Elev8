import { configureStore } from "@reduxjs/toolkit";
import { isMainSlice } from "./slices/walletSlice";

export const store = configureStore({
  reducer: {
    MainSlice: isMainSlice.reducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
