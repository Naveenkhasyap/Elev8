import { createSlice } from "@reduxjs/toolkit";

export interface initialStateI {
  isDarkMode: boolean;
}
const initialState = {
  modelOpen: true,
  walletAddress: "",
};

export const isMainSlice = createSlice({
  name: "MainSlice",
  initialState,
  reducers: {
    setIsModelOpen: (state, action) => {
      state.modelOpen = action.payload;
    },
    setWalletAddress: (state, action) => {
      state.walletAddress = action.payload;
    },
  },
});
export const { setIsModelOpen } = isMainSlice.actions;
export default isMainSlice;
