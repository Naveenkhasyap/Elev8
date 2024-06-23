import { createSlice } from "@reduxjs/toolkit";

export interface initialStateI {
  isDarkMode: boolean;
}
const initialState = {
  modelOpen: true,
  cryptoCoins: [],
};

export const isMainSlice = createSlice({
  name: "MainSlice",
  initialState,
  reducers: {
    setIsModelOpen: (state, action) => {
      state.modelOpen = action.payload;
    },
    setCryptoCoins: (state, action) => {
      state.cryptoCoins = action.payload;
    },
  },
});
export const { setIsModelOpen, setCryptoCoins } = isMainSlice.actions;
export default isMainSlice;
